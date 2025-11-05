package property

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &DomainsResource{}
	_ resource.ResourceWithImportState = &DomainsResource{}
	_ resource.ResourceWithConfigure   = &DomainsResource{}
	_ resource.ResourceWithModifyPlan  = &DomainsResource{}
)

// DomainsResource represents akamai_property_domainownership_domains resource.
type DomainsResource struct {
	meta meta.Meta
}

type (
	domainsResourceModel struct {
		Domains types.Set `tfsdk:"domains"`
	}

	domainResourceModel struct {
		DomainName              types.String `tfsdk:"domain_name"`
		ValidationScope         types.String `tfsdk:"validation_scope"`
		AccountID               types.String `tfsdk:"account_id"`
		DomainStatus            types.String `tfsdk:"domain_status"`
		ValidationMethod        types.String `tfsdk:"validation_method"`
		ValidationRequestedBy   types.String `tfsdk:"validation_requested_by"`
		ValidationRequestedDate types.String `tfsdk:"validation_requested_date"`
		ValidationCompletedDate types.String `tfsdk:"validation_completed_date"`
		ValidationChallenge     types.Object `tfsdk:"validation_challenge"`
	}
)

// NewDomainsResource returns new domains resource.
func NewDomainsResource() resource.Resource {
	return &DomainsResource{}
}

// Metadata implements resource.Resource.
func (r *DomainsResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_property_domainownership_domains"
}

// Schema implements resource's Schema.
func (r *DomainsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domains": domainsSchema(),
		},
	}
}

func domainsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Required:    true,
		Description: "List of domains.",
		Validators: []validator.Set{
			setvalidator.SizeBetween(1, 1000),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"domain_name": schema.StringAttribute{
					Required:    true,
					Description: "Your domain's name.",
				},
				"validation_scope": schema.StringAttribute{
					Required: true,
					MarkdownDescription: "Your domain's validation scope. Possible values are: \n" +
						"* `HOST` - The scope is only the exactly specified domain.\n" +
						"* `WILDCARD` - The scope covers any hostname within one subdomain level.\n" +
						"* `DOMAIN` - The scope covers any hostnames under the domain, regardless of the level of subdomains.",
					Validators: []validator.String{
						stringvalidator.OneOf("HOST", "WILDCARD", "DOMAIN"),
					},
				},
				"account_id": schema.StringAttribute{
					Computed:    true,
					Description: "Your account's ID.",
				},
				"domain_status": schema.StringAttribute{
					Computed: true,
					MarkdownDescription: "The domain's validation status. Possible values are: \n" +
						"* `REQUEST_ACCEPTED` - When you successfully submit the domain for validation.\n" +
						"* `VALIDATION_IN_PROGRESS` - When the DOM background jobs are trying to validate the domain.\n" +
						"* `VALIDATED` - When the validation is completed successfully. Akamai recognizes you as the domain owner.\n" +
						"* `TOKEN_EXPIRED` - When you haven't completed the validation in the requested time frame and the challenge token is not valid anymore. You need to generate new validation challenges for the domain.\n" +
						"* `INVALIDATED` - When the domain was invalidated and Akamai doesn't recognize you as its owner.",
				},
				"validation_method": schema.StringAttribute{
					Computed: true,
					MarkdownDescription: "The method used to validate the domain. Possible values are: \n" +
						"* `DNS_CNAME` - For this method, Akamai generates a `cname_record` that you copy as the `target` to a `CNAME` record of your DNS configuration. The record's name needs to be in the `_acme-challenge.domain-name` format.\n" +
						"* `DNS_TXT` - For this method, Akamai generates a `txt_record` with a token `value` that you copy as the `target` to a `TXT` record of your DNS configuration. The record's name needs to be in the `_akamai-{host|wildcard|domain}-challenge.domainName` format based on the validation scope.\n" +
						"* `HTTP` - Applies only to domains with the `HOST` validation scope. For this method, you create the file containing a token and place it on your HTTP server in the location specified by the `validation_challenge.http_file.path` or use a redirect to the `validation_challenge.http_redirect.to` with the token.\n" +
						"* `SYSTEM` - This method refers to domains that were automatically validated before Domain Validation Manager (DOM) was introduced.\n" +
						"* `MANUAL` - For this method, the DOM team manually performed the validation.",
				},
				"validation_requested_by": schema.StringAttribute{
					Computed:    true,
					Description: "The name of the user who requested the domain validation.",
				},
				"validation_requested_date": schema.StringAttribute{
					Computed:    true,
					Description: "The timestamp indicating when the domain validation was requested.",
				},
				"validation_completed_date": schema.StringAttribute{
					Computed:    true,
					Description: "The timestamp indicating when the domain validation was completed.",
				},
				"validation_challenge": validationChallengeSchema(),
			},
		},
	}
}

