package property

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &HostnameBucketResource{}
	_ resource.ResourceWithImportState = &HostnameBucketResource{}
	_ resource.ResourceWithConfigure   = &HostnameBucketResource{}
	_ resource.ResourceWithModifyPlan  = &HostnameBucketResource{}
)

// HostnameBucketResource represents akamai_property_hostname_bucket resource.
type HostnameBucketResource struct {
	meta meta.Meta
}

// HostnameBucketResourceModel is a model for akamai_property_hostname_bucket resource.
type HostnameBucketResourceModel struct {
	PropertyID           types.String `tfsdk:"property_id"`
	ContractID           types.String `tfsdk:"contract_id"`
	GroupID              types.String `tfsdk:"group_id"`
	Network              types.String `tfsdk:"network"`
	Note                 types.String `tfsdk:"note"`
	NotifyEmails         types.List   `tfsdk:"notify_emails"`
	ActivationID         types.String `tfsdk:"activation_id"`
	ID                   types.String `tfsdk:"id"`
	Hostnames            types.Map    `tfsdk:"hostnames"`
	TimeoutForActivation types.Int64  `tfsdk:"timeout_for_activation"`
	HostnameCount        types.Int64  `tfsdk:"hostname_count"`
	PendingDefaultCerts  types.Int64  `tfsdk:"pending_default_certs"`
}

func (m *HostnameBucketResourceModel) haveHostnamesChanged(plan *HostnameBucketResourceModel) bool {
	return !m.Hostnames.Equal(plan.Hostnames)
}

func (m *HostnameBucketResourceModel) hasTimeoutChanged(plan *HostnameBucketResourceModel) bool {
	return !m.TimeoutForActivation.Equal(plan.TimeoutForActivation)
}

func (m *HostnameBucketResourceModel) hasNoteChanged(plan *HostnameBucketResourceModel) bool {
	return !m.Note.Equal(plan.Note)
}

func (m *HostnameBucketResourceModel) haveEmailsChanged(plan *HostnameBucketResourceModel) bool {
	return !m.NotifyEmails.Equal(plan.NotifyEmails)
}

func (m *HostnameBucketResourceModel) doesPendingDefaultCertsTriggerDiff(plan *HostnameBucketResourceModel) bool {
	return plan.PendingDefaultCerts.IsUnknown() && !m.PendingDefaultCerts.IsNull()
}

func (m *HostnameBucketResourceModel) doesActivationIDTriggerDiff(plan *HostnameBucketResourceModel) bool {
	return plan.ActivationID.IsUnknown() && m.ActivationID.ValueString() != ""
}

func (m *HostnameBucketResourceModel) isGroupIDDefined() bool {
	return !m.GroupID.IsUnknown() && !m.GroupID.IsNull()
}

// setIDs sets the GroupID and ContractID from the response of ListActivePropertyHostnames, if the values
// were not provided in the configuration.
func (m *HostnameBucketResourceModel) setIDs(responses []papi.ListActivePropertyHostnamesResponse) error {
	var groupID, contractID string
	if len(responses) == 0 {
		return fmt.Errorf("there are no responses from ListActivePropertyHostnames endpoint")
	}
	groupID = responses[0].GroupID
	contractID = responses[0].ContractID

	m.GroupID = getStringValueWithPrefixOrNull(groupID, "grp_")
	m.ContractID = getStringValueWithPrefixOrNull(contractID, "ctr_")

	return nil
}

func (m *HostnameBucketResourceModel) sendRequests(ctx context.Context, client papi.PAPI, requests []papi.PatchPropertyHostnameBucketRequest,
	waitFunc func(context.Context, papi.PAPI, HostnameBucketResourceModel) error) error {
	for i, r := range requests {
		tflog.Debug(ctx, "sending patch requests", map[string]interface{}{
			fmt.Sprintf("%d", i): r,
		})
		response, err := client.PatchPropertyHostnameBucket(ctx, r)
		if err != nil {
			return err
		}

		// Overwrite activationID with the latest successful result.
		tflog.Debug(ctx, "set new activation_id value", map[string]interface{}{
			"previous_activation_id": m.ActivationID.ValueString(),
			"new_activation_id":      response.ActivationID,
			"new_activation_link":    response.ActivationLink,
		})
		m.ActivationID = types.StringValue(response.ActivationID)

		if err = waitFunc(ctx, client, *m); err != nil {
			return err
		}
	}
	return nil
}

