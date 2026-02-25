package property

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/domainownership"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &DomainOwnershipLateValidationResource{}
	_ resource.ResourceWithImportState = &DomainOwnershipLateValidationResource{}
	_ resource.ResourceWithConfigure   = &DomainOwnershipLateValidationResource{}
)

const defaultLateValidationPollTimeout = 30 * time.Minute

type (
	// DomainOwnershipLateValidationResource represents akamai_domainownership_late_validation resource.
	DomainOwnershipLateValidationResource struct {
		meta.Resource
		domainOwnershipLateValidationResourceConfig
	}

	domainOwnershipLateValidationResourceConfig struct {
		searchInterval time.Duration
	}

	domainOwnershipLateValidationResourceModel struct {
		PropertyID       types.String   `tfsdk:"property_id"`
		Version          types.Int64    `tfsdk:"version"`
		ContractID       types.String   `tfsdk:"contract_id"`
		GroupID          types.String   `tfsdk:"group_id"`
		ValidationMethod types.String   `tfsdk:"validation_method"`
		Timeouts         timeouts.Value `tfsdk:"timeouts"`
	}
)

func defaultDomainOwnershipLateValidationResourceConfig() domainOwnershipLateValidationResourceConfig {
	return domainOwnershipLateValidationResourceConfig{
		searchInterval: 30 * time.Second,
	}
}

// NewDomainOwnershipLateValidationResource returns new domain ownership late validation resource.
func NewDomainOwnershipLateValidationResource(config domainOwnershipLateValidationResourceConfig) func() resource.Resource {
	return func() resource.Resource {
		return &DomainOwnershipLateValidationResource{
			domainOwnershipLateValidationResourceConfig: config,
		}
	}
}

// Metadata implements resource.Resource.
func (d *DomainOwnershipLateValidationResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_property_domainownership_late_validation"
}

