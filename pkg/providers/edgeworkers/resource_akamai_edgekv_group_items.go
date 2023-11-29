package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/collections"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEdgeKVGroupItems() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEdgeKVGroupItemsCreate,
		ReadContext:   resourceEdgeKVGroupItemsRead,
		UpdateContext: resourceEdgeKVGroupItemsUpdate,
		DeleteContext: resourceEdgeKVGroupItemsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: &timeouts.SDKDefaultTimeout,
		},
		Schema: map[string]*schema.Schema{
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name for the EdgeKV namespace.",
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					string(edgeworkers.NamespaceStagingNetwork),
					string(edgeworkers.NamespaceProductionNetwork),
				}, false)),
				Description: "The network against which to execute the API request.",
			},
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the EdgeKV group.",
				ForceNew:    true,
			},
			"items": {
				Type:             schema.TypeMap,
				Required:         true,
				ValidateDiagFunc: tf.ValidateMapMinimalLength(1),
				Description:      "A map of items within the specified group. Each item consists of an item key and a value.",
				Elem:             &schema.Schema{Type: schema.TypeString},
			},
			"timeouts": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enables to set timeout for processing",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: timeouts.ValidateDurationFormat,
						},
					},
				},
			},
		},
	}
}

func resourceEdgeKVGroupItemsCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVGroupItemsCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Creating EdgeKV group items")

	attrs, err := getAttributes(rd)
	if err != nil {
		return diag.Errorf("could not get attributes: %s", err)
	}

	for key, valueRaw := range attrs.items {
		value, ok := valueRaw.(string)
		if !ok {
			return diag.Errorf("could not cast value of type %T into string", value)
		}
		_, err = client.UpsertItem(ctx, edgeworkers.UpsertItemRequest{
			ItemID:   key,
			ItemData: edgeworkers.Item(value),
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				Network:     attrs.network,
				NamespaceID: attrs.namespace,
				GroupID:     attrs.groupName,
			},
		})
		if err != nil {
			return diag.Errorf("could not upsert an item with key '%s': %s", key, err)
		}
	}

	if err = waitForEdgeKVGroupCreation(ctx, client, attrs.groupName, attrs); err != nil {
		return diag.Errorf("waitForEdgeKVGroupCreation error: %s", err)
	}

	if err = waitForConsistentEdgeKVDatabase(ctx, client, nil, attrs); err != nil {
		return diag.Errorf("waitForConsistentEdgeKVDatabase error: %s", err)
	}
	rd.SetId(fmt.Sprintf("%s:%s:%s", attrs.namespace, attrs.network, attrs.groupName))

	return resourceEdgeKVGroupItemsRead(ctx, rd, m)
}

func resourceEdgeKVGroupItemsRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVGroupItemsRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Reading EdgeKV group items")

	rdID := rd.Id()
	parts := strings.Split(rdID, ":")
	if len(parts) != 3 {
		return diag.Errorf("incorrect resource id format: must be `namespace_name:network:group_name`")
	}

	namespace, network, groupName := parts[0], parts[1], parts[2]

	items, err := client.ListItems(ctx, edgeworkers.ListItemsRequest{
		ItemsRequestParams: edgeworkers.ItemsRequestParams{
			NamespaceID: namespace,
			GroupID:     groupName,
			Network:     edgeworkers.ItemNetwork(network),
		},
	})
	if err != nil {
		return diag.Errorf("could not list items: %s", err)
	}

	itemsMap, err := getItems(ctx, items, client, network, namespace, groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := make(map[string]interface{})
	attrs["namespace_name"] = namespace
	attrs["network"] = network
	attrs["group_name"] = groupName
	attrs["items"] = itemsMap

	if err = tf.SetAttrs(rd, attrs); err != nil {
		return diag.Errorf("could not set attributes: %s", err)
	}

	return nil
}

func resourceEdgeKVGroupItemsUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVGroupItemsUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Updating EdgeKV group items")

	if !rd.HasChangeExcept("timeouts") {
		logger.Debug("Only timeouts were updated, skipping")
		return nil
	}

	if !rd.HasChanges("items") {
		return resourceEdgeKVGroupItemsRead(ctx, rd, m)
	}

	attrs, err := getAttributes(rd)
	if err != nil {
		return diag.Errorf("could not get attributes: %s", err)
	}

	remoteStateItems, err := client.ListItems(ctx, edgeworkers.ListItemsRequest{
		ItemsRequestParams: edgeworkers.ItemsRequestParams{
			NamespaceID: attrs.namespace,
			GroupID:     attrs.groupName,
			Network:     attrs.network,
		},
	})
	if err != nil {
		return diag.Errorf("could not list items: %s", err)
	}

	remoteStateItemsArray := []string(*remoteStateItems)
	var deletedItems []string

	// first loop for creating items, where the loop iterates through items specified in config
	for key, valueRaw := range attrs.items {
		value, ok := valueRaw.(string)
		if !ok {
			return diag.Errorf("could not cast value of type %T to string", valueRaw)
		}
		if !collections.StringInSlice(remoteStateItemsArray, key) {
			_, err = client.UpsertItem(ctx, edgeworkers.UpsertItemRequest{
				ItemID:   key,
				ItemData: edgeworkers.Item(value),
				ItemsRequestParams: edgeworkers.ItemsRequestParams{
					Network:     attrs.network,
					NamespaceID: attrs.namespace,
					GroupID:     attrs.groupName,
				},
			})
			if err != nil {
				return diag.Errorf("could not upsert an item with key '%s': %s", key, err)
			}
		}
	}

	// second loop updates or deletes items, where the loop iterates through items present in the remote state
	for _, remoteStateItemKey := range remoteStateItemsArray {
		if val, ok := attrs.items[remoteStateItemKey]; ok {
			strVal, ok := val.(string)
			if !ok {
				return diag.Errorf("could not cast value of type %T to string", val)
			}

			_, err = client.UpsertItem(ctx, edgeworkers.UpsertItemRequest{
				ItemID:   remoteStateItemKey,
				ItemData: edgeworkers.Item(strVal),
				ItemsRequestParams: edgeworkers.ItemsRequestParams{
					Network:     attrs.network,
					NamespaceID: attrs.namespace,
					GroupID:     attrs.groupName,
				},
			})
			if err != nil {
				return diag.Errorf("could not upsert an item with key '%s': %s", remoteStateItemKey, err)
			}
		} else {
			_, err = client.DeleteItem(ctx, edgeworkers.DeleteItemRequest{
				ItemID: remoteStateItemKey,
				ItemsRequestParams: edgeworkers.ItemsRequestParams{
					Network:     attrs.network,
					NamespaceID: attrs.namespace,
					GroupID:     attrs.groupName,
				},
			})
			if err != nil {
				return diag.Errorf("could not delete an item with key '%s': %s", remoteStateItemKey, err)
			}
			deletedItems = append(deletedItems, remoteStateItemKey)
		}
	}

	if err = waitForConsistentEdgeKVDatabase(ctx, client, deletedItems, attrs); err != nil {
		return diag.Errorf("waitForConsistentEdgeKVDatabase error: %s", err)
	}

	return resourceEdgeKVGroupItemsRead(ctx, rd, m)
}

