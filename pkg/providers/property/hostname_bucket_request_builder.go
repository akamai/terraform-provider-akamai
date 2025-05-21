package property

import (
	"context"
	"fmt"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/str"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	// hostnameLimit is the maximum amount of hostnames that can be passed to the Patch request. It includes both the added
	// and removed hostnames.
	hostnameLimit int = 1000

	// actionAdd describes the hostname to be added.
	actionAdd actionType = "ADD"

	// actionRemove describes the hostname to be removed.
	actionRemove actionType = "REMOVE"
)

type (
	// hostnameRequestBuilder builds the PATCH requests based on the diff between state and plan hostnames.
	hostnameRequestBuilder struct {
		hostnameRequestData hostnameRequestData
		diags               diag.Diagnostics
		ctx                 context.Context
	}

	hostnameRequestData struct {
		planHostnames  map[string]Hostname
		stateHostnames map[string]Hostname
		propertyID     string
		contractID     string
		groupID        string
		network        string
		note           string
		emails         []string
		requests       []papi.PatchPropertyHostnameBucketRequest
	}

	// actionType represents whether the hostname is added or removed.
	actionType string

	// hostnameOp gathers details about changes to the hostname.
	hostnameOp struct {
		cnameFrom            string
		certProvisioningType string
		edgeHostnameID       string
		action               actionType
	}
)

// toLog is a helper function that convert hostnameOp to more readable format for the logs.
func (h hostnameOp) toLog() map[string]any {
	return map[string]any{
		"cname_from":             h.cnameFrom,
		"edge_hostname_id":       h.edgeHostnameID,
		"cert_provisioning_type": h.certProvisioningType,
		"action":                 string(h.action),
	}
}

func newRequestBuilder(ctx context.Context, data HostnameBucketResourceModel) *hostnameRequestBuilder {
	var emails []string
	diags := data.NotifyEmails.ElementsAs(ctx, &emails, false)

	return &hostnameRequestBuilder{
		hostnameRequestData: hostnameRequestData{
			propertyID:     data.PropertyID.ValueString(),
			contractID:     data.ContractID.ValueString(),
			groupID:        data.GroupID.ValueString(),
			network:        data.Network.ValueString(),
			note:           data.Note.ValueString(),
			emails:         emails,
			stateHostnames: map[string]Hostname{},
			planHostnames:  map[string]Hostname{},
		},
		diags: diags,
		ctx:   ctx,
	}
}

func (h *hostnameRequestBuilder) setPlanHostnames(p types.Map) *hostnameRequestBuilder {
	h.diags.Append(p.ElementsAs(h.ctx, &h.hostnameRequestData.planHostnames, false)...)

	return h
}

func (h *hostnameRequestBuilder) setStateHostnames(s types.Map) *hostnameRequestBuilder {
	h.diags.Append(s.ElementsAs(h.ctx, &h.hostnameRequestData.stateHostnames, false)...)

	return h
}

func (h *hostnameRequestBuilder) build() (hostnameRequestData, diag.Diagnostics) {
	h.hostnameRequestData.requests = h.buildRequests()
	return h.hostnameRequestData, h.diags
}