// Hostname represents a hostname object.
type Hostname struct {
	CertProvisioningType types.String `tfsdk:"cert_provisioning_type"`
	EdgeHostnameID       types.String `tfsdk:"edge_hostname_id"`
	CnameTo              types.String `tfsdk:"cname_to"`
}

// equal compares two hostnames for their equality. The hostnames are stores as a map, where a key is the `cname_from` attribute.
// To identify if two hostnames are the same, we need to compare the rest of required attributes, that being `edge_hostname_id`
// and `cert_provisioning_type`.
func (h Hostname) equal(other Hostname) bool {
	return h.CertProvisioningType == other.CertProvisioningType && h.EdgeHostnameID == other.EdgeHostnameID
}

func (h Hostname) toLog() map[string]any {
	return map[string]any{
		"cname_to":               h.CnameTo.ValueString(),
		"edge_hostname_id":       h.EdgeHostnameID.ValueString(),
		"cert_provisioning_type": h.CertProvisioningType.ValueString(),
	}
}

func getStringValueWithPrefixOrNull(val, pre string) basetypes.StringValue {
	if val == "" {
		return types.StringNull()
	}

	return types.StringValue(str.AddPrefix(val, pre))
}

var (
	// getHostnameBucketActivationInterval is the time interval after which consecutive requests are being sent.
	getHostnameBucketActivationInterval = time.Second * 30
	// forceTimeoutDuration is used to overwrite `timeout_for_activation` for unit tests.
	forceTimeoutDuration time.Duration
	// errCancelActivation is returned when an activation has been cancelled.
	errCancelActivation = errors.New("timeout has been reached; activation has been cancelled")
	// hostnameObjectType represents the object inside the 'hostnames' attribute.
	hostnameObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"cert_provisioning_type": types.StringType,
			"edge_hostname_id":       types.StringType,
			"cname_to":               types.StringType,
		},
	}
)

const (
	// maxHostnamesNumber is the maximum amount of hostnames that can be configured for a hostname bucket.
	maxHostnamesNumber int = 99999
	// activationTimeout is the default timeout value for the hostname activation.
	activationTimeout int64 = 50
)

// NewHostnameBucketResource returns new property hostname bucket resource.
func NewHostnameBucketResource() resource.Resource {
	return &HostnameBucketResource{}
}

// Metadata implements resource.Resource.
func (h *HostnameBucketResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_property_hostname_bucket"
}

// Configure implements resource.ResourceWithConfigure.
func (h *HostnameBucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	h.meta = meta.Must(req.ProviderData)
}

