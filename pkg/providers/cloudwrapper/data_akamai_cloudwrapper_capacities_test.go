package cloudwrapper

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudwrapper"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCapacitiesDataSource(t *testing.T) {
	t.Parallel()

	contractIDs := []string{"ctr_123", "345"}
	contractIDsTrimmed := []string{"123", "345"}
	request := cloudwrapper.ListCapacitiesRequest{
		ContractIDs: contractIDsTrimmed,
	}

	testCases := map[string]struct {
		respData      []cloudwrapper.LocationCapacity
		init          func(*testing.T, *cloudwrapper.Mock, []cloudwrapper.LocationCapacity)
		expectedError *regexp.Regexp
	}{
		"listing capacities successful": {
			respData: []cloudwrapper.LocationCapacity{
				{
					LocationID:   1,
					LocationName: "US West",
					Type:         cloudwrapper.CapacityTypeMedia,
					ContractID:   "ctr_123",
					ApprovedCapacity: cloudwrapper.Capacity{
						Value: 2000,
						Unit:  cloudwrapper.UnitGB,
					},
					UnassignedCapacity: cloudwrapper.Capacity{
						Value: 2000,
						Unit:  cloudwrapper.UnitGB,
					},
					AssignedCapacity: cloudwrapper.Capacity{
						Value: 0,
						Unit:  cloudwrapper.UnitGB,
					},
				},
				{
					LocationID:   2,
					LocationName: "US East",
					Type:         cloudwrapper.CapacityTypeMedia,
					ContractID:   "ctr_345",
					ApprovedCapacity: cloudwrapper.Capacity{
						Value: 4000,
						Unit:  cloudwrapper.UnitGB,
					},
					UnassignedCapacity: cloudwrapper.Capacity{
						Value: 2000,
						Unit:  cloudwrapper.UnitGB,
					},
					AssignedCapacity: cloudwrapper.Capacity{
						Value: 2000,
						Unit:  cloudwrapper.UnitGB,
					},
				},
			},
			init: func(t *testing.T, m *cloudwrapper.Mock, capacities []cloudwrapper.LocationCapacity) {
				resp := cloudwrapper.ListCapacitiesResponse{
					Capacities: capacities,
				}
				m.On("ListCapacities", mock.Anything, request).Return(&resp, nil).Times(5)
			},
		},
		"no capacities found": {
			respData: []cloudwrapper.LocationCapacity{},
			init: func(t *testing.T, m *cloudwrapper.Mock, capacities []cloudwrapper.LocationCapacity) {
				resp := cloudwrapper.ListCapacitiesResponse{
					Capacities: capacities,
				}
				m.On("ListCapacities", mock.Anything, request).Return(&resp, nil).Times(5)
			},
		},
		"listing capacities failed": {
			init: func(t *testing.T, m *cloudwrapper.Mock, _ []cloudwrapper.LocationCapacity) {
				err := fmt.Errorf("listing capacities failed")
				m.On("ListCapacities", mock.Anything, request).Return(nil, err).Once()
			},
			expectedError: regexp.MustCompile("listing capacities failed"),
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := &cloudwrapper.Mock{}
			if tc.init != nil {
				tc.init(t, client, tc.respData)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: newProviderFactory(withMockClient(client)),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      dataCloudWrapperCapacityConfig(contractIDs),
						Check:       checkDataCloudWrapperCapacityAttrs(contractIDs, tc.respData),
						ExpectError: tc.expectedError,
					},
				},
			})
			client.AssertExpectations(t)
		})
	}
}

