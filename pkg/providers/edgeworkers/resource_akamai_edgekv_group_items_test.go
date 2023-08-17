package edgeworkers

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestCreateEdgeKVGroupItems(t *testing.T) {
	// decrease interval for testing
	pollForConsistentEdgeKVDatabaseInterval = time.Microsecond

	tests := map[string]struct {
		configPath string
		attrs      edgeKVConfigurationForTests
		init       func(*edgeworkers.Mock, edgeKVConfigurationForTests)
		withError  *regexp.Regexp
	}{
		"create edgeKV items": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			attrs: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			init: func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrs, "key1", "value1", 1)
				mockUpsertItem(m, attrs, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrs, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrs, "key1", "value1", 1)
				mockGetItem(m, attrs, "key2", "value2", 1)
				// read
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2"}, 2)
				mockGetItem(m, attrs, "key1", "value1", 2)
				mockGetItem(m, attrs, "key2", "value2", 2)
				// delete
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrs, "key1", "message1", 1)
				mockDeleteItem(m, attrs, "key2", "message2", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrs, []string{}, 1)
			},
		},
		"create edgeKV items, group not yet created": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			attrs: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			init: func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrs, "key1", "value1", 1)
				mockUpsertItem(m, attrs, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation - no group
				mockErrListGroupsWithinNamespace(m, attrs, 1)
				// waitForEdgeKVGroupCreation - loop #1 - group created
				mockListGroupsWithinNamespace(m, attrs, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrs, "key1", "value1", 1)
				mockGetItem(m, attrs, "key2", "value2", 1)
				// read
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2"}, 2)
				mockGetItem(m, attrs, "key1", "value1", 2)
				mockGetItem(m, attrs, "key2", "value2", 2)
				// delete
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrs, "key1", "message1", 1)
				mockDeleteItem(m, attrs, "key2", "message2", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrs, []string{}, 1)
			},
		},
		"create edgeKV items, items not yet populated": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			attrs: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			init: func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrs, "key1", "value1", 1)
				mockUpsertItem(m, attrs, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrs, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				// item with key1 not populated yet
				mockErrGetItem(m, attrs, "key1", 1)
				// item with key1 now populated
				mockGetItem(m, attrs, "key1", "value1", 1)
				mockGetItem(m, attrs, "key2", "value2", 1)
				// read
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2"}, 2)
				mockGetItem(m, attrs, "key1", "value1", 2)
				mockGetItem(m, attrs, "key2", "value2", 2)
				// delete
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrs, "key1", "message1", 1)
				mockDeleteItem(m, attrs, "key2", "message2", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrs, []string{}, 1)
			},
		},
		"create with no group items - error": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/empty_items.tf",
			init:       func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {},
			withError:  regexp.MustCompile("Error: map must contain at least 1 element\\(s\\), but has 0"),
		},
		"no namespace_name - error": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/no_namespace.tf",
			init:       func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {},
			withError:  regexp.MustCompile("Error: Missing required argument"),
		},
		"no group_name - error": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/no_group.tf",
			init:       func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {},
			withError:  regexp.MustCompile("Error: Missing required argument"),
		},
		"no network - error": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/no_network.tf",
			init:       func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {},
			withError:  regexp.MustCompile("Error: Missing required argument"),
		},
		"no items attribute - error": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/no_items.tf",
			init:       func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {},
			withError:  regexp.MustCompile("Error: Missing required argument"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			test.init(client, test.attrs)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkEdgeKVGroupItemsAttrs(test.attrs),
							ExpectError: test.withError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestReadEdgeKVGroupItems(t *testing.T) {
	// decrease interval for testing
	pollForConsistentEdgeKVDatabaseInterval = time.Microsecond

	tests := map[string]struct {
		configPathForCreate string
		configPathForUpdate string
		attrsForCreate      edgeKVConfigurationForTests
		attrsForUpdate      edgeKVConfigurationForTests
		init                func(*edgeworkers.Mock, edgeKVConfigurationForTests, edgeKVConfigurationForTests)
		errorForCreate      *regexp.Regexp
		errorForUpdate      *regexp.Regexp
	}{
		"create, remove 1 item and read inconsistent state - item was not deleted yet from remote": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_3_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/remove_1_item_of_3.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key3": "value3",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				mockUpsertItem(m, attrsForCreate, "key3", "value3", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				mockGetItem(m, attrsForCreate, "key3", "value3", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				mockGetItem(m, attrsForCreate, "key3", "value3", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message2", 1)
				mockUpsertItem(m, attrsForUpdate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForUpdate, "key3", "value3", 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key3", "value3", 1)
				// item with key2 still exists
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				// waitForConsistentEdgeKVDatabase: wait for database to be consistent - loop #1, item with key2 deleted
				mockErrGetItem(m, attrsForUpdate, "key2", 1)
				// read
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key3"}, 2)
				mockGetItem(m, attrsForCreate, "key1", "value1", 2)
				mockGetItem(m, attrsForCreate, "key3", "value3", 2)
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key3"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockDeleteItem(m, attrsForUpdate, "key3", "message3", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{}, 1)
			},
		},
		"remove, add and upsert items - read inconsistent state - item was not deleted yet from remote and new items were not yet added - check logic for counting deleted and added items": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/check_counting_logic.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key2": "updatedValue",
					"key3": "value3",
					"key4": "value4",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockUpsertItem(m, attrsForUpdate, "key2", "updatedValue", 1)
				mockUpsertItem(m, attrsForUpdate, "key3", "value3", 1)
				mockUpsertItem(m, attrsForUpdate, "key4", "value4", 1)
				// waitForConsistentEdgeKVDatabase: read dirty state - item was not yet deleted from remote state and items were not yet added
				mockGetItem(m, attrsForCreate, "key2", "updatedValue", 1)
				// items with key3 and key4 not yet created
				mockErrGetItem(m, attrsForUpdate, "key3", 1)
				mockErrGetItem(m, attrsForUpdate, "key4", 1)
				// items with key3 and key4 created
				mockGetItem(m, attrsForCreate, "key3", "value3", 1)
				mockGetItem(m, attrsForCreate, "key4", "value4", 1)
				// item with key1 not yet deleted
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				// item with key1 deleted
				mockErrGetItem(m, attrsForUpdate, "key1", 1)
				// read
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key2", "key3", "key4"}, 2)
				mockGetItem(m, attrsForCreate, "key2", "updatedValue", 2)
				mockGetItem(m, attrsForCreate, "key3", "value3", 2)
				mockGetItem(m, attrsForCreate, "key4", "value4", 2)
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key2", "key3", "key4"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message2", 1)
				mockDeleteItem(m, attrsForUpdate, "key3", "message3", 1)
				mockDeleteItem(m, attrsForUpdate, "key4", "message4", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{}, 1)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			test.init(client, test.attrsForCreate, test.attrsForUpdate)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForCreate),
							Check:       checkEdgeKVGroupItemsAttrs(test.attrsForCreate),
							ExpectError: test.errorForCreate,
						},
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForUpdate),
							Check:       checkEdgeKVGroupItemsAttrs(test.attrsForUpdate),
							ExpectError: test.errorForUpdate,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestUpdateEdgeKVGroupItems(t *testing.T) {
	// decrease interval for testing
	pollForConsistentEdgeKVDatabaseInterval = time.Microsecond

	tests := map[string]struct {
		configPathForCreate string
		configPathForUpdate string
		attrsForCreate      edgeKVConfigurationForTests
		attrsForUpdate      edgeKVConfigurationForTests
		init                func(*edgeworkers.Mock, edgeKVConfigurationForTests, edgeKVConfigurationForTests)
		planOnly            bool
		expectNonEmptyPlan  bool
		errorForCreate      *regexp.Regexp
		errorForUpdate      *regexp.Regexp
	}{
		"update - add 1 item": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/add_1_item.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockUpsertItem(m, attrsForUpdate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForUpdate, "key2", "value2", 1)
				mockUpsertItem(m, attrsForUpdate, "key3", "value3", 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForUpdate, "key1", "value1", 1)
				mockGetItem(m, attrsForUpdate, "key2", "value2", 1)
				mockGetItem(m, attrsForUpdate, "key3", "value3", 1)
				// read
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 2)
				mockGetItem(m, attrsForUpdate, "key1", "value1", 2)
				mockGetItem(m, attrsForUpdate, "key2", "value2", 2)
				mockGetItem(m, attrsForUpdate, "key3", "value3", 2)
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message1", 1)
				mockDeleteItem(m, attrsForUpdate, "key3", "message1", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForUpdate, []string{}, 1)
			},
		},
		"update - remove 1 item": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/remove_1_item.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key2": "value2",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockUpsertItem(m, attrsForUpdate, "key2", "value2", 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForUpdate, "key2", "value2", 1)
				mockErrGetItem(m, attrsForUpdate, "key1", 1)
				// read
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key2"}, 2)
				mockGetItem(m, attrsForUpdate, "key2", "value2", 2)
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key2"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message1", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForUpdate, []string{}, 1)
			},
		},
		"update - upsert 1 item": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/upsert_1_item.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "updatedValue",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockUpsertItem(m, attrsForUpdate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForUpdate, "key2", "updatedValue", 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForUpdate, "key1", "value1", 1)
				mockGetItem(m, attrsForUpdate, "key2", "updatedValue", 1)
				// read
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 2)
				mockGetItem(m, attrsForUpdate, "key1", "value1", 2)
				mockGetItem(m, attrsForUpdate, "key2", "updatedValue", 2)
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message1", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForUpdate, []string{}, 1)
			},
		},
		"update - upsert 1 item by changing the key - delete and create new item": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/basic_upsert_key.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key3": "value1",
					"key2": "value2",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockUpsertItem(m, attrsForUpdate, "key2", "value2", 1)
				mockUpsertItem(m, attrsForUpdate, "key3", "value1", 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForUpdate, "key2", "value2", 1)
				mockGetItem(m, attrsForUpdate, "key3", "value1", 1)
				mockErrGetItem(m, attrsForUpdate, "key1", 1)
				// read
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key2", "key3"}, 2)
				mockGetItem(m, attrsForUpdate, "key2", "value2", 2)
				mockGetItem(m, attrsForUpdate, "key3", "value1", 2)
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key2", "key3"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message2", 1)
				mockDeleteItem(m, attrsForUpdate, "key3", "message3", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForUpdate, []string{}, 1)
			},
		},
		"update - update all items keys": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/change_all_keys.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key3": "value1",
					"key4": "value2",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockUpsertItem(m, attrsForUpdate, "key3", "value1", 1)
				mockUpsertItem(m, attrsForUpdate, "key4", "value2", 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message1", 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForUpdate, "key3", "value1", 1)
				mockGetItem(m, attrsForUpdate, "key4", "value2", 1)
				mockErrGetItem(m, attrsForUpdate, "key1", 1)
				mockErrGetItem(m, attrsForUpdate, "key2", 1)
				// read
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key3", "key4"}, 2)
				mockGetItem(m, attrsForUpdate, "key3", "value1", 2)
				mockGetItem(m, attrsForUpdate, "key4", "value2", 2)
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key3", "key4"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key3", "message3", 1)
				mockDeleteItem(m, attrsForUpdate, "key4", "message4", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForUpdate, []string{}, 1)
			},
		},
		"update - add, remove and upsert items together": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/add_remove_upsert_items.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key2": "updatedValue",
					"key3": "value3",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockUpsertItem(m, attrsForUpdate, "key2", "updatedValue", 1)
				mockUpsertItem(m, attrsForUpdate, "key3", "value3", 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForUpdate, "key2", "updatedValue", 1)
				mockGetItem(m, attrsForUpdate, "key3", "value3", 1)
				mockErrGetItem(m, attrsForUpdate, "key1", 1)
				// read
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key2", "key3"}, 2)
				mockGetItem(m, attrsForUpdate, "key2", "updatedValue", 2)
				mockGetItem(m, attrsForUpdate, "key3", "value3", 2)
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key2", "key3"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message1", 1)
				mockDeleteItem(m, attrsForUpdate, "key3", "message3", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForUpdate, []string{}, 1)
			},
		},
		"update - changed order - no diff": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_3_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/changed_order.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key2": "value2",
					"key3": "value3",
					"key1": "value1",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				mockUpsertItem(m, attrsForCreate, "key3", "value3", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				mockGetItem(m, attrsForCreate, "key3", "value3", 1)
				// read x4
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 4)
				mockGetItem(m, attrsForCreate, "key1", "value1", 4)
				mockGetItem(m, attrsForCreate, "key2", "value2", 4)
				mockGetItem(m, attrsForCreate, "key3", "value3", 4)
				// no update
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message2", 1)
				mockDeleteItem(m, attrsForUpdate, "key3", "message3", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{}, 1)
			},
		},
		"update - delete 2 items with one last failing": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_3_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/remove_fail.tf",
			attrsForCreate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
			},
			attrsForUpdate: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
				},
			},
			init: func(m *edgeworkers.Mock, attrsForCreate, attrsForUpdate edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrsForCreate, "key1", "value1", 1)
				mockUpsertItem(m, attrsForCreate, "key2", "value2", 1)
				mockUpsertItem(m, attrsForCreate, "key3", "value3", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrsForCreate, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrsForCreate, "key1", "value1", 1)
				mockGetItem(m, attrsForCreate, "key2", "value2", 1)
				mockGetItem(m, attrsForCreate, "key3", "value3", 1)
				// read (2x after create, 1x before update)
				mockListItems(m, attrsForCreate, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 3)
				mockGetItem(m, attrsForCreate, "key1", "value1", 3)
				mockGetItem(m, attrsForCreate, "key2", "value2", 3)
				mockGetItem(m, attrsForCreate, "key3", "value3", 3)
				// update
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 1)
				mockUpsertItem(m, attrsForUpdate, "key1", "value1", 1)
				mockDeleteItem(m, attrsForUpdate, "key2", "message2", 1)
				// return an error on delete call
				m.On("DeleteItem", mock.Anything, edgeworkers.DeleteItemRequest{
					ItemID: "key3",
					ItemsRequestParams: edgeworkers.ItemsRequestParams{
						NamespaceID: attrsForUpdate.namespaceID,
						Network:     attrsForUpdate.network,
						GroupID:     attrsForUpdate.groupID,
					},
				}).Return(nil, fmt.Errorf("error deleting an item")).Once()
				// delete
				mockListItems(m, attrsForUpdate, edgeworkers.ListItemsResponse{"key1"}, 1)
				mockDeleteItem(m, attrsForUpdate, "key1", "message1", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrsForUpdate, []string{}, 1)
			},
			errorForUpdate: regexp.MustCompile("error deleting an item"),
		},
		"update - remove 2 out of 2 items - error": {
			configPathForCreate: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			configPathForUpdate: "testdata/TestResourceEdgeKVGroupItems/update/remove_2_items.tf",
			init:                func(_ *edgeworkers.Mock, _, _ edgeKVConfigurationForTests) {},
			planOnly:            true,
			expectNonEmptyPlan:  true,
			errorForUpdate:      regexp.MustCompile("Error: map must contain at least 1 element\\(s\\), but has 0"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			test.init(client, test.attrsForCreate, test.attrsForUpdate)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							PlanOnly:           test.planOnly,
							ExpectNonEmptyPlan: test.expectNonEmptyPlan,
							Config:             testutils.LoadFixtureString(t, test.configPathForCreate),
							Check:              checkEdgeKVGroupItemsAttrs(test.attrsForCreate),
							ExpectError:        test.errorForCreate,
						},
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForUpdate),
							Check:       checkEdgeKVGroupItemsAttrs(test.attrsForUpdate),
							ExpectError: test.errorForUpdate,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestDeleteEdgeKVGroupItems(t *testing.T) {
	// decrease interval for testing
	pollForConsistentEdgeKVDatabaseInterval = time.Microsecond

	tests := map[string]struct {
		configPath string
		attrs      edgeKVConfigurationForTests
		init       func(*edgeworkers.Mock, edgeKVConfigurationForTests)
		error      *regexp.Regexp
	}{
		"create, wait for deletion of group": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/basic_2_items.tf",
			attrs: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			init: func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrs, "key1", "value1", 1)
				mockUpsertItem(m, attrs, "key2", "value2", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrs, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrs, "key1", "value1", 1)
				mockGetItem(m, attrs, "key2", "value2", 1)
				// read
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2"}, 2)
				mockGetItem(m, attrs, "key1", "value1", 2)
				mockGetItem(m, attrs, "key2", "value2", 2)
				// delete
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2"}, 1)
				mockDeleteItem(m, attrs, "key1", "message1", 1)
				mockDeleteItem(m, attrs, "key2", "message2", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrs, []string{"1234"}, 1) // group still exists
				mockListGroupsWithinNamespace(m, attrs, []string{"1234"}, 1) // group still exists
				mockListGroupsWithinNamespace(m, attrs, []string{}, 1)       // group removed from remote state
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			test.init(client, test.attrs)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkEdgeKVGroupItemsAttrs(test.attrs),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestImportEdgeKVGroupItems(t *testing.T) {
	tests := map[string]struct {
		configPath string
		attrs      edgeKVConfigurationForTests
		init       func(*edgeworkers.Mock, edgeKVConfigurationForTests)
		error      *regexp.Regexp
	}{
		"import": {
			configPath: "testdata/TestResourceEdgeKVGroupItems/create/basic_3_items.tf",
			attrs: edgeKVConfigurationForTests{
				namespaceID: "test_namespace",
				network:     "staging",
				groupID:     "1234",
				items: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
			},
			init: func(m *edgeworkers.Mock, attrs edgeKVConfigurationForTests) {
				// create
				mockUpsertItem(m, attrs, "key1", "value1", 1)
				mockUpsertItem(m, attrs, "key2", "value2", 1)
				mockUpsertItem(m, attrs, "key3", "value3", 1)
				// waitForEdgeKVGroupCreation
				mockListGroupsWithinNamespace(m, attrs, []string{"1234"}, 1)
				// waitForConsistentEdgeKVDatabase
				mockGetItem(m, attrs, "key1", "value1", 1)
				mockGetItem(m, attrs, "key2", "value2", 1)
				mockGetItem(m, attrs, "key3", "value3", 1)
				// read
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 2)
				mockGetItem(m, attrs, "key1", "value1", 2)
				mockGetItem(m, attrs, "key2", "value2", 2)
				mockGetItem(m, attrs, "key3", "value3", 2)
				// import
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 1)
				mockGetItem(m, attrs, "key1", "value1", 1)
				mockGetItem(m, attrs, "key2", "value2", 1)
				mockGetItem(m, attrs, "key3", "value3", 1)
				// delete
				mockListItems(m, attrs, edgeworkers.ListItemsResponse{"key1", "key2", "key3"}, 1)
				mockDeleteItem(m, attrs, "key1", "message1", 1)
				mockDeleteItem(m, attrs, "key2", "message2", 1)
				mockDeleteItem(m, attrs, "key3", "message2", 1)
				// waitForEdgeKVGroupDeletion
				mockListGroupsWithinNamespace(m, attrs, []string{}, 1)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			test.init(client, test.attrs)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,

					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, test.configPath),
						},
						{
							ImportState:       true,
							ImportStateId:     "test_namespace:staging:1234",
							ResourceName:      "akamai_edgekv_group_items.test",
							ImportStateVerify: true,
							ExpectError:       nil,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// edgeKVConfigurationForTests contains configuration attributes for edgeKV_group_items resource used in tests
type edgeKVConfigurationForTests struct {
	namespaceID string
	network     edgeworkers.ItemNetwork
	groupID     string
	items       map[string]string
}

var (
	// mockUpsertItem mocks 'UpsertItem' call with provided data
	mockUpsertItem = func(client *edgeworkers.Mock, attrs edgeKVConfigurationForTests, itemID, itemData string, timesToRun int) {
		client.On("UpsertItem", mock.Anything, edgeworkers.UpsertItemRequest{
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				NamespaceID: attrs.namespaceID,
				Network:     attrs.network,
				GroupID:     attrs.groupID,
			},
			ItemID:   itemID,
			ItemData: edgeworkers.Item(itemData),
		}).Return(tools.StringPtr("value1"), nil).Times(timesToRun)
	}

	// mockListItems mocks 'ListItems' call with provided data
	mockListItems = func(client *edgeworkers.Mock, attrs edgeKVConfigurationForTests, items edgeworkers.ListItemsResponse, timesToRun int) {
		client.On("ListItems", mock.Anything, edgeworkers.ListItemsRequest{
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				NamespaceID: attrs.namespaceID,
				Network:     attrs.network,
				GroupID:     attrs.groupID,
			},
		}).Return(&items, nil).Times(timesToRun)
	}

	// mockGetItem mocks 'GetItem' call with provided data
	mockGetItem = func(client *edgeworkers.Mock, attrs edgeKVConfigurationForTests, itemID, itemData string, timesToRun int) {
		client.On("GetItem", mock.Anything, edgeworkers.GetItemRequest{
			ItemID: itemID,
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				NamespaceID: attrs.namespaceID,
				Network:     attrs.network,
				GroupID:     attrs.groupID,
			},
		}).Return(getEdgeKVItemPtr(itemData), nil).Times(timesToRun)
	}

	// mockDeleteItem mocks 'DeleteItem' call with provided data
	mockDeleteItem = func(client *edgeworkers.Mock, attrs edgeKVConfigurationForTests, itemID, responseMessage string, timesToRun int) {
		client.On("DeleteItem", mock.Anything, edgeworkers.DeleteItemRequest{
			ItemID: itemID,
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				NamespaceID: attrs.namespaceID,
				Network:     attrs.network,
				GroupID:     attrs.groupID,
			},
		}).Return(tools.StringPtr(responseMessage), nil).Times(timesToRun)
	}

	// mockListGroupsWithinNamespace mocks 'ListGroupsWithinNamespace' call with provided data
	mockListGroupsWithinNamespace = func(client *edgeworkers.Mock, attrs edgeKVConfigurationForTests, groups []string, timesToRun int) {
		client.On("ListGroupsWithinNamespace", mock.Anything, edgeworkers.ListGroupsWithinNamespaceRequest{
			Network:     edgeworkers.NamespaceNetwork(attrs.network),
			NamespaceID: attrs.namespaceID,
		}).Return(groups, nil).Times(timesToRun)
	}

	// mockErrListGroupsWithinNamespace mocks 'ListGroupsWithinNamespace' call and returns 'ErrGroupNotFound' error
	mockErrListGroupsWithinNamespace = func(client *edgeworkers.Mock, attrs edgeKVConfigurationForTests, timesToRun int) {
		client.On("ListGroupsWithinNamespace", mock.Anything, edgeworkers.ListGroupsWithinNamespaceRequest{
			Network:     edgeworkers.NamespaceNetwork(attrs.network),
			NamespaceID: attrs.namespaceID,
		}).Return(nil, edgeworkers.ErrNotFound).Times(timesToRun)
	}

	// mockErrGetItem mocks 'GetItem' call and returns 'ErrItemNotFound' error
	mockErrGetItem = func(client *edgeworkers.Mock, attrs edgeKVConfigurationForTests, itemKey string, timesToRun int) {
		client.On("GetItem", mock.Anything, edgeworkers.GetItemRequest{
			ItemID: itemKey,
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				NamespaceID: attrs.namespaceID,
				Network:     attrs.network,
				GroupID:     attrs.groupID,
			},
		}).Return(nil, edgeworkers.ErrNotFound).Times(timesToRun)
	}
)

// checkEdgeKVGroupItemsAttrs creates resource.TestCheckFunc functions that check the resource' state based on the provided data
func checkEdgeKVGroupItemsAttrs(data edgeKVConfigurationForTests) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edgekv_group_items.test", "namespace_name", data.namespaceID))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edgekv_group_items.test", "network", string(data.network)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edgekv_group_items.test", "group_name", data.groupID))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edgekv_group_items.test", "items.%", strconv.Itoa(len(data.items))))
	for key, val := range data.items {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edgekv_group_items.test", fmt.Sprintf("items.%s", key), val))
	}
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edgekv_group_items.test", "id", fmt.Sprintf("%s:%s:%s", data.namespaceID, data.network, data.groupID)))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