func validationChallengeType() map[string]attr.Type {
	return validationChallengeSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func domainsType() types.ObjectType {
	return domainsSchema().NestedObject.Type().(types.ObjectType)
}

// Configure implements resource.ResourceWithConfigure.
func (r *DomainsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.meta = meta.Must(req.ProviderData)
}

// ModifyPlan is a plan modifier for domainownership_domains resource.
func (r *DomainsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if modifiers.IsUpdate(req) {
		var plan, state domainsResourceModel

		if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
			return
		}
		if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
			return
		}

		var domainsFromPlan, domainsFromState []domainResourceModel
		if resp.Diagnostics.Append(plan.Domains.ElementsAs(ctx, &domainsFromPlan, false)...); resp.Diagnostics.HasError() {
			return
		}
		if resp.Diagnostics.Append(state.Domains.ElementsAs(ctx, &domainsFromState, false)...); resp.Diagnostics.HasError() {
			return
		}

		domains, diags := useStateForUnknownForDomains(ctx, domainsFromPlan, domainsFromState)
		if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
			return
		}
		plan.Domains = domains
		if resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...); resp.Diagnostics.HasError() {
			return
		}

		planMap := map[domainKey]domainResourceModel{}
		for _, domain := range domainsFromPlan {
			planMap[domainKey{domainName: domain.DomainName.ValueString(), validationScope: domain.ValidationScope.ValueString()}] = domain
		}
		domainsToRemove := buildRemoveDomainsRequest(domainsFromState, planMap)
		resp.Diagnostics.Append(warnAboutDroppingValidatedDomains(domainsToRemove, domainsFromState)...)
	}

	if modifiers.IsDelete(req) {
		var state domainsResourceModel

		if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
			return
		}

		var domainsFromState []domainResourceModel
		if resp.Diagnostics.Append(state.Domains.ElementsAs(ctx, &domainsFromState, false)...); resp.Diagnostics.HasError() {
			return
		}

		domainsToRemove := make([]domainownership.Domain, 0, len(domainsFromState))
		for _, domain := range domainsFromState {
			domainsToRemove = append(domainsToRemove, domainownership.Domain{
				DomainName:      domain.DomainName.ValueString(),
				ValidationScope: domainownership.ValidationScope(domain.ValidationScope.ValueString()),
			})
		}
		resp.Diagnostics.Append(warnAboutDroppingValidatedDomains(domainsToRemove, domainsFromState)...)
	}
}

func warnAboutDroppingValidatedDomains(domainsToRemove []domainownership.Domain, domainsFromState []domainResourceModel) diag.Diagnostics {
	if len(domainsToRemove) > 0 {
		validatedDomains := calculateValidatedDomains(domainsToRemove, domainsFromState)
		if len(validatedDomains) > 0 {
			var listOfDomains []string
			for _, domain := range validatedDomains {
				listOfDomains = append(listOfDomains, fmt.Sprintf("%s:%s", domain.DomainName, domain.ValidationScope))
			}
			diags := diag.Diagnostics{}
			diags.AddWarning("VALIDATED domains are planned for removal", fmt.Sprintf("The following VALIDATED domains are planned for removal. They will be invalidated during the apply phase before removal: [%s]", strings.Join(listOfDomains, ",")))
			return diags
		}
	}
	return nil
}

