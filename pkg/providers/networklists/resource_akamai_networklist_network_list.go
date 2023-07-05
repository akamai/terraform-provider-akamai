package networklists

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// network_lists v2
//
// https://techdocs.akamai.com/network-lists/reference/api
func resourceNetworkList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkListCreate,
		ReadContext:   resourceNetworkListRead,
		UpdateContext: resourceNetworkListUpdate,
		DeleteContext: resourceNetworkListDelete,
		CustomizeDiff: customdiff.All(
			verifyContractGroupUnchanged,
			markSyncPointComputedIfListModified,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name to be assigned to the network list",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the network list; must be either 'IP' or 'GEO'",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					IP,
					Geo,
				}, false)),
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of the network list",
			},
			"list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of IP addresses or locations to be included in the list, added to an existing list, or removed from an existing list",
			},
			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A string specifying the interpretation of the `list` parameter. Must be 'APPEND', 'REPLACE', or 'REMOVE'",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					Append,
					Replace,
					Remove,
				}, false)),
			},
			"uniqueid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "unique ID",
			},
			"network_list_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "network list ID",
			},
			"sync_point": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "sync point",
			},
			"contract_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "contract ID",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "group ID",
			},
		},
	}
}

func resourceNetworkListCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListCreate")

	attrs, err := getAttributes(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createNetworkList := networklists.CreateNetworkListRequest{
		Name:        attrs.name,
		Type:        attrs.listType,
		Description: attrs.description,
		ContractID:  attrs.contractID,
		GroupID:     attrs.groupID,
	}

	if len(attrs.contractID) > 0 || attrs.groupID > 0 {
		if len(attrs.contractID) == 0 || attrs.groupID == 0 {
			return diag.Errorf("If either a contract_id or group_id is provided, both must be provided")
		}
	}

	getNetworkLists := networklists.GetNetworkListsRequest{}
	getNetworkLists.Name = attrs.name
	getNetworkLists.Type = attrs.listType

	networklists, err := client.GetNetworkLists(ctx, getNetworkLists)
	if err != nil {
		logger.Errorf("calling 'getNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	netlist := d.Get("list").(*schema.Set)
	networkListElements := make([]string, 0, len(netlist.List()))

	for _, h := range netlist.List() {
		networkListElements = append(networkListElements, strings.ToLower(h.(string)))
	}

	finallist := make([]string, 0, len(netlist.List()))

	switch attrs.mode {
	case Remove:
		for _, hl := range netlist.List() {
			for _, h := range networklists.NetworkLists {

				if h.Name == hl.(string) {
					finallist = append(finallist, strings.ToLower(h.Name))
				}
			}
		}
	case Append:
		var oneShot bool

		for _, h := range networklists.NetworkLists {
			finallist = appendIfMissing(finallist, strings.ToLower(h.Name))
			for _, hl := range netlist.List() {
				finallist = appendIfMissing(finallist, strings.ToLower(hl.(string)))
			}
			oneShot = true
		}

		if !oneShot {
			finallist = networkListElements
		}

	default:
		finallist = networkListElements
	}

	createNetworkList.List = finallist

	spcr, err := client.CreateNetworkList(ctx, createNetworkList)
	if err != nil {
		logger.Errorf("calling 'createNetworkList': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("name", spcr.Name); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("sync_point", spcr.SyncPoint); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("uniqueid", spcr.UniqueID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("network_list_id", spcr.UniqueID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("mode", attrs.mode); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("contract_id", attrs.contractID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("group_id", attrs.groupID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(spcr.UniqueID)

	return resourceNetworkListRead(ctx, d, m)
}

func resourceNetworkListUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListUpdate")

	attrs, err := getAttributes(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateNetworkList := networklists.UpdateNetworkListRequest{
		Name:        attrs.name,
		Type:        attrs.listType,
		Description: attrs.description,
		ContractID:  attrs.contractID,
		GroupID:     attrs.groupID,
		UniqueID:    d.Id(),
	}

	if len(attrs.contractID) > 0 || attrs.groupID > 0 {
		if len(attrs.contractID) == 0 || attrs.groupID == 0 {
			return diag.Errorf("If either a contract_id or group_id is provided, both must be provided")
		}
	}

	listRequest := networklists.GetNetworkListRequest{}
	listRequest.UniqueID = d.Id()

	networkLists, err := client.GetNetworkList(ctx, listRequest)
	if err != nil {
		logger.Errorf("calling 'getNetworkList': %s", err.Error())
		return diag.FromErr(err)
	}

	netlist := d.Get("list").(*schema.Set)
	nru := make([]string, 0, len(netlist.List()))

	for _, h := range netlist.List() {
		nru = append(nru, strings.ToLower(h.(string)))
	}

	finallist := make([]string, 0, len(netlist.List()))

	switch attrs.mode {
	case Remove:
		for _, hl := range netlist.List() {

			for idx, h := range networkLists.List {
				if strings.EqualFold(h, hl.(string)) {
					networkLists.List = RemoveIndex(networkLists.List, idx)
				}
			}
		}
		finallist = networkLists.List

	case Append:
		for _, h := range networkLists.List {
			finallist = append(finallist, strings.ToLower(h))
		}
		for _, hl := range netlist.List() {
			finallist = appendIfMissing(finallist, strings.ToLower(hl.(string)))
		}
	case Replace:
		finallist = nru
	default:
		finallist = nru
	}

	updateNetworkList.List = finallist

	syncPoint, err := tf.GetIntValue("sync_point", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateNetworkList.SyncPoint = syncPoint

	_, err = client.UpdateNetworkList(ctx, updateNetworkList)
	if err != nil {
		logger.Errorf("calling 'updateNetworkList': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("contract_id", attrs.contractID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("group_id", attrs.groupID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return resourceNetworkListRead(ctx, d, m)
}

func resourceNetworkListDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListRemove")

	removeNetworkList := networklists.RemoveNetworkListRequest{}
	removeNetworkList.UniqueID = d.Id()
	_, errd := client.RemoveNetworkList(ctx, removeNetworkList)
	if errd != nil {
		logger.Errorf("calling 'removeNetworkList': %s", errd.Error())
	}

	d.SetId("")

	return nil
}

func resourceNetworkListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListRead")

	getNetworkList := networklists.GetNetworkListRequest{}
	getNetworkList.UniqueID = d.Id()

	mode, err := tf.GetStringValue("mode", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	var detectCase string

	netlist := d.Get("list").(*schema.Set)
	for _, hl := range netlist.List() {
		if hl.(string) == strings.ToLower(hl.(string)) {
			detectCase = "LOWER"
		} else {
			detectCase = "UPPER"
		}
	}
	finalldata := make([]string, 0, len(netlist.List()))

	networklist, err := client.GetNetworkList(ctx, getNetworkList)
	if err != nil {
		logger.Errorf("calling 'getNetworkList': %s", err.Error())
		return diag.FromErr(err)
	}
	finalldata = resolveNetworkList(mode, netlist, networklist, finalldata)
	sort.Strings(finalldata)

	if err := d.Set("sync_point", networklist.SyncPoint); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("name", networklist.Name); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("type", networklist.Type); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if detectCase == "LOWER" {
		for index, value := range finalldata {
			finalldata[index] = strings.ToLower(value)
		}
	} else {
		for index, value := range finalldata {
			finalldata[index] = strings.ToUpper(value)
		}
	}

	if err := d.Set("description", networklist.Description); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("list", nil); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("list", finalldata); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if mode == "" {
		mode = "REPLACE"
	}

	if err := d.Set("mode", mode); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("uniqueid", networklist.UniqueID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("network_list_id", networklist.UniqueID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(networklist.UniqueID)

	return nil
}

// networkListAttrs represent networkList attributes
type networkListAttrs struct {
	name        string
	listType    string
	description string
	contractID  string
	groupID     int
	mode        string
}

// getAttributes fetches some attributes grouped together
func getAttributes(d *schema.ResourceData) (*networkListAttrs, error) {
	name, err := tf.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	listType, err := tf.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	description, err := tf.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	groupID, err := tf.GetIntValue("group_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	mode, err := tf.GetStringValue("mode", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}

	return &networkListAttrs{
		name:        name,
		listType:    listType,
		description: description,
		contractID:  contractID,
		groupID:     groupID,
		mode:        mode,
	}, nil
}

// resolveNetworkList resolves network list based on its mode
func resolveNetworkList(mode string, netlist *schema.Set, networklist *networklists.GetNetworkListResponse, finalldata []string) []string {
	switch mode {
	case Remove:
		for _, hl := range netlist.List() {
			for _, h := range networklist.List {

				if strings.EqualFold(h, hl.(string)) {
					finalldata = append(finalldata, strings.ToLower(h))
				}
			}
		}

		if len(finalldata) == 0 {
			for _, hl := range netlist.List() {
				finalldata = append(finalldata, strings.ToLower(hl.(string)))
			}
		}

	case Append:
		for _, h := range networklist.List {

			for _, hl := range netlist.List() {
				if strings.EqualFold(h, hl.(string)) {
					finalldata = append(finalldata, strings.ToLower(h))
				}
			}
		}
	case Replace:
		for _, h := range networklist.List {
			finalldata = append(finalldata, strings.ToLower(h))
		}
	default:
		for _, h := range networklist.List {
			finalldata = append(finalldata, strings.ToLower(h))
		}
	}
	return finalldata
}

func appendIfMissing(slice []string, s string) []string {
	for _, element := range slice {
		if element == s {
			return slice
		}
	}
	return append(slice, s)
}

// RemoveIndex removes an element from the slice and returns it
func RemoveIndex(hl []string, index int) []string {
	return append(hl[:index], hl[index+1:]...)
}

// verifyContractGroupUnchanged compares the configuration's value for the contract_id and group_id with the resource's
// value specified in the resources's ID, to ensure that the user has not inadvertently modified the configuration's
// value; any such modifications indicate an incorrect understanding of the Update operation.
func verifyContractGroupUnchanged(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("NETWORKLIST", "VerifyContractGroupUnchanged")

	if d.HasChange("contract_id") {
		oldContract, newContract := d.GetChange("contract_id")
		oldvalue := oldContract.(string)
		newvalue := newContract.(string)
		if len(oldvalue) > 0 {
			logger.Errorf("%s value %s specified in configuration differs from resource ID's value %s", "contract_id", newvalue, oldvalue)
			return fmt.Errorf("%s value %s specified in configuration differs from resource ID's value %s", "contract_id", newvalue, oldvalue)
		}
	}

	if d.HasChange("group_id") {
		oldGroup, newGroup := d.GetChange("group_id")
		oldvalue := oldGroup.(int)
		newvalue := newGroup.(int)
		if oldvalue > 0 {
			logger.Errorf("%s value %d specified in configuration differs from resource ID's value %d", "group_id", newvalue, oldvalue)
			return fmt.Errorf("%s value %d specified in configuration differs from resource ID's value %d", "group_id", newvalue, oldvalue)
		}
	}

	return nil
}

// markSyncPointComputedIfListModified sets 'sync_point' field as new computed
// if a new version of network list is expected to be created.
func markSyncPointComputedIfListModified(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("NETWORKLIST", "MarkSyncPointComputedIfListModified")
	if d.HasChange("list") {
		logger.Debug("setting sync_point as new computed")
		return d.SetNewComputed("sync_point")
	}
	return nil
}

// Append Replace Remove mode flags
const (
	Append  = "APPEND"
	Replace = "REPLACE"
	Remove  = "REMOVE"

	IP  = "IP"
	Geo = "GEO"
)
