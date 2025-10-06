package property

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &domainsDataSource{}
	_ datasource.DataSourceWithConfigure = &domainsDataSource{}
)

type (
	domainsDataSource struct {
		meta meta.Meta
	}
	domainsDataSourceModel struct {
		Domains []domainDetails `tfsdk:"domains"`
	}

	domainDetails struct {
		Name                    types.String         `tfsdk:"domain_name"`
		ValidationScope         types.String         `tfsdk:"validation_scope"`
		AccountID               types.String         `tfsdk:"account_id"`
		DomainStatus            types.String         `tfsdk:"domain_status"`
		ValidationMethod        types.String         `tfsdk:"validation_method"`
		ValidationRequestedBy   types.String         `tfsdk:"validation_requested_by"`
		ValidationRequestedDate types.String         `tfsdk:"validation_requested_date"`
		ValidationCompletedDate types.String         `tfsdk:"validation_completed_date"`
		ValidationChallenge     *validationChallenge `tfsdk:"validation_challenge"`
	}

	validationChallenge struct {
		DNSCname                  types.String `tfsdk:"dns_cname"`
		ChallengeToken            types.String `tfsdk:"challenge_token"`
		ChallengeTokenExpiresDate types.String `tfsdk:"challenge_token_expires_date"`
		HTTPRedirectFrom          types.String `tfsdk:"http_redirect_from"`
		HTTPRedirectTo            types.String `tfsdk:"http_redirect_to"`
	}
)

// NewDomainOwnershipDomainsDataSource returns a new domainDataSource.
func NewDomainOwnershipDomainsDataSource() datasource.DataSource {
	return &domainsDataSource{}
}

// Metadata configures data source's meta information.
func (d domainsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_domainownership_domains"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *domainsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					req.ProviderData))
		}
	}()
	d.meta = meta.Must(req.ProviderData)
}

// Schema is used to define data source's terraform schema.
func (d domainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details of Domains.",
		Attributes: map[string]schema.Attribute{
			"domains": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of domains",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain_name": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Name of the domain.",
						},
						"validation_scope": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Scope of the domain validation, either HOST, WILDCARD, or DOMAIN.",
						},
						"account_id": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "ID of an account.",
						},
						"domain_status": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Validation status of the domain, either REQUEST_ACCEPTED, VALIDATION_IN_PROGRESS, VALIDATED, TOKEN_EXPIRED, or INVALIDATED.",
						},
						"validation_method": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Method of the domain validation, either DNS_CNAME, DNS_TXT, HTTP, SYSTEM, or MANUAL.",
						},
						"validation_requested_by": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Name of the user who requested domain validation.",
						},
						"validation_requested_date": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Timestamp of the request.",
						},
						"validation_completed_date": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Timestamp of completing the validation.",
						},
						"validation_challenge": schema.SingleNestedAttribute{
							Computed:    true,
							Optional:    true,
							Description: "Validation challenge of the domain.",
							Attributes: map[string]schema.Attribute{
								"dns_cname": schema.StringAttribute{
									Computed:    true,
									Optional:    true,
									Description: "DNS CNAME you need to use for DNS CNAME domain validation.",
								},
								"challenge_token": schema.StringAttribute{
									Computed:    true,
									Optional:    true,
									Description: "Challenge token you need to use for domain validation.",
								},
								"challenge_token_expires_date": schema.StringAttribute{
									Computed:    true,
									Optional:    true,
									Description: "An ISO 8601 timestamp indicating when the domain validation token expires.",
								},
								"http_redirect_from": schema.StringAttribute{
									Computed:    true,
									Optional:    true,
									Description: "HTTP URL for checking the challenge token during HTTP validation.",
								},
								"http_redirect_to": schema.StringAttribute{
									Computed:    true,
									Optional:    true,
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

func (d *domainsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "DomainOwnership Domains DataSource Read")
	var data domainsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := DomainOwnershipClient(d.meta)

	//Pagination parameters
	pageSize := int64(1000)
	page := int64(1)
	paginate := true

	var domains domainownership.ListDomainsResponse

	for {
		tflog.Debug(ctx, fmt.Sprintf("Fetching domains page %d with size %d", page, pageSize)) // Good for debugging

		domainsResp, err := client.ListDomains(ctx, domainownership.ListDomainsRequest{
			Paginate: &paginate,
			Page:     &page,
			PageSize: &pageSize,
		})

		if err != nil {
			resp.Diagnostics.AddError("Read DomainOwnership Domains failed", err.Error())
			return
		}
		if domainsResp == nil {
			break
		}

		domains.Domains = append(domains.Domains, domainsResp.Domains...)
		if !domainsResp.Metadata.HasNext {
			break
		}
		page++
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully retrieved a total of %d domains.", len(domains.Domains)))
	data.convertDomainsToModel(domains)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *domainsDataSourceModel) convertDomainsToModel(domains domainownership.ListDomainsResponse) {
	m.Domains = make([]domainDetails, len(domains.Domains))
	for i, domain := range domains.Domains {
		currentDomain := domainDetails{
			AccountID:               types.StringValue(domain.AccountID),
			DomainStatus:            types.StringValue(domain.DomainStatus),
			Name:                    types.StringValue(domain.DomainName),
			ValidationScope:         types.StringValue(domain.ValidationScope),
			ValidationMethod:        types.StringPointerValue(domain.ValidationMethod),
			ValidationRequestedBy:   types.StringValue(domain.ValidationRequestedBy),
			ValidationRequestedDate: date.TimeRFC3339NanoValue(domain.ValidationRequestedDate),
			ValidationCompletedDate: date.TimeRFC3339NanoPointerValue(domain.ValidationCompletedDate),
			ValidationChallenge:     &validationChallenge{},
		}

		if domain.ValidationChallenge != nil {
			vc := &validationChallenge{
				ChallengeToken:            types.StringValue(domain.ValidationChallenge.ChallengeToken),
				ChallengeTokenExpiresDate: date.TimeRFC3339NanoValue(domain.ValidationChallenge.ChallengeTokenExpiresDate),
				DNSCname:                  types.StringValue(domain.ValidationChallenge.DNSCname),
				HTTPRedirectFrom:          types.StringPointerValue(domain.ValidationChallenge.HTTPRedirectFrom),
				HTTPRedirectTo:            types.StringPointerValue(domain.ValidationChallenge.HTTPRedirectTo),
			}

			if vc.isEmpty() {
				currentDomain.ValidationChallenge = nil
			} else {
				currentDomain.ValidationChallenge = vc
			}
		}

		m.Domains[i] = currentDomain
	}
}

func (vc *validationChallenge) isEmpty() bool {
	return vc.ChallengeToken.IsNull() &&
		vc.ChallengeTokenExpiresDate.IsNull() &&
		vc.DNSCname.IsNull() &&
		vc.HTTPRedirectFrom.IsNull() &&
		vc.HTTPRedirectTo.IsNull()
}