// Schema implements resource's Schema.
func (h *HostnameBucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"property_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("prp_")),
					modifiers.PreventStringUpdate(),
				},
				Description: "The unique identifier for the property.",
			},
			"contract_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					stringplanmodifier.UseStateForUnknown(),
					modifiers.PreventStringUpdateIfKnown("contract_id"),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{path.MatchRoot("group_id")}...),
				},
				Description: "The unique identifier for the contract. Provide it if resolving the property without 'contract_id' and 'group_id' is not possible",
			},
			"group_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("grp_")),
					stringplanmodifier.UseStateForUnknown(),
					modifiers.PreventStringUpdateIfKnown("group_id"),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{path.MatchRoot("contract_id")}...),
				},
				Description: "The unique identifier for the group. Provide it if resolving the property without 'contract_id' and 'group_id' is not possible",
			},
			"network": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(papi.ActivationNetworkStaging), string(papi.ActivationNetworkProduction)),
				},
				Description: "The network to activate on, either `STAGING` or `PRODUCTION`.",
			},
			"note": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// Default value of '   ' (3 spaces) is assigned by the API when we don't send any value.
				Default:     stringdefault.StaticString("   "),
				Description: "Assigns a log message to the request.",
			},
			"notify_emails": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType, []attr.Value{types.StringValue("nomail@akamai.com")})),
				Description: "Email addresses to notify when the activation status changes.",
			},
			"activation_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the latest hostname bucket activation.",
			},
			"timeout_for_activation": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(activationTimeout),
				Description: "The timeout value in minutes after which a single hostname activation will be canceled. " +
					"Defaults to 50 minutes.",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The resource ID in the <property_id:network> format.",
			},
			"hostname_count": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "The computed number of hostnames after applying desired modifications. Used only to inform" +
					"during the plan phase about the number of hostnames that will be active after making the changes.",
			},
			"hostnames": schema.MapNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cert_provisioning_type": schema.StringAttribute{
							Required: true,
							Description: "Indicates the type of the certificate used in the property hostname. " +
								"Either `CPS_MANAGED` for certificates you create with the Certificate Provisioning System (CPS) API, " +
								"or `DEFAULT` for Domain Validation (DV) certificates deployed automatically.",
							Validators: []validator.String{
								stringvalidator.OneOf(string(papi.CertTypeDefault), string(papi.CertTypeCPSManaged)),
							},
						},
						"edge_hostname_id": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ehn_")),
							},
							Description: "Identifies the edge hostname you mapped your traffic to on the production network.",
						},
						"cname_to": schema.StringAttribute{
							Computed: true,
							Description: "The edge hostname you point the property hostname to so that you can start serving traffic through Akamai servers. " +
								"This member corresponds to the edge hostname object's `edgeHostnameDomain` member.",
						},
					},
				},
				PlanModifiers: []planmodifier.Map{
					newHostnamesPlanModifier(),
				},
				Validators: []validator.Map{
					mapvalidator.SizeBetween(1, maxHostnamesNumber),
				},
				Description: "The hostnames mapping. The key represents 'cname_from' and the value contains hostnames details, " +
					"consisting of certificate provisioning type and edge hostname.",
			},
			"pending_default_certs": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of hostnames with a `DEFAULT` certificate type that are still in the `PENDING` state.",
			},
		},
	}
}

// ModifyPlan performs plan modification on a resource level.
func (h *HostnameBucketResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if request.Plan.Raw.IsNull() || request.State.Raw.IsNull() {
		return
	}

	var state, plan *HostnameBucketResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Get the hostnames to calculate their quantity and set in during planning in order to be visible.
	hostnames := plan.Hostnames.Elements()
	response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("hostname_count"), len(hostnames))...)
	if response.Diagnostics.HasError() {
		return
	}

	// Suppress changes to `timeout_for_activation` if it has changed, but the hostnames did not.
	if !state.haveHostnamesChanged(plan) && state.hasTimeoutChanged(plan) {
		tflog.Debug(ctx, "only 'timeout_for_activation' changed, using state value instead")
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("timeout_for_activation"), state.TimeoutForActivation.ValueInt64())...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	// Suppress changes to `pending_default_certs` on the terraform plan, if there are no changes to the hostnames.
	// For other updates, this field should be marked as "known after apply".
	if !state.haveHostnamesChanged(plan) && state.doesPendingDefaultCertsTriggerDiff(plan) {
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("pending_default_certs"), state.PendingDefaultCerts.ValueInt64())...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	// Suppress changes to `activation_id` on the terraform plan, if there are no changes to the hostnames.
	// For other updates, this field should be marked as "known after apply".
	if !state.haveHostnamesChanged(plan) && state.doesActivationIDTriggerDiff(plan) {
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("activation_id"), state.ActivationID.ValueString())...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	// Suppress changes to `notify_emails` if they have changes, but the hostnames did not. It will not be shown in the plan
	// and such change won't trigger any update.
	if !state.haveHostnamesChanged(plan) && state.haveEmailsChanged(plan) {
		tflog.Debug(ctx, "only 'notify_emails' changed, using state value instead")
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("notify_emails"), state.NotifyEmails)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	// Suppress changes to `note` if it has changed attribute, but the hostnames did not. It will not be shown in the plan
	// and such change won't trigger any update.
	if !state.haveHostnamesChanged(plan) && state.hasNoteChanged(plan) {
		tflog.Debug(ctx, "only 'note' changed, using state value instead")
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("note"), state.Note)...)
		if response.Diagnostics.HasError() {
			return
		}
	}
}