// useStateForUnknownForDomains handles setting the computed attributes. We cannot use library one because of there are nullable fields which cause problems there.
func useStateForUnknownForDomains(ctx context.Context, domainsFromPlan, domainsFromState []domainResourceModel) (basetypes.SetValue, diag.Diagnostics) {
	domainsMapFromState := make(map[domainKey]domainResourceModel, len(domainsFromState))
	for _, domain := range domainsFromState {
		domainsMapFromState[domainKey{domainName: domain.DomainName.ValueString(), validationScope: domain.ValidationScope.ValueString()}] = domain
	}
	var domainsWithComputedFields []domainResourceModel
	for _, domain := range domainsFromPlan {
		if domainFromState, ok := domainsMapFromState[domainKey{domainName: domain.DomainName.ValueString(), validationScope: domain.ValidationScope.ValueString()}]; ok {
			domainsWithComputedFields = append(domainsWithComputedFields, domainFromState)
		} else {
			// this is new domains so we don't have computed fields yet
			domainsWithComputedFields = append(domainsWithComputedFields, domain)
		}
	}

	domainsModel, diags := types.SetValueFrom(ctx, domainsType(), domainsWithComputedFields)
	if diags.Append(diags...); diags.HasError() {
		return basetypes.SetValue{}, diags
	}
	return domainsModel, diags
}