func resourceEdgeKVGroupItemsDelete(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "resourceEdgeKVGroupItemsDelete")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Deleting EdgeKV group items")

	attrs, err := getAttributes(rd)
	if err != nil {
		return diag.Errorf("could not get attributes: %s", err)
	}

	remoteStateItems, err := client.ListItems(ctx, edgeworkers.ListItemsRequest{
		ItemsRequestParams: edgeworkers.ItemsRequestParams{
			Network:     attrs.network,
			NamespaceID: attrs.namespace,
			GroupID:     attrs.groupName,
		},
	})
	if err != nil {
		return diag.Errorf("could not list items: %s", err)
	}

	remoteStateItemsArray := []string(*remoteStateItems)
	if len(attrs.items) != len(remoteStateItemsArray) {
		return diag.Errorf("in order to delete whole group of items, number of items in the configuration and remote state should be the same")
	}

	for key := range attrs.items {
		if !collections.StringInSlice(remoteStateItemsArray, key) {
			return diag.Errorf("item with key '%s' does not exist in the remote state of the database", key)
		}
		_, err = client.DeleteItem(ctx, edgeworkers.DeleteItemRequest{
			ItemID: key,
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				Network:     attrs.network,
				NamespaceID: attrs.namespace,
				GroupID:     attrs.groupName,
			},
		})
		if err != nil {
			return diag.Errorf("could not delete an item with key '%s': %s", key, err)
		}
	}
	if err = waitForEdgeKVGroupDeletion(ctx, client, attrs.groupName, attrs); err != nil {
		return diag.Errorf("waitForEdgeKVGroupDeletion error: %s", err)
	}

	rd.SetId("")
	return nil
}

var (
	// pollForConsistentEdgeKVDatabaseInterval defines retry interval for listing items or getting and item
	pollForConsistentEdgeKVDatabaseInterval = 5 * time.Second
)

// waitForEdgeKVGroupCreation waits for the group to be created in the remote state
func waitForEdgeKVGroupCreation(ctx context.Context, client edgeworkers.Edgeworkers, groupName string, attrs *edgeKVGroupItemsAttrs) error {
	var groupExists bool

	for !groupExists {
		select {
		case <-time.After(pollForConsistentEdgeKVDatabaseInterval):
			groups, err := client.ListGroupsWithinNamespace(ctx, edgeworkers.ListGroupsWithinNamespaceRequest{
				Network:     edgeworkers.NamespaceNetwork(attrs.network),
				NamespaceID: attrs.namespace,
			})
			if err != nil && !errors.Is(err, edgeworkers.ErrNotFound) {
				return fmt.Errorf("could not list groups within network: `%s` and namespace_name: `%s`: %s", attrs.network, attrs.namespace, err)
			}

			if collections.StringInSlice(groups, groupName) {
				groupExists = true
			}
		case <-ctx.Done():
			return fmt.Errorf("retry timeout reached for list groups")
		}
	}

	return nil
}

// waitForEdgeKVGroupDeletion waits for the group to be deleted from the remote state
func waitForEdgeKVGroupDeletion(ctx context.Context, client edgeworkers.Edgeworkers, groupName string, attrs *edgeKVGroupItemsAttrs) error {
	groupExists := true

	for groupExists {
		select {
		case <-time.After(pollForConsistentEdgeKVDatabaseInterval):
			groups, err := client.ListGroupsWithinNamespace(ctx, edgeworkers.ListGroupsWithinNamespaceRequest{
				Network:     edgeworkers.NamespaceNetwork(attrs.network),
				NamespaceID: attrs.namespace,
			})
			if err != nil && !errors.Is(err, edgeworkers.ErrNotFound) {
				return fmt.Errorf("could not list groups within network: `%s` and namespace_name: `%s`: %s", attrs.network, attrs.namespace, err)
			}

			if !collections.StringInSlice(groups, groupName) {
				groupExists = false
			}
		case <-ctx.Done():
			return fmt.Errorf("retry timeout reached for list groups")
		}
	}

	return nil
}

