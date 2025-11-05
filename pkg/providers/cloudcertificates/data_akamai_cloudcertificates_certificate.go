package cloudcertificates

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &certificateDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateDataSource{}
)

type (
	certificateDataSource struct {
		meta meta.Meta
	}

	certificateDataSourceModel struct {
		CertificateID                       types.String   `tfsdk:"certificate_id"`
		IncludeHostnameBindings             types.Bool     `tfsdk:"include_hostname_bindings"`
		AccountID                           types.String   `tfsdk:"account_id"`
		CertificateName                     types.String   `tfsdk:"certificate_name"`
		CertificateStatus                   types.String   `tfsdk:"certificate_status"`
		CertificateType                     types.String   `tfsdk:"certificate_type"`
		ContractID                          types.String   `tfsdk:"contract_id"`
		CreatedBy                           types.String   `tfsdk:"created_by"`
		CreatedDate                         types.String   `tfsdk:"created_date"`
		CSRExpirationDate                   types.String   `tfsdk:"csr_expiration_date"`
		CSRPEM                              types.String   `tfsdk:"csr_pem"`
		KeySize                             types.String   `tfsdk:"key_size"`
		KeyType                             types.String   `tfsdk:"key_type"`
		ModifiedBy                          types.String   `tfsdk:"modified_by"`
		ModifiedDate                        types.String   `tfsdk:"modified_date"`
		SANs                                types.Set      `tfsdk:"sans"`
		SecureNetwork                       types.String   `tfsdk:"secure_network"`
		SignedCertificateIssuer             types.String   `tfsdk:"signed_certificate_issuer"`
		SignedCertificateNotValidAfterDate  types.String   `tfsdk:"signed_certificate_not_valid_after_date"`
		SignedCertificateNotValidBeforeDate types.String   `tfsdk:"signed_certificate_not_valid_before_date"`
		SignedCertificatePEM                types.String   `tfsdk:"signed_certificate_pem"`
		SignedCertificateSHA256Fingerprint  types.String   `tfsdk:"signed_certificate_sha256_fingerprint"`
		SignedCertificateSerialNumber       types.String   `tfsdk:"signed_certificate_serial_number"`
		Subject                             *subjectModel  `tfsdk:"subject"`
		TrustChainPEM                       types.String   `tfsdk:"trust_chain_pem"`
		Bindings                            []bindingModel `tfsdk:"bindings"`
	}
)

const defaultPageSize int64 = 100

// NewCertificateDataSource returns a new CloudCertificates Certificate data source.
func NewCertificateDataSource() datasource.DataSource {
	return &certificateDataSource{}
}

func (d *certificateDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_certificate"
}