// Create implements resource's Create method.
func (r *DomainsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Domain Ownership Domains Resource")
	var data domainsResourceModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var domainList []domainResourceModel
	if resp.Diagnostics.Append(data.Domains.ElementsAs(ctx, &domainList, false)...); resp.Diagnostics.HasError() {
		return
	}

	sortDomainList(domainList)

	client := DomainOwnershipClient(r.meta)
	foundDomains, err := client.SearchDomains(ctx, buildSearchDomainsRequest(domainList))
	if err != nil {
		resp.Diagnostics.AddError("Error checking status of domains in configuration", err.Error())
		return
	}

	domainsToAdd, _, diags := buildAddAndRemoveDomainsRequest(ctx, domainList, nil, foundDomains)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	if len(domainsToAdd) > 0 {
		if resp.Diagnostics.Append(addDomains(ctx, client, domainsToAdd)...); resp.Diagnostics.HasError() {
			return
		}
	}
	tflog.Debug(ctx, "Domains were successfully created; fetching necessary data")
	foundDomains, err = client.SearchDomains(ctx, buildSearchDomainsRequest(domainList))
	if err != nil {
		resp.Diagnostics.AddError("Error getting status of added domains", err.Error())
		return
	}
	if resp.Diagnostics.Append(data.setDomains(ctx, foundDomains.Domains)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(resp.State.Set(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
}

func sortDomainList(domainList []domainResourceModel) {
	sort.Slice(domainList, func(i, j int) bool {
		if domainList[i].DomainName == domainList[j].DomainName {
			return domainList[i].ValidationScope.ValueString() < domainList[j].ValidationScope.ValueString()
		}
		return domainList[i].DomainName.ValueString() < domainList[j].DomainName.ValueString()
	})
}

func buildSearchDomainsRequest(domainList []domainResourceModel) domainownership.SearchDomainsRequest {
	domains := make([]domainownership.Domain, 0, len(domainList))
	for _, domain := range domainList {
		domains = append(domains, domainownership.Domain{
			DomainName:      domain.DomainName.ValueString(),
			ValidationScope: domainownership.ValidationScope(domain.ValidationScope.ValueString()),
		})
	}

	return domainownership.SearchDomainsRequest{
		IncludeAll: true,
		Body: domainownership.SearchDomainsBody{
			Domains: domains,
		},
	}
}

func buildAddAndRemoveDomainsRequest(ctx context.Context, planDomains, stateDomains []domainResourceModel, serverDomains *domainownership.SearchDomainsResponse) ([]domainownership.Domain, []domainownership.Domain, diag.Diagnostics) {
	domainsToAdd := make([]domainownership.Domain, 0, len(planDomains))
	domainsToImport := make([]domainownership.Domain, 0, len(planDomains))

	planMap := make(map[domainKey]domainResourceModel, len(planDomains))
	for _, domain := range planDomains {
		planMap[domainKey{domainName: domain.DomainName.ValueString(), validationScope: domain.ValidationScope.ValueString()}] = domain
	}
	stateMap := make(map[domainKey]domainResourceModel, len(stateDomains))
	for _, domain := range stateDomains {
		stateMap[domainKey{domainName: domain.DomainName.ValueString(), validationScope: domain.ValidationScope.ValueString()}] = domain
	}
	serverMap := make(map[domainKey]domainownership.SearchDomainItem, len(serverDomains.Domains))
	for _, domain := range serverDomains.Domains {
		serverMap[domainKey{domainName: domain.DomainName, validationScope: domain.ValidationScope}] = domain
	}

	diags := diag.Diagnostics{}

	// We need to calculate which domains are to be added and which are already present on the server (to be imported).
	for _, modelDomain := range planDomains {
		domainName := modelDomain.DomainName.ValueString()
		validationScope := modelDomain.ValidationScope.ValueString()
		// Was domains added in one of the previous apply?
		_, alreadyInState := stateMap[domainKey{domainName: domainName, validationScope: validationScope}]
		if !alreadyInState {
			// Was domain added to the server outside of terraform?
			serverItem, ok := serverMap[domainKey{domainName: domainName, validationScope: validationScope}]
			if ok {
				if serverItem.ValidationLevel == "ROOT/WILDCARD" {
					// Domains with ROOT/WILDCARD validationLevel can be added only if they are not yet VALIDATED.
					if serverItem.DomainStatus == "VALIDATED" {
						diags.Append(diag.NewErrorDiagnostic("error adding domains", fmt.Sprintf("domain %s with validation scope %s is already part of other, already validated domain/wildcard and cannot be added again", domainName, validationScope)))
					} else {
						domainsToAdd = append(domainsToAdd, domainownership.Domain{
							DomainName:      domainName,
							ValidationScope: domainownership.ValidationScope(validationScope),
						})
					}
					// Domain in TOKEN_EXPIRED or INVALIDATED status cannot be moved to VALIDATED directly, so we need to re-add them (if they were just added to plan).
				} else if serverItem.DomainStatus != "TOKEN_EXPIRED" && serverItem.DomainStatus != "INVALIDATED" {
					domainsToImport = append(domainsToImport, domainownership.Domain{
						DomainName:      domainName,
						ValidationScope: domainownership.ValidationScope(validationScope),
					})
				} else {
					// They are FQDN and in the status that requires re-adding.
					domainsToAdd = append(domainsToAdd, domainownership.Domain{
						DomainName:      domainName,
						ValidationScope: domainownership.ValidationScope(validationScope),
					})
				}
			} else {
				// Not found on the server, we need to add it.
				domainsToAdd = append(domainsToAdd, domainownership.Domain{
					DomainName:      domainName,
					ValidationScope: domainownership.ValidationScope(validationScope),
				})
			}
		}
	}
	if diags.HasError() {
		return nil, nil, diags
	}
	importMsg := strings.Builder{}
	for i, domain := range domainsToImport {
		if i > 0 {
			importMsg.WriteString(",")
		}
		importMsg.WriteString(fmt.Sprintf("%s:%s", domain.DomainName, domain.ValidationScope))
	}
	if importMsg.String() != "" {
		tflog.Debug(ctx, fmt.Sprintf("Some domains were already found on the server. Imported: [%s]", importMsg.String()))
	}
	if stateDomains == nil {
		return domainsToAdd, nil, nil
	}

	domainsToRemove := buildRemoveDomainsRequest(stateDomains, planMap)
	return domainsToAdd, domainsToRemove, nil
}

func buildRemoveDomainsRequest(stateDomains []domainResourceModel, planMap map[domainKey]domainResourceModel) []domainownership.Domain {
	domainsToRemove := make([]domainownership.Domain, 0, len(stateDomains))
	// Now we need to find which domains are to be removed.
	for _, domainFromState := range stateDomains {
		domainName := domainFromState.DomainName.ValueString()
		validationScope := domainFromState.ValidationScope.ValueString()
		if _, found := planMap[domainKey{domainName: domainName, validationScope: validationScope}]; !found {
			domainsToRemove = append(domainsToRemove, domainownership.Domain{
				DomainName:      domainName,
				ValidationScope: domainownership.ValidationScope(validationScope),
			})
		}
	}
	return domainsToRemove
}

func addDomains(ctx context.Context, client domainownership.DomainOwnership, domainsToAdd []domainownership.Domain) diag.Diagnostics {
	addedDomains, err := client.AddDomains(ctx, domainownership.AddDomainsRequest{Domains: domainsToAdd})
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("error adding domains", err.Error())}
	}
	if len(addedDomains.Errors) > 0 {
		tflog.Debug(ctx, "Adding domains resulted in errors, attempting rollback")
		domainsToRollback := make([]domainownership.Domain, 0, len(addedDomains.Successes))
		for _, domain := range addedDomains.Successes {
			domainsToRollback = append(domainsToRollback, domainownership.Domain{
				DomainName:      domain.DomainName,
				ValidationScope: domainownership.ValidationScope(domain.ValidationScope),
			})
		}
		if err := client.DeleteDomains(ctx, domainownership.DeleteDomainsRequest{Domains: domainsToRollback}); err != nil {
			return diag.Diagnostics{formatErrorMessageForAddDomains(addedDomains, err)}
		}
		return diag.Diagnostics{formatErrorMessageForAddDomains(addedDomains, nil)}
	}

	successMsg := strings.Builder{}
	for i, domain := range addedDomains.Successes {
		if i > 0 {
			successMsg.WriteString(",")
		}
		successMsg.WriteString(fmt.Sprintf("%s:%s", domain.DomainName, domain.ValidationScope))
	}
	tflog.Debug(ctx, fmt.Sprintf("Domains added successfully: [%s]", successMsg.String()))
	return nil
}

func formatErrorMessageForAddDomains(addedDomains *domainownership.AddDomainsResponse, err error) diag.Diagnostic {
	messageBuffer := strings.Builder{}
	for i, e := range addedDomains.Errors {
		if i > 0 {
			messageBuffer.WriteString(",\n")
		}
		messageBuffer.WriteString(fmt.Sprintf("{\n\tdomainName: %s,\n\tvalidationScope: %s,\n\ttitle: %s,\n\tdetail: %s\n}", e.DomainName, e.ValidationScope, e.Title, e.Detail))
	}
	if err != nil {
		return diag.NewErrorDiagnostic("error adding domains", fmt.Sprintf("%v\nRollback was not successful: %s", messageBuffer.String(), err))
	}
	return diag.NewErrorDiagnostic("error adding domains", fmt.Sprintf("%v\nRollback was successful", messageBuffer.String()))
}

func (m *domainsResourceModel) setDomains(ctx context.Context, domains []domainownership.SearchDomainItem) diag.Diagnostics {
	diags := diag.Diagnostics{}
	domainModels, d := getDomainsModelForSearch(ctx, domains)
	if diags.Append(d...); diags.HasError() {
		return diags
	}
	domainsModel, d := types.SetValueFrom(ctx, domainsType(), domainModels)
	if diags.Append(d...); diags.HasError() {
		return diags
	}
	m.Domains = domainsModel
	return diags
}

func getDomainsModelForSearch(ctx context.Context, domains []domainownership.SearchDomainItem) ([]domainResourceModel, diag.Diagnostics) {
	domainModels := make([]domainResourceModel, 0, len(domains))
	diags := diag.Diagnostics{}
	for _, domain := range domains {
		domainModel := domainResourceModel{
			DomainName:              types.StringValue(domain.DomainName),
			ValidationScope:         types.StringValue(domain.ValidationScope),
			AccountID:               types.StringPointerValue(domain.AccountID),
			DomainStatus:            types.StringValue(domain.DomainStatus),
			ValidationMethod:        types.StringPointerValue(domain.ValidationMethod),
			ValidationRequestedBy:   types.StringPointerValue(domain.ValidationRequestedBy),
			ValidationRequestedDate: date.TimeRFC3339PointerValue(domain.ValidationRequestedDate),
			ValidationCompletedDate: date.TimeRFC3339PointerValue(domain.ValidationCompletedDate),
		}

		if domain.ValidationChallenge != nil {
			validationModel := validationChallengeModel{
				CnameRecord: cnameRecordModel{
					Name:   types.StringValue(domain.ValidationChallenge.CnameRecord.Name),
					Target: types.StringValue(domain.ValidationChallenge.CnameRecord.Target),
				},
				TXTRecord: txtRecordModel{
					Name:  types.StringValue(domain.ValidationChallenge.TXTRecord.Name),
					Value: types.StringValue(domain.ValidationChallenge.TXTRecord.Value),
				},
				ExpirationDate: date.TimeRFC3339Value(domain.ValidationChallenge.ExpirationDate),
			}
			if domain.ValidationChallenge.HTTPRedirect != nil {
				validationModel.HTTPRedirect = &httpRedirectModel{
					From: types.StringValue(domain.ValidationChallenge.HTTPRedirect.From),
					To:   types.StringValue(domain.ValidationChallenge.HTTPRedirect.To),
				}
			}
			if domain.ValidationChallenge.HTTPFile != nil {
				validationModel.HTTPFile = &httpFileModel{
					Path:        types.StringValue(domain.ValidationChallenge.HTTPFile.Path),
					Content:     types.StringValue(domain.ValidationChallenge.HTTPFile.Content),
					ContentType: types.StringValue(domain.ValidationChallenge.HTTPFile.ContentType),
				}
			}

			challenge, d := types.ObjectValueFrom(ctx, validationChallengeType(), validationModel)
			if diags.Append(d...); diags.HasError() {
				return nil, diags
			}

			domainModel.ValidationChallenge = challenge
		} else {
			domainModel.ValidationChallenge = types.ObjectNull(validationChallengeType())
		}

		domainModels = append(domainModels, domainModel)
	}
	return domainModels, diags
}

// Read implements resource's Read method.
func (r *DomainsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Domain Ownership Domains Resource")
	var data domainsResourceModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := DomainOwnershipClient(r.meta)
	var domainList []domainResourceModel
	if resp.Diagnostics.Append(data.Domains.ElementsAs(ctx, &domainList, false)...); resp.Diagnostics.HasError() {
		return
	}
	sortDomainList(domainList)

	searchDomainsRequest := buildSearchDomainsRequest(domainList)
	foundDomains, err := client.SearchDomains(ctx, searchDomainsRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching status of domains from state", err.Error())
		return
	}
	if len(foundDomains.Domains) == 0 {
		tflog.Warn(ctx, "No domains from state found on the server, removing from state")
		resp.State.RemoveResource(ctx)
		return
	}
	foundDomainsMap := make(map[domainKey]struct{})
	for _, domain := range foundDomains.Domains {
		foundDomainsMap[domainKey{domainName: domain.DomainName, validationScope: domain.ValidationScope}] = struct{}{}
	}

	droppedDomains := make([]domainownership.Domain, 0, len(searchDomainsRequest.Body.Domains))
	for _, domain := range searchDomainsRequest.Body.Domains {
		if _, ok := foundDomainsMap[domainKey{domainName: domain.DomainName, validationScope: string(domain.ValidationScope)}]; !ok {
			droppedDomains = append(droppedDomains, domain)
		}
	}
	if len(droppedDomains) > 0 {
		droppedMsg := strings.Builder{}
		for i, domain := range droppedDomains {
			if i > 0 {
				droppedMsg.WriteString(",")
			}
			droppedMsg.WriteString(fmt.Sprintf("%s:%s", domain.DomainName, domain.ValidationScope))
		}
		tflog.Info(ctx, fmt.Sprintf("Some domains from state were not found on the server, removing from state: [%s]", droppedMsg.String()))
	}

	if resp.Diagnostics.Append(data.setDomains(ctx, foundDomains.Domains)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(resp.State.Set(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
}

// Update implements resource's Update method.
func (r *DomainsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Domain Ownership Domains Resource")
	var plan, state domainsResourceModel

	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	client := DomainOwnershipClient(r.meta)
	var domainsFromPlan, domainsFromState []domainResourceModel
	if resp.Diagnostics.Append(plan.Domains.ElementsAs(ctx, &domainsFromPlan, false)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(state.Domains.ElementsAs(ctx, &domainsFromState, false)...); resp.Diagnostics.HasError() {
		return
	}
	sortDomainList(domainsFromPlan)
	sortDomainList(domainsFromState)

	// Domains status may have changed since the read
	domainsFromState, diags := fetchLatestDomainsValue(ctx, client, buildSearchDomainsRequest(domainsFromState))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	foundDomains, err := client.SearchDomains(ctx, buildSearchDomainsRequest(domainsFromPlan))
	if err != nil {
		resp.Diagnostics.AddError("Error checking status of domains in configuration", err.Error())
		return
	}

	domainsToAdd, domainsToRemove, diags := buildAddAndRemoveDomainsRequest(ctx, domainsFromPlan, domainsFromState, foundDomains)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	if len(domainsToRemove) > 0 {
		if err := invalidateDomainsToDelete(ctx, client, domainsToRemove, domainsFromState); err != nil {
			resp.Diagnostics.AddError("Error deleting domains", err.Error())
			return
		}
	}
	if len(domainsToAdd) > 0 {
		if resp.Diagnostics.Append(addDomains(ctx, client, domainsToAdd)...); resp.Diagnostics.HasError() {
			return
		}
		tflog.Debug(ctx, "Domains were successfully added")
	}
	if len(domainsToRemove) > 0 {
		if err := client.DeleteDomains(ctx, domainownership.DeleteDomainsRequest{Domains: domainsToRemove}); err != nil {
			resp.Diagnostics.AddError("Error deleting domains", err.Error())
			return
		}
		tflog.Debug(ctx, "Domains were successfully deleted")
	}

	foundDomains, err = client.SearchDomains(ctx, buildSearchDomainsRequest(domainsFromPlan))
	if err != nil {
		resp.Diagnostics.AddError("Error getting status of added domains", err.Error())
		return
	}
	if resp.Diagnostics.Append(plan.setDomains(ctx, foundDomains.Domains)...); resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource's Delete method.
func (r *DomainsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting DomainOwnershipDomains Resource")
	var state domainsResourceModel

	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	var domainList []domainResourceModel
	if resp.Diagnostics.Append(state.Domains.ElementsAs(ctx, &domainList, false)...); resp.Diagnostics.HasError() {
		return
	}
	sortDomainList(domainList)

	client := DomainOwnershipClient(r.meta)

	// Domains status may have changed since the read
	domainList, diags := fetchLatestDomainsValue(ctx, client, buildSearchDomainsRequest(domainList))
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	domainsToRemove := make([]domainownership.Domain, 0, len(domainList))
	for _, domain := range domainList {
		domainsToRemove = append(domainsToRemove, domainownership.Domain{
			DomainName:      domain.DomainName.ValueString(),
			ValidationScope: domainownership.ValidationScope(domain.ValidationScope.ValueString()),
		})
	}

	if err := invalidateDomainsToDelete(ctx, client, domainsToRemove, domainList); err != nil {
		resp.Diagnostics.AddError("Error deleting domains", err.Error())
		return
	}

	if err := client.DeleteDomains(ctx, domainownership.DeleteDomainsRequest{Domains: domainsToRemove}); err != nil {
		resp.Diagnostics.AddError("Error deleting domains", err.Error())
		return
	}
}

func fetchLatestDomainsValue(ctx context.Context, client domainownership.DomainOwnership, searchRequest domainownership.SearchDomainsRequest) ([]domainResourceModel, diag.Diagnostics) {
	domains, err := client.SearchDomains(ctx, searchRequest)
	if err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("error fetching domains", err.Error())}
	}
	return getDomainsModelForSearch(ctx, domains.Domains)
}

func invalidateDomainsToDelete(ctx context.Context, client domainownership.DomainOwnership, domainList []domainownership.Domain, allDomains []domainResourceModel) error {
	validatedDomains := calculateValidatedDomains(domainList, allDomains)

	if len(validatedDomains) > 0 {
		var domainsList []string
		for _, domain := range validatedDomains {
			domainsList = append(domainsList, fmt.Sprintf("%s:%s", domain.DomainName, domain.ValidationScope))
		}

		tflog.Debug(ctx, fmt.Sprintf("Some domains requested for deletion are in VALIDATED status, invalidating them first: [%s]", strings.Join(domainsList, ",")))
		_, err := client.InvalidateDomains(ctx, domainownership.InvalidateDomainsRequest{Domains: validatedDomains})
		if err != nil {
			return err
		}
	}
	return nil
}

func calculateValidatedDomains(domainList []domainownership.Domain, allDomains []domainResourceModel) []domainownership.Domain {
	validatedDomains := make([]domainownership.Domain, 0, len(domainList))
	allDomainsMap := make(map[domainKey]string, len(allDomains))
	for _, domain := range allDomains {
		allDomainsMap[domainKey{domainName: domain.DomainName.ValueString(), validationScope: domain.ValidationScope.ValueString()}] = domain.DomainStatus.ValueString()
	}

	for _, domain := range domainList {
		if status, ok := allDomainsMap[domainKey{domainName: domain.DomainName, validationScope: string(domain.ValidationScope)}]; ok {
			if status == "VALIDATED" {
				validatedDomains = append(validatedDomains, domain)
			}
		}
	}
	return validatedDomains
}

// ImportState implements resource's ImportState method.
func (r *DomainsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	client := DomainOwnershipClient(r.meta)
	domains, diags := parseDomains(ctx, client, id, false)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var data domainsResourceModel
	var domainModels []domainResourceModel
	for _, domain := range domains {
		domainModels = append(domainModels, domainResourceModel{
			DomainName:              types.StringValue(domain.DomainName),
			ValidationScope:         types.StringValue(string(domain.ValidationScope)),
			AccountID:               types.StringNull(),
			DomainStatus:            types.StringNull(),
			ValidationMethod:        types.StringNull(),
			ValidationRequestedBy:   types.StringNull(),
			ValidationRequestedDate: types.StringNull(),
			ValidationCompletedDate: types.StringNull(),
			ValidationChallenge:     types.ObjectNull(validationChallengeType()),
		})
	}
	domainsModel, diags := types.SetValueFrom(ctx, domainsType(), domainModels)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
	data.Domains = domainsModel
	if resp.Diagnostics.Append(resp.State.Set(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
}
