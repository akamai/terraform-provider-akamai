package property

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/papi"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type requestStructure struct {
	add    int
	remove int
}

func TestHostnameRequestBuilder(t *testing.T) {
	tests := map[string]struct {
		plan                     map[string]Hostname
		state                    map[string]Hostname
		expectedRequestStructure map[int]requestStructure
	}{
		"only to add": {
			plan: map[string]Hostname{
				"test.cnameFrom.1.com": {
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
				},
				"test.cnameFrom.2.com": {
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
				},
			},
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    2,
					remove: 0,
				},
			},
		},
		"only to remove": {
			state: map[string]Hostname{
				"test.cnameFrom.1.com": {
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
				},
				"test.cnameFrom.2.com": {
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
				},
			},
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    0,
					remove: 2,
				},
			},
		},
		"add and remove, below 1000": {
			state: map[string]Hostname{
				"test.cnameFrom.1.com": {
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
				},
				"test.cnameFrom.2.com": {
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
				},
			},
			plan: map[string]Hostname{
				"test.cnameFrom.3.com": {
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
				},
				"test.cnameFrom.4.com": {
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
				},
			},
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    2,
					remove: 2,
				},
			},
		},
		"add 1000 - expect 1 request": {
			plan: generateHostnames(1000, "CPS_MANAGED", "ehn_444"),
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    1000,
					remove: 0,
				},
			},
		},
		"add 1001 - expect 2 requests": {
			plan: generateHostnames(1001, "CPS_MANAGED", "ehn_444"),
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    1000,
					remove: 0,
				},
				1: {
					add:    1,
					remove: 0,
				},
			},
		},
		"remove 1000 - expect 1 request": {
			state: generateHostnames(1000, "CPS_MANAGED", "ehn_444"),
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    0,
					remove: 1000,
				},
			},
		},
		"remove 1001 - expect 2 requests": {
			state: generateHostnames(1001, "CPS_MANAGED", "ehn_444"),
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    0,
					remove: 1000,
				},
				1: {
					add:    0,
					remove: 1,
				},
			},
		},
		"have 3000, update 2000, remove 1000 - expect 5 requests": {
			state: generateHostnames(3000, "CPS_MANAGED", "ehn_444"),
			plan:  generateHostnames(2000, "DEFAULT", "ehn_444"),
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    0,
					remove: 1000,
				},
				1: {
					add:    0,
					remove: 1000,
				},
				2: {
					add:    1000,
					remove: 0,
				},
				3: {
					add:    1000,
					remove: 0,
				},
				4: {
					add:    0,
					remove: 1000,
				},
			},
		},
		"have 3000, update 2500, remove 1000 - expect 5 requests": {
			state: generateHostnames(3000, "CPS_MANAGED", "ehn_444"),
			plan:  generateHostnames(2500, "DEFAULT", "ehn_444"),
			expectedRequestStructure: map[int]requestStructure{
				0: {
					add:    0,
					remove: 1000,
				},
				1: {
					add:    0,
					remove: 1000,
				},
				2: {
					add:    0,
					remove: 500,
				},
				3: {
					add:    1000,
					remove: 0,
				},
				4: {
					add:    1000,
					remove: 0,
				},
				5: {
					add:    500,
					remove: 500,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b := hostnameRequestBuilder{
				hostnameRequestData: hostnameRequestData{
					planHostnames:  tc.plan,
					stateHostnames: tc.state,
				},
				ctx: context.Background(),
			}
			requestData, _ := b.build()

			assert.Equal(t, len(tc.expectedRequestStructure), len(requestData.requests))
			for k, v := range tc.expectedRequestStructure {
				assert.Equal(t, v.add, len(requestData.requests[k].Body.Add), "expected the request %d to have %d additions, but have %d", k, v.add, len(requestData.requests[k].Body.Add))
				assert.Equal(t, v.remove, len(requestData.requests[k].Body.Remove), "expected the request %d to have %d removals, but have %d", k, v.remove, len(requestData.requests[k].Body.Remove))
			}
		})
	}

	t.Run("check request content", func(t *testing.T) {
		b := hostnameRequestBuilder{
			hostnameRequestData: hostnameRequestData{
				propertyID: "prp_111",
				contractID: "ctr_222",
				groupID:    "grp_333",
				network:    "STAGING",
				note:       "Note",
				emails:     []string{"test@mail.com"},
				planHostnames: map[string]Hostname{
					"test.cnameFrom.1.com": {
						CertProvisioningType: types.StringValue("CPS_MANAGED"),
						EdgeHostnameID:       types.StringValue("ehn_444"),
					},
				},
			},
			ctx: context.Background(),
		}
		rd, _ := b.build()
		assert.Equal(t, "prp_111", rd.requests[0].PropertyID)
		assert.Equal(t, "ctr_222", rd.requests[0].ContractID)
		assert.Equal(t, "grp_333", rd.requests[0].GroupID)
		assert.Equal(t, papi.ActivationNetworkStaging, rd.requests[0].Body.Network)
		assert.Equal(t, "Note", rd.requests[0].Body.Note)
		assert.Equal(t, "test@mail.com", rd.requests[0].Body.NotifyEmails[0])
		assert.Equal(t, papi.CertTypeCPSManaged, rd.requests[0].Body.Add[0].CertProvisioningType)
		assert.Equal(t, "ehn_444", rd.requests[0].Body.Add[0].EdgeHostnameID)
	})
}