// waitForConsistentEdgeKVDatabase waits until all items specified in the config are propagated in the remote state, as well as all
// the deleted items from the config are removed from the remote state.
func waitForConsistentEdgeKVDatabase(ctx context.Context, client edgeworkers.Edgeworkers, deletedItems []string, attrs *edgeKVGroupItemsAttrs) error {
	for itemKey, itemRaw := range attrs.items {
		itemVal, ok := itemRaw.(string)
		if !ok {
			return fmt.Errorf("could not cast value of type '%T' into string", itemRaw)
		}

		var isPresent bool
		for !isPresent {
			select {
			case <-time.After(pollForConsistentEdgeKVDatabaseInterval):
				stateVal, err := client.GetItem(ctx, edgeworkers.GetItemRequest{
					ItemID: itemKey,
					ItemsRequestParams: edgeworkers.ItemsRequestParams{
						Network:     attrs.network,
						NamespaceID: attrs.namespace,
						GroupID:     attrs.groupName,
					},
				})
				if err != nil && !errors.Is(err, edgeworkers.ErrNotFound) {
					return fmt.Errorf("could not get an item with key '%s': %s", itemKey, err)
				} else if err == nil && string(*stateVal) == itemVal {
					isPresent = true
				}
			case <-ctx.Done():
				return fmt.Errorf("retry timeout reached for get an item")
			}
		}
	}

	for _, itemKey := range deletedItems {
		var isDeleted bool
		for !isDeleted {
			select {
			case <-time.After(pollForConsistentEdgeKVDatabaseInterval):
				_, err := client.GetItem(ctx, edgeworkers.GetItemRequest{
					ItemID: itemKey,
					ItemsRequestParams: edgeworkers.ItemsRequestParams{
						Network:     attrs.network,
						NamespaceID: attrs.namespace,
						GroupID:     attrs.groupName,
					},
				})
				if err != nil && !errors.Is(err, edgeworkers.ErrNotFound) {
					return fmt.Errorf("could not get an item with key '%s': %s", itemKey, err)
				} else if err != nil && errors.Is(err, edgeworkers.ErrNotFound) {
					isDeleted = true
				}
			case <-ctx.Done():
				return fmt.Errorf("retry timeout reached for get an item")
			}
		}
	}

	return nil
}

// edgeKVGroupItemsAttrs represents attributes for edgeKV_group_items resource
type edgeKVGroupItemsAttrs struct {
	namespace, groupName string
	network              edgeworkers.ItemNetwork
	items                map[string]interface{}
}

// getAttributes retrieves edgeKV_group_items attributes from config
func getAttributes(rd *schema.ResourceData) (*edgeKVGroupItemsAttrs, error) {
	namespace, err := tf.GetStringValue("namespace_name", rd)
	if err != nil {
		return nil, fmt.Errorf("could not get 'namespace_name' attribute: %s", err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return nil, fmt.Errorf("could not get 'network' attribute: %s", err)
	}

	groupName, err := tf.GetStringValue("group_name", rd)
	if err != nil {
		return nil, fmt.Errorf("could not get 'group_name' attribute: %s", err)
	}

	items, err := tf.GetMapValue("items", rd)
	if err != nil {
		return nil, fmt.Errorf("could not get 'items' attribute: %s", err)
	}

	return &edgeKVGroupItemsAttrs{
		namespace: namespace,
		groupName: groupName,
		network:   edgeworkers.ItemNetwork(network),
		items:     items,
	}, nil
}

// getEdgeKVItemPtr returns pointer to the edgeworkers.Item
func getEdgeKVItemPtr(value string) *edgeworkers.Item {
	itemVal := edgeworkers.Item(value)
	return &itemVal
}

func getItems(ctx context.Context, items *edgeworkers.ListItemsResponse, client edgeworkers.Edgeworkers, network, namespace, groupName string) (map[string]string, error) {
	itemsMap := make(map[string]string)
	for _, itemKey := range *items {
		itemValue, err := client.GetItem(ctx, edgeworkers.GetItemRequest{
			ItemID: itemKey,
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				Network:     edgeworkers.ItemNetwork(network),
				NamespaceID: namespace,
				GroupID:     groupName,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("could not get an item with key '%s': %s", itemKey, err)
		}
		itemsMap[itemKey] = string(*itemValue)
	}
	return itemsMap, nil
}
