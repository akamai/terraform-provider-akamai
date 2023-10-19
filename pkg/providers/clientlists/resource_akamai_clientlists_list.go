package clientlists

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceClientList() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceClientListRead,
		CreateContext: resourceClientListCreate,
		UpdateContext: resourceClientListUpdate,
		DeleteContext: resourceClientListDelete,
		CustomizeDiff: customdiff.All(
			markVersionComputedIfListModified,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the client list.",
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      fmt.Sprintf("The type of the client list. Valid types: %s", getValidListTypes()),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(getValidListTypes(), false)),
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The client list notes.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The client list tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Contract ID for which client list is assigned.",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Group ID for which client list is assigned.",
			},
			"list_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the client list.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The current version of the client list.",
			},
			"items_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of items that a client list contains.",
			},
			"items": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of items containing item information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the item. (i.e. IP address, AS Number, GEO, ...etc)",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A description of the item.",
							Default:     "",
						},
						"tags": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The item tags.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"expiration_date": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The item expiration date.",
							Default:     "",
						},
					},
				},
			},
		},
	}
}

func resourceClientListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "resourceClientListRead")
	logger.Debug("Reading client list")

	getCLientList := clientlists.GetClientListRequest{
		ListID:       d.Id(),
		IncludeItems: true,
	}

	list, err := client.GetClientList(ctx, getCLientList)
	if e, ok := err.(*clientlists.Error); ok && e.StatusCode == http.StatusNotFound || list.Deprecated {
		d.SetId("")
		return nil
	} else if err != nil {
		logger.Errorf("calling 'getClientList' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	items := make([]interface{}, 0, len(list.Items))
	for _, v := range list.Items {
		i := map[string]interface{}{
			"value":           v.Value,
			"description":     v.Description,
			"expiration_date": v.ExpirationDate,
			"tags":            v.Tags,
		}

		items = append(items, i)
	}

	fields := map[string]interface{}{
		"contract_id": list.ContractID,
		"group_id":    list.GroupID,
		"name":        list.Name,
		"type":        list.Type,
		"notes":       list.Notes,
		"tags":        list.Tags,
		"list_id":     list.ListID,
		"version":     list.Version,
		"items_count": list.ItemsCount,
		"items":       items,
	}

	if err = tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceClientListCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "resourceClientListCreate")
	logger.Debug("Creating client list")

	if err := validateItemsUniqueness(d); err != nil {
		return diag.FromErr(err)
	}

	listAttrs, err := getClientListAttr(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createCLientListRequest := clientlists.CreateClientListRequest{
		Name:       listAttrs.Name,
		Type:       clientlists.ClientListType(listAttrs.ListType),
		Notes:      listAttrs.Notes,
		Tags:       listAttrs.Tags,
		ContractID: listAttrs.ContractID,
		GroupID:    listAttrs.GroupID,
		Items:      listAttrs.Items,
	}

	list, err := client.CreateClientList(ctx, createCLientListRequest)
	if err != nil {
		logger.Errorf("calling 'createClientList' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(list.ListID)

	return resourceClientListRead(ctx, d, m)
}

func resourceClientListUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "resourceClientListUpdate")
	logger.Debug("Updating client list")

	if err := validateItemsUniqueness(d); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("items") {
		getListRes, err := client.GetClientList(ctx, clientlists.GetClientListRequest{
			ListID:       d.Id(),
			IncludeItems: true,
		})
		if err != nil {
			logger.Errorf("calling 'getClientList' failed: %s", err.Error())
			return diag.FromErr(err)
		}

		itemsUpdateReq, err := getListItemsUpdateReq(*getListRes, d)
		if err != nil {
			logger.Errorf("constructing items update request failed: %s", err.Error())
			return diag.FromErr(err)
		}

		_, err = client.UpdateClientListItems(ctx, *itemsUpdateReq)
		if err != nil {
			logger.Errorf("calling 'UpdateClientListItems' failed: %s", err.Error())
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("name", "notes", "tags") {
		listAttrs, err := getClientListAttr(d)
		if err != nil {
			return diag.FromErr(err)
		}

		updateClientListRequest := clientlists.UpdateClientListRequest{
			ListID: d.Id(),
			UpdateClientList: clientlists.UpdateClientList{
				Name:  listAttrs.Name,
				Notes: listAttrs.Notes,
				Tags:  listAttrs.Tags,
			},
		}

		_, err = client.UpdateClientList(ctx, updateClientListRequest)
		if err != nil {
			logger.Errorf("calling 'updateClientList' failed: %s", err.Error())
			return diag.FromErr(err)
		}
	}

	return resourceClientListRead(ctx, d, m)
}

func resourceClientListDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "resourceClientListDelete")
	logger.Debug("Deleting client list")

	deleteCLientListRequest := clientlists.DeleteClientListRequest{
		ListID: d.Id(),
	}

	err := client.DeleteClientList(ctx, deleteCLientListRequest)
	if err != nil {
		logger.Errorf("calling 'deleteClientList' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

type clientListAttrs struct {
	Name       string
	ListType   string
	Notes      string
	Tags       []string
	ContractID string
	GroupID    int64
	Items      []clientlists.ListItemPayload
}

func getClientListAttr(d *schema.ResourceData) (*clientListAttrs, error) {
	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return nil, err
	}
	listType, err := tf.GetStringValue("type", d)
	if err != nil {
		return nil, err
	}
	notes, err := tf.GetStringValue("notes", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}

	ts, err := tf.GetSetValue("tags", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	tags := make([]string, 0, len(ts.List()))
	for _, t := range ts.List() {
		tags = append(tags, t.(string))
	}

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return nil, err
	}
	groupID, err := tf.GetIntValue("group_id", d)
	if err != nil {
		return nil, err
	}

	itemsSet, err := tf.GetSetValue("items", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	items := make([]clientlists.ListItemPayload, 0, itemsSet.Len())
	for _, v := range itemsSet.List() {
		itemMap := v.(map[string]interface{})

		t := itemMap["tags"].(*schema.Set)
		items = append(items, clientlists.ListItemPayload{
			Value:          itemMap["value"].(string),
			Description:    itemMap["description"].(string),
			Tags:           tf.SetToStringSlice(t),
			ExpirationDate: itemMap["expiration_date"].(string),
		})
	}

	return &clientListAttrs{
		Name:       name,
		ListType:   listType,
		Notes:      notes,
		Tags:       tags,
		ContractID: contractID,
		GroupID:    int64(groupID),
		Items:      items,
	}, nil
}

func getListItemsUpdateReq(list clientlists.GetClientListResponse, d *schema.ResourceData) (*clientlists.UpdateClientListItemsRequest, error) {
	itemsSet, err := tf.GetSetValue("items", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	// Map of item value to ListItemPayload representing items in the config
	configItemsMap := make(map[string]clientlists.ListItemPayload)
	for _, v := range itemsSet.List() {
		itemMap := v.(map[string]interface{})

		configItemsMap[itemMap["value"].(string)] = clientlists.ListItemPayload{
			Value:          itemMap["value"].(string),
			Description:    itemMap["description"].(string),
			Tags:           tf.SetToStringSlice(itemMap["tags"].(*schema.Set)),
			ExpirationDate: itemMap["expiration_date"].(string),
		}
	}

	// Map of item value to item representing list of item in remote state
	listItemsMap := make(map[string]clientlists.ListItemContent)
	for _, v := range list.Items {
		listItemsMap[v.Value] = v
	}

	res := &clientlists.UpdateClientListItemsRequest{
		ListID: list.ListID,
		UpdateClientListItems: clientlists.UpdateClientListItems{
			Append: []clientlists.ListItemPayload{},
			Update: []clientlists.ListItemPayload{},
			Delete: []clientlists.ListItemPayload{},
		},
	}

	for _, configItem := range configItemsMap {
		if listItem, ok := listItemsMap[configItem.Value]; ok {
			if shouldUpdateItem(configItem, listItem) {
				res.UpdateClientListItems.Update = append(res.UpdateClientListItems.Update, configItem)
			}
		} else {
			res.UpdateClientListItems.Append = append(res.UpdateClientListItems.Append, configItem)
		}
	}

	for _, listItem := range listItemsMap {
		if _, ok := configItemsMap[listItem.Value]; !ok {
			res.UpdateClientListItems.Delete = append(res.UpdateClientListItems.Delete, clientlists.ListItemPayload{
				Value: listItem.Value,
			})
		}
	}

	return res, nil
}

func shouldUpdateItem(a clientlists.ListItemPayload, b clientlists.ListItemContent) bool {
	if a.Value == b.Value &&
		a.Description == b.Description &&
		a.ExpirationDate == b.ExpirationDate &&
		isEqualTags(a.Tags, b.Tags) {
		return false
	}
	return true
}

func isEqualTags(t1, t2 []string) bool {
	if len(t1) != len(t2) {
		return false
	}

	a := make([]string, len(t1))
	b := make([]string, len(t2))

	copy(a, t1)
	copy(b, t2)

	sort.Strings(a)
	sort.Strings(b)

	return reflect.DeepEqual(a, b)
}

func validateItemsUniqueness(d *schema.ResourceData) error {
	itemsSet, err := tf.GetSetValue("items", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	values := map[string]interface{}{}
	for _, v := range itemsSet.List() {
		itemMap := v.(map[string]interface{})
		value := itemMap["value"].(string)

		if _, ok := values[value]; ok {
			return fmt.Errorf("'Items' collection contains duplicate values for 'value' field. Duplicate value: %s", value)
		}
		values[value] = itemMap
	}

	return nil
}

// markVersionComputedIfListModified sets 'version' field as new computed
// if a new version of client list is expected to be created.
func markVersionComputedIfListModified(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("CLIENTLIST", "markVersionComputedIfListModified")

	itemsHasChange := d.HasChange("items")
	oldItems, newItems := d.GetChange("items")

	isVersionUpdateRequired, err := isVersionUpdateRequired(oldItems, newItems)
	if err != nil {
		return err
	}

	if itemsHasChange && isVersionUpdateRequired {
		logger.Debug("setting version as new computed")
		if err := d.SetNewComputed("version"); err != nil {
			return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
		}
	}

	return nil
}

// isVersionUpdateRequired determines if list version update is required based on items changes
func isVersionUpdateRequired(oldValue, newValue interface{}) (bool, error) {
	if oldValue == nil || newValue == nil {
		return oldValue != newValue, nil
	}

	o, ok := oldValue.(*schema.Set)
	if !ok {
		return false, fmt.Errorf("'items' old value is not of type schema.Set")
	}
	n, ok := newValue.(*schema.Set)
	if !ok {
		return false, fmt.Errorf("'items' new value is not of type schema.Set")
	}

	if o.Len() != n.Len() {
		return true, nil
	}

	oldMap := mapExpirationDateToValue(o)
	newMap := mapExpirationDateToValue(n)

	for newValue, newExpDate := range newMap {
		// if value does not exist or expiration dates are different,
		// then version update is required
		if oldExpDate, ok := oldMap[newValue]; !ok || oldExpDate != newExpDate {
			return true, nil
		}
	}

	return false, nil
}

func mapExpirationDateToValue(items *schema.Set) map[string]string {
	res := make(map[string]string)

	for _, v := range items.List() {
		item := v.(map[string]interface{})
		res[item["value"].(string)] = item["expiration_date"].(string)
	}

	return res
}