func TestSortHostnameOps(t *testing.T) {
	ops := []hostnameOp{
		{action: actionAdd, cnameFrom: "kilo.edgesuite.net"},
		{action: actionRemove, cnameFrom: "foxtrot.edgesuite.net"},
		{action: actionAdd, cnameFrom: "golf.edgesuite.net"},
		{action: actionRemove, cnameFrom: "mike.edgesuite.net"},
		{action: actionRemove, cnameFrom: "oscar.edgesuite.net"},
		{action: actionAdd, cnameFrom: "lima.edgesuite.net"},
		{action: actionRemove, cnameFrom: "charlie.edgesuite.net"},
		{action: actionAdd, cnameFrom: "november.edgesuite.net"},
		{action: actionAdd, cnameFrom: "echo.edgesuite.net"},
		{action: actionAdd, cnameFrom: "india.edgesuite.net"},
		{action: actionRemove, cnameFrom: "bravo.edgesuite.net"},
		{action: actionRemove, cnameFrom: "alpha.edgesuite.net"},
		{action: actionAdd, cnameFrom: "delta.edgesuite.net"},
		{action: actionRemove, cnameFrom: "hotel.edgesuite.net"},
		{action: actionAdd, cnameFrom: "juliett.edgesuite.net"},
	}

	expected := []hostnameOp{
		{action: actionAdd, cnameFrom: "delta.edgesuite.net"},
		{action: actionAdd, cnameFrom: "echo.edgesuite.net"},
		{action: actionAdd, cnameFrom: "golf.edgesuite.net"},
		{action: actionAdd, cnameFrom: "india.edgesuite.net"},
		{action: actionAdd, cnameFrom: "juliett.edgesuite.net"},
		{action: actionAdd, cnameFrom: "kilo.edgesuite.net"},
		{action: actionAdd, cnameFrom: "lima.edgesuite.net"},
		{action: actionAdd, cnameFrom: "november.edgesuite.net"},
		{action: actionRemove, cnameFrom: "alpha.edgesuite.net"},
		{action: actionRemove, cnameFrom: "bravo.edgesuite.net"},
		{action: actionRemove, cnameFrom: "charlie.edgesuite.net"},
		{action: actionRemove, cnameFrom: "foxtrot.edgesuite.net"},
		{action: actionRemove, cnameFrom: "hotel.edgesuite.net"},
		{action: actionRemove, cnameFrom: "mike.edgesuite.net"},
		{action: actionRemove, cnameFrom: "oscar.edgesuite.net"},
	}

	t.Run("sorting hostname ops", func(t *testing.T) {
		sortHostnameOps(ops)
		assert.Equal(t, expected, ops)
	})
}