func (h *hostnameRequestBuilder) buildRequests() []papi.PatchPropertyHostnameBucketRequest {
	// We need to compare plan hostnames with state hostnames to determine changes.
	// Loop through the hostnames in the plan, if such key is not present in the state, append it to the add queue.
	// If there is already such key present in the state, compare the two hostnames.
	// If they are not equal, it means that the first one needs to be
	// removed in a separate request prior to adding the same hostname with different values.
	var priorityHostnameOps []hostnameOp
	var hostnameOps []hostnameOp

	for key, value := range h.hostnameRequestData.planHostnames {
		if stateHostname, ok := h.hostnameRequestData.stateHostnames[key]; !ok {
			op := newHostnameOp(key, value, actionAdd)
			hostnameOps = append(hostnameOps, op)
			tflog.Debug(h.ctx, "buildRequests: adding new hostname", op.toLog())
		} else if !value.equal(stateHostname) {
			opRemove := newHostnameOp(key, value, actionRemove)
			opAdd := newHostnameOp(key, value, actionAdd)
			tflog.Debug(h.ctx, "buildRequests: removing existing hostname", opRemove.toLog())
			tflog.Debug(h.ctx, "buildRequests: adding new hostname", opAdd.toLog())
			priorityHostnameOps = append(priorityHostnameOps, opRemove)
			hostnameOps = append(hostnameOps, opAdd)
		}
	}

	// Loop through the hostnames in the state, if such key is not present in the plan, append it to the remove queue.
	for key, value := range h.hostnameRequestData.stateHostnames {
		if _, ok := h.hostnameRequestData.planHostnames[key]; !ok {
			op := newHostnameOp(key, value, actionRemove)
			tflog.Debug(h.ctx, "buildRequests: removing old hostname", op.toLog())
			hostnameOps = append(hostnameOps, op)
		}
	}

	var requests []papi.PatchPropertyHostnameBucketRequest
	// Handle the hostnames that need to be removed first in order to be updated.
	requests = append(requests, h.buildPatchRequests(priorityHostnameOps)...)
	// Handle the rest of removals and additions.
	requests = append(requests, h.buildPatchRequests(hostnameOps)...)

	return requests
}

func newHostnameOp(cnameFrom string, h Hostname, action actionType) hostnameOp {
	return hostnameOp{
		cnameFrom:            cnameFrom,
		certProvisioningType: h.CertProvisioningType.ValueString(),
		edgeHostnameID:       h.EdgeHostnameID.ValueString(),
		action:               action,
	}
}

// buildPatchRequests creates Patch requests based on the hostnames to add and hostnames to remove. It minimizes the
// number of requests, combining both operations into single request if possible.
func (h *hostnameRequestBuilder) buildPatchRequests(ops []hostnameOp) []papi.PatchPropertyHostnameBucketRequest {
	tflog.Debug(h.ctx, fmt.Sprintf("in 'buildPatchRequests' method with %d remaining operations", len(ops)))

	sortHostnameOps(ops)
	var requests []papi.PatchPropertyHostnameBucketRequest

	for len(ops) > 0 {
		var hostnamesToAdd []papi.PatchPropertyHostnameBucketAdd
		var hostnamesToRemove []string

		addCount := min(hostnameLimit, len(ops))
		for _, op := range ops[:addCount] {
			if op.action == actionAdd {
				hostnamesToAdd = append(hostnamesToAdd, papi.PatchPropertyHostnameBucketAdd{
					EdgeHostnameID:       str.AddPrefix(op.edgeHostnameID, "ehn_"),
					CertProvisioningType: papi.CertType(op.certProvisioningType),
					CnameType:            papi.HostnameCnameTypeEdgeHostname,
					CnameFrom:            op.cnameFrom,
				})
			} else {
				hostnamesToRemove = append(hostnamesToRemove, op.cnameFrom)
			}
		}
		req := papi.PatchPropertyHostnameBucketRequest{
			PropertyID: str.AddPrefix(h.hostnameRequestData.propertyID, "prp_"),
			ContractID: str.AddPrefix(h.hostnameRequestData.contractID, "ctr_"),
			GroupID:    str.AddPrefix(h.hostnameRequestData.groupID, "grp_"),
			Body: papi.PatchPropertyHostnameBucketBody{
				Add:          hostnamesToAdd,
				Remove:       hostnamesToRemove,
				Network:      papi.ActivationNetwork(h.hostnameRequestData.network),
				NotifyEmails: h.hostnameRequestData.emails,
				Note:         h.hostnameRequestData.note,
			},
		}
		requests = append(requests, req)
		ops = ops[addCount:]
		tflog.Debug(h.ctx, "buildPatchRequests: request built", map[string]interface{}{
			"request":              req,
			"remaining_operations": ops,
		})
	}

	return requests
}

// sortHostnameOps sorts hostname operations using lexicographical ordering. The key is logically
// a tuple (action, cnameFrom) where add operations have precedence over remove operations. To achieve continuity,
// it might be beneficial to add new hostnames first and then delete the old ones, as they may point to the same domains
// via edge hostnames.
func sortHostnameOps(ops []hostnameOp) {
	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].cnameFrom < ops[j].cnameFrom
	})
	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].action != ops[j].action && ops[i].action == actionAdd
	})
}