func newHostnamesPlanModifier() planmodifier.Map {
	return hostnamesPlanModifier{}
}

type hostnamesPlanModifier struct{}

func (m hostnamesPlanModifier) Description(_ context.Context) string {
	return "The hostnamesPlanModifier populates the `cname_to` attribute from the state value, in order to limit the diff for that field."
}

func (m hostnamesPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m hostnamesPlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.StateValue.IsNull() || req.PlanValue.IsNull() {
		return
	}

	var planHostnames map[string]Hostname
	resp.Diagnostics.Append(req.PlanValue.ElementsAs(ctx, &planHostnames, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateHostnames map[string]Hostname
	resp.Diagnostics.Append(req.StateValue.ElementsAs(ctx, &stateHostnames, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for cnameFrom, planHostname := range planHostnames {
		if stateHostname, ok := stateHostnames[cnameFrom]; ok {
			if planHostname.equal(stateHostname) {
				planHostname.CnameTo = stateHostname.CnameTo
				planHostnames[cnameFrom] = planHostname
			}
		}
	}

	planHostnamesValue, diags := types.MapValueFrom(ctx, hostnameObjectType, planHostnames)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = planHostnamesValue
}

// Create implements resource's Create method.
func (h *HostnameBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Property Hostname Bucket Resource")

	var plan HostnameBucketResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.PropertyID.ValueString(), plan.Network.ValueString()))
	ctx = tflog.SetField(ctx, "id", plan.ID.ValueString())

	requestsData, diags := newRequestBuilder(ctx, plan).
		setPlanHostnames(plan.Hostnames).
		build()
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(h.meta)

	// Send the PATCH requests to add the hostnames and wait for their activation.
	if err := plan.sendRequests(ctx, client, requestsData.requests, waitForHostnameBucketActivation); err != nil {
		resp.Diagnostics.AddError("Create Property Hostname Bucket error", err.Error())
		return
	}

	// After all PATCH requests, list the hostnames and fill the `cname_to` attributes.
	responses, err := listHostnamesResponses(ctx, client, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Create Property Hostname Bucket error", err.Error())
		return
	}
	hostnames := extractConcatenatedHostnames(responses)

	// Fill out `cname_to` attributes for each of the hostname.
	planHostnames := setHostnameDetails(ctx, hostnames, requestsData.planHostnames, plan.Network.ValueString(), false)
	planHostnamesValue, diags := types.MapValueFrom(ctx, hostnameObjectType, planHostnames)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// `group_id` and `contract_id` are not needed if it's possible to resolve the property without them. If it's not, they
	// must be provided. After Create, the values for those fields are taken from the API responses and set in state.
	if !plan.isGroupIDDefined() {
		if err = plan.setIDs(responses); err != nil {
			resp.Diagnostics.AddError("Create Property Hostname Bucket Resource error", err.Error())
			return
		}
	}
	plan.Hostnames = planHostnamesValue
	plan.HostnameCount = types.Int64Value(int64(len(planHostnames)))
	plan.PendingDefaultCerts = types.Int64Value(countPendingDefaultCerts(hostnames))

	if plan.PendingDefaultCerts.ValueInt64() > 0 {
		resp.Diagnostics.Append(formatWarning(plan.PendingDefaultCerts.ValueInt64()))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read implements resource's Read method.
func (h *HostnameBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Property Hostname Bucket Resource")

	var state HostnameBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "id", state.ID.ValueString())
	client := Client(h.meta)

	act, err := findCurrentActivation(ctx, client, &state)
	if err != nil {
		resp.Diagnostics.AddError("Read Property Hostname Bucket Resource error", err.Error())
		return
	}
	state.ActivationID = types.StringValue(act.HostnameActivationID)
	state.Note = types.StringValue(act.Note)
	emails, diags := types.ListValueFrom(ctx, types.StringType, act.NotifyEmails)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	state.NotifyEmails = emails

	responses, err := listHostnamesResponses(ctx, client, &state)
	if err != nil {
		resp.Diagnostics.AddError("Read Property Hostname Bucket Resource error", err.Error())
		return
	}
	hostnames := extractConcatenatedHostnames(responses)

	// If there are no hostnames in the API, remove the resource.
	if len(hostnames) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	var stateHostnames map[string]Hostname
	resp.Diagnostics.Append(state.Hostnames.ElementsAs(ctx, &stateHostnames, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize the stateHostnames map if empty.
	if stateHostnames == nil {
		stateHostnames = make(map[string]Hostname)
	}

	// Fill out `cname_to` attribute for each of the hostname and add missing hostname to the state if not present.
	stateHostnames = setHostnameDetails(ctx, hostnames, stateHostnames, state.Network.ValueString(), true)

	// We need to create a map from the API hostnames and check if there are some hostnames in the current state, but not in the API.
	// If they are, we need to delete them.
	apiHostnamesMap := make(map[string]bool)
	for _, hn := range hostnames {
		apiHostnamesMap[hn.CnameFrom] = true
	}
	for stateCnameFrom := range stateHostnames {
		if !apiHostnamesMap[stateCnameFrom] {
			// If there is no hostname in the API, but it is in the state, delete it.
			delete(stateHostnames, stateCnameFrom)
		}
	}

	stateHostnamesValue, diags := types.MapValueFrom(ctx, hostnameObjectType, stateHostnames)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if !state.isGroupIDDefined() {
		if err = state.setIDs(responses); err != nil {
			resp.Diagnostics.AddError("Read Property Hostname Bucket Resource error", err.Error())
			return
		}
	}
	state.Hostnames = stateHostnamesValue
	state.HostnameCount = types.Int64Value(int64(len(stateHostnames)))
	state.PendingDefaultCerts = types.Int64Value(countPendingDefaultCerts(hostnames))

	if state.PendingDefaultCerts.ValueInt64() > 0 {
		resp.Diagnostics.Append(formatWarning(state.PendingDefaultCerts.ValueInt64()))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements resource's Update method.
func (h *HostnameBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Property Hostname Bucket Resource")

	var plan, state HostnameBucketResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "id", state.ID.ValueString())

	requestsData, diags := newRequestBuilder(ctx, plan).
		setStateHostnames(state.Hostnames).
		setPlanHostnames(plan.Hostnames).
		build()
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(h.meta)

	if err := plan.sendRequests(ctx, client, requestsData.requests, waitForHostnameBucketActivation); err != nil {
		resp.Diagnostics.AddError("Update Property Hostname Bucket error", err.Error())
		return
	}

	// After all PATCH requests, list the hostnames and fill the `cname_to` attributes.
	responses, err := listHostnamesResponses(ctx, client, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Update Property Hostname Bucket error", err.Error())
		return
	}
	hostnames := extractConcatenatedHostnames(responses)

	planHostnames := setHostnameDetails(ctx, hostnames, requestsData.planHostnames, plan.Network.ValueString(), false)
	planHostnamesValue, diags := types.MapValueFrom(ctx, hostnameObjectType, planHostnames)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if !plan.isGroupIDDefined() {
		if err = plan.setIDs(responses); err != nil {
			resp.Diagnostics.AddError("Update Property Hostname Bucket Resource error", err.Error())
			return
		}
	}
	plan.Hostnames = planHostnamesValue
	plan.HostnameCount = types.Int64Value(int64(len(planHostnames)))
	plan.PendingDefaultCerts = types.Int64Value(countPendingDefaultCerts(hostnames))

	if plan.PendingDefaultCerts.ValueInt64() > 0 {
		resp.Diagnostics.Append(formatWarning(plan.PendingDefaultCerts.ValueInt64()))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource's Delete method.
func (h *HostnameBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Property Hostname Bucket Resource")

	var state HostnameBucketResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "id", state.ID.ValueString())

	requestsData, diags := newRequestBuilder(ctx, state).
		setPlanHostnames(types.MapNull(hostnameObjectType)).
		setStateHostnames(state.Hostnames).
		build()
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(h.meta)
	if err := state.sendRequests(ctx, client, requestsData.requests, waitForHostnameBucketDeletion); err != nil {
		resp.Diagnostics.AddError("Delete Property Hostname Bucket Resource error", err.Error())
		return
	}
}

// ImportState implements resource's ImportState method. The importID has a format of <propertyID:network[:contractID:groupID]>
// for example: "prp_123:STAGING", "prp_123:PRODUCTION:ctr_456:grp_789".
func (h *HostnameBucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Property Hostname Bucket Resource")

	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 && len(parts) != 4 {
		resp.Diagnostics.AddError("Import Property Hostname Bucket error", "importID must be of format 'property_id:network[:contract_id:group_id]'")
		return
	}

	prpID := parts[0]
	net := parts[1]
	if net != "STAGING" && net != "PRODUCTION" {
		resp.Diagnostics.AddError("Import Property Hostname Bucket error", "network must have correct value of 'STAGING' or 'PRODUCTION'")
		return
	}
	var ctrID, grpID string
	if len(parts) == 4 {
		ctrID = parts[2]
		grpID = parts[3]
	}

	// Set the state attributes so the Read method can process the resource.
	state := HostnameBucketResourceModel{
		PropertyID:           types.StringValue(prpID),
		Network:              types.StringValue(net),
		NotifyEmails:         types.ListUnknown(types.StringType),
		ID:                   types.StringValue(fmt.Sprintf("%s:%s", prpID, net)),
		Hostnames:            types.MapNull(hostnameObjectType),
		TimeoutForActivation: types.Int64Value(50),
		PendingDefaultCerts:  types.Int64Unknown(),
	}
	if ctrID != "" && grpID != "" {
		state.ContractID = types.StringValue(ctrID)
		state.GroupID = types.StringValue(grpID)
	} else {
		state.ContractID = types.StringUnknown()
		state.GroupID = types.StringUnknown()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// listHostnamesResponses returns a slice of ListActivePropertyHostnames responses.
func listHostnamesResponses(ctx context.Context, client papi.PAPI, data *HostnameBucketResourceModel) ([]papi.ListActivePropertyHostnamesResponse, error) {
	var result []papi.ListActivePropertyHostnamesResponse
	offset, limit := 0, 999
	for {
		res, err := client.ListActivePropertyHostnames(ctx, papi.ListActivePropertyHostnamesRequest{
			Offset:            offset,
			Limit:             limit,
			PropertyID:        data.PropertyID.ValueString(),
			Network:           papi.ActivationNetwork(data.Network.ValueString()),
			ContractID:        data.ContractID.ValueString(),
			GroupID:           data.GroupID.ValueString(),
			Sort:              papi.SortAscending,
			IncludeCertStatus: true,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, *res)

		offset += limit
		if offset >= res.Hostnames.TotalItems {
			break
		}
	}

	return result, nil
}

// extractConcatenatedHostnames extracts the hostnames from a slice of ListActivePropertyHostnames responses.
func extractConcatenatedHostnames(responses []papi.ListActivePropertyHostnamesResponse) []papi.HostnameItem {
	var hostnames []papi.HostnameItem
	for _, r := range responses {
		hostnames = append(hostnames, r.Hostnames.Items...)
	}

	return hostnames
}

// countPendingDefaultCerts counts the hostnames with a default certificate type that has a PENDING certificate status.
func countPendingDefaultCerts(hostnames []papi.HostnameItem) int64 {
	var counter int64
	for _, h := range hostnames {
		if isPendingDefaultCert(h) {
			counter++
		}
	}
	return counter
}

func isPendingDefaultCert(item papi.HostnameItem) bool {
	if item.CertStatus != nil && (item.ProductionCertType == papi.CertTypeDefault || item.StagingCertType == papi.CertTypeDefault) {
		return (len(item.CertStatus.Staging) > 0 && item.CertStatus.Staging[0].Status == "PENDING") ||
			(len(item.CertStatus.Production) > 0 && item.CertStatus.Production[0].Status == "PENDING")
	}

	return false
}

// setHostnameDetails fills out cname_to attributes for each of the hostnames and updates the state with missing hostnames
// if they are not present.
func setHostnameDetails(ctx context.Context, hostnames []papi.HostnameItem, result map[string]Hostname, network string, addMissing bool) map[string]Hostname {
	for _, hn := range hostnames {
		switch network {
		case "STAGING":
			if val, ok := result[hn.CnameFrom]; ok {
				val.CnameTo = types.StringValue(hn.StagingCnameTo)
				result[hn.CnameFrom] = val
				tflog.Debug(ctx, "read 'cname_to' value from the API", val.toLog(), map[string]any{"cname_from": hn.CnameFrom})
			} else if addMissing {
				// If there is a hostname in the API but not in the current state, add it.
				h := Hostname{
					CertProvisioningType: types.StringValue(string(hn.StagingCertType)),
					EdgeHostnameID:       types.StringValue(hn.StagingEdgeHostnameID),
					CnameTo:              types.StringValue(hn.StagingCnameTo),
				}
				result[hn.CnameFrom] = h
				tflog.Debug(ctx, "added hostname from API to the state", h.toLog(), map[string]any{"cname_from": hn.CnameFrom})
			}
		case "PRODUCTION":
			if val, ok := result[hn.CnameFrom]; ok {
				val.CnameTo = types.StringValue(hn.ProductionCnameTo)
				result[hn.CnameFrom] = val
				tflog.Debug(ctx, "read 'cname_to' value from the API", val.toLog(), map[string]any{"cname_from": hn.CnameFrom})
			} else if addMissing {
				// If there is a hostname in the API but not in the current state, add it.
				h := Hostname{
					CertProvisioningType: types.StringValue(string(hn.ProductionCertType)),
					EdgeHostnameID:       types.StringValue(hn.ProductionEdgeHostnameID),
					CnameTo:              types.StringValue(hn.ProductionCnameTo),
				}
				result[hn.CnameFrom] = h
				tflog.Debug(ctx, "added hostname from API to the state", h.toLog(), map[string]any{"cname_from": hn.CnameFrom})
			}
		}
	}

	return result
}

// findCurrentActivation returns the latest ACTIVE activation.
func findCurrentActivation(ctx context.Context, client papi.PAPI, data *HostnameBucketResourceModel) (papi.HostnameActivationListItem, error) {
	offset, limit := 0, 999
	for {
		activations, err := client.ListPropertyHostnameActivations(ctx, papi.ListPropertyHostnameActivationsRequest{
			PropertyID: data.PropertyID.ValueString(),
			ContractID: data.ContractID.ValueString(),
			GroupID:    data.GroupID.ValueString(),
			Offset:     offset,
			Limit:      limit,
		})
		if err != nil {
			return papi.HostnameActivationListItem{}, err
		}

		for _, act := range activations.HostnameActivations.Items {
			if act.Status == "ACTIVE" && string(act.Network) == data.Network.ValueString() {
				return act, nil
			}
		}

		offset += limit
		if offset >= activations.HostnameActivations.TotalItems {
			break
		}
	}

	return papi.HostnameActivationListItem{}, fmt.Errorf("there is no active hostname activation for given property")
}

func waitForHostnameBucketActivation(ctx context.Context, client papi.PAPI, data HostnameBucketResourceModel) error {
	timeout := time.Duration(data.TimeoutForActivation.ValueInt64()) * time.Minute
	// Overwrite the timeout value if the forceTimeoutDuration has been configured in unit tests.
	if forceTimeoutDuration != 0 {
		timeout = forceTimeoutDuration
	}
	deadline := time.Now().Add(timeout)

	for {
		activation, err := client.GetPropertyHostnameActivation(ctx, papi.GetPropertyHostnameActivationRequest{
			PropertyID:           data.PropertyID.ValueString(),
			HostnameActivationID: data.ActivationID.ValueString(),
			ContractID:           data.ContractID.ValueString(),
			GroupID:              data.GroupID.ValueString(),
		})
		if err != nil {
			return err
		}
		if activation.HostnameActivation.Status == "ACTIVE" {
			return nil
		}

		if time.Now().After(deadline) {
			cancelResp, err := client.CancelPropertyHostnameActivation(ctx, papi.CancelPropertyHostnameActivationRequest{
				PropertyID:           data.PropertyID.ValueString(),
				ContractID:           data.ContractID.ValueString(),
				GroupID:              data.GroupID.ValueString(),
				HostnameActivationID: data.ActivationID.ValueString(),
			})
			if err != nil && (errors.Is(err, papi.ErrActivationTooFar) || errors.Is(err, papi.ErrActivationAlreadyActive)) {
				deadline = time.Now().Add(timeout)
				continue
			} else if err != nil && (!errors.Is(err, papi.ErrActivationTooFar) || !errors.Is(err, papi.ErrActivationAlreadyActive)) {
				return err
			} else if err == nil {
				innerDeadline := time.Now().Add(timeout)
				cancelActivationID := cancelResp.HostnameActivation.HostnameActivationID
				for {
					cancelActivation, err := client.GetPropertyHostnameActivation(ctx, papi.GetPropertyHostnameActivationRequest{
						PropertyID:           data.PropertyID.ValueString(),
						HostnameActivationID: cancelActivationID,
						ContractID:           data.ContractID.ValueString(),
						GroupID:              data.GroupID.ValueString(),
					})
					if err != nil {
						return err
					}
					if cancelActivation.HostnameActivation.Status == "ABORTED" {
						return errCancelActivation
					}

					if time.Now().After(innerDeadline) {
						return fmt.Errorf("sent cancel request for the activation: %s, but reached the timeout for waiting until the change is active. Please remove local state and import the resource", cancelActivationID)
					}

					time.Sleep(getHostnameBucketActivationInterval)
				}
			}
		}

		time.Sleep(getHostnameBucketActivationInterval)
	}
}

// In delete, we need to wait until the activation is ACTIVE, in other way if we would cancel such type of activation,
// the state of the resource would enter weird scenario - API would preserve the hostnames,
// but the resource would not be deleted locally.
// Moreover, we would need to check and verify if the cancel request went through and was actually successful,
// so let's not overcomplicate already complicated system.
func waitForHostnameBucketDeletion(ctx context.Context, client papi.PAPI, data HostnameBucketResourceModel) error {
	for {
		activation, err := client.GetPropertyHostnameActivation(ctx, papi.GetPropertyHostnameActivationRequest{
			PropertyID:           data.PropertyID.ValueString(),
			HostnameActivationID: data.ActivationID.ValueString(),
			ContractID:           data.ContractID.ValueString(),
			GroupID:              data.GroupID.ValueString(),
		})
		if err != nil {
			return err
		}
		if activation.HostnameActivation.Status == "ACTIVE" {
			return nil
		}

		time.Sleep(getHostnameBucketActivationInterval)
	}
}

func formatWarning(numberOfPendingCerts int64) diag.WarningDiagnostic {
	var verb, noun string
	if numberOfPendingCerts == 1 {
		verb = "is"
		noun = "certificate"
	} else {
		verb = "are"
		noun = "certificates"
	}
	detail := "use 'akamai_property_hostnames' data source with the 'filter_pending_default_certs' set to true to get more details"
	summary := fmt.Sprintf("there %s %d default %s that %s not active", verb, numberOfPendingCerts, noun, verb)

	return diag.NewWarningDiagnostic(summary, detail)
}