func (d *certificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *certificateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details of a specific cloud certificate.",
		Attributes: map[string]schema.Attribute{
			"certificate_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the certificate.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"include_hostname_bindings": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to include hostname bindings from the certificate-bindings API. By default, hostname bindings are not included.",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Account associated with 'contract_id'.",
			},
			"certificate_name": schema.StringAttribute{
				Computed:    true,
				Description: "The certificate name.",
			},
			"certificate_status": schema.StringAttribute{
				Computed:    true,
				Description: "The current status of the certificate.",
			},
			"certificate_type": schema.StringAttribute{
				Computed:    true,
				Description: "The certificate type.",
			},
			"contract_id": schema.StringAttribute{
				Computed:    true,
				Description: "Contract ID under which this certificate was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who created the certificate.",
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the certificate (UTC).",
			},
			"csr_expiration_date": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration time of the certificate signing request (UTC).",
			},
			"csr_pem": schema.StringAttribute{
				Computed:    true,
				Description: "PEM-encoded certificate signing request (CSR) generated by Akamai for your selected key type.",
			},
			"key_size": schema.StringAttribute{
				Computed:    true,
				Description: "The certificate key size.",
			},
			"key_type": schema.StringAttribute{
				Computed:    true,
				Description: "The certificate key type.",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who last updated the certificate.",
			},
			"modified_date": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the certificate (UTC).",
			},
			"sans": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The list of Subject Alternative Names (SANs) for the certificate.",
			},
			"secure_network": schema.StringAttribute{
				Computed:    true,
				Description: "The secure network type.",
			},
			"signed_certificate_issuer": schema.StringAttribute{
				Description: "The issuer of the signed certificate.",
				Computed:    true,
			},
			"signed_certificate_not_valid_after_date": schema.StringAttribute{
				Description: "The date after which the signed certificate is no longer valid.",
				Computed:    true,
			},
			"signed_certificate_not_valid_before_date": schema.StringAttribute{
				Description: "The date before which the signed certificate is not valid.",
				Computed:    true,
			},
			"signed_certificate_pem": schema.StringAttribute{
				Description: "PEM-encoded signed certificate uploaded for the selected key type.",
				Computed:    true,
			},
			"signed_certificate_sha256_fingerprint": schema.StringAttribute{
				Description: "The SHA256 fingerprint of the signed certificate.",
				Computed:    true,
			},
			"signed_certificate_serial_number": schema.StringAttribute{
				Description: "The serial number of the signed certificate in hex format.",
				Computed:    true,
			},
			"subject": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Subject fields as defined in X.509 certificates (RFC 5280).",
				Attributes: map[string]schema.Attribute{
					"common_name": schema.StringAttribute{
						Computed:    true,
						Description: "Fully qualified domain name (FQDN) or other name associated with the subject. If specified, this value must also be included in the SANs list.",
					},
					"organization": schema.StringAttribute{
						Computed:    true,
						Description: "Legal name of the organization.",
					},
					"country": schema.StringAttribute{
						Computed:    true,
						Description: "Two-letter ISO 3166 country code.",
					},
					"state": schema.StringAttribute{
						Computed:    true,
						Description: "Full name of the state or province.",
					},
					"locality": schema.StringAttribute{
						Computed:    true,
						Description: "City or locality name.",
					},
				},
			},
			"trust_chain_pem": schema.StringAttribute{
				Description: "The trust chain PEM content uploaded by end user.",
				Computed:    true,
			},
			"bindings": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of hostname bindings for the certificate identified by certificate_id.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"certificate_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the third-party certificate.",
						},
						"hostname": schema.StringAttribute{
							Computed:    true,
							Description: "Hostname on the Akamai CDN the certificate applies to.",
						},
						"network": schema.StringAttribute{
							Computed:    true,
							Description: "The deployment network, either STAGING or PRODUCTION, on which the certificate is active for a property version.",
						},
						"resource_type": schema.StringAttribute{
							Computed:    true,
							Description: "Resource type this binding applies to. Currently, only CDN_HOSTNAME is available.",
						},
					},
				},
			},
		},
	}
}

func (d *certificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CCM Certificate DataSource Read")
	var data certificateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)

	cert, err := client.GetCertificate(ctx, cloudcertificates.GetCertificateRequest{
		CertificateID: data.CertificateID.ValueString(),
	})
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrCertificateNotFound) {
			resp.Diagnostics.AddError("Certificate Not Found", fmt.Sprintf("No certificate found with ID: %s", data.CertificateID.ValueString()))
		} else {
			resp.Diagnostics.AddError("Failed to retrieve certificate", err.Error())
		}
		return
	}

	// IncludeHostnameBindings defaults to false - only include bindings if explicitly set to true
	if data.IncludeHostnameBindings.ValueBool() {
		allBindings, err := getAllBindings(ctx, client, data.CertificateID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to retrieve bindings", err.Error())
			return
		}
		data.convertBindingsToModel(allBindings)
	}

	resp.Diagnostics.Append(data.convertCertificateToModel(ctx, *cert)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *certificateDataSourceModel) convertBindingsToModel(bindings []cloudcertificates.CertificateBinding) {
	m.Bindings = make([]bindingModel, len(bindings))
	for i, b := range bindings {
		m.Bindings[i] = bindingModel{
			CertificateID: types.StringValue(b.CertificateID.String()),
			Hostname:      types.StringValue(b.Hostname),
			Network:       types.StringValue(b.Network),
			ResourceType:  types.StringValue(b.ResourceType),
		}
	}
}

