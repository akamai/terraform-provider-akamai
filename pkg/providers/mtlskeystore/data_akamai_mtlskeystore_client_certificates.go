package mtlskeystore

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clientCertificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &clientCertificatesDataSource{}
)

type (
	clientCertificatesDataSource struct {
		meta meta.Meta
	}

	clientCertificatesDataSourceModel struct {
		ClientCertificates []clientCertificateModel `tfsdk:"certificates"`
	}

	clientCertificateModel struct {
		CertificateID      types.Int64  `tfsdk:"certificate_id"`
		CertificateName    types.String `tfsdk:"certificate_name"`
		CreatedBy          types.String `tfsdk:"created_by"`
		CreatedDate        types.String `tfsdk:"created_date"`
		Geography          types.String `tfsdk:"geography"`
		KeyAlgorithm       types.String `tfsdk:"key_algorithm"`
		NotificationEmails types.List   `tfsdk:"notification_emails"`
		SecureNetwork      types.String `tfsdk:"secure_network"`
		Signer             types.String `tfsdk:"signer"`
		Subject            types.String `tfsdk:"subject"`
	}
)

// NewClientCertificatesDataSource returns a new mtls keystore client certificates data source.
func NewClientCertificatesDataSource() datasource.DataSource {
	return &clientCertificatesDataSource{}
}

func (d *clientCertificatesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlskeystore_client_certificates"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *clientCertificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *clientCertificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve client certificates under MTLS Keystore.",
		Attributes: map[string]schema.Attribute{
			"certificates": schema.ListNestedAttribute{
				Description: "A list of client certificates under the account.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"certificate_id": schema.Int64Attribute{
							Description: "The unique identifier of the client certificate.",
							Computed:    true,
						},
						"certificate_name": schema.StringAttribute{
							Description: "The name of the client certificate.",
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
						"geography": schema.StringAttribute{
							Description: "Specifies the type of network to deploy the client certificate. Possible values: `CORE`, `RUSSIA_AND_CORE`, or `CHINA_AND_CORE`.",
							Computed:    true,
						},
						"key_algorithm": schema.StringAttribute{
							Description: "Identifies the CA certificate's encryption algorithm. Possible values: `RSA` or `ECDSA`.",
							Computed:    true,
						},
						"notification_emails": schema.ListAttribute{
							Description: "The email addresses to notify for client certificate-related issues.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"secure_network": schema.StringAttribute{
							Description: "Identifies the network deployment type. Possible values: `STANDARD_TLS` or `ENHANCED_TLS`.",
							Computed:    true,
						},
						"signer": schema.StringAttribute{
							Description: "The signing entity of the client certificate. Possible values: `AKAMAI` or `THIRD_PARTY`.",
							Computed:    true,
						},
						"subject": schema.StringAttribute{
							Description: "The CA certificate’s key value details.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *clientCertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS Keystore Client Certificates Data Source Read")

	var data clientCertificatesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := Client(d.meta)

	clientCertificates, err := client.ListClientCertificates(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to List Client Certificates", err.Error())
		return
	}
	if clientCertificates == nil {
		resp.Diagnostics.AddError("No Client Certificates", "Received nil response from client")
		return
	}

	if resp.Diagnostics.Append(data.convertClientCertificatesToModel(ctx, clientCertificates)...); resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *clientCertificatesDataSourceModel) convertClientCertificatesToModel(ctx context.Context, clientCertificates *mtlskeystore.ListClientCertificatesResponse) diag.Diagnostics {
	var model clientCertificateModel
	var diag diag.Diagnostics
	for _, cert := range clientCertificates.Certificates {
		model = clientCertificateModel{
			CertificateID:   types.Int64Value(cert.CertificateID),
			CertificateName: types.StringValue(cert.CertificateName),
			CreatedBy:       types.StringValue(cert.CreatedBy),
			CreatedDate:     date.TimeRFC3339Value(cert.CreatedDate),
			Geography:       types.StringValue(cert.Geography),
			KeyAlgorithm:    types.StringValue(cert.KeyAlgorithm),
			SecureNetwork:   types.StringValue(cert.SecureNetwork),
			Signer:          types.StringValue(cert.Signer),
			Subject:         types.StringValue(cert.Subject),
		}

		model.NotificationEmails, diag = types.ListValueFrom(ctx, types.StringType, cert.NotificationEmails)
		if diag.HasError() {
			return diag
		}
		m.ClientCertificates = append(m.ClientCertificates, model)
	}
	return nil
}
