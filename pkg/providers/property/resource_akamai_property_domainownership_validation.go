package property

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &DomainOwnershipValidationResource{}
	_ resource.ResourceWithImportState = &DomainOwnershipValidationResource{}
	_ resource.ResourceWithConfigure   = &DomainOwnershipValidationResource{}

	searchInterval     = 30 * time.Second
	defaultPollTimeout = 30 * time.Minute
)

type (
	// DomainOwnershipValidationResource represents akamai_domainownership_validation resource.
	DomainOwnershipValidationResource struct {
		meta meta.Meta
	}

	domainOwnershipValidationResourceModel struct {
		Domains  types.Set      `tfsdk:"domains"`
		Timeouts timeouts.Value `tfsdk:"timeouts"`
	}

	domainModel struct {
		DomainName       types.String `tfsdk:"domain_name"`
		ValidationScope  types.String `tfsdk:"validation_scope"`
		ValidationMethod types.String `tfsdk:"validation_method"`
	}

	domainKey struct {
		domainName      string
		validationScope string
	}

	domainDetails struct {
		validationMethod *string
		validationStatus string
		validationLevel  string
	}
)

// NewDomainOwnershipValidationResource returns new domain ownership validation resource.
func NewDomainOwnershipValidationResource() resource.Resource {
	return &DomainOwnershipValidationResource{}
}

// Metadata implements resource.Resource.
func (d *DomainOwnershipValidationResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_property_domainownership_validation"
}

// Configure implements resource.ResourceWithConfigure.
func (d *DomainOwnershipValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.meta = meta.Must(req.ProviderData)
}

// Schema implements resource's Schema.
func (d *DomainOwnershipValidationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domains": domainsSchema(),
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

func domainsType() types.ObjectType {
	return domainsSchema().NestedObject.Type().(types.ObjectType)
}

func domainsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Required:    true,
		Description: "List of domains to be validated.",
		Validators: []validator.Set{
			setvalidator.SizeBetween(1, 1000),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"domain_name": schema.StringAttribute{
					Required:    true,
					Description: "Name of the domain.",
				},
				"validation_scope": schema.StringAttribute{
					Required:    true,
					Description: "Scope of the domain validation, either 'HOST', 'WILDCARD', or 'DOMAIN'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							string(domainownership.ValidationScopeHost),
							string(domainownership.ValidationScopeWildcard),
							string(domainownership.ValidationScopeDomain)),
					},
				},
				"validation_method": schema.StringAttribute{
					Optional:    true,
					Description: "If it is not provided, the default validation method will be used.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							string(domainownership.ValidationMethodDNSCNAME),
							string(domainownership.ValidationMethodDNSTXT),
							string(domainownership.ValidationMethodHTTP)),
					},
				},
			},
		},
	}
}

