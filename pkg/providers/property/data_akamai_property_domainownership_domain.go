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
				Description: "Name of the domain.",
			},
			"validation_scope": schema.StringAttribute{
				Required: true,
				Description: "Validation scope of the domain. For HOST, the scope is only the exactly specified domain. " +
					"For WILDCARD, the scope covers any hostname within one subdomain level. " +
					"For DOMAIN, the scope covers any hostnames under the domain, regardless of the level of subdomains.",
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
				Description: "Validation status of the domain, either `REQUEST_ACCEPTED`, `VALIDATION_IN_PROGRESS`, `VALIDATED`, `TOKEN_EXPIRED`, or `INVALIDATED`.",
			},
			"validation_method": schema.StringAttribute{
				Computed:    true,
				Description: "Method of the domain validation, either `DNS_CNAME`, `DNS_TXT`, `HTTP`, `SYSTEM`, or `MANUAL`.",
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
					"cname_record": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "CNAME record details for domain validation.",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:    true,
								Description: "The name of the CNAME record.",
							},
							"target": schema.StringAttribute{
								Computed:    true,
								Description: "The target value of the CNAME record.",
							},
						},
					},
					"txt_record": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "TXT record details for domain validation.",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:    true,
								Description: "The name of the TXT record.",
							},
							"value": schema.StringAttribute{
								Computed:    true,
								Description: "The value of the TXT record.",
							},
						},
					},
					"http_file": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "HTTP file details for domain validation.",
						Attributes: map[string]schema.Attribute{
							"path": schema.StringAttribute{
								Computed:    true,
								Description: "The path where the file should be accessible.",
							},
							"content": schema.StringAttribute{
								Computed:    true,
								Description: "The content of the file.",
							},
							"content_type": schema.StringAttribute{
								Computed:    true,
								Description: "The content type of the file.",
							},
						},
					},
					"http_redirect": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "HTTP redirect details for domain validation.",
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								Computed:    true,
								Description: "HTTP URL for checking the challenge token during HTTP validation.",
							},
							"to": schema.StringAttribute{
								Computed:    true,
								Description: "HTTP redirect URL for HTTP validation.",
							},
						},
					},
					"expiration_date": schema.StringAttribute{
						Computed:    true,
						Description: "The ISO 8601 timestamp indicating when the validation challenge expires.",
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
							Description: "Status of the domain, either `REQUEST_ACCEPTED`, `VALIDATION_IN_PROGRESS`, `VALIDATED`, `TOKEN_EXPIRED`, or `INVALIDATED`.",
						},
						"modified_user": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the user who requested the status change.",
						},
						"modified_date": schema.StringAttribute{
							Computed:    true,
							Description: "An ISO 8601 timestamp indicating when the domain status changed.",
						},
						"message": schema.StringAttribute{
							Computed:    true,
							Description: "Information about the status change.",
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
