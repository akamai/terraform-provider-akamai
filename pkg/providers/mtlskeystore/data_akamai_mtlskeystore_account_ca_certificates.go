package mtlskeystore

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &accountCACertificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &accountCACertificatesDataSource{}
)

type (
	accountCACertificatesDataSource struct {
		meta meta.Meta
	}

	accountCACertificatesDataSourceModel struct {
		Status       types.List         `tfsdk:"status"`
		Certificates []certificateModel `tfsdk:"certificates"`
	}

	certificateModel struct {
		AccountID          types.String `tfsdk:"account_id"`
		Certificate        types.String `tfsdk:"certificate"`
		CommonName         types.String `tfsdk:"common_name"`
		CreatedBy          types.String `tfsdk:"created_by"`
		CreatedDate        types.String `tfsdk:"created_date"`
		ExpiryDate         types.String `tfsdk:"expiry_date"`
		ID                 types.Int64  `tfsdk:"id"`
		IssuedDate         types.String `tfsdk:"issued_date"`
		KeyAlgorithm       types.String `tfsdk:"key_algorithm"`
		KeySizeInBytes     types.Int64  `tfsdk:"key_size_in_bytes"`
		QualificationDate  types.String `tfsdk:"qualification_date"`
		SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
		Status             types.String `tfsdk:"status"`
		Subject            types.String `tfsdk:"subject"`
		Version            types.Int64  `tfsdk:"version"`
	}
)

// NewAccountCACertificatesDataSource returns a new mtls keystore account ca certificates data source.
func NewAccountCACertificatesDataSource() datasource.DataSource {
	return &accountCACertificatesDataSource{}
}

// Metadata configures data source's meta information.
func (d *accountCACertificatesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlskeystore_account_ca_certificates"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *accountCACertificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *accountCACertificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve CA certificates under MTLS Keystore account.",
		Attributes: map[string]schema.Attribute{
			"status": schema.ListAttribute{
				Description: "A list of statuses separated by commas used to filter account CA certificates. Only certificates matching the provided statuses are included in the results." +
					" Possible values: QUALIFYING, CURRENT, PREVIOUS, or EXPIRED.",
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.OneOf(
						string(mtlskeystore.CertificateStatusExpired),
						string(mtlskeystore.CertificateStatusCurrent),
						string(mtlskeystore.CertificateStatusQualifying),
						string(mtlskeystore.CertificateStatusPrevious),
					)),
				},
			},
			"certificates": schema.ListNestedAttribute{
				Description: "A list of account CA certificates.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Description: "The account the CA certificate is under.",
							Computed:    true,
						},
						"certificate": schema.StringAttribute{
							Description: "The certificate block of the CA certificate.",
							Computed:    true,
						},
						"common_name": schema.StringAttribute{
							Description: "The common name of the CA certificate.",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "The user who created the CA certificate.",
							Computed:    true,
						},
						"created_date": schema.StringAttribute{
							Description: "An ISO 8601 timestamp indicating the CA certificate's creation.",
							Computed:    true,
						},
						"expiry_date": schema.StringAttribute{
							Description: "An ISO 8601 timestamp indicating when the CA certificate expires.",
							Computed:    true,
						},
						"id": schema.Int64Attribute{
							Description: "The unique identifier of the CA certificate.",
							Computed:    true,
						},
						"issued_date": schema.StringAttribute{
							Description: "An ISO 8601 timestamp indicating the CA certificate's availability.",
							Computed:    true,
						},
						"key_algorithm": schema.StringAttribute{
							Description: "Identifies the CA certificate's encryption algorithm. Possible values: `RSA` or `ECDSA`.",
							Computed:    true,
						},
						"key_size_in_bytes": schema.Int64Attribute{
							Description: "The private key length of the CA certificate.",
							Computed:    true,
						},
						"qualification_date": schema.StringAttribute{
							Description: "An ISO 8601 timestamp indicating when the CA certificate's status moved from QUALIFYING to CURRENT.",
							Computed:    true,
						},
						"signature_algorithm": schema.StringAttribute{
							Description: "Specifies the algorithm that secures the data exchange between the edge server and origin.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The status of the CA certificate. Possible values: QUALIFYING, CURRENT, PREVIOUS, or EXPIRED.",
							Computed:    true,
						},
						"subject": schema.StringAttribute{
							Description: "The public key's entity stored in the CA certificate's subject public key field.",
							Computed:    true,
						},
						"version": schema.Int64Attribute{
							Description: "The version of the CA certificate.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *accountCACertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS Keystore Account CA Certificates Data Source Read")

	var data accountCACertificatesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := Client(d.meta)

	var status []mtlskeystore.CertificateStatus
	diags := data.Status.ElementsAs(ctx, &status, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	caCerts, err := client.ListAccountCACertificates(ctx, mtlskeystore.ListAccountCACertificatesRequest{
		Status: status,
	})
	if err != nil {
		resp.Diagnostics.AddError("Read Account CA Certificates failed", err.Error())
		return
	}

	data.convertResponseToModel(caCerts)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (a *accountCACertificatesDataSourceModel) convertResponseToModel(caCerts *mtlskeystore.ListAccountCACertificatesResponse) {
	for _, cert := range caCerts.Certificates {
		certificates := certificateModel{
			AccountID:          types.StringValue(cert.AccountID),
			Certificate:        types.StringValue(cert.Certificate),
			CommonName:         types.StringValue(cert.CommonName),
			CreatedBy:          types.StringValue(cert.CreatedBy),
			CreatedDate:        date.TimeRFC3339Value(cert.CreatedDate),
			ExpiryDate:         date.TimeRFC3339Value(cert.ExpiryDate),
			ID:                 types.Int64Value(cert.ID),
			IssuedDate:         date.TimeRFC3339Value(cert.IssuedDate),
			KeyAlgorithm:       types.StringValue(cert.KeyAlgorithm),
			KeySizeInBytes:     types.Int64Value(cert.KeySizeInBytes),
			QualificationDate:  date.TimeRFC3339PointerValue(cert.QualificationDate),
			SignatureAlgorithm: types.StringValue(cert.SignatureAlgorithm),
			Status:             types.StringValue(cert.Status),
			Subject:            types.StringValue(cert.Subject),
			Version:            types.Int64Value(cert.Version),
		}

		a.Certificates = append(a.Certificates, certificates)
	}
}