// Create implements resource's Create method.
func (d *DomainOwnershipValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Domain Ownership Validation Resource")
	ctx = tflog.SetField(ctx, "method", "create")

	var plan domainOwnershipValidationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planDomains []domainModel
	if diags := plan.Domains.ElementsAs(ctx, &planDomains, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Debug(ctx, "planned domains", map[string]any{
		"domains": planDomains,
	})

	client := DomainOwnershipClient(d.meta)
	apiDomains, err := fetchDomainsFromAPI(ctx, client, planDomains)
	if err != nil {
		resp.Diagnostics.AddError("Error Searching Domains", err.Error())
		return
	}

	planDomainsMap := domainsToMap(planDomains)
	apiDomainsMap := apiDomainsToMap(apiDomains)

	validationHandler := newValidationHandler(ctx).
		setPlanDomains(planDomainsMap).
		setAPIDomains(apiDomainsMap)

	validationHandler, err = validationHandler.calculateDomainsToValidate()
	if err != nil {
		resp.Diagnostics.AddError("Error Validating Domains", err.Error())
		return
	}

	requests := validationHandler.buildValidateRequests()
	domainsToPoll, err := validateDomains(ctx, client, requests)
	if err != nil {
		resp.Diagnostics.AddError("Error Validating Domains", err.Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, defaultPollTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if len(domainsToPoll) != 0 {
		tflog.Debug(ctx, "need to wait for domains to become validated", map[string]any{
			"domains": domainsToPoll,
			"timeout": createTimeout,
		})
		resp.Diagnostics.Append(waitForDomains(ctx, client, domainsToPoll, createTimeout)...)
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddWarning("Partial success of create",
				"Some domains scheduled for validation may not have been validated. "+
					"Rerun 'terraform apply' to retry validating the remaining domains.")
			return
		}
	}

	var state domainOwnershipValidationResourceModel
	state.Domains = plan.Domains
	state.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read implements resource's Read method.
func (d *DomainOwnershipValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Domain Ownership Validation Resource")
	ctx = tflog.SetField(ctx, "method", "read")

	var state domainOwnershipValidationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateDomains []domainModel
	if diags := state.Domains.ElementsAs(ctx, &stateDomains, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Debug(ctx, "state domains", map[string]any{
		"domains": stateDomains,
	})

	stateDomainsMap := domainsToMap(stateDomains)
	tflog.Debug(ctx, "mapped state domains", map[string]any{
		"state_domains_map": stateDomainsMap,
	})

	client := DomainOwnershipClient(d.meta)
	apiDomains, err := fetchDomainsFromAPI(ctx, client, stateDomains)
	if err != nil {
		resp.Diagnostics.AddError("Error Searching Domains", err.Error())
		return
	}

	apiDomainsMap := apiDomainsToMap(apiDomains)
	tflog.Debug(ctx, "mapped API domains", map[string]any{
		"api_domains_map": apiDomainsMap,
	})

	for stateDomain := range stateDomainsMap {
		if apiDomain, foundInAPI := apiDomainsMap[stateDomain]; !foundInAPI {
			tflog.Debug(ctx, "domain not found in API, removing from state", map[string]any{
				"domain_name":      stateDomain.domainName,
				"validation_scope": stateDomain.validationScope,
			})
			delete(stateDomainsMap, stateDomain)
		} else {
			if apiDomain.validationStatus != "VALIDATED" {
				tflog.Debug(ctx, fmt.Sprintf("domain with '%s' status in API, removing from state",
					apiDomain.validationStatus), map[string]any{
					"domain_name":      stateDomain.domainName,
					"validation_scope": stateDomain.validationScope,
				})
				delete(stateDomainsMap, stateDomain)
			}
		}
	}

	if len(stateDomainsMap) == 0 {
		tflog.Debug(ctx, "no domains left in state after read, removing resource from state")
		resp.State.RemoveResource(ctx)
		return
	}

	var refreshedStateDomains []domainModel
	for domain, domainDetails := range stateDomainsMap {
		refreshedStateDomains = append(refreshedStateDomains, domainModel{
			DomainName:       types.StringValue(domain.domainName),
			ValidationScope:  types.StringValue(domain.validationScope),
			ValidationMethod: types.StringPointerValue(domainDetails.validationMethod),
		})
	}
	tflog.Debug(ctx, "refreshed state domains", map[string]any{
		"refreshed_state_domains": refreshedStateDomains,
	})

	domainsSet, diags := types.SetValueFrom(ctx, domainsType(), refreshedStateDomains)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.Domains = domainsSet

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update implements resource's Update method.
func (d *DomainOwnershipValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Domain Ownership Validation Resource")
	ctx = tflog.SetField(ctx, "method", "update")

	var plan domainOwnershipValidationResourceModel
	var state domainOwnershipValidationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planDomains []domainModel
	if diags := plan.Domains.ElementsAs(ctx, &planDomains, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var stateDomains []domainModel
	if diags := state.Domains.ElementsAs(ctx, &stateDomains, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Debug(ctx, "read state and plan domains", map[string]any{
		"state_domains": stateDomains,
		"plan_domains":  planDomains,
	})

	planDomainsMap := domainsToMap(planDomains)
	stateDomainsMap := domainsToMap(stateDomains)

	// Read domains from state and plan, to determine which domains
	// need to be fetched from API to check their statuses.
	domainKeys := make(map[domainKey]struct{})
	for k := range stateDomainsMap {
		domainKeys[k] = struct{}{}
	}
	for k := range planDomainsMap {
		domainKeys[k] = struct{}{}
	}
	var domainsToSearch []domainModel
	for k := range domainKeys {
		domainsToSearch = append(domainsToSearch, domainModel{
			DomainName:      types.StringValue(k.domainName),
			ValidationScope: types.StringValue(k.validationScope),
		})
	}

	client := DomainOwnershipClient(d.meta)
	apiDomains, err := fetchDomainsFromAPI(ctx, client, domainsToSearch)
	if err != nil {
		resp.Diagnostics.AddError("Error Searching Domains", err.Error())
		return
	}
	apiDomainsMap := apiDomainsToMap(apiDomains)

	validationHandler := newValidationHandler(ctx).
		setPlanDomains(planDomainsMap).
		setAPIDomains(apiDomainsMap).
		setStateDomains(stateDomainsMap)

	invalidateRequest := validationHandler.
		calculateDomainsToInvalidate().
		buildInvalidateRequest()
	if invalidateRequest != nil {
		invalidateResponse, err := client.InvalidateDomains(ctx, *invalidateRequest)
		if err != nil {
			resp.Diagnostics.AddError("Error Invalidating Domains", err.Error())
			return
		}
		tflog.Debug(ctx, "domains invalidated", map[string]any{
			"response": invalidateResponse,
		})
	}

	validationHandler, err = validationHandler.calculateDomainsToValidate()
	if err != nil {
		resp.Diagnostics.AddError("Error Validating Domains", err.Error())
		return
	}

	validateRequests := validationHandler.buildValidateRequests()
	domainsToPoll, err := validateDomains(ctx, client, validateRequests)
	if err != nil {
		resp.Diagnostics.AddError("Error Validating Domains", err.Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, defaultPollTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if len(domainsToPoll) != 0 {
		tflog.Debug(ctx, "need to wait for domains to become validated", map[string]any{
			"domains": domainsToPoll,
			"timeout": updateTimeout,
		})
		resp.Diagnostics.Append(waitForDomains(ctx, client, domainsToPoll, updateTimeout)...)
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddWarning("Partial success of update",
				"Domains scheduled for invalidation have been successfully invalidated while "+
					"some domains scheduled for validation may not have been validated")
			return
		}
	}

	state.Domains = plan.Domains
	state.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Delete implements resource's Delete method.
func (d *DomainOwnershipValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Domain Ownership Validation Resource")
	ctx = tflog.SetField(ctx, "method", "delete")

	var state domainOwnershipValidationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateDomains []domainModel
	if diags := state.Domains.ElementsAs(ctx, &stateDomains, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Debug(ctx, "state domains", map[string]any{
		"domains": stateDomains,
	})

	client := DomainOwnershipClient(d.meta)
	apiDomains, err := fetchDomainsFromAPI(ctx, client, stateDomains)
	if err != nil {
		resp.Diagnostics.AddError("Error Searching Domains", err.Error())
		return
	}

	apiDomainsMap := apiDomainsToMap(apiDomains)
	stateDomainsMap := domainsToMap(stateDomains)

	validationHandler := newValidationHandler(ctx).
		setAPIDomains(apiDomainsMap).
		setStateDomains(stateDomainsMap).
		setPlanDomains(make(map[domainKey]domainDetails))

	invalidateRequest := validationHandler.
		calculateDomainsToInvalidate().
		buildInvalidateRequest()
	if invalidateRequest != nil {
		invalidateResponse, err := client.InvalidateDomains(ctx, *invalidateRequest)
		if err != nil {
			resp.Diagnostics.AddError("Error Invalidating Domains", err.Error())
			return
		}
		tflog.Debug(ctx, "domains invalidated", map[string]any{
			"response": invalidateResponse,
		})
	}
}

// ImportState implements resource's ImportState method.
func (d *DomainOwnershipValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Domain Ownership Validation Resource")
	ctx = tflog.SetField(ctx, "method", "import")

	id := req.ID
	tflog.Debug(ctx, fmt.Sprintf("importID: %s", id))

	client := DomainOwnershipClient(d.meta)
	domains, diags := parseDomains(ctx, client, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var domainModels []domainModel
	for _, domain := range domains {
		domainModels = append(domainModels, domainModel{
			DomainName:      types.StringValue(domain.DomainName),
			ValidationScope: types.StringValue(string(domain.ValidationScope)),
		})
	}
	domainsModel, diags := types.SetValueFrom(ctx, domainsType(), domainModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data domainOwnershipValidationResourceModel
	data.Domains = domainsModel
	data.Timeouts = timeouts.Value{
		Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"update": types.StringType,
		}),
	}
	if resp.Diagnostics.Append(resp.State.Set(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
}

func parseDomains(ctx context.Context, client domainownership.DomainOwnership, importID string) ([]domainownership.Domain, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	domainsToImport, err := importIDToDomains(importID)
	if err != nil {
		diags.AddError("Error parsing import ID", err.Error())
		return nil, diags
	}
	tflog.Debug(ctx, "calculated domains to import", map[string]any{
		"domains": domainsToImport,
	})

	finalDomains, err := verifyDomainsToImport(ctx, client, domainsToImport)
	if err != nil {
		diags.AddError("Error verifying domains", err.Error())
		return nil, diags
	}

	return finalDomains, diags
}

func importIDToDomains(importID string) (map[domainKey]domainDetails, error) {
	domainsToImport := make(map[domainKey]domainDetails)
	splittedImportID := strings.Split(importID, ",")
	for _, domain := range splittedImportID {
		domainParts := strings.Split(domain, ":")
		switch len(domainParts) {
		case 1:
			if _, ok := domainsToImport[domainKey{domainName: domainParts[0], validationScope: ""}]; ok {
				return nil, fmt.Errorf("domain '%s' was already provided in the importID. Please remove duplicate domain entries", domainParts[0])
			}
			domainsToImport[domainKey{
				domainName:      domainParts[0],
				validationScope: "",
			}] = domainDetails{}
		case 2:
			if _, ok := domainsToImport[domainKey{domainName: domainParts[0], validationScope: domainParts[1]}]; ok {
				return nil, fmt.Errorf("domain '%s' with validation scope '%s' was already provided in the importID. Please remove duplicate domain entries", domainParts[0], domainParts[1])
			}
			if domainParts[1] == "HOST" || domainParts[1] == "WILDCARD" || domainParts[1] == "DOMAIN" {
				domainsToImport[domainKey{
					domainName:      domainParts[0],
					validationScope: domainParts[1],
				}] = domainDetails{}
			} else {
				return nil, fmt.Errorf("invalid validation scope '%s' for domain '%s'. Expected one of 'HOST', 'WILDCARD', 'DOMAIN'", domainParts[1], domainParts[0])
			}
		default:
			return nil, fmt.Errorf("invalid import ID format. Expected format is a list of: domain[:scope], got %v", domain)
		}
	}

	return domainsToImport, nil
}

func verifyDomainsToImport(ctx context.Context, client domainownership.DomainOwnership, domains map[domainKey]domainDetails) ([]domainownership.Domain, error) {
	var searchDomainsBody []domainownership.Domain
	for d := range domains {
		if d.validationScope != "" {
			searchDomainsBody = append(searchDomainsBody, newDomain(d.domainName, d.validationScope))
		} else {
			searchDomainsBody = append(searchDomainsBody, newDomain(d.domainName, "HOST"))
			searchDomainsBody = append(searchDomainsBody, newDomain(d.domainName, "DOMAIN"))
			searchDomainsBody = append(searchDomainsBody, newDomain(d.domainName, "WILDCARD"))
		}
	}
	sortDomains(searchDomainsBody)
	tflog.Debug(ctx, "built search domains request body", map[string]any{
		"search_domains_body": searchDomainsBody,
	})

	apiDomains, err := client.SearchDomains(ctx, domainownership.SearchDomainsRequest{
		IncludeAll: true,
		Body: domainownership.SearchDomainsBody{
			Domains: searchDomainsBody,
		},
	})
	if err != nil {
		return nil, err
	}

	apiDomainsMap := apiDomainsToMap(apiDomains.Domains)
	tflog.Debug(ctx, "mapped API domains", map[string]any{
		"api_domains_map": apiDomainsMap,
	})

	finalDomains := make([]domainownership.Domain, 0, len(domains))
	for d := range domains {
		if d.validationScope == "" {
			// If the validation scope is not provided, we need to check
			// if the domain exists with multiple validation scopes.
			// If it does, we cannot determine which one to import and return an error.
			// If it exists with only one scope, we can import that one.
			// If it doesn't exist, we return an error.
			var counter int
			var candidateDomain domainownership.Domain
			var candidateDomainStatus string
			for _, scope := range []string{"HOST", "WILDCARD", "DOMAIN"} {
				if apiDomainDetails, ok := apiDomainsMap[domainKey{domainName: d.domainName, validationScope: scope}]; ok {
					if apiDomainDetails.validationLevel == "FQDN" {
						counter++
						candidateDomain = newDomain(d.domainName, scope)
						candidateDomainStatus = apiDomainDetails.validationStatus
					}
				}
			}
			if counter == 0 {
				return nil, fmt.Errorf("the domain '%s' was not found", d.domainName)
			} else if counter > 1 {
				return nil, fmt.Errorf("the domain '%s' exists with multiple validation scopes. Please re-import specifying the validation scope for the domain", d.domainName)
			}
			if candidateDomainStatus != "VALIDATED" {
				return nil, fmt.Errorf("the domain '%s' is in '%s' status and cannot be imported", d.domainName, candidateDomainStatus)
			}
			finalDomains = append(finalDomains, candidateDomain)
		} else {
			// If the validation scope is provided, we need to check
			// if the domain exists with that validation scope in the API,
			// and if its validation level is FQDN. FQDN is the validation level for the domain,
			// which was explicitly submitted for validation. Other validation levels is
			// ROOT/WILDCARD, which can be found without prior explicit submitting,
			// so we don't import those as they are part of the WILDCARD or DOMAIN domain.
			if apiDomainDetails, ok := apiDomainsMap[d]; ok {
				if apiDomainDetails.validationStatus != "VALIDATED" {
					return nil, fmt.Errorf("the domain '%s' with validation scope '%s' is in '%s' status and cannot be imported", d.domainName, d.validationScope, apiDomainDetails.validationStatus)
				}
				if apiDomainDetails.validationLevel == "FQDN" {
					finalDomains = append(finalDomains, newDomain(d.domainName, d.validationScope))
				} else {
					return nil, fmt.Errorf("only domains with validation level FQDN can be imported, the requested domain '%s' with validation scope '%s' has validation level '%s'",
						d.domainName, d.validationScope, apiDomainDetails.validationLevel)
				}
			} else {
				return nil, fmt.Errorf("the domain '%s' with validation scope '%s' was not found", d.domainName, d.validationScope)
			}
		}
	}

	return finalDomains, nil
}

func domainsToMap(domains []domainModel) map[domainKey]domainDetails {
	domainMap := make(map[domainKey]domainDetails)
	for _, d := range domains {
		domainMap[domainKey{
			domainName:      d.DomainName.ValueString(),
			validationScope: d.ValidationScope.ValueString(),
		}] = domainDetails{
			validationMethod: d.ValidationMethod.ValueStringPointer(),
		}
	}
	return domainMap
}

func apiDomainsToMap(domains []domainownership.SearchDomainItem) map[domainKey]domainDetails {
	domainMap := make(map[domainKey]domainDetails)
	for _, d := range domains {
		domainMap[domainKey{
			domainName:      d.DomainName,
			validationScope: d.ValidationScope,
		}] = domainDetails{
			validationMethod: d.ValidationMethod,
			validationStatus: d.DomainStatus,
			validationLevel:  d.ValidationLevel,
		}
	}
	return domainMap
}

func fetchDomainsFromAPI(ctx context.Context, client domainownership.DomainOwnership, domains []domainModel) ([]domainownership.SearchDomainItem, error) {
	var searchDomainsBody []domainownership.Domain
	for _, d := range domains {
		searchDomainsBody = append(searchDomainsBody,
			newDomain(d.DomainName.ValueString(), d.ValidationScope.ValueString()))
	}

	// Sort the slice to have consistent order.
	sortDomains(searchDomainsBody)

	if len(searchDomainsBody) > 1000 {
		var allDomains []domainownership.SearchDomainItem

		for batch := range slices.Chunk(searchDomainsBody, 1000) {
			tflog.Debug(ctx, "built search domains request batch", map[string]any{
				"batch": batch,
			})

			batchDomains, err := client.SearchDomains(ctx, domainownership.SearchDomainsRequest{
				IncludeAll: true,
				Body: domainownership.SearchDomainsBody{
					Domains: batch,
				},
			})
			if err != nil {
				return nil, err
			}
			allDomains = append(allDomains, batchDomains.Domains...)
		}

		tflog.Debug(ctx, "domains fetched from API", map[string]any{
			"domains": allDomains,
		})
		return allDomains, nil
	}

	tflog.Debug(ctx, "built search domains request body", map[string]any{
		"body": searchDomainsBody,
	})

	apiDomains, err := client.SearchDomains(ctx, domainownership.SearchDomainsRequest{
		IncludeAll: true,
		Body: domainownership.SearchDomainsBody{
			Domains: searchDomainsBody,
		},
	})
	if err != nil {
		return nil, err
	}
	tflog.Debug(ctx, "domains fetched from API", map[string]any{
		"domains": apiDomains.Domains,
	})

	return apiDomains.Domains, nil
}

func validateDomains(ctx context.Context, client domainownership.DomainOwnership, requests []domainownership.ValidateDomainsRequest) (map[domainKey]domainDetails, error) {
	domainsToPoll := make(map[domainKey]domainDetails)
	for _, request := range requests {
		validatedDomains, err := client.ValidateDomains(ctx, request)
		if err != nil {
			return nil, err
		}
		tflog.Debug(ctx, "successfully validated domains", map[string]any{
			"validated_domains": validatedDomains,
		})

		for _, d := range validatedDomains.Domains {
			if d.DomainStatus != "VALIDATED" {
				domainsToPoll[domainKey{
					domainName:      d.DomainName,
					validationScope: d.ValidationScope,
				}] = domainDetails{
					validationStatus: d.DomainStatus,
				}
			}
		}
	}
	return domainsToPoll, nil
}

func formatDomainsForDiag(domains map[domainKey]domainDetails) string {
	var domainInfos []string
	for k := range domains {
		domainInfos = append(domainInfos, fmt.Sprintf("%s:%s", k.domainName, k.validationScope))
	}
	return strings.Join(domainInfos, ", ")
}

func waitForDomains(ctx context.Context, client domainownership.DomainOwnership, domainsToPoll map[domainKey]domainDetails, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	searchPollIntervalTicker := time.NewTicker(searchInterval)
	defer searchPollIntervalTicker.Stop()

	for {
		if len(domainsToPoll) == 0 {
			tflog.Debug(ctx, "all domains validated")
			break
		}

		select {
		case <-ctx.Done():
			diags.AddError("Timeout while waiting for domain validation",
				fmt.Sprintf("Please make sure that challenges for the domains: %s are correctly set up and try again.", formatDomainsForDiag(domainsToPoll)))
			return diags
		case <-searchPollIntervalTicker.C:
			var domainsToSearch []domainownership.Domain
			for domain := range domainsToPoll {
				domainsToSearch = append(domainsToSearch, newDomain(domain.domainName, domain.validationScope))
			}

			sortDomains(domainsToSearch)

			tflog.Debug(ctx, "polling domains", map[string]any{
				"domains": domainsToSearch,
			})

			apiDomains, err := client.SearchDomains(ctx, domainownership.SearchDomainsRequest{
				IncludeAll: true,
				Body: domainownership.SearchDomainsBody{
					Domains: domainsToSearch,
				},
			})
			if err != nil {
				if ctx.Err() != nil {
					diags.AddError(
						"Context terminated while waiting for domain validation",
						fmt.Sprintf("Please make sure that challenges for the domains: %s are correctly set up and try again.", formatDomainsForDiag(domainsToPoll)))
					return diags
				}
				diags.AddError("Error Searching Domains", err.Error())
				return diags
			}

			for _, d := range apiDomains.Domains {
				if d.DomainStatus == "VALIDATED" {
					delete(domainsToPoll, domainKey{
						domainName:      d.DomainName,
						validationScope: d.ValidationScope,
					})
				} else {
					domainsToPoll[domainKey{
						domainName:      d.DomainName,
						validationScope: d.ValidationScope,
					}] = domainDetails{
						validationStatus: d.DomainStatus,
					}
				}
			}
		}
	}
	return diags
}

func newDomain(name, validationScope string) domainownership.Domain {
	return domainownership.Domain{
		DomainName:      name,
		ValidationScope: domainownership.ValidationScope(validationScope),
	}
}

func sortDomains(domains []domainownership.Domain) {
	sort.Slice(domains, func(i, j int) bool {
		if domains[i].DomainName == domains[j].DomainName {
			return domains[i].ValidationScope < domains[j].ValidationScope
		}
		return domains[i].DomainName < domains[j].DomainName
	})
}
