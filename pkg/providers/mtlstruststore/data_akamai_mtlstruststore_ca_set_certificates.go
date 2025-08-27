package mtlstruststore

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/internal/customtypes"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
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
	_ datasource.DataSource                   = &caSetCertificatesDataSource{}
	_ datasource.DataSourceWithConfigure      = &caSetCertificatesDataSource{}
	_ datasource.DataSourceWithValidateConfig = &caSetCertificatesDataSource{}
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
		IncludeActive         types.Bool         `tfsdk:"include_active"`
		IncludeExpired        types.Bool         `tfsdk:"include_expired"`
		IncludeExpiringInDays types.Int64        `tfsdk:"include_expiring_in_days"`
		IncludeExpiringByDate types.String       `tfsdk:"include_expiring_by_date"`
	}
)

func (d *caSetCertificatesDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	tflog.Debug(ctx, "MTLS TrustStore CA Set Certificates DataSource ValidateConfig")
	var data caSetCertificatesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.IncludeActive.IsUnknown() && !data.IncludeExpired.IsUnknown() &&
		!data.IncludeExpiringInDays.IsUnknown() && !data.IncludeExpiringByDate.IsUnknown() {
		nulledExpiring := data.IncludeExpiringInDays.IsNull() && data.IncludeExpiringByDate.IsNull()
		includeActive := !data.IncludeActive.IsNull() && data.IncludeActive.ValueBool()
		includeExpired := !data.IncludeExpired.IsNull() && data.IncludeExpired.ValueBool()
		if !includeActive && !includeExpired && nulledExpiring {
			resp.Diagnostics.AddError(
				"At least one include option must be set",
				"At least one attribute out of 'include_active', 'include_expired', 'include_expiring_in_days', or 'include_expiring_by_date' must be specified with 'true' value for booleans, or some value for the rest")
			return
		}
	}

	if tf.IsKnown(data.IncludeExpiringByDate) {
		timestamp := data.IncludeExpiringByDate.ValueString()
		parsedTimestamp, err := parseExpiringTimestamp(timestamp)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("include_expiring_by_date"),
				"Invalid expiring timestamp",
				fmt.Sprintf("The provided expiring timestamp '%s' is not a valid RFC3339 or RFC3339Nano formatted date", timestamp),
			)
			return
		}
		if parsedTimestamp.Before(time.Now()) {
			resp.Diagnostics.AddAttributeError(
				path.Root("include_expiring_by_date"),
				"Invalid expiring timestamp",
				fmt.Sprintf("The provided expiring threshold timestamp '%s' cannot be in the past", timestamp),
			)
		}
	}
}

func parseExpiringTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, timestamp)
}

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
				Description: "Identifies each CA set. Either 'id' or 'name' must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the CA set. Either 'id' or 'name' must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S{3,}`), "must not be empty or only whitespace"),
				},
			},
			"version": schema.Int64Attribute{
				Description: "Version identifier of the CA Set. If not provided, the latest version is used.",
				Optional:    true,
				Computed:    true,
			},
			"include_active": schema.BoolAttribute{
				Description: "When true, returns all active (not expired) certificates.",
				Optional:    true,
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.MatchRoot("include_expiring_in_days"),
						path.MatchRoot("include_expiring_by_date")),
				},
			},
			"include_expired": schema.BoolAttribute{
				Description: "When true, returns all expired certificates.",
				Optional:    true,
			},
			"include_expiring_in_days": schema.Int64Attribute{
				Description: "When provided it returns certificates that will expire within the specified number of days.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.ConflictsWith(path.MatchRoot("include_expiring_by_date")),
				},
			},
			"include_expiring_by_date": schema.StringAttribute{
				Description: "When provided it returns certificates that will expire by the specified date.",
				Optional:    true,
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
	requests, err := prepareGetCertificatesRequests(data)
	if err != nil {
		return nil, err
	}
	var certificateResponses []*mtlstruststore.GetCASetVersionCertificatesResponse
	for _, req := range requests {
		certificates, err := client.GetCASetVersionCertificates(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch CA set version certificates: %w", err)
		}
		certificateResponses = append(certificateResponses, certificates)
	}

	return repackCertsToSingleResponse(certificateResponses)
}