// Schema implements resource's Schema.
func (d *DomainOwnershipLateValidationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"property_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("prp_")),
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(prp_)?\d+$`),
						"must start with 'prp_' prefix followed by digits, or be digits only",
					),
				},
				Description: "Property ID of the Property which domains will be validated.",
			},
			"version": schema.Int64Attribute{
				Required:    true,
				Description: "Property version containing domains to be validated.",
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "Group ID of the Property.",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("grp_")),
				},
			},
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Contract ID of the Property.",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					modifiers.PreventStringUpdate(),
				},
			},
			"validation_method": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The method used to validate the domain. Possible values are: \n" +
					"* `DNS_CNAME` - For this method, Akamai generates a `cname_record` that you copy as the `target` to a `CNAME` record of your DNS configuration. The record's name needs to be in the `_acme-challenge.domain-name` format.\n" +
					"* `DNS_TXT` - For this method, Akamai generates a `txt_record` with a token `value` that you copy as the `target` to a `TXT` record of your DNS configuration. The record's name needs to be in the `_akamai-{host|wildcard|domain}-challenge.domainName` format based on the validation scope.\n" +
					"* `HTTP` - Applies only to domains with the `HOST` validation scope. For this method, you create the file containing a token and place it on your HTTP server in the location specified by the `validation_challenge.http_file.path` or use a redirect to the `validation_challenge.http_redirect.to` with the token.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(domainownership.ValidationMethodDNSCNAME),
						string(domainownership.ValidationMethodDNSTXT),
						string(domainownership.ValidationMethodHTTP)),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create:            true,
				CreateDescription: "Optional configurable domains validation timeout to be used on resource create. By default it's 30m.",
				Update:            true,
				UpdateDescription: "Optional configurable domains validation timeout to be used on resource update. By default it's 30m.",
			}),
		},
	}
}

// Create implements resource's Create method.
func (d *DomainOwnershipLateValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Domain Ownership Late Validation Resource")
	ctx = tflog.SetField(ctx, "method", "create")

	var plan domainOwnershipLateValidationResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, timeoutDiags := plan.Timeouts.Create(ctx, defaultLateValidationPollTimeout)
	if resp.Diagnostics.Append(timeoutDiags...); resp.Diagnostics.HasError() {
		return
	}

	diags := d.validateDomains(ctx, &plan, timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read implements resource's Read method.
func (d *DomainOwnershipLateValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Domain Ownership Late Validation Resource")
	ctx = tflog.SetField(ctx, "method", "read")

	var state domainOwnershipLateValidationResourceModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	papiClient := d.Client.GetPAPI()
	prop := papi.Property{
		ContractID: state.ContractID.ValueString(),
		PropertyID: state.PropertyID.ValueString(),
		GroupID:    state.GroupID.ValueString(),
	}
	version := int(state.Version.ValueInt64())

	hostnames, err := fetchPropertyVersionHostnames(ctx, papiClient, prop, version)
	if err != nil {
		if errors.Is(err, papi.ErrNotFound) {
			tflog.Warn(ctx, "Property or version not found, removing resource from state")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error on Read", fmt.Sprintf("Failed to fetch property hostnames: %s", err))
		return
	}

	var invalidHostnames []string
	for _, hostname := range hostnames {
		if hostname.DomainOwnershipVerification != nil && hostname.DomainOwnershipVerification.Status != "VALIDATED" {
			invalidHostnames = append(invalidHostnames, hostname.CnameFrom)
		}
	}

	if len(invalidHostnames) > 0 {
		tflog.Debug(ctx, "Hostnames are not validated, removing resource from state", map[string]any{
			"hostnames": invalidHostnames,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	tflog.Debug(ctx, "Resource is valid")
}

// Update implements resource's Update method.
func (d *DomainOwnershipLateValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Domain Ownership Late Validation Resource")
	ctx = tflog.SetField(ctx, "method", "update")

	var plan domainOwnershipLateValidationResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	timeout, timeoutDiags := plan.Timeouts.Update(ctx, defaultLateValidationPollTimeout)
	if resp.Diagnostics.Append(timeoutDiags...); resp.Diagnostics.HasError() {
		return
	}

	diags := d.validateDomains(ctx, &plan, timeout)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource's Delete method.
func (d *DomainOwnershipLateValidationResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Domain Ownership Late Validation Resource. This is a no-op.")
}

// ImportState implements resource's ImportState method.
func (d *DomainOwnershipLateValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Domain Ownership Late Validation Resource")
	ctx = tflog.SetField(ctx, "method", "import")

	id := req.ID
	tflog.Debug(ctx, fmt.Sprintf("importID: %s", id))

	parts := strings.Split(req.ID, ",")

	if len(parts) != 5 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected a comma-separated import ID with 5 parts: property_id,property_version,contract_id,group_id,validation_method. Got: %q", req.ID),
		)
		return
	}

	propertyID := parts[0]
	versionStr := parts[1]
	contractID := parts[2]
	groupID := parts[3]
	validationMethod := parts[4]

	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Property Version in Import ID",
			fmt.Sprintf("The version part of the import ID could not be parsed as an integer. Got: %q. Error: %s", versionStr, err),
		)
		return
	}

	if validationMethod != string(domainownership.ValidationMethodDNSCNAME) &&
		validationMethod != string(domainownership.ValidationMethodDNSTXT) &&
		validationMethod != string(domainownership.ValidationMethodHTTP) {
		resp.Diagnostics.AddError(
			"Invalid Validation Method in Import ID",
			fmt.Sprintf("The validation method part of the import ID is invalid. Expected one of: %q, %q, %q. Got: %q",
				domainownership.ValidationMethodDNSCNAME,
				domainownership.ValidationMethodDNSTXT,
				domainownership.ValidationMethodHTTP,
				validationMethod),
		)
		return
	}

	var data domainOwnershipLateValidationResourceModel
	data.PropertyID = types.StringValue(propertyID)
	data.Version = types.Int64Value(version)
	data.ContractID = types.StringValue(contractID)
	data.GroupID = types.StringValue(groupID)
	data.ValidationMethod = types.StringValue(validationMethod)

	data.Timeouts = timeouts.Value{
		Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"update": types.StringType,
		}),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// validateDomains contains the shared logic for Create and Update.
func (d *DomainOwnershipLateValidationResource) validateDomains(ctx context.Context, plan *domainOwnershipLateValidationResourceModel, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	papiClient := d.Client.GetPAPI()
	prop := papi.Property{
		ContractID: plan.ContractID.ValueString(),
		PropertyID: plan.PropertyID.ValueString(),
		GroupID:    plan.GroupID.ValueString(),
	}
	version := int(plan.Version.ValueInt64())

	hostnames, err := fetchPropertyVersionHostnames(ctx, papiClient, prop, version)
	if err != nil {
		diags.AddError("API Error", fmt.Sprintf("Failed to fetch property version hostnames: %s", err))
		return diags
	}

	var domainsToValidate []string
	for _, item := range hostnames {
		if item.DomainOwnershipVerification != nil && item.DomainOwnershipVerification.Status != "VALIDATED" {
			domainsToValidate = append(domainsToValidate, item.CnameFrom)
		}
	}

	if len(domainsToValidate) == 0 {
		tflog.Info(ctx, "All domains are already validated.")
		return diags
	}
	tflog.Info(ctx, "Domains identified for ownership validation", map[string]any{"domains": domainsToValidate})

	domClient := d.Client.GetDomainOwnership()
	domainsToValidateDOM := make([]domainownership.ValidateDomain, 0, len(domainsToValidate))
	for _, domain := range domainsToValidate {
		var scope domainownership.ValidationScope
		actualDomain := domain
		if strings.HasPrefix(domain, "*.") {
			scope = domainownership.ValidationScopeDomain
			actualDomain = strings.TrimPrefix(domain, "*.")
		} else {
			scope = domainownership.ValidationScopeHost
		}

		domainsToValidateDOM = append(domainsToValidateDOM, domainownership.ValidateDomain{DomainName: actualDomain, ValidationScope: scope, ValidationMethod: domainownership.ValidationMethod(plan.ValidationMethod.ValueString())})
	}
	validateDomainsResponse, err := domClient.ValidateDomains(ctx, domainownership.ValidateDomainsRequest{Domains: domainsToValidateDOM})
	if err != nil {
		diags.AddError("API Error", fmt.Sprintf("Failed to initiate domain ownership validation: %s", err))
		return diags
	}

	var domainsToPoll []string
	for _, hostname := range validateDomainsResponse.Domains {
		if hostname.DomainStatus != "VALIDATED" {
			domainsToPoll = append(domainsToPoll, hostname.DomainName)
		}
	}

	if len(domainsToPoll) == 0 {
		tflog.Info(ctx, "All domains were successfully validated.")
		return diags
	}

	tflog.Info(ctx, "Polling required for domains", map[string]any{"domains": domainsToPoll, "timeout": timeout})
	waitDiags := d.waitForDomainsValidation(ctx, papiClient, prop, version, timeout)
	diags.Append(waitDiags...)

	return diags
}

// waitForDomainsValidation polls for domain validation status until all domains are validated or a timeout occurs.
func (d *DomainOwnershipLateValidationResource) waitForDomainsValidation(ctx context.Context, client papi.PAPI, prop papi.Property, version int, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	searchPollIntervalTicker := time.NewTicker(d.searchInterval)
	defer searchPollIntervalTicker.Stop()

	var stillNotValidated []string

	for {
		select {
		case <-ctx.Done():
			diags.AddError("Timeout while waiting for domain validation",
				fmt.Sprintf("Please make sure that challenges for the domains: %s are correctly set up and try again.", strings.Join(stillNotValidated, ", ")))
			return diags
		case <-searchPollIntervalTicker.C:
			tflog.Debug(ctx, "Polling for domain validation status...")

			hostnamesResp, err := fetchPropertyVersionHostnames(ctx, client, prop, version)
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Polling check failed, will retry on next tick: %s", err))
				continue
			}

			var pendingDomains []string
			for _, item := range hostnamesResp {
				if item.DomainOwnershipVerification != nil && item.DomainOwnershipVerification.Status != "VALIDATED" {
					pendingDomains = append(pendingDomains, item.CnameFrom)
				}
			}

			stillNotValidated = pendingDomains

			if len(stillNotValidated) == 0 {
				tflog.Info(ctx, "All polled domains have been successfully validated.")
				return diags
			}

			tflog.Debug(ctx, "Still waiting for domains to be validated", map[string]any{"remaining_domains": stillNotValidated})
		}
	}
}
