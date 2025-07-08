package mtlstruststore

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &caSetCertificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &caSetCertificatesDataSource{}
)

type (
	caSetCertificatesDataSource struct {
		meta meta.Meta
	}

	caSetCertificatesDataSourceModel struct {
		ID                    types.String       `tfsdk:"id"`
		Name                  types.String       `tfsdk:"name"`
		Version               types.Int64        `tfsdk:"version"`
		Certificates          []certificateModel `tfsdk:"certificates"`
		Expired               types.Bool         `tfsdk:"expired"`
		ExpiryThresholdInDays types.Int64        `tfsdk:"expiry_threshold_in_days"`
	}
)

// NewCASetCertificatesDataSource returns a new mtls truststore ca set certificates data source.
func NewCASetCertificatesDataSource() datasource.DataSource {
	return &caSetCertificatesDataSource{}
}

// Metadata configures data source's meta information.
func (d *caSetCertificatesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set_certificates"
}

func (d *caSetCertificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *caSetCertificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve CA certificates for a specific MTLS Truststore CA Set version.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifies each CA set.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the CA set.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S`), "must not be empty or only whitespace"),
				},
			},
			"version": schema.Int64Attribute{
				Description: "Version identifier of the CA Set. If not provided, the latest version is used.",
				Optional:    true,
				Computed:    true,
			},
			"expired": schema.BoolAttribute{
				Description: "When true, returns certificates that expired within the past N days, where N is from the `expiry_threshold_in_days` (if provided). If `expiry_threshold_in_days` is not set, all expired certificates are returned.",
				Optional:    true,
			},
			"expiry_threshold_in_days": schema.Int64Attribute{
				Description: "When provided it filters certificates that will expire within the specified number of days. If `expired` is also set, it returns certificates that expired within the past specified number of days.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"certificates": schema.ListNestedAttribute{
				Description: "The CA certificates.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"certificate_pem": schema.StringAttribute{
							Description: "The certificate in PEM format (Base64 ASCII encoded).",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the CA certificate.",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "The user who created this CA certificate.",
							Computed:    true,
						},
						"created_date": schema.StringAttribute{
							Description: "When the CA certificate was created.",
							Computed:    true,
						},
						"end_date": schema.StringAttribute{
							Description: "The ISO 8601 formatted expiration date of the certificate.",
							Computed:    true,
						},
						"fingerprint": schema.StringAttribute{
							Description: "The fingerprint of the certificate.",
							Computed:    true,
						},
						"issuer": schema.StringAttribute{
							Description: "The certificate's issuer.",
							Computed:    true,
						},
						"serial_number": schema.StringAttribute{
							Description: "The unique serial number of the certificate.",
							Computed:    true,
						},
						"signature_algorithm": schema.StringAttribute{
							Description: "The signature algorithm of the CA certificate.",
							Computed:    true,
						},
						"start_date": schema.StringAttribute{
							Description: "The start date of the certificate.",
							Computed:    true,
						},
						"subject": schema.StringAttribute{
							Description: "The subject field of the certificate.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *caSetCertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS TrustStore CA Set Certificates DataSource Read")

	var data caSetCertificatesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client = Client(d.meta)

	if err := data.resolveDefaults(ctx, client); err != nil {
		resp.Diagnostics.AddError("Resolving CA set inputs failed", err.Error())
		return
	}

	certificates, err := data.getCertificates(ctx, client)
	if err != nil {
		resp.Diagnostics.AddError("Read CA set certificates failed", err.Error())
		return
	}

	data.setData(mapCertificatesResponseToModel(certificates))
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (data *caSetCertificatesDataSourceModel) resolveDefaults(ctx context.Context, client mtlstruststore.MTLSTruststore) error {

	if !data.ID.IsNull() && !data.Version.IsNull() {
		return nil
	}

	var caSet *mtlstruststore.CASetResponse

	if !data.Name.IsNull() {
		resp, err := findNotDeletedCASet(ctx, client, data.Name.ValueString())
		if err != nil {
			return fmt.Errorf("failed to find CA set by name '%s': %w", data.Name.ValueString(), err)
		}
		data.ID = types.StringValue(resp.CASetID)
		caSet = &resp
	}

	if data.Version.IsNull() {
		if caSet == nil {
			resp, err := client.GetCASet(ctx, mtlstruststore.GetCASetRequest{
				CASetID: data.ID.ValueString(),
			})
			if err != nil {
				return fmt.Errorf("failed to get CA set '%s': %w", data.ID.ValueString(), err)
			}
			caSet = (*mtlstruststore.CASetResponse)(resp)
		}

		if caSet.LatestVersion == nil {
			return fmt.Errorf("no version provided and CA set has no latest version available")
		}
		data.Version = types.Int64Value(*caSet.LatestVersion)
	}

	return nil
}

func (data *caSetCertificatesDataSourceModel) getCertificates(ctx context.Context, client mtlstruststore.MTLSTruststore) (*mtlstruststore.GetCASetVersionCertificatesResponse, error) {
	var (
		expiryThresholdInDays *int
		certificateStatus     *mtlstruststore.CertificateStatus
	)

	hasExpiry := !data.ExpiryThresholdInDays.IsNull()
	hasExpired := !data.Expired.IsNull() && data.Expired.ValueBool()

	switch {
	case hasExpired && hasExpiry:
		days := int(data.ExpiryThresholdInDays.ValueInt64())
		expiryThresholdInDays = &days
		status := mtlstruststore.ExpiredOrExpiringCert
		certificateStatus = &status
	case hasExpired:
		status := mtlstruststore.ExpiredCert
		certificateStatus = &status
	case hasExpiry:
		days := int(data.ExpiryThresholdInDays.ValueInt64())
		expiryThresholdInDays = &days
		status := mtlstruststore.ExpiringCert
		certificateStatus = &status
	}

	certificates, err := client.GetCASetVersionCertificates(ctx, mtlstruststore.GetCASetVersionCertificatesRequest{
		CASetID:               data.ID.ValueString(),
		Version:               data.Version.ValueInt64(),
		ExpiryThresholdInDays: expiryThresholdInDays,
		CertificateStatus:     certificateStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CA set version certificates: %w", err)
	}

	return certificates, nil
}

func mapCertificatesResponseToModel(certificatesResp *mtlstruststore.GetCASetVersionCertificatesResponse) caSetCertificatesDataSourceModel {
	data := caSetCertificatesDataSourceModel{
		ID:      types.StringValue(certificatesResp.CASetID),
		Name:    types.StringValue(certificatesResp.CASetName),
		Version: types.Int64Value(certificatesResp.Version),
	}

	certs := make([]certificateModel, len(certificatesResp.Certificates))
	for i, cert := range certificatesResp.Certificates {
		certs[i] = certificateModel{
			CertificatePEM:     types.StringValue(cert.CertificatePEM),
			CreatedBy:          types.StringValue(cert.CreatedBy),
			CreatedDate:        types.StringValue(cert.CreatedDate.Format(time.RFC3339)),
			Description:        types.StringPointerValue(cert.Description),
			EndDate:            types.StringValue(cert.EndDate.Format(time.RFC3339)),
			Fingerprint:        types.StringValue(cert.Fingerprint),
			Issuer:             types.StringValue(cert.Issuer),
			SerialNumber:       types.StringValue(cert.SerialNumber),
			SignatureAlgorithm: types.StringValue(cert.SignatureAlgorithm),
			StartDate:          types.StringValue(cert.StartDate.Format(time.RFC3339)),
			Subject:            types.StringValue(cert.Subject),
		}
	}

	data.Certificates = certs
	return data
}

func (data *caSetCertificatesDataSourceModel) setData(certificates caSetCertificatesDataSourceModel) {
	data.ID = certificates.ID
	data.Name = certificates.Name
	data.Version = certificates.Version
	data.Certificates = certificates.Certificates
}
