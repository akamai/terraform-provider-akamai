package property

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &domainDataSource{}
	_ datasource.DataSourceWithConfigure = &domainDataSource{}
)

type (
	domainDataSource struct {
		meta meta.Meta
	}

	domainDataSourceModel struct {
		Name                    types.String              `tfsdk:"domain_name"`
		ValidationScope         types.String              `tfsdk:"validation_scope"`
		AccountID               types.String              `tfsdk:"account_id"`
		DomainStatus            types.String              `tfsdk:"domain_status"`
		ValidationMethod        types.String              `tfsdk:"validation_method"`
		ValidationRequestedBy   types.String              `tfsdk:"validation_requested_by"`
		ValidationRequestedDate types.String              `tfsdk:"validation_requested_date"`
		ValidationCompletedDate types.String              `tfsdk:"validation_completed_date"`
		ValidationChallenge     *validationChallengeModel `tfsdk:"validation_challenge"`
		DomainStatusHistory     []domainStatusHistory     `tfsdk:"domain_status_history"`
	}

	validationChallengeModel struct {
		CnameRecord    cnameRecordModel   `tfsdk:"cname_record"`
		TXTRecord      txtRecordModel     `tfsdk:"txt_record"`
		HTTPFile       *httpFileModel     `tfsdk:"http_file"`
		HTTPRedirect   *httpRedirectModel `tfsdk:"http_redirect"`
		ExpirationDate types.String       `tfsdk:"expiration_date"`
	}

	// cnameRecordModel represents a CNAME record for domain validation cnameRecordModel.
	cnameRecordModel struct {
		Name   types.String `tfsdk:"name"`
		Target types.String `tfsdk:"target"`
	}

	// txtRecordModel represents a TXT record for domain validation txtRecordModel.
	txtRecordModel struct {
		Name  types.String `tfsdk:"name"`
		Value types.String `tfsdk:"value"`
	}

	// httpFileModel represents an HTTP file for domain validation httpFileModel.
	httpFileModel struct {
		Path        types.String `tfsdk:"path"`
		Content     types.String `tfsdk:"content"`
		ContentType types.String `tfsdk:"content_type"`
	}

	// httpRedirectModel represents an HTTP redirect for domain validation httpRedirectModel.
	httpRedirectModel struct {
		From types.String `tfsdk:"from"`
		To   types.String `tfsdk:"to"`
	}

	domainStatusHistory struct {
		DomainStatus types.String `tfsdk:"domain_status"`
		ModifiedDate types.String `tfsdk:"modified_date"`
		ModifiedUser types.String `tfsdk:"modified_user"`
		Message      types.String `tfsdk:"message"`
	}
)

// NewDomainOwnershipDomainDataSource returns a new domainDataSource.
func NewDomainOwnershipDomainDataSource() datasource.DataSource {
	return &domainDataSource{}
}