func repackCertsToSingleResponse(responses []*mtlstruststore.GetCASetVersionCertificatesResponse) (*mtlstruststore.GetCASetVersionCertificatesResponse, error) {
	var allCerts []mtlstruststore.CertificateResponse

	if len(responses) == 0 {
		return nil, fmt.Errorf("no responses received for CA set version certificates; please report this issue to the provider developers")
	}

	for _, response := range responses {
		allCerts = append(allCerts, response.Certificates...)
	}

	responses[0].Certificates = allCerts
	return responses[0], nil
}

func prepareGetCertificatesRequests(data *caSetCertificatesDataSourceModel) ([]mtlstruststore.GetCASetVersionCertificatesRequest, error) {
	var requests []mtlstruststore.GetCASetVersionCertificatesRequest

	if tf.IsKnown(data.IncludeExpiringInDays) || tf.IsKnown(data.IncludeExpiringByDate) {
		var expiringTimestamp time.Time
		if tf.IsKnown(data.IncludeExpiringByDate) {
			var err error
			expiringTimestamp, err = parseExpiringTimestamp(data.IncludeExpiringByDate.ValueString())
			if err != nil {
				return nil, fmt.Errorf("invalid expiring timestamp: %w", err)
			}
		}

		var expiringInDays *int
		if tf.IsKnown(data.IncludeExpiringInDays) {
			expiringInDays = ptr.To(int(data.IncludeExpiringInDays.ValueInt64()))
		}

		requests = append(requests, mtlstruststore.GetCASetVersionCertificatesRequest{
			CASetID:                  data.ID.ValueString(),
			Version:                  data.Version.ValueInt64(),
			ExpiryThresholdInDays:    expiringInDays,
			ExpiryThresholdTimestamp: expiringTimestamp,
			CertificateStatus:        ptr.To(mtlstruststore.ExpiringCert),
		})
	}

	includeActive := !data.IncludeActive.IsNull() && data.IncludeActive.ValueBool()
	includeExpired := !data.IncludeExpired.IsNull() && data.IncludeExpired.ValueBool()
	if includeExpired || includeActive {
		requests = append(requests, mtlstruststore.GetCASetVersionCertificatesRequest{
			CASetID:           data.ID.ValueString(),
			Version:           data.Version.ValueInt64(),
			CertificateStatus: ptr.To(calculateNotExpiringCertificateStatus(data)),
		})
	}

	return requests, nil
}

func calculateNotExpiringCertificateStatus(data *caSetCertificatesDataSourceModel) mtlstruststore.CertificateStatus {
	includeActive := !data.IncludeActive.IsNull() && data.IncludeActive.ValueBool()
	includeExpired := !data.IncludeExpired.IsNull() && data.IncludeExpired.ValueBool()
	if includeActive && includeExpired {
		return mtlstruststore.ActiveOrExpiredCert
	}
	if includeActive {
		return mtlstruststore.ActiveCert
	}
	// Both includeActive and includeExpired as false are blocked by the schema validation.
	return mtlstruststore.ExpiredCert
}

func mapCertificatesResponseToModel(certificatesResp *mtlstruststore.GetCASetVersionCertificatesResponse) caSetCertificatesDataSourceModel {
	data := caSetCertificatesDataSourceModel{
		ID:      types.StringValue(certificatesResp.CASetID),
		Name:    types.StringValue(certificatesResp.CASetName),
		Version: types.Int64Value(certificatesResp.Version),
	}

	certs := make([]certificateModel, 0, len(certificatesResp.Certificates))
	for _, cert := range certificatesResp.Certificates {
		certs = append(certs, certificateModel{
			CertificatePEM:     customtypes.NewIgnoreTrailingWhitespaceValue(cert.CertificatePEM),
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
		})
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