func getAllBindings(ctx context.Context, client cloudcertificates.CloudCertificates, certificateID string) ([]cloudcertificates.CertificateBinding, error) {
	pageSize := defaultPageSize
	var page int64 = 1
	var allBindings []cloudcertificates.CertificateBinding

	for {
		tflog.Debug(ctx, fmt.Sprintf("Fetching bindings page %d with size %d", page, pageSize))
		bindingsResp, err := client.ListCertificateBindings(ctx, cloudcertificates.ListCertificateBindingsRequest{
			CertificateID: certificateID,
			PageSize:      pageSize,
			Page:          page,
		})
		if err != nil {
			return nil, err
		}

		allBindings = append(allBindings, bindingsResp.Bindings...)
		if bindingsResp.Links.Next == nil {
			break
		}
		page++
	}
	return allBindings, nil
}

func (m *certificateDataSourceModel) convertCertificateToModel(ctx context.Context, certificate cloudcertificates.GetCertificateResponse) diag.Diagnostics {
	m.AccountID = types.StringValue(certificate.Certificate.AccountID)
	m.ContractID = types.StringValue(certificate.Certificate.ContractID)
	m.CertificateName = types.StringValue(certificate.Certificate.CertificateName)
	m.CertificateStatus = types.StringValue(certificate.Certificate.CertificateStatus)
	m.CertificateType = types.StringValue(certificate.Certificate.CertificateType)
	m.CreatedDate = date.TimeRFC3339NanoValue(certificate.Certificate.CreatedDate)
	m.CreatedBy = types.StringValue(certificate.Certificate.CreatedBy)
	m.CSRExpirationDate = date.TimeRFC3339Value(certificate.Certificate.CSRExpirationDate)
	m.CSRPEM = types.StringValue(certificate.Certificate.CSRPEM)
	m.KeyType = types.StringValue(string(certificate.Certificate.KeyType))
	m.KeySize = types.StringValue(string(certificate.Certificate.KeySize))
	m.SecureNetwork = types.StringValue(certificate.Certificate.SecureNetwork)
	m.SignedCertificatePEM = types.StringPointerValue(certificate.Certificate.SignedCertificatePEM)
	m.SignedCertificateIssuer = types.StringPointerValue(certificate.Certificate.SignedCertificateIssuer)
	m.SignedCertificateNotValidBeforeDate = date.TimeRFC3339NanoPointerValue(certificate.Certificate.SignedCertificateNotValidBeforeDate)
	m.SignedCertificateNotValidAfterDate = date.TimeRFC3339NanoPointerValue(certificate.Certificate.SignedCertificateNotValidAfterDate)
	m.SignedCertificateSerialNumber = types.StringPointerValue(certificate.Certificate.SignedCertificateSerialNumber)
	m.SignedCertificateSHA256Fingerprint = types.StringPointerValue(certificate.Certificate.SignedCertificateSHA256Fingerprint)
	m.TrustChainPEM = types.StringPointerValue(certificate.Certificate.TrustChainPEM)
	m.ModifiedDate = date.TimeRFC3339NanoValue(certificate.Certificate.ModifiedDate)
	m.ModifiedBy = types.StringValue(certificate.Certificate.ModifiedBy)

	if certificate.Certificate.Subject != nil {
		m.Subject = &subjectModel{
			CommonName:   types.StringValue(certificate.Certificate.Subject.CommonName),
			Organization: types.StringValue(certificate.Certificate.Subject.Organization),
			Country:      types.StringValue(certificate.Certificate.Subject.Country),
			State:        types.StringValue(certificate.Certificate.Subject.State),
			Locality:     types.StringValue(certificate.Certificate.Subject.Locality),
		}
	}

	var diags diag.Diagnostics
	sans, sanDiags := types.SetValueFrom(ctx, types.StringType, certificate.Certificate.SANs)
	diags.Append(sanDiags...)
	if !diags.HasError() {
		m.SANs = sans
	}
	return diags
}
