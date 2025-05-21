package property

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataPropertyIncludes(t *testing.T) {
	tests := map[string]struct {
		attrs      attributes
		init       func(*papi.Mock, attributes)
		configPath string
		error      *regexp.Regexp
	}{
		"ListIncludes - without type filter": {
			attrs: attributes{
				contractID:     contractForTests,
				groupID:        groupForTests,
				includeType:    "",
				parentProperty: nil,
				includes:       createIncludes(10, contractForTests, groupForTests, false, false),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListIncludes(m, attrs.contractID, attrs.groupID, attrs.includes, 3)
			},
			configPath: "testdata/TestDataPropertyIncludes/without_parent_property/list_includes_no_filters.tf",
		},
		"ListIncludes - with MICROSERVICES include type": {
			attrs: attributes{
				contractID:     contractForTests,
				groupID:        groupForTests,
				includeType:    string(papi.IncludeTypeMicroServices),
				parentProperty: nil,
				includes:       createIncludes(10, contractForTests, groupForTests, false, false),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListIncludes(m, attrs.contractID, attrs.groupID, attrs.includes, 3)
			},
			configPath: "testdata/TestDataPropertyIncludes/without_parent_property/list_includes_type_microservices.tf",
		},
		"ListIncludes - with COMMON_SETTINGS include type": {
			attrs: attributes{
				contractID:     contractForTests,
				groupID:        groupForTests,
				includeType:    string(papi.IncludeTypeCommonSettings),
				parentProperty: nil,
				includes:       createIncludes(10, contractForTests, groupForTests, false, false),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListIncludes(m, attrs.contractID, attrs.groupID, attrs.includes, 3)
			},
			configPath: "testdata/TestDataPropertyIncludes/without_parent_property/list_includes_type_common_settings.tf",
		},
		"ListIncludes - no includes": {
			attrs: attributes{
				contractID:     contractForTests,
				groupID:        groupForTests,
				includeType:    "",
				parentProperty: nil,
				includes:       nil,
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListIncludes(m, attrs.contractID, attrs.groupID, attrs.includes, 3)
			},
			configPath: "testdata/TestDataPropertyIncludes/without_parent_property/list_includes_no_filters.tf",
		},
		"ListIncludes - nil production and staging version": {
			attrs: attributes{
				contractID:     contractForTests,
				groupID:        groupForTests,
				includeType:    "",
				parentProperty: nil,
				includes:       createIncludes(1, contractForTests, groupForTests, true, true),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListIncludes(m, attrs.contractID, attrs.groupID, attrs.includes, 3)
			},
			configPath: "testdata/TestDataPropertyIncludes/without_parent_property/list_includes_no_filters.tf",
		},
		"ParentProperty provided - ListAvailableIncludes - without type filter": {
			attrs: attributes{
				contractID:  contractForTests,
				groupID:     groupForTests,
				includeType: "",
				parentProperty: &parentPropertyAttr{
					id:      "property_id_123",
					version: 2,
				},
				externalIncludes: createExternalIncludeData(2),
				includes:         createIncludes(2, contractForTests, groupForTests, false, false),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListAvailableIncludes(m, attrs.contractID, attrs.groupID, attrs.parentProperty.id, attrs.parentProperty.version, len(attrs.includes), 3)
				for i, include := range attrs.externalIncludes {
					mockGetInclude(m, attrs.contractID, attrs.groupID, include.IncludeID, 3, attrs.includes[i])
				}
			},
			configPath: "testdata/TestDataPropertyIncludes/with_parent_property/list_available_includes_no_filters.tf",
		},
		"ParentProperty provided - ListAvailableIncludes - with `MICROSERVICES` include type": {
			attrs: attributes{
				contractID:  contractForTests,
				groupID:     groupForTests,
				includeType: string(papi.IncludeTypeMicroServices),
				parentProperty: &parentPropertyAttr{
					id:      "property_id_123",
					version: 10,
				},
				externalIncludes: createExternalIncludeData(15),
				includes:         createIncludes(15, contractForTests, groupForTests, false, false),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListAvailableIncludes(m, attrs.contractID, attrs.groupID, attrs.parentProperty.id, attrs.parentProperty.version, len(attrs.includes), 3)
				for i, include := range attrs.externalIncludes {
					mockGetInclude(m, attrs.contractID, attrs.groupID, include.IncludeID, 3, attrs.includes[i])
				}
			},
			configPath: "testdata/TestDataPropertyIncludes/with_parent_property/list_available_includes_type_microservices.tf",
		},
		"ParentProperty provided - ListAvailableIncludes - with `COMMON_SETTINGS` include type": {
			attrs: attributes{
				contractID:  contractForTests,
				groupID:     groupForTests,
				includeType: string(papi.IncludeTypeCommonSettings),
				parentProperty: &parentPropertyAttr{
					id:      "property_id_123",
					version: 47,
				},
				externalIncludes: createExternalIncludeData(30),
				includes:         createIncludes(30, contractForTests, groupForTests, false, false),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListAvailableIncludes(m, attrs.contractID, attrs.groupID, attrs.parentProperty.id, attrs.parentProperty.version, len(attrs.includes), 3)
				for i, include := range attrs.externalIncludes {
					mockGetInclude(m, attrs.contractID, attrs.groupID, include.IncludeID, 3, attrs.includes[i])
				}
			},
			configPath: "testdata/TestDataPropertyIncludes/with_parent_property/list_available_includes_type_common_settings.tf",
		},
		"ParentProperty provided - ListAvailableIncludes - no includes": {
			attrs: attributes{
				contractID:  contractForTests,
				groupID:     groupForTests,
				includeType: "",
				parentProperty: &parentPropertyAttr{
					id:      "property_id_123",
					version: 2,
				},
				externalIncludes: createExternalIncludeData(0),
				includes:         createIncludes(0, contractForTests, groupForTests, false, false),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListAvailableIncludes(m, attrs.contractID, attrs.groupID, attrs.parentProperty.id, attrs.parentProperty.version, len(attrs.includes), 3)
				for i, include := range attrs.externalIncludes {
					mockGetInclude(m, attrs.contractID, attrs.groupID, include.IncludeID, 5, attrs.includes[i])
				}
			},
			configPath: "testdata/TestDataPropertyIncludes/with_parent_property/list_available_includes_no_filters.tf",
		},
		"ParentProperty provided - ListAvailableIncludes - with no production version": {
			attrs: attributes{
				contractID:  contractForTests,
				groupID:     groupForTests,
				includeType: string(papi.IncludeTypeMicroServices),
				parentProperty: &parentPropertyAttr{
					id:      "property_id_123",
					version: 10,
				},
				externalIncludes: createExternalIncludeData(15),
				includes:         createIncludes(15, contractForTests, groupForTests, false, true),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListAvailableIncludes(m, attrs.contractID, attrs.groupID, attrs.parentProperty.id, attrs.parentProperty.version, len(attrs.includes), 3)
				for i, include := range attrs.externalIncludes {
					mockGetInclude(m, attrs.contractID, attrs.groupID, include.IncludeID, 3, attrs.includes[i])
				}
			},
			configPath: "testdata/TestDataPropertyIncludes/with_parent_property/list_available_includes_type_microservices.tf",
		},
		"ListIncludes - API error": {
			attrs: attributes{
				contractID:     contractForTests,
				groupID:        groupForTests,
				includeType:    "",
				parentProperty: nil,
				includes:       nil,
			},
			init: func(m *papi.Mock, attrs attributes) {
				m.On("ListIncludes", testutils.MockContext, papi.ListIncludesRequest{
					ContractID: attrs.contractID,
					GroupID:    attrs.groupID,
				}).Return(nil, fmt.Errorf("fetching includes failed")).Once()
			},
			configPath: "testdata/TestDataPropertyIncludes/without_parent_property/list_includes_no_filters.tf",
			error:      regexp.MustCompile("Error: sendListIncludesRequest error: could not list includes: fetching includes failed"),
		},
		"ListAvailableIncludes - API error": {
			attrs: attributes{
				contractID:  contractForTests,
				groupID:     groupForTests,
				includeType: "",
				parentProperty: &parentPropertyAttr{
					id:      "property_id_123",
					version: 2,
				},
			},
			init: func(m *papi.Mock, attrs attributes) {
				m.On("ListAvailableIncludes", testutils.MockContext, papi.ListAvailableIncludesRequest{
					PropertyID:      attrs.parentProperty.id,
					PropertyVersion: attrs.parentProperty.version,
					ContractID:      attrs.contractID,
					GroupID:         attrs.groupID,
				}).Return(nil, fmt.Errorf("could not list available includes")).Once()
			},
			configPath: "testdata/TestDataPropertyIncludes/with_parent_property/list_available_includes_no_filters.tf",
			error:      regexp.MustCompile("could not list available includes"),
		},
		"GetInclude - API error": {
			attrs: attributes{
				contractID:  contractForTests,
				groupID:     groupForTests,
				includeType: "",
				parentProperty: &parentPropertyAttr{
					id:      "property_id_123",
					version: 2,
				},
				includes: createIncludes(8, contractForTests, groupForTests, false, false),
			},
			init: func(m *papi.Mock, attrs attributes) {
				mockListAvailableIncludes(m, attrs.contractID, attrs.groupID, attrs.parentProperty.id, attrs.parentProperty.version, len(attrs.includes), 1)
				m.On("GetInclude", testutils.MockContext, papi.GetIncludeRequest{
					ContractID: attrs.contractID,
					GroupID:    attrs.groupID,
					IncludeID:  attrs.includes[0].IncludeID,
				}).Return(nil, fmt.Errorf("could not get include")).Once()
			},
			configPath: "testdata/TestDataPropertyIncludes/with_parent_property/list_available_includes_no_filters.tf",
			error:      regexp.MustCompile("Error: sendListIncludesRequest error: could not get an include with ID: 0, could not get include"),
		},
		"missing required argument - contractID": {
			configPath: "testdata/TestDataPropertyIncludes/no_contract_id.tf",
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"missing required argument - groupID": {
			configPath: "testdata/TestDataPropertyIncludes/no_group_id.tf",
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"missing required argument - property ID": {
			configPath: "testdata/TestDataPropertyIncludes/no_property_id.tf",
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"missing required argument - property version": {
			configPath: "testdata/TestDataPropertyIncludes/no_property_version.tf",
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"invalid include type": {
			configPath: "testdata/TestDataPropertyIncludes/invalid_include_type.tf",
			error:      regexp.MustCompile(`Error: expected type to be one of \['MICROSERVICES', 'COMMON_SETTINGS'], got WRONG TYPE`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			if test.init != nil {
				test.init(client, test.attrs)
			}
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, test.configPath),
						Check:       checkPropertyIncludesAttrs(test.attrs),
						ExpectError: test.error,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

type attributes struct {
	contractID       string
	groupID          string
	includeType      string
	parentProperty   *parentPropertyAttr
	externalIncludes []papi.ExternalIncludeData
	includes         []papi.Include
}

var (
	contractForTests = "contract_123"
	groupForTests    = "group_321"

	// checkPropertyIncludesAttrs creates resource.TestCheckFunc that checks the data source attributes
	checkPropertyIncludesAttrs = func(attrs attributes) resource.TestCheckFunc {
		var testCheckFuncs []resource.TestCheckFunc
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_includes.test", "id", createDataSourceIDForTests(attrs)))
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_includes.test", "contract_id", attrs.contractID))
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_includes.test", "group_id", attrs.groupID))

		parentPropertyCheck := "0"
		if attrs.parentProperty != nil {
			parentPropertyCheck = "1"
		}
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_includes.test", "parent_property.#", parentPropertyCheck))

		if attrs.includeType == "" {
			testCheckFuncs = append(testCheckFuncs, resource.TestCheckNoResourceAttr("data.akamai_property_includes.test", "type"))
		} else {
			testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_includes.test", "type", attrs.includeType))
		}

		filteredIncludes := filterIncludes(attrs.includes, attrs.includeType)
		testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_includes.test", "includes.#", strconv.Itoa(len(filteredIncludes))))
		testCheckFuncs = append(testCheckFuncs, createIncludeTestCheckFuncs(filteredIncludes)...)

		return resource.ComposeAggregateTestCheckFunc(testCheckFuncs...)
	}

	// createIncludes creates Includes with provided parameters
	createIncludes = func(includesNumber int, contractID, groupID string, nilStagVer, nilProdVer bool) []papi.Include {
		var includes []papi.Include
		includeTypes := []string{string(papi.IncludeTypeMicroServices), string(papi.IncludeTypeCommonSettings)}

		var stagingVersion *int
		if !nilStagVer {
			stagingVersion = ptr.To(10)
		}

		var productionVersion *int
		if !nilProdVer {
			productionVersion = ptr.To(11)
		}

		for i := 0; i < includesNumber; i++ {
			index := rand.Int() % 2
			includes = append(includes, papi.Include{
				AccountID:         "account_123",
				AssetID:           "asset_123",
				ContractID:        contractID,
				GroupID:           groupID,
				IncludeID:         strconv.Itoa(i),
				IncludeName:       fmt.Sprintf("Name %d", i),
				IncludeType:       papi.IncludeType(includeTypes[index]),
				LatestVersion:     i,
				ProductionVersion: productionVersion,
				PropertyType:      nil,
				StagingVersion:    stagingVersion,
			})
		}

		return includes
	}

	// createExternalIncludeData creates ExternalIncludeData
	createExternalIncludeData = func(includesNumber int) []papi.ExternalIncludeData {
		var includes []papi.ExternalIncludeData
		includeTypes := []string{string(papi.IncludeTypeMicroServices), string(papi.IncludeTypeCommonSettings)}

		for i := 0; i < includesNumber; i++ {
			index := rand.Int() % 2
			includes = append(includes, papi.ExternalIncludeData{
				IncludeID:   strconv.Itoa(i),
				IncludeName: fmt.Sprintf("Name %d", i),
				IncludeType: papi.IncludeType(includeTypes[index]),
				FileName:    "test_fileName",
				ProductName: "test_productName",
				RuleFormat:  "v.2020-01-01",
			})
		}

		return includes
	}

	// mockListIncludes mocks ListIncludes call with provided parameters
	mockListIncludes = func(m *papi.Mock, contractID, groupID string, includes []papi.Include, timesToRun int) {
		m.On("ListIncludes", testutils.MockContext, papi.ListIncludesRequest{
			ContractID: contractID,
			GroupID:    groupID,
		}).Return(&papi.ListIncludesResponse{
			Includes: papi.IncludeItems{
				Items: includes,
			},
		}, nil).Times(timesToRun)
	}

	// mockListAvailableIncludes mocks ListAvailableIncludes call with provided parameters
	mockListAvailableIncludes = func(m *papi.Mock, contractID, groupID, propertyID string, propertyVersion, includesNumber, timesToRun int) {
		m.On("ListAvailableIncludes", testutils.MockContext, papi.ListAvailableIncludesRequest{
			ContractID:      contractID,
			GroupID:         groupID,
			PropertyID:      propertyID,
			PropertyVersion: propertyVersion,
		}).Return(&papi.ListAvailableIncludesResponse{
			AvailableIncludes: createExternalIncludeData(includesNumber)}, nil).Times(timesToRun)
	}

	// mockGetInclude mocks GetInclude call with provided parameters
	mockGetInclude = func(m *papi.Mock, contractID, groupID, includeID string, timesToRun int, include papi.Include) {
		m.On("GetInclude", testutils.MockContext, papi.GetIncludeRequest{
			ContractID: contractID,
			GroupID:    groupID,
			IncludeID:  includeID,
		}).Return(&papi.GetIncludeResponse{
			Includes: papi.IncludeItems{
				Items: []papi.Include{include},
			},
		}, nil).Times(timesToRun)
	}
)

// createIncludeTestCheckFuncs creates resource.TestCheckFunc for particular includes
func createIncludeTestCheckFuncs(filteredIncludes []papi.Include) []resource.TestCheckFunc {
	var testCheckIncludeFunc []resource.TestCheckFunc
	for i, include := range filteredIncludes {
		testCheckIncludeFunc = append(testCheckIncludeFunc, resource.TestCheckResourceAttr("data.akamai_property_includes.test", fmt.Sprintf("includes.%d.latest_version", i), strconv.Itoa(include.LatestVersion)))
		testCheckIncludeFunc = append(testCheckIncludeFunc, resource.TestCheckResourceAttr("data.akamai_property_includes.test", fmt.Sprintf("includes.%d.id", i), include.IncludeID))
		testCheckIncludeFunc = append(testCheckIncludeFunc, resource.TestCheckResourceAttr("data.akamai_property_includes.test", fmt.Sprintf("includes.%d.name", i), include.IncludeName))
		testCheckIncludeFunc = append(testCheckIncludeFunc, resource.TestCheckResourceAttr("data.akamai_property_includes.test", fmt.Sprintf("includes.%d.type", i), string(include.IncludeType)))
		if include.StagingVersion != nil {
			testCheckIncludeFunc = append(testCheckIncludeFunc, resource.TestCheckResourceAttr("data.akamai_property_includes.test", fmt.Sprintf("includes.%d.staging_version", i), strconv.Itoa(*include.StagingVersion)))
		} else {
			testCheckIncludeFunc = append(testCheckIncludeFunc, resource.TestCheckResourceAttr("data.akamai_property_includes.test", fmt.Sprintf("includes.%d.staging_version", i), ""))
		}
		if include.ProductionVersion != nil {
			testCheckIncludeFunc = append(testCheckIncludeFunc, resource.TestCheckResourceAttr("data.akamai_property_includes.test", fmt.Sprintf("includes.%d.production_version", i), strconv.Itoa(*include.ProductionVersion)))
		} else {
			testCheckIncludeFunc = append(testCheckIncludeFunc, resource.TestCheckResourceAttr("data.akamai_property_includes.test", fmt.Sprintf("includes.%d.production_version", i), ""))
		}
	}

	return testCheckIncludeFunc
}

// createDataSourceIDForTests creates ID of a data source based on provided attributes
func createDataSourceIDForTests(attrs attributes) string {
	var idElements []string
	idElements = append(idElements, attrs.contractID)
	idElements = append(idElements, attrs.groupID)

	if attrs.includeType != "" {
		idElements = append(idElements, attrs.includeType)
	}

	if attrs.parentProperty != nil {
		idElements = append(idElements, attrs.parentProperty.id)
		idElements = append(idElements, strconv.Itoa(attrs.parentProperty.version))
	}

	return strings.Join(idElements, ":")
}
