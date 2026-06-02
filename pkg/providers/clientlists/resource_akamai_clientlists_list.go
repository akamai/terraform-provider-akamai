package clientlists

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/google/uuid"
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
			validateItems,
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
				Set:         itemsHashFunc,
				Description: "Set of items containing item information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value of the item (e.g. IP address, AS Number, GEO, domain, TLS fingerprint, file hash, user ID). Not applicable for REQUEST_HEADER_NAME_VALUE list type.",
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key of the item (e.g. request header name). Applicable only for REQUEST_HEADER_NAME_VALUE list type.",
						},
						"values": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Values of the item (e.g. request header name values). Applicable only for REQUEST_HEADER_NAME_VALUE list type.",
							Elem:        &schema.Schema{Type: schema.TypeString},
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
	metaInfo := meta.Must(m)
	client := inst.Client(metaInfo)
	logger := metaInfo.Log("CLIENTLIST", "resourceClientListRead")
	logger.Debug("Reading client list")

	getClientListReq := clientlists.GetClientListRequest{
		ListID:       d.Id(),
		IncludeItems: true,
	}

	list, err := client.GetClientList(ctx, getClientListReq)
	var clientListErr *clientlists.Error
	if errors.As(err, &clientListErr) && clientListErr.StatusCode == http.StatusNotFound || (list != nil && list.Deprecated) {
		d.SetId("")
		return nil
	} else if err != nil {
		logger.Errorf("calling 'getClientList' failed: %s", err.Error())
		return diag.FromErr(err)
	}
	if list == nil {
		return diag.Errorf("unexpected nil response from GetClientList")
	}

	if list.Type == clientlists.USER {
		getClientListItems := clientlists.GetClientListItemsRequest{
			ListID: list.ListID,
		}
		items, err := client.GetClientListItems(ctx, getClientListItems)
		if err != nil {
			logger.Errorf("calling 'GetClientListItems': %s", err.Error())
			return diag.FromErr(err)
		} else if len(items.Items) > 0 {
			list.Items = items.Items
		}
	}

	items, diags := extractItems(ctx, d, client, list)
	if diags.HasError() {
		return diags
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
	metaInfo := meta.Must(m)
	client := inst.Client(metaInfo)
	logger := metaInfo.Log("CLIENTLIST", "resourceClientListCreate")
	logger.Debug("Creating client list")

	diags := validateItemsUniqueness(ctx, d, client)
	if diags.HasError() {
		return diags
	}

	listAttrs, err := getClientListAttr(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createClientListRequest := clientlists.CreateClientListRequest{
		Name:       listAttrs.Name,
		Type:       clientlists.ClientListType(listAttrs.ListType),
		Notes:      listAttrs.Notes,
		Tags:       listAttrs.Tags,
		ContractID: listAttrs.ContractID,
		GroupID:    listAttrs.GroupID,
		Items:      listAttrs.Items,
	}

	list, err := client.CreateClientList(ctx, createClientListRequest)
	if err != nil {
		logger.Errorf("calling 'createClientList' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(list.ListID)

	return resourceClientListRead(ctx, d, m)
}

func resourceClientListUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	metaInfo := meta.Must(m)
	client := inst.Client(metaInfo)
	logger := metaInfo.Log("CLIENTLIST", "resourceClientListUpdate")
	logger.Debug("Updating client list")

	diags := validateItemsUniqueness(ctx, d, client)
	if diags.HasError() {
		return diags
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
	metaInfo := meta.Must(m)
	client := inst.Client(metaInfo)
	logger := metaInfo.Log("CLIENTLIST", "resourceClientListDelete")
	logger.Debug("Deleting client list")

	deleteClientListRequest := clientlists.DeleteClientListRequest{
		ListID: d.Id(),
	}

	err := client.DeleteClientList(ctx, deleteClientListRequest)
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
		items = append(items, buildItemPayload(itemMap, listType))
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
	listType, err := tf.GetStringValue("type", d)
	if err != nil {
		return nil, err
	}
	isKeyValuesItem := isKeyValuesItem(listType)

	configItemsMap := buildConfigItemsMap(itemsSet, listType, isKeyValuesItem)
	listItemsMap := buildListItemsMap(list.Items, isKeyValuesItem)

	res := &clientlists.UpdateClientListItemsRequest{
		ListID: list.ListID,
		UpdateClientListItems: clientlists.UpdateClientListItems{
			Append: []clientlists.ListItemPayload{},
			Update: []clientlists.ListItemPayload{},
			Delete: []clientlists.ListItemPayload{},
		},
	}

	for id, configItem := range configItemsMap {
		if listItem, ok := listItemsMap[id]; ok {
			if shouldUpdateItem(configItem, listItem, isKeyValuesItem) {
				res.Update = append(res.Update, configItem)
			}
		} else {
			res.Append = append(res.Append, configItem)
		}
	}

	for id, listItem := range listItemsMap {
		if _, ok := configItemsMap[id]; !ok {
			res.Delete = append(res.Delete, toDeletePayload(listItem, isKeyValuesItem))
		}
	}

	return res, nil
}

func buildConfigItemsMap(itemsSet *schema.Set, listType string, isRequestHeaderNameValue bool) map[string]clientlists.ListItemPayload {
	result := make(map[string]clientlists.ListItemPayload, itemsSet.Len())
	for _, v := range itemsSet.List() {
		itemMap := v.(map[string]interface{})
		payload := buildItemPayload(itemMap, listType)
		result[configItemID(itemMap, isRequestHeaderNameValue)] = payload
	}
	return result
}

func buildListItemsMap(items []clientlists.ListItemContent, isRequestHeaderNameValue bool) map[string]clientlists.ListItemContent {
	result := make(map[string]clientlists.ListItemContent, len(items))
	for _, item := range items {
		result[listItemID(item, isRequestHeaderNameValue)] = item
	}
	return result
}

func configItemID(item map[string]interface{}, isRequestHeaderNameValue bool) string {
	if isRequestHeaderNameValue {
		return item["key"].(string)
	}
	return item["value"].(string)
}

func listItemID(item clientlists.ListItemContent, isRequestHeaderNameValue bool) string {
	if isRequestHeaderNameValue {
		return item.Key
	}
	return item.Value
}

func toDeletePayload(item clientlists.ListItemContent, isRequestHeaderNameValue bool) clientlists.ListItemPayload {
	if isRequestHeaderNameValue {
		return clientlists.ListItemPayload{Key: item.Key}
	}
	return clientlists.ListItemPayload{Value: item.Value}
}

func shouldUpdateItem(a clientlists.ListItemPayload, b clientlists.ListItemContent, isKeyValuesItem bool) bool {
	if isKeyValuesItem {
		return a.Key != b.Key ||
			!isEqual(a.Values, b.Values) ||
			a.Description != b.Description ||
			a.ExpirationDate != b.ExpirationDate ||
			!isEqual(a.Tags, b.Tags)
	}
	return a.Value != b.Value ||
		a.Description != b.Description ||
		a.ExpirationDate != b.ExpirationDate ||
		!isEqual(a.Tags, b.Tags)
}

func isEqual(t1, t2 []string) bool {
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

func validateItemsUniqueness(ctx context.Context, d *schema.ResourceData, client clientlists.ClientLists) diag.Diagnostics {
	itemsSet, err := tf.GetSetValue("items", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	listType := d.Get("type").(string)
	translatedUsernames, diags := translateUsernames(ctx, d, client)
	if diags.HasError() {
		return diags
	}

	values := map[string]interface{}{}
	for _, v := range itemsSet.List() {
		itemMap := v.(map[string]interface{})

		var value, duplicateFieldName string
		if isKeyValuesItem(listType) {
			value = itemMap["key"].(string)
			duplicateFieldName = "key"
		} else {
			value = itemMap["value"].(string)
			duplicateFieldName = "value"
		}
		originalValue := value

		if listType == string(clientlists.USER) && len(translatedUsernames) > 0 {
			if userID := translatedUsernames[value]; userID != "" {
				value = userID
			}
		}

		if _, ok := values[value]; ok {
			return diag.FromErr(fmt.Errorf("'Items' collection contains duplicate values for '%s' field. Duplicate value: %s", duplicateFieldName, originalValue))
		}
		values[value] = itemMap
	}

	return nil
}

func validateItems(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	listType := d.Get("type").(string)
	isKeyValuesItem := isKeyValuesItem(listType)

	var errs []error
	for _, v := range d.Get("items").(*schema.Set).List() {
		item := v.(map[string]interface{})
		if isKeyValuesItem {
			errs = append(errs, validateKeyValuesItem(item))
		} else {
			errs = append(errs, validateValueItem(item))
		}
	}

	return errors.Join(errs...)
}

func validateKeyValuesItem(item map[string]interface{}) error {
	key := item["key"].(string)
	value := item["value"].(string)
	valuesSet := item["values"].(*schema.Set)

	var errs []error
	if key == "" {
		errs = append(errs, fmt.Errorf("invalid item: missing required field 'key'"))
	}
	if valuesSet.Len() == 0 {
		errs = append(errs, fmt.Errorf("invalid item: missing required field 'values'"))
	}
	if value != "" {
		errs = append(errs, fmt.Errorf("invalid item: unsupported field 'value'"))
	}
	return errors.Join(errs...)
}

func validateValueItem(item map[string]interface{}) error {
	key := item["key"].(string)
	value := item["value"].(string)
	valuesSet := item["values"].(*schema.Set)

	var errs []error
	if value == "" {
		errs = append(errs, fmt.Errorf("invalid item: missing required field 'value'"))
	}
	if key != "" {
		errs = append(errs, fmt.Errorf("invalid item: unsupported field 'key'"))
	}
	if valuesSet.Len() > 0 {
		errs = append(errs, fmt.Errorf("invalid item: unsupported field 'values'"))
	}
	return errors.Join(errs...)
}

func markVersionComputedIfListModified(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	metaInfo := meta.Must(m)
	logger := metaInfo.Log("CLIENTLIST", "markVersionComputedIfListModified")

	itemsHasChange := d.HasChange("items")
	oldItems, newItems := d.GetChange("items")
	isKeyValues := isKeyValuesItem(d.Get("type").(string))

	isVersionUpdateRequired, err := isVersionUpdateRequired(oldItems, newItems, isKeyValues)
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
func isVersionUpdateRequired(oldValue, newValue interface{}, isKeyValues bool) (bool, error) {
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

	oldMap := buildVersionComparisonMap(o, isKeyValues)
	newMap := buildVersionComparisonMap(n, isKeyValues)

	for id, newSignature := range newMap {
		// if item does not exist (new key for key-value items, new value for value items)
		// or signature is different (expiration_date or values changed),
		// then version update is required
		if oldSignature, ok := oldMap[id]; !ok || oldSignature != newSignature {
			return true, nil
		}
	}

	return false, nil
}

// buildVersionComparisonMap creates a map of item ID to version-affecting fields signature.
// This map is used to determine if list version update is required by comparing old vs new items.
//
// For key-value items (e.g. REQUEST_HEADER_NAME_VALUE):
//   - ID: item's "key" field (e.g., "User-Agent")
//   - Signature: "expiration_date|sorted_values" (e.g., "2026-12-31|Chrome,Mozilla")
//
// For value items (IP, GEO, ASN, etc.):
//   - ID: item's "value" field (e.g., "1.2.3.4")
//   - Signature: "expiration_date" (e.g., "2026-12-31")
//
// The signature includes only fields that, when changed, should trigger a list version update.
// Changes to description or tags do not affect the version.
func buildVersionComparisonMap(items *schema.Set, isKeyValues bool) map[string]string {
	res := make(map[string]string, items.Len())
	for _, v := range items.List() {
		item := v.(map[string]interface{})
		var id, signature string
		if isKeyValues {
			id, _ = item["key"].(string)
			// Include sorted values to detect changes in values field for key-value items (e.g. REQUEST_HEADER_NAME_VALUE)
			values := tf.SetToStringSlice(item["values"].(*schema.Set))
			sort.Strings(values)
			signature = fmt.Sprintf("%s|%s", item["expiration_date"].(string), strings.Join(values, ","))
		} else {
			id, _ = item["value"].(string)
			signature = item["expiration_date"].(string)
		}
		res[id] = signature
	}

	return res
}

func createTranslateUsernamesRequest(d *schema.ResourceData) (clientlists.TranslateUsernamesRequest, error) {
	itemsSet, err := tf.GetSetValue("items", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	usernames := make(clientlists.TranslateUsernamesRequest, 0, itemsSet.Len())
	for _, v := range itemsSet.List() {
		itemMap := v.(map[string]interface{})
		if isUserID(itemMap["value"].(string)) {
			continue
		}
		usernames = append(usernames, itemMap["value"].(string))
	}
	return usernames, nil
}

func translateValueToUsername(m clientlists.TranslateUsernamesResponse, value string) string {
	for k, v := range m {
		if v == value {
			return k
		}
	}
	return ""
}

func isUserID(s string) bool {
	u, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return u.Version() == uuid.Version(4)
}

func shouldUseUsername(d *schema.ResourceData, value string) bool {
	items := d.Get("items").(*schema.Set).List()

	if len(items) == 0 {
		return true
	}

	for _, item := range items {
		v := item.(map[string]interface{})["value"].(string)
		if value == v && !isUserID(v) {
			return true
		}
	}
	return false
}

func extractItems(ctx context.Context, d *schema.ResourceData, client clientlists.ClientLists, list *clientlists.GetClientListResponse) ([]interface{}, diag.Diagnostics) {
	items := make([]interface{}, 0, len(list.Items))
	translatedUsernames, diags := translateUsernames(ctx, d, client)
	if diags.HasError() {
		return nil, diags
	}

	for _, v := range list.Items {
		if list.Type == clientlists.USER {
			if v.Username != "" && shouldUseUsername(d, v.Username) {
				v.Value = v.Username
			} else if len(translatedUsernames) > 0 {
				username := translateValueToUsername(translatedUsernames, v.Value)
				if username != "" && shouldUseUsername(d, username) {
					v.Value = username
				}
			}
		}

		items = append(items, buildItem(v, list.Type))
	}

	return items, nil
}

func buildItem(itemContent clientlists.ListItemContent, listType clientlists.ClientListType) interface{} {
	if isKeyValuesItem(string(listType)) {
		return map[string]interface{}{
			"key":             itemContent.Key,
			"values":          itemContent.Values,
			"description":     itemContent.Description,
			"expiration_date": itemContent.ExpirationDate,
			"tags":            itemContent.Tags,
		}
	}

	return map[string]interface{}{
		"value":           itemContent.Value,
		"description":     itemContent.Description,
		"expiration_date": itemContent.ExpirationDate,
		"tags":            itemContent.Tags,
	}
}

// valueItemHashFn hashes value items (IP, GEO, ASN, TLS_FINGERPRINT, FILE_HASH, USER_ID, DOMAIN)
var valueItemHashFn = schema.HashResource(&schema.Resource{Schema: map[string]*schema.Schema{
	"value":           {Type: schema.TypeString, Optional: true, Default: ""},
	"description":     {Type: schema.TypeString, Optional: true, Default: ""},
	"tags":            {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
	"expiration_date": {Type: schema.TypeString, Optional: true, Default: ""},
}})

// keyValuesItemHashFn hashes key-values items (REQUEST_HEADER_NAME_VALUE)
var keyValuesItemHashFn = schema.HashResource(&schema.Resource{Schema: map[string]*schema.Schema{
	"key":             {Type: schema.TypeString, Optional: true, Default: ""},
	"values":          {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
	"description":     {Type: schema.TypeString, Optional: true, Default: ""},
	"tags":            {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
	"expiration_date": {Type: schema.TypeString, Optional: true, Default: ""},
}})

// itemsHashFunc is a custom hash function for the items TypeSet.
func itemsHashFunc(v interface{}) int {
	m := v.(map[string]interface{})
	key, _ := m["key"].(string)
	if key != "" {
		return keyValuesItemHashFn(m)
	}
	return valueItemHashFn(m)
}

func buildItemPayload(itemMap map[string]interface{}, listType string) clientlists.ListItemPayload {
	if isKeyValuesItem(listType) {
		return clientlists.ListItemPayload{
			Key:            itemMap["key"].(string),
			Values:         tf.SetToStringSlice(itemMap["values"].(*schema.Set)),
			Description:    itemMap["description"].(string),
			Tags:           tf.SetToStringSlice(itemMap["tags"].(*schema.Set)),
			ExpirationDate: itemMap["expiration_date"].(string),
		}
	}

	return clientlists.ListItemPayload{
		Value:          itemMap["value"].(string),
		Description:    itemMap["description"].(string),
		Tags:           tf.SetToStringSlice(itemMap["tags"].(*schema.Set)),
		ExpirationDate: itemMap["expiration_date"].(string),
	}
}

func isKeyValuesItem(listType string) bool {
	return listType == string(clientlists.RequestHeaderNameValue)
}

func translateUsernames(ctx context.Context, d *schema.ResourceData, client clientlists.ClientLists) (clientlists.TranslateUsernamesResponse, diag.Diagnostics) {
	listType := d.Get("type").(string)
	items := d.Get("items").(*schema.Set).List()

	if len(items) == 0 || listType != string(clientlists.USER) {
		return nil, nil
	}

	translationRequest, err := createTranslateUsernamesRequest(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if len(translationRequest) == 0 {
		return clientlists.TranslateUsernamesResponse{}, nil
	}

	translatedUsernames, err := client.TranslateUsernames(ctx, translationRequest)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return *translatedUsernames, nil
}