// Metadata configures data source's meta information.
func (d *domainDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_property_domainownership_domain"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *domainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *domainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details of a Domain.",
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
			"validation_challenge": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The domain's validation challenge details.",
				Attributes: map[string]schema.Attribute{
					"cname_record": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "The details of the 'CNAME' record you copy to your DNS configuration to prove you own the domain. You should use the 'DNS_CNAME' method in most cases.",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:    true,
								Description: "The 'CNAME' record for your domain that you add to the DNS configuration.",
							},
							"target": schema.StringAttribute{
								Computed:    true,
								Description: "The 'target' value you set in the 'CNAME' record that validates the domain ownership.",
							},
						},
					},
					"txt_record": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "The details of the 'TXT' record with the challenge token that you copy to your DNS configuration to prove you own the domain.",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:    true,
								Description: "The hostname where you should add the 'TXT' record to validate the domain ownership.",
							},
							"value": schema.StringAttribute{
								Computed:    true,
								Description: "The token you need to copy to the DNS 'TXT' record that validates the domain ownership.",
							},
						},
					},
					"http_file": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Available only for the 'HOST' validation scope. The details for the HTTP validation method in which you create a file containing a token and save it on your HTTP server at the provided URL. Alternatively, you can use the 'http_redirect' method.",
						Attributes: map[string]schema.Attribute{
							"path": schema.StringAttribute{
								Computed:    true,
								Description: "The URL where you should place the file containing the challenge token.",
							},
							"content": schema.StringAttribute{
								Computed:    true,
								Description: "The content of the file that you should place at the specified URL.",
							},
							"content_type": schema.StringAttribute{
								Computed:    true,
								Description: "The content type of the file containing the token.",
							},
						},
					},
					"http_redirect": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Available only for the 'HOST' validation scope. The details for the HTTP validation method in which you use a redirect URL with the token. Alternatively, you can use the 'http_file' method.",
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								Computed:    true,
								Description: "The location on your HTTP server where you set up the redirect.",
							},
							"to": schema.StringAttribute{
								Computed:    true,
								Description: "The redirect URL with the token that you place on your HTTP server.",
							},
						},
					},
					"expiration_date": schema.StringAttribute{
						Computed:    true,
						Description: "The timestamp indicating when the challenge data expires.",
					},
				},
			},
			"domain_status_history": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of domain status history changes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain_status": schema.StringAttribute{
							Computed:    true,
							Description: "The domain's validation status. Possible values are: `REQUEST_ACCEPTED`, `VALIDATION_IN_PROGRESS`, `VALIDATED`, `TOKEN_EXPIRED`, or `INVALIDATED`.",
						},
						"modified_user": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the user who requested the status change.",
						},
						"modified_date": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp indicating when the domain status changed.",
						},
						"message": schema.StringAttribute{
							Computed:    true,
							Description: "Additional information about the status change.",
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *domainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Domain Ownership Domain DataSource Read")

	var data domainDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	domainownershipClient = DomainOwnershipClient(d.meta)

	domain, err := domainownershipClient.GetDomain(ctx, domainownership.GetDomainRequest{
		DomainName:                 data.Name.ValueString(),
		ValidationScope:            domainownership.ValidationScope(data.ValidationScope.ValueString()),
		IncludeDomainStatusHistory: true,
	})
	if err != nil {
		resp.Diagnostics.AddError("Read Domain Ownership Domain failed", err.Error())
		return
	}

	data.convertDomainToModel(*domain)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *domainDataSourceModel) convertDomainToModel(domain domainownership.GetDomainResponse) {
	m.AccountID = types.StringValue(domain.AccountID)
	m.DomainStatus = types.StringValue(domain.DomainStatus)
	m.Name = types.StringValue(domain.DomainName)
	m.ValidationScope = types.StringValue(domain.ValidationScope)
	m.ValidationMethod = types.StringPointerValue(domain.ValidationMethod)
	m.ValidationRequestedBy = types.StringValue(domain.ValidationRequestedBy)
	m.ValidationRequestedDate = date.TimeRFC3339NanoValue(domain.ValidationRequestedDate)
	m.ValidationCompletedDate = date.TimeRFC3339NanoPointerValue(domain.ValidationCompletedDate)

	if domain.ValidationChallenge != nil {
		challenge := domain.ValidationChallenge
		challengeModel := &validationChallengeModel{}

		challengeModel.CnameRecord = cnameRecordModel{
			Name:   types.StringValue(challenge.CnameRecord.Name),
			Target: types.StringValue(challenge.CnameRecord.Target),
		}

		challengeModel.TXTRecord = txtRecordModel{
			Name:  types.StringValue(challenge.TXTRecord.Name),
			Value: types.StringValue(challenge.TXTRecord.Value),
		}

		if challenge.HTTPFile != nil {
			challengeModel.HTTPFile = &httpFileModel{
				Path:        types.StringValue(challenge.HTTPFile.Path),
				Content:     types.StringValue(challenge.HTTPFile.Content),
				ContentType: types.StringValue(challenge.HTTPFile.ContentType),
			}
		}

		if challenge.HTTPRedirect != nil {
			challengeModel.HTTPRedirect = &httpRedirectModel{
				From: types.StringValue(challenge.HTTPRedirect.From),
				To:   types.StringValue(challenge.HTTPRedirect.To),
			}
		}

		challengeModel.ExpirationDate = date.TimeRFC3339NanoValue(challenge.ExpirationDate)

		m.ValidationChallenge = challengeModel
	}

	for _, history := range domain.DomainStatusHistory {
		m.DomainStatusHistory = append(m.DomainStatusHistory, domainStatusHistory{
			DomainStatus: types.StringValue(history.DomainStatus),
			ModifiedDate: date.TimeRFC3339NanoValue(history.ModifiedDate),
			ModifiedUser: types.StringValue(history.ModifiedUser),
			Message:      types.StringPointerValue(history.Message),
		})
	}
}
