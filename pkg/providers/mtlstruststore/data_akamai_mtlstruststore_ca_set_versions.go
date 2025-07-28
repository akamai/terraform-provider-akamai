package mtlstruststore

import (
	"context"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/internal/customtypes"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &caSetVersionsDataSource{}
	_ datasource.DataSourceWithConfigure = &caSetVersionsDataSource{}
)

type (
	caSetVersionsDataSource struct {
		meta meta.Meta
	}

	caSetVersionsDataSourceModel struct {
		ID                  types.String   `tfsdk:"id"`
		Name                types.String   `tfsdk:"name"`
		ActiveVersionsOnly  types.Bool     `tfsdk:"active_versions_only"`
		IncludeCertificates types.Bool     `tfsdk:"include_certificates"`
		Versions            []versionModel `tfsdk:"versions"`
	}

	versionModel struct {
		Version            types.Int64        `tfsdk:"version"`
		AllowInsecureSHA1  types.Bool         `tfsdk:"allow_insecure_sha1"`
		VersionDescription types.String       `tfsdk:"version_description"`
		CreatedBy          types.String       `tfsdk:"created_by"`
		CreatedDate        types.String       `tfsdk:"created_date"`
		ModifiedBy         types.String       `tfsdk:"modified_by"`
		ModifiedDate       types.String       `tfsdk:"modified_date"`
		ProductionStatus   types.String       `tfsdk:"production_status"`
		StagingStatus      types.String       `tfsdk:"staging_status"`
		Certificates       []certificateModel `tfsdk:"certificates"`
	}
)

// NewCASetVersionsDataSource returns a new mtls truststore ca set versions data source.
func NewCASetVersionsDataSource() datasource.DataSource {
	return &caSetVersionsDataSource{}
}

