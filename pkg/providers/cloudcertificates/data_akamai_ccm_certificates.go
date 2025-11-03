package cloudcertificates

import (
	"context"
	"fmt"
	"regexp"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ccm"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &certificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesDataSource{}
)

type (
	certificatesDataSource struct {
		meta meta.Meta
	}

	certificatesDataSourceModel struct {
		ContractID                  types.String       `tfsdk:"contract_id"`
		GroupID                     types.String       `tfsdk:"group_id"`
		CertificateStatus           types.Set          `tfsdk:"certificate_status"`
		ExpiringInDays              types.Int64        `tfsdk:"expiring_in_days"`
		Domain                      types.String       `tfsdk:"domain"`
		CertificateName             types.String       `tfsdk:"certificate_name"`
		KeyType                     types.String       `tfsdk:"key_type"`
		Issuer                      types.String       `tfsdk:"issuer"`
		IncludeCertificateMaterials types.Bool         `tfsdk:"include_certificate_materials"`
		Sort                        types.String       `tfsdk:"sort"`
		Certificates                []certificateModel `tfsdk:"certificates"`
	}

	certificateModel struct {
		CertificateID                       types.String  `tfsdk:"certificate_id"`
		CertificateName                     types.String  `tfsdk:"certificate_name"`
		SANs                                types.Set     `tfsdk:"sans"`
		Subject                             *subjectModel `tfsdk:"subject"`
		CertificateType                     types.String  `tfsdk:"certificate_type"`
		KeyType                             types.String  `tfsdk:"key_type"`
		KeySize                             types.String  `tfsdk:"key_size"`
		SecureNetwork                       types.String  `tfsdk:"secure_network"`
		ContractID                          types.String  `tfsdk:"contract_id"`
		AccountID                           types.String  `tfsdk:"account_id"`
		CreatedDate                         types.String  `tfsdk:"created_date"`
		CreatedBy                           types.String  `tfsdk:"created_by"`
		ModifiedDate                        types.String  `tfsdk:"modified_date"`
		ModifiedBy                          types.String  `tfsdk:"modified_by"`
		CertificateStatus                   types.String  `tfsdk:"certificate_status"`
		CSRPEM                              types.String  `tfsdk:"csr_pem"`
		CSRExpirationDate                   types.String  `tfsdk:"csr_expiration_date"`
		SignedCertificatePEM                types.String  `tfsdk:"signed_certificate_pem"`
		SignedCertificateNotValidAfterDate  types.String  `tfsdk:"signed_certificate_not_valid_after_date"`
		SignedCertificateNotValidBeforeDate types.String  `tfsdk:"signed_certificate_not_valid_before_date"`
		SignedCertificateSerialNumber       types.String  `tfsdk:"signed_certificate_serial_number"`
		SignedCertificateSHA256Fingerprint  types.String  `tfsdk:"signed_certificate_sha256_fingerprint"`
		SignedCertificateIssuer             types.String  `tfsdk:"signed_certificate_issuer"`
		TrustChainPEM                       types.String  `tfsdk:"trust_chain_pem"`
	}
)

// ccmCertificatesPageSize defines the maximum number of items to be retrieved per page from the CCM API.
const ccmCertificatesPageSize int64 = 100

// NewCertificatesDataSource returns a new CloudCertificates Certificates data source.
func NewCertificatesDataSource() datasource.DataSource {
	return &certificatesDataSource{}
}

