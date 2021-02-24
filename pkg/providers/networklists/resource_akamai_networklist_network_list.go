package networklists

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html
func resourceNetworkList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkListCreate,
		ReadContext:   resourceNetworkListRead,
		UpdateContext: resourceNetworkListUpdate,
		DeleteContext: resourceNetworkListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					IP,
					Geo,
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Append,
					Replace,
					Remove,
				}, false),
			},
			"uniqueid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "uniqueId",
			},
			"sync_point": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "sync point",
			},
		},
	}
}

func resourceNetworkListCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListCreate")

	createNetworkList := networklists.CreateNetworkListRequest{}

	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createNetworkList.Name = name

	listType, err := tools.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createNetworkList.Type = listType

	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createNetworkList.Description = description

	//createNetworkList.List = tools.SetToStringSlice(d.Get("list").(*schema.Set))

	netlist := d.Get("list").([]interface{})
	nru := make([]string, 0, len(netlist))

	for _, h := range netlist {
		nru = append(nru, h.(string))

	}
	createNetworkList.List = nru

	spcr, errc := client.CreateNetworkList(ctx, createNetworkList)
	if errc != nil {
		logger.Errorf("calling 'createNetworkList': %s", errc.Error())
		return diag.FromErr(errc)
	}

	d.Set("name", spcr.Name)

	d.Set("sync_point", strconv.Itoa(spcr.SyncPoint))

	if err := d.Set("uniqueid", spcr.UniqueID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(spcr.UniqueID)

	return resourceNetworkListRead(ctx, d, m)
}

func resourceNetworkListUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListUpdate")

	updateNetworkList := networklists.UpdateNetworkListRequest{}
	updateNetworkList.Name = d.Id()

	listType, err := tools.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateNetworkList.Type = listType

	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateNetworkList.Description = description

	//	updateNetworkList.List = tools.SetToStringSlice(d.Get("list").(*schema.Set))
	netlist := d.Get("list").([]interface{})
	nru := make([]string, 0, len(netlist))

	for _, h := range netlist {
		nru = append(nru, h.(string))

	}
	updateNetworkList.List = nru

	syncPoint, err := tools.GetIntValue("sync_point", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateNetworkList.SyncPoint = syncPoint

	_, erru := client.UpdateNetworkList(ctx, updateNetworkList)
	if erru != nil {
		logger.Errorf("calling 'updateNetworkList': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceNetworkListRead(ctx, d, m)
}

func resourceNetworkListDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListRemove")

	removeNetworkList := networklists.RemoveNetworkListRequest{}
	removeNetworkList.Name = d.Id()
	_, errd := client.RemoveNetworkList(ctx, removeNetworkList)
	if errd != nil {
		logger.Errorf("calling 'removeNetworkList': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceNetworkListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLISTs", "resourceNetworkListRead")

	getNetworkList := networklists.GetNetworkListRequest{}
	getNetworkList.Name = d.Id()

	networklist, err := client.GetNetworkList(ctx, getNetworkList)
	if err != nil {
		logger.Errorf("calling 'getNetworkList': %s", err.Error())
		return diag.FromErr(err)
	}

	logger.Errorf("calling 'getNetworkList': SYNC POINT %d", networklist.SyncPoint)
	/*
		syncPoint, errconv := strconv.Atoi(networklist.SyncPoint)
		if errconv != nil {
			return diag.FromErr(errconv)
		}
	*/
	if err := d.Set("sync_point", networklist.SyncPoint); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	//d.Set("sync_point", strconv.Itoa(networklist.SyncPoint))

	d.SetId(networklist.UniqueID)

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
