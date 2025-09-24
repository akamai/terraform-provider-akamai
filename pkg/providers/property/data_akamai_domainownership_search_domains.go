package property

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &domainOwnershipSearchDomains{}
var _ datasource.DataSourceWithConfigure = &domainOwnershipSearchDomains{}

// NewDomainOwnershipSearchDomains returns a new domain ownership search domains data source.
func NewDomainOwnershipSearchDomains() datasource.DataSource {
	return &domainOwnershipSearchDomains{}
}

// domainOwnershipSearchDomains defines the data source implementation for domain ownership search domains.
type domainOwnershipSearchDomains struct {
	meta meta.Meta
}

// domainOwnershipSearchDomainsModel describes the data source data model for PropertyDomainOwnershipSearchDomains.
type domainOwnershipSearchDomainsModel struct {
	Domains []domainOwnershipDomainModel `tfsdk:"domains"`
}

// domainOwnershipDomainModel models each domain in the domains set.
type domainOwnershipDomainModel struct {
	DomainName              types.String                             `tfsdk:"domain_name"`
	ValidationScope         types.String                             `tfsdk:"validation_scope"`
	AccountID               types.String                             `tfsdk:"account_id"`
	DomainStatus            types.String                             `tfsdk:"domain_status"`
	ValidationLevel         types.String                             `tfsdk:"validation_level"`
	ValidationMethod        types.String                             `tfsdk:"validation_method"`
	ValidationRequestedBy   types.String                             `tfsdk:"validation_requested_by"`
	ValidationRequestedDate types.String                             `tfsdk:"validation_requested_date"`
	ValidationCompletedDate types.String                             `tfsdk:"validation_completed_date"`
	ValidationChallenge     *domainOwnershipValidationChallengeModel `tfsdk:"validation_challenge"`
}

// domainOwnershipValidationChallengeModel models the nested validation_challenge attribute.
type domainOwnershipValidationChallengeModel struct {
	DNSCName                  types.String `tfsdk:"dns_cname"`
	ChallengeToken            types.String `tfsdk:"challenge_token"`
	ChallengeTokenExpiresDate types.String `tfsdk:"challenge_token_expires_date"`
	HTTPRedirectFrom          types.String `tfsdk:"http_redirect_from"`
	HTTPRedirectTo            types.String `tfsdk:"http_redirect_to"`
}

// Metadata configures data source's meta information
func (d *domainOwnershipSearchDomains) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_domainownership_search_domains"
}

// Schema is used to define data source's terraform schema
func (d *domainOwnershipSearchDomains) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Domain Ownership - Search Domains data source",
		Attributes: map[string]schema.Attribute{
			"domains": schema.SetNestedAttribute{
				Required:    true,
				Description: "List of domains.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain_name": schema.StringAttribute{
							Required:    true,
							Description: "Name of the domain.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 253),
							},
						},
						"validation_scope": schema.StringAttribute{
							Required:    true,
							Description: "Scope of the domain validation, either HOST, WILDCARD, or DOMAIN.",
							Validators: []validator.String{
								stringvalidator.OneOf("HOST", "WILDCARD", "DOMAIN"),
							},
						},
						"account_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of an account.",
						},
						"domain_status": schema.StringAttribute{
							Computed:    true,
							Description: "Validation status of the domain, either REQUEST_ACCEPTED, VALIDATION_IN_PROGRESS, VALIDATED, TOKEN_EXPIRED, or INVALIDATED.",
						},
						"validation_level": schema.StringAttribute{
							Computed:    true,
							Description: "Level of the domain validation, either FQDN or WILDCARD.",
						},
						"validation_method": schema.StringAttribute{
							Computed:    true,
							Description: "Method of the domain validation, either DNS_CNAME, DNS_TXT, HTTP, SYSTEM, or MANUAL.",
						},
						"validation_requested_by": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the user who requested the domain validation.",
						},
						"validation_requested_date": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp of the request.",
						},
						"validation_completed_date": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp of completing the validation.",
						},
						"validation_challenge": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Validation challenge of the domain.",
							Attributes: map[string]schema.Attribute{
								"dns_cname": schema.StringAttribute{
									Computed:    true,
									Description: "DNS CNAME you need to use for DNS CNAME domain validation.",
								},
								"challenge_token": schema.StringAttribute{
									Computed:    true,
									Description: "Challenge token you need to use for domain validation.",
								},
								"challenge_token_expires_date": schema.StringAttribute{
									Computed:    true,
									Description: "An ISO 8601 timestamp indicating when the domain validation token expires.",
								},
								"http_redirect_from": schema.StringAttribute{
									Computed:    true,
									Description: "HTTP URL for checking the challenge token during HTTP validation.",
								},
								"http_redirect_to": schema.StringAttribute{
									Computed:    true,
									Description: "HTTP redirect URL for HTTP validation.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure  configures data source at the beginning of the lifecycle
func (d *domainOwnershipSearchDomains) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	d.meta = meta.Must(req.ProviderData)
}

// Read is called when the provider must read data source values in order to update state
func (d *domainOwnershipSearchDomains) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "DomainOwnershipSearchDomains Read")

	var data domainOwnershipSearchDomainsModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	domains := make([]domainownership.Domain, 0, len(data.Domains))
	for _, domain := range data.Domains {
		domains = append(domains, domainownership.Domain{
			DomainName:      domain.DomainName.ValueString(),
			ValidationScope: domainownership.ValidationScope(domain.ValidationScope.ValueString()),
		})
	}
	client := DomainOwnershipClient(d.meta)
	response, err := client.SearchDomains(ctx, domainownership.SearchDomainsRequest{
		IncludeAll: true,
		Body: domainownership.SearchDomainsBody{
			Domains: domains,
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("searching domains failed", err.Error())
		return
	}

	data.Domains = make([]domainOwnershipDomainModel, 0, len(response.Domains))
	for _, domain := range response.Domains {
		dm := domainOwnershipDomainModel{
			DomainName:              types.StringValue(domain.DomainName),
			ValidationScope:         types.StringValue(domain.ValidationScope),
			AccountID:               types.StringPointerValue(domain.AccountID),
			DomainStatus:            types.StringValue(domain.DomainStatus),
			ValidationLevel:         types.StringValue(domain.ValidationLevel),
			ValidationMethod:        types.StringPointerValue(domain.ValidationMethod),
			ValidationRequestedBy:   types.StringPointerValue(domain.ValidationRequestedBy),
			ValidationRequestedDate: date.TimeRFC3339PointerValue(domain.ValidationRequestedDate),
			ValidationCompletedDate: date.TimeRFC3339PointerValue(domain.ValidationCompletedDate),
		}
		if domain.ValidationChallenge != nil {
			dm.ValidationChallenge = &domainOwnershipValidationChallengeModel{
				DNSCName:                  types.StringValue(domain.ValidationChallenge.DNSCname),
				ChallengeToken:            types.StringValue(domain.ValidationChallenge.ChallengeToken),
				ChallengeTokenExpiresDate: date.TimeRFC3339Value(domain.ValidationChallenge.ChallengeTokenExpiresDate),
				HTTPRedirectFrom:          types.StringPointerValue(domain.ValidationChallenge.HTTPRedirectFrom),
				HTTPRedirectTo:            types.StringPointerValue(domain.ValidationChallenge.HTTPRedirectTo),
			}
		}
		data.Domains = append(data.Domains, dm)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