// Metadata configures data source's meta information.
func (d *certificatesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_certificates"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *certificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *certificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for retrieving all certificates available to you in your account.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Description: "Filter by contract identifier and only return CCM certificates associated with that contract.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"group_id": schema.StringAttribute{
				Description: "Filter by group identifier and only return CCM certificates associated with that group.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"certificate_status": schema.SetAttribute{
				Description: "Filter certificates by status. Valid values are `ACTIVE`, `READY_FOR_USE`, `CSR_READY`.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.OneOf("ACTIVE", "READY_FOR_USE", "CSR_READY")),
				},
			},
			"expiring_in_days": schema.Int64Attribute{
				Description: "Filter certificates that are expiring in the specified number of days. " +
					"A value of 0 returns only expired certificates.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"domain": schema.StringAttribute{
				Description: "Filter certificates by domain in the certificate's SANs or subject CN. Supports partial matches. " +
					"Matches are case-insensitive, and support wildcards.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"certificate_name": schema.StringAttribute{
				Description: "Filter certificates by name. Supports partial matches.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"key_type": schema.StringAttribute{
				Description: "Filter certificates by key type. Valid values are `RSA` and `ECDSA`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("RSA", "ECDSA"),
				},
			},
			"issuer": schema.StringAttribute{
				Description: "Filter certificates by issuer. Supports partial matches.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"include_certificate_materials": schema.BoolAttribute{
				Description: "Include certificate materials in the response for each certificate (CSR, signed certificate, and trust chain).",
				Optional:    true,
			},
			"sort": schema.StringAttribute{
				Description: "Sorts results by one or more comma-separated certificate fields, in order from left to right. " +
					"Valid values are `certificateName`, `createdDate`, `expirationDate`, and `modifiedDate`. " +
					"Prefix a field with a plus sign (+) for ascending order or a minus sign (-) for descending order. " +
					"By default, results are sorted by `modifiedDate` in descending order.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^`+ccm.SortFieldPat+`(,`+ccm.SortFieldPat+
						`)*$`), "must be a comma-separated list of fields with optional '+' or '-' prefix. "+
						"Valid fields are \"certificateName\", \"createdDate\", \"expirationDate\", and \"modifiedDate\""),
				},
			},
			"certificates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"certificate_id": schema.StringAttribute{
							Description: "The unique identifier for the certificate.",
							Computed:    true,
						},
						"certificate_name": schema.StringAttribute{
							Description: "The name of the certificate.",
							Computed:    true,
						},
						"sans": schema.SetAttribute{
							Description: "The list of SAN (Subject Alternative Name) domains included in the certificate.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"subject": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Subject fields as defined in X.509 certificates (RFC 5280).",
							Attributes: map[string]schema.Attribute{
								"common_name": schema.StringAttribute{
									Description: "Fully qualified domain name (FQDN) or other name associated with the subject.",
									Computed:    true,
								},
								"organization": schema.StringAttribute{
									Description: "Legal name of the organization.",
									Computed:    true,
								},
								"country": schema.StringAttribute{
									Description: "Two-letter ISO 3166 country code.",
									Computed:    true,
								},
								"state": schema.StringAttribute{
									Description: "Full name of the state or province.",
									Computed:    true,
								},
								"locality": schema.StringAttribute{
									Description: "City or locality name.",
									Computed:    true,
								},
							},
						},
						"certificate_type": schema.StringAttribute{
							Description: "The type of the certificate.",
							Computed:    true,
						},
						"key_type": schema.StringAttribute{
							Description: "The key type of the algorithm used in the certificate signing request (CSR).",
							Computed:    true,
						},
						"key_size": schema.StringAttribute{
							Description: "Size of the key used in the certificate signing request (CSR) in bits.",
							Computed:    true,
						},
						"secure_network": schema.StringAttribute{
							Description: "The secure network associated with the certificate.",
							Computed:    true,
						},
						"contract_id": schema.StringAttribute{
							Description: "The contract identifier associated with the certificate.",
							Computed:    true,
						},
						"account_id": schema.StringAttribute{
							Description: "The account identifier associated with the certificate.",
							Computed:    true,
						},
						"created_date": schema.StringAttribute{
							Description: "The date when the certificate was created.",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "The user who created the certificate.",
							Computed:    true,
						},
						"modified_date": schema.StringAttribute{
							Description: "The date when the certificate was last modified.",
							Computed:    true,
						},
						"modified_by": schema.StringAttribute{
							Description: "The user who last modified the certificate.",
							Computed:    true,
						},
						"certificate_status": schema.StringAttribute{
							Description: "The status of the certificate.",
							Computed:    true,
						},
						"csr_pem": schema.StringAttribute{
							Description: "PEM-encoded certificate signing request (CSR) generated by Akamai for your selected key type.",
							Computed:    true,
						},
						"csr_expiration_date": schema.StringAttribute{
							Description: "The expiration date of the CSR.",
							Computed:    true,
						},
						"signed_certificate_pem": schema.StringAttribute{
							Description: "PEM-encoded signed certificate you uploaded for your selected key type.",
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
						"signed_certificate_serial_number": schema.StringAttribute{
							Description: "Signed certificate serial number in hex format.",
							Computed:    true,
						},
						"signed_certificate_sha256_fingerprint": schema.StringAttribute{
							Description: "The SHA256 fingerprint of the signed certificate.",
							Computed:    true,
						},
						"signed_certificate_issuer": schema.StringAttribute{
							Description: "The issuer of the signed certificate.",
							Computed:    true,
						},
						"trust_chain_pem": schema.StringAttribute{
							Description: "The trust chain PEM content uploaded by end user.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *certificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CCM Certificates DataSource Read")

	var data certificatesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)

	request := ccm.ListCertificatesRequest{}

	if !data.ContractID.IsNull() {
		request.ContractID = data.ContractID.ValueString()
	}
	if !data.GroupID.IsNull() {
		request.GroupID = data.GroupID.ValueString()
	}
	if !data.CertificateStatus.IsNull() {
		var statusList []ccm.CertificateStatus
		data.CertificateStatus.ElementsAs(ctx, &statusList, false)
		request.CertificateStatus = statusList
	}
	if !data.Domain.IsNull() {
		request.Domain = data.Domain.ValueString()
	}
	if !data.CertificateName.IsNull() {
		request.CertificateName = data.CertificateName.ValueString()
	}
	if !data.KeyType.IsNull() {
		request.KeyType = ccm.CryptographicAlgorithm(data.KeyType.ValueString())
	}
	if !data.Issuer.IsNull() {
		request.Issuer = data.Issuer.ValueString()
	}
	if !data.IncludeCertificateMaterials.IsNull() {
		request.IncludeCertificateMaterials = data.IncludeCertificateMaterials.ValueBool()
	}
	if !data.Sort.IsNull() {
		request.Sort = data.Sort.ValueString()
	}
	request.ExpiringInDays = data.ExpiringInDays.ValueInt64Pointer()

	cert, err := getAllCertificates(ctx, client, request)
	if err != nil {
		resp.Diagnostics.AddError("Read Certificates failed", err.Error())
		return
	}

	models, diags := certificatesToModels(ctx, cert)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.Certificates = models

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getAllCertificates(ctx context.Context, client ccm.CCM, request ccm.ListCertificatesRequest) (*ccm.ListCertificatesResponse, error) {
	var allCertificates ccm.ListCertificatesResponse

	request.PageSize = ccmCertificatesPageSize
	request.Page = 1

	for {
		certificatesResponse, err := client.ListCertificates(ctx, request)
		if err != nil {
			return nil, err
		}

		allCertificates.Certificates = append(allCertificates.Certificates, certificatesResponse.Certificates...)

		if certificatesResponse.Links.Next == nil {
			break
		}
		request.Page++
	}

	return &allCertificates, nil
}