func checkDataCloudWrapperCapacityAttrs(contractIDs []string, capacities []cloudwrapper.LocationCapacity) resource.TestCheckFunc {
	name := "data.akamai_cloudwrapper_capacities.test"
	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, "contract_ids.#", strconv.Itoa(len(contractIDs))))
	for i, ctr := range contractIDs {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("contract_ids.%d", i), ctr))
	}

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(name, "capacities.#", strconv.Itoa(len(capacities))))
	for i, cap := range capacities {
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.location_id", i), strconv.Itoa(cap.LocationID)),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.location_name", i), cap.LocationName),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.contract_id", i), cap.ContractID),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.type", i), string(cap.Type)),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.approved.value", i), strconv.FormatInt(cap.ApprovedCapacity.Value, 10)),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.approved.unit", i), string(cap.ApprovedCapacity.Unit)),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.assigned.value", i), strconv.FormatInt(cap.AssignedCapacity.Value, 10)),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.assigned.unit", i), string(cap.AssignedCapacity.Unit)),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.unassigned.value", i), strconv.FormatInt(cap.UnassignedCapacity.Value, 10)),
			resource.TestCheckResourceAttr(name, fmt.Sprintf("capacities.%d.unassigned.unit", i), string(cap.UnassignedCapacity.Unit)),
		)
	}

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttrSet(name, "id"))
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func dataCloudWrapperCapacityConfig(contracts []string) string {
	var contractIDs string
	if len(contracts) > 0 {
		contractIDs = fmt.Sprintf(`"%s"`, strings.Join(contracts, `", "`))
	}

	return fmt.Sprintf(`
data "akamai_cloudwrapper_capacities" "test" {
	contract_ids = [%s]
}`, contractIDs)
}

func TestCapacitiesDataSourceModel_getContractIDs(t *testing.T) {
	t.Parallel()

	newCapacitiesDataSourceModelWithContractIDs := func(t *testing.T, ctrIDs []string) capacitiesDataSourceModel {
		t.Helper()

		ctrList, diags := types.ListValueFrom(context.Background(), types.StringType, ctrIDs)
		require.False(t, diags.HasError(), "converting []string to types.List failed with: %s", diags)

		return capacitiesDataSourceModel{
			ContractIDs: ctrList,
		}
	}

	testCases := map[string]struct {
		createModel func(*testing.T) capacitiesDataSourceModel
		want        []string
		withError   bool
	}{
		"returns all contract IDs": {
			createModel: func(t *testing.T) capacitiesDataSourceModel {
				contractIDs := []string{"123", "345", "567"}
				return newCapacitiesDataSourceModelWithContractIDs(t, contractIDs)
			},
			want: []string{"123", "345", "567"},
		},
		"trims 'ctr_' prefix": {
			createModel: func(t *testing.T) capacitiesDataSourceModel {
				contractIDs := []string{"ctr_123", "345", "ctr_567"}
				return newCapacitiesDataSourceModelWithContractIDs(t, contractIDs)
			},
			want: []string{"123", "345", "567"},
		},
		"returns empty slice when contract IDs empty": {
			createModel: func(t *testing.T) capacitiesDataSourceModel {
				return newCapacitiesDataSourceModelWithContractIDs(t, []string{})
			},
			want: []string{},
		},
		"returns nil when contract IDs are nil": {
			createModel: func(t *testing.T) capacitiesDataSourceModel {
				return newCapacitiesDataSourceModelWithContractIDs(t, nil)
			},
			want: nil,
		},
		"fails when contractIDs are of wrong type": {
			createModel: func(t *testing.T) capacitiesDataSourceModel {
				contractIDs := []int{123, 345, 567}
				ctrList, diags := types.ListValueFrom(context.Background(), types.Int64Type, contractIDs)
				require.False(t, diags.HasError(), "converting []int to types.List failed with: %s", diags)

				return capacitiesDataSourceModel{
					ContractIDs: ctrList,
				}
			},
			withError: true,
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			model := tc.createModel(t)

			got, diags := model.getContractIDs(context.Background())
			if tc.withError {
				assert.True(t, diags.HasError(), "diag.Diagnostics does not contain error")
				return
			}

			require.False(t, diags.HasError(), "returned diag.Diagnostics contain error")
			assert.Equal(t, tc.want, got)
		})
	}
}