// Metadata configures data source's meta information.
func (d *caSetVersionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set_versions"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *caSetVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *caSetVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve versions for a specific MTLS Truststore CA Set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifies each CA set. Either `id` or `name` must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the CA set. Either `id` or `name` must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"include_certificates": schema.BoolAttribute{
				Description: "If this option is set to true, the response includes certificates belonging to the version. The default is true.",
				Optional:    true,
			},
			"active_versions_only": schema.BoolAttribute{
				Description: "If true, only the active versions of the CA set will be returned. The default is false.",
				Optional:    true,
			},
			"versions": schema.ListNestedAttribute{
				Description: "List of CA set versions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"version": schema.Int64Attribute{
							Description: "Version identifier on which to perform the desired operation.",
							Computed:    true,
						},
						"allow_insecure_sha1": schema.BoolAttribute{
							Description: "By default, all certificates in the version need a signature algorithm of SHA-256 or better. Enabling this allows certificates with SHA-1 signatures.",
							Computed:    true,
						},
						"version_description": schema.StringAttribute{
							Description: "Any additional description you can provide while creating or updating the CA set version.",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "The user who created the CA set version.",
							Computed:    true,
						},
						"created_date": schema.StringAttribute{
							Description: "When the CA set version was created.",
							Computed:    true,
						},
						"modified_by": schema.StringAttribute{
							Description: "The user who last modified the CA set version.",
							Computed:    true,
						},
						"modified_date": schema.StringAttribute{
							Description: "When the CA set version was last modified.",
							Computed:    true,
						},
						"production_status": schema.StringAttribute{
							Description: "The CA set version's status on the production network, either `ACTIVE` or `INACTIVE`.",
							Computed:    true,
						},
						"staging_status": schema.StringAttribute{
							Description: "The CA set version's status on the staging network, either `ACTIVE` or `INACTIVE`.",
							Computed:    true,
						},
						"certificates": schema.ListNestedAttribute{
							Description: "List of certificate objects in the version, with each element corresponding to one root or intermediate certificate.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"certificate_pem": schema.StringAttribute{
										Description: "The certificate in PEM format, as found in a Base64 ASCII encoded file.",
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
										Description: "The certificate's ISO 8601 formatted expiration date.",
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
										Description: "The certificate's subject field.",
										Computed:    true,
									},
									"description": schema.StringAttribute{
										Description: "Description for the certificate.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *caSetVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS TrustStore CA Set Versions DataSource Read")

	var data caSetVersionsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client = Client(d.meta)

	if !data.Name.IsNull() {
		tflog.Debug(ctx, "'name' provided, attempting to find CA set ID")
		setID, err := findCASetID(ctx, client, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Read CA set versions failed", err.Error())
			return
		}
		data.ID = types.StringValue(setID)
	}

	if data.IncludeCertificates.IsNull() {
		data.IncludeCertificates = types.BoolValue(true)
	}
	if data.ActiveVersionsOnly.IsNull() {
		data.ActiveVersionsOnly = types.BoolValue(false)
	}

	versions, err := client.ListCASetVersions(ctx, mtlstruststore.ListCASetVersionsRequest{
		CASetID:             data.ID.ValueString(),
		IncludeCertificates: data.IncludeCertificates.ValueBool(),
		ActiveVersionsOnly:  data.ActiveVersionsOnly.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read CA set versions failed", err.Error())
		return
	}

	modelData := convertCASetVersionsDataToModel(*versions)
	data.setData(modelData)

	if len(versions.Versions) == 0 {
		if data.Name.IsNull() {
			caSet, err := client.GetCASet(ctx, mtlstruststore.GetCASetRequest{
				CASetID: data.ID.ValueString(),
			})
			if err != nil {
				resp.Diagnostics.AddError("Read CA set versions failed", err.Error())
				return
			}
			data.Name = types.StringValue(caSet.CASetName)
		}
	} else {
		data.Name, data.ID = extractCASetNameAndID(*versions)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func convertCASetVersionsDataToModel(versions mtlstruststore.ListCASetVersionsResponse) caSetVersionsDataSourceModel {
	var model caSetVersionsDataSourceModel
	for _, version := range versions.Versions {
		var certModel []certificateModel
		for _, cert := range version.Certificates {
			certModel = append(certModel, certificateModel{
				CertificatePEM:     customtypes.NewIgnoreTrailingWhitespaceValue(cert.CertificatePEM),
				CreatedBy:          types.StringValue(cert.CreatedBy),
				CreatedDate:        types.StringValue(cert.CreatedDate.Format(time.RFC3339Nano)),
				EndDate:            types.StringValue(cert.EndDate.Format(time.RFC3339Nano)),
				Fingerprint:        types.StringValue(cert.Fingerprint),
				Issuer:             types.StringValue(cert.Issuer),
				SerialNumber:       types.StringValue(cert.SerialNumber),
				SignatureAlgorithm: types.StringValue(cert.SignatureAlgorithm),
				StartDate:          types.StringValue(cert.StartDate.Format(time.RFC3339Nano)),
				Subject:            types.StringValue(cert.Subject),
				Description:        types.StringPointerValue(cert.Description),
			})
		}

		versionModel := versionModel{
			Version:            types.Int64Value(version.Version),
			AllowInsecureSHA1:  types.BoolValue(version.AllowInsecureSHA1),
			VersionDescription: types.StringPointerValue(version.Description),
			CreatedBy:          types.StringValue(version.CreatedBy),
			CreatedDate:        types.StringValue(version.CreatedDate.Format(time.RFC3339Nano)),
			ProductionStatus:   types.StringValue(version.ProductionStatus),
			StagingStatus:      types.StringValue(version.StagingStatus),
			ModifiedBy:         types.StringPointerValue(version.ModifiedBy),
			Certificates:       certModel,
		}

		if version.ModifiedDate != nil {
			versionModel.ModifiedDate = types.StringValue(version.ModifiedDate.Format(time.RFC3339Nano))
		}

		model.Versions = append(model.Versions, versionModel)
	}

	return model
}

func (m *caSetVersionsDataSourceModel) setData(data caSetVersionsDataSourceModel) {
	m.Versions = make([]versionModel, len(data.Versions))

	for i, version := range data.Versions {
		m.Versions[i] = versionModel{
			Version:            version.Version,
			AllowInsecureSHA1:  version.AllowInsecureSHA1,
			VersionDescription: version.VersionDescription,
			CreatedBy:          version.CreatedBy,
			CreatedDate:        version.CreatedDate,
			ModifiedBy:         version.ModifiedBy,
			ModifiedDate:       version.ModifiedDate,
			ProductionStatus:   version.ProductionStatus,
			StagingStatus:      version.StagingStatus,
			Certificates:       make([]certificateModel, len(version.Certificates)),
		}

		for j, cert := range version.Certificates {
			m.Versions[i].Certificates[j] = certificateModel{
				CertificatePEM:     cert.CertificatePEM,
				CreatedBy:          cert.CreatedBy,
				CreatedDate:        cert.CreatedDate,
				EndDate:            cert.EndDate,
				Fingerprint:        cert.Fingerprint,
				Issuer:             cert.Issuer,
				SerialNumber:       cert.SerialNumber,
				SignatureAlgorithm: cert.SignatureAlgorithm,
				StartDate:          cert.StartDate,
				Subject:            cert.Subject,
				Description:        cert.Description,
			}
		}
	}
}

func extractCASetNameAndID(versions mtlstruststore.ListCASetVersionsResponse) (types.String, types.String) {
	if len(versions.Versions) == 0 {
		return types.StringNull(), types.StringNull()
	}

	return types.StringValue(versions.Versions[0].CASetName), types.StringValue(versions.Versions[0].CASetID)
}