func certificatesToModels(ctx context.Context, certificates *ccm.ListCertificatesResponse) ([]certificateModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	var certs []certificateModel
	for _, cert := range certificates.Certificates {
		c := certificateModel{
			CertificateID:                       types.StringValue(cert.CertificateID),
			CertificateName:                     types.StringValue(cert.CertificateName),
			CertificateType:                     types.StringValue(cert.CertificateType),
			KeyType:                             types.StringValue(string(cert.KeyType)),
			KeySize:                             types.StringValue(string(cert.KeySize)),
			SecureNetwork:                       types.StringValue(cert.SecureNetwork),
			ContractID:                          types.StringValue(cert.ContractID),
			AccountID:                           types.StringValue(cert.AccountID),
			CreatedDate:                         date.TimeRFC3339NanoValue(cert.CreatedDate),
			CreatedBy:                           types.StringValue(cert.CreatedBy),
			ModifiedDate:                        date.TimeRFC3339NanoValue(cert.ModifiedDate),
			ModifiedBy:                          types.StringValue(cert.ModifiedBy),
			CertificateStatus:                   types.StringValue(cert.CertificateStatus),
			CSRPEM:                              types.StringPointerValue(cert.CSRPEM),
			CSRExpirationDate:                   date.TimeRFC3339NanoValue(cert.CSRExpirationDate),
			SignedCertificatePEM:                types.StringPointerValue(cert.SignedCertificatePEM),
			SignedCertificateNotValidAfterDate:  date.TimeRFC3339NanoPointerValue(cert.SignedCertificateNotValidAfterDate),
			SignedCertificateNotValidBeforeDate: date.TimeRFC3339NanoPointerValue(cert.SignedCertificateNotValidBeforeDate),
			SignedCertificateSerialNumber:       types.StringPointerValue(cert.SignedCertificateSerialNumber),
			SignedCertificateSHA256Fingerprint:  types.StringPointerValue(cert.SignedCertificateSHA256Fingerprint),
			SignedCertificateIssuer:             types.StringPointerValue(cert.SignedCertificateIssuer),
			TrustChainPEM:                       types.StringPointerValue(cert.TrustChainPEM),
		}

		if cert.Subject != nil {
			c.Subject = &subjectModel{
				CommonName:   types.StringValue(cert.Subject.CommonName),
				Organization: types.StringValue(cert.Subject.Organization),
				Country:      types.StringValue(cert.Subject.Country),
				State:        types.StringValue(cert.Subject.State),
				Locality:     types.StringValue(cert.Subject.Locality),
			}
		}

		sans, dd := types.SetValueFrom(ctx, types.StringType, cert.SANs)
		diags.Append(dd...)
		if diags.HasError() {
			return nil, diags
		}
		c.SANs = sans

		certs = append(certs, c)
	}

	return certs, diags
}
