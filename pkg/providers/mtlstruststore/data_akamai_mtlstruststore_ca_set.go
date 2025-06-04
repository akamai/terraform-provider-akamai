package mtlstruststore

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &caSetDataSource{}
	_ datasource.DataSourceWithConfigure = &caSetDataSource{}
)

type (
	caSetDataSource struct {
		meta meta.Meta
	}

	caSetDataSourceModel struct {
		ID                  types.String       `tfsdk:"id"`
		Name                types.String       `tfsdk:"name"`
		Version             types.Int64        `tfsdk:"version"`
		Description         types.String       `tfsdk:"description"`
		AccountID           types.String       `tfsdk:"account_id"`
		CreatedBy           types.String       `tfsdk:"created_by"`
		CreatedDate         types.String       `tfsdk:"created_date"`
		DeletedBy           types.String       `tfsdk:"deleted_by"`
		DeletedDate         types.String       `tfsdk:"deleted_date"`
		VersionCreatedBy    types.String       `tfsdk:"version_created_by"`
		VersionCreatedDate  types.String       `tfsdk:"version_created_date"`
		VersionModifiedBy   types.String       `tfsdk:"version_modified_by"`
		VersionModifiedDate types.String       `tfsdk:"version_modified_date"`
		AllowInsecureSHA1   types.Bool         `tfsdk:"allow_insecure_sha1"`
		VersionDescription  types.String       `tfsdk:"version_description"`
		StagingVersion      types.Int64        `tfsdk:"staging_version"`
		ProductionVersion   types.Int64        `tfsdk:"production_version"`
		Certificates        []certificateModel `tfsdk:"certificates"`
	}

	certificateModel struct {
		CertificatePEM     types.String `tfsdk:"certificate_pem"`
		Description        types.String `tfsdk:"description"`
		CreatedBy          types.String `tfsdk:"created_by"`
		CreatedDate        types.String `tfsdk:"created_date"`
		StartDate          types.String `tfsdk:"start_date"`
		EndDate            types.String `tfsdk:"end_date"`
		Fingerprint        types.String `tfsdk:"fingerprint"`
		Issuer             types.String `tfsdk:"issuer"`
		SerialNumber       types.String `tfsdk:"serial_number"`
		SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
		Subject            types.String `tfsdk:"subject"`
	}
)

// NewCASetDataSource returns a new mtls truststore ca set data source.
func NewCASetDataSource() datasource.DataSource {
	return &caSetDataSource{}
}

// Metadata configures data source's meta information.
func (d *caSetDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *caSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *caSetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details of a specific MTLS Truststore CA Set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the CA set.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the CA set.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"version": schema.Int64Attribute{
				Description: "When provided, it defines which version's details to provide. If not, the latest version is used",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Any additional comments you can add to the CA set.",
				Computed:    true,
			},
			"account_id": schema.StringAttribute{
				Description: "Identifies the account the CA set belongs to.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The user who created the CA set.",
				Computed:    true,
			},
			"created_date": schema.StringAttribute{
				Description: "The date when the CA set was created.",
				Computed:    true,
			},
			"deleted_by": schema.StringAttribute{
				Description: "The user who deleted the CA set.",
				Computed:    true,
			},
			"deleted_date": schema.StringAttribute{
				Description: "The date when the CA set was deleted.",
				Computed:    true,
			},
			"version_created_by": schema.StringAttribute{
				Description: "The user who created the CA set version.",
				Computed:    true,
			},
			"version_created_date": schema.StringAttribute{
				Description: "When the CA set version was created.",
				Computed:    true,
			},
			"version_modified_by": schema.StringAttribute{
				Description: "The user who modified the CA set version.",
				Computed:    true,
			},
			"version_modified_date": schema.StringAttribute{
				Description: "When the CA set version was modified.",
				Computed:    true,
			},
			"allow_insecure_sha1": schema.BoolAttribute{
				Description: "By default, the version's certificates need a signature algorithm of SHA-256 or better. Enabling this allows certificates with SHA-1 signatures.",
				Computed:    true,
			},
			"version_description": schema.StringAttribute{
				Description: "Any additional description you can provide while creating or updating the CA set version.",
				Computed:    true,
			},
			"staging_version": schema.Int64Attribute{
				Description: "Version number of the CA set that is active on staging.",
				Computed:    true,
			},
			"production_version": schema.Int64Attribute{
				Description: "Version  of the CA set that is active on production.",
				Computed:    true,
			},
			"certificates": schema.ListNestedAttribute{
				Description: "The certificates that are valid, non-expired, root, or intermediate.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"certificate_pem": schema.StringAttribute{
							Description: "The certificate in PEM format, as found in a Base64 ASCII encoded file.",
							Computed:    true,
							Sensitive:   true,
						},
						"description": schema.StringAttribute{
							Description: "Optional description for the certificate.",
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
						"start_date": schema.StringAttribute{
							Description: "The start date of the certificate.",
							Computed:    true,
						},
						"end_date": schema.StringAttribute{
							Description: "The certificate's ISO 8601 formatted expiration date.",
							Computed:    true,
						},
						"fingerprint": schema.StringAttribute{
							Description: "The fingerprint of the certificate.",
							Computed:    true,
							Sensitive:   true,
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
						"subject": schema.StringAttribute{
							Description: "The certificate's subject field.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *caSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS TrustStore CA Set DataSource Read")

	var data caSetDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := Client(d.meta)

	if !data.Name.IsNull() {
		tflog.Debug(ctx, "'name' provided, attempting to find CA set ID")
		setID, err := findCASetID(ctx, client, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Read CA set failed", err.Error())
			return
		}
		data.ID = types.StringValue(setID)
	}

	caSet, err := client.GetCASet(ctx, mtlstruststore.GetCASetRequest{
		CASetID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read CA set failed", err.Error())
		return
	}

	var caSetVersion *mtlstruststore.GetCASetVersionResponse
	var version int64
	if !data.Version.IsNull() && !data.Version.IsUnknown() {
		tflog.Debug(ctx, "'version' provided, fetching specific version of CA set")
		version = data.Version.ValueInt64()
	} else if caSet.LatestVersion != nil {
		tflog.Debug(ctx, "No 'version' provided, using latest version of CA set")
		version = *caSet.LatestVersion
	} else {
		tflog.Debug(ctx, "No 'version' provided and CA Set has no latest version")
	}

	if version != 0 {
		tflog.Debug(ctx, "Fetching specific version of CA set", map[string]interface{}{
			"ca_set_id": data.ID.ValueString(),
			"version":   version,
		})
		caSetVersion, err = client.GetCASetVersion(ctx, mtlstruststore.GetCASetVersionRequest{
			CASetID: data.ID.ValueString(),
			Version: version,
		})
		if err != nil {
			resp.Diagnostics.AddError("Read CA set failed", err.Error())
			return
		}
	}

	modelData := convertCASetDataToModel(caSet, caSetVersion)
	data.setData(modelData)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func convertCASetDataToModel(caSet *mtlstruststore.GetCASetResponse, caSetVersion *mtlstruststore.GetCASetVersionResponse) caSetDataSourceModel {
	model := caSetDataSourceModel{
		ID:          types.StringValue(caSet.CASetID),
		Name:        types.StringValue(caSet.CASetName),
		Description: types.StringValue(caSet.Description),
		AccountID:   types.StringValue(caSet.AccountID),
		CreatedBy:   types.StringValue(caSet.CreatedBy),
		CreatedDate: types.StringValue(caSet.CreatedDate.String()),
	}

	if caSet.LatestVersion != nil {
		model.Version = types.Int64Value(*caSet.LatestVersion)
	} else {
		model.Version = types.Int64Null()
	}

	if caSet.StagingVersion != nil {
		model.StagingVersion = types.Int64Value(*caSet.StagingVersion)
	} else {
		model.StagingVersion = types.Int64Null()
	}

	if caSet.ProductionVersion != nil {
		model.ProductionVersion = types.Int64Value(*caSet.ProductionVersion)
	} else {
		model.ProductionVersion = types.Int64Null()
	}

	if caSet.DeletedBy != nil {
		model.DeletedBy = types.StringValue(*caSet.DeletedBy)
	} else {
		model.DeletedBy = types.StringNull()
	}

	if caSet.DeletedDate != nil {
		model.DeletedDate = types.StringValue(caSet.DeletedDate.String())
	} else {
		model.DeletedDate = types.StringNull()
	}

	if caSetVersion != nil {
		model.setCASetVersionData(caSetVersion)
	}

	return model
}

func (m *caSetDataSourceModel) setCASetVersionData(v *mtlstruststore.GetCASetVersionResponse) {
	m.AllowInsecureSHA1 = types.BoolValue(v.AllowInsecureSHA1)
	m.VersionDescription = types.StringValue(v.Description)
	m.VersionCreatedBy = types.StringValue(v.CreatedBy)
	m.VersionCreatedDate = types.StringValue(v.CreatedDate.String())

	if v.ModifiedBy != nil {
		m.VersionModifiedBy = types.StringValue(*v.ModifiedBy)
	} else {
		m.VersionModifiedBy = types.StringNull()
	}

	if v.ModifiedDate != nil {
		m.VersionModifiedDate = types.StringValue(v.ModifiedDate.String())
	} else {
		m.VersionModifiedDate = types.StringNull()
	}

	certificates := make([]certificateModel, len(v.Certificates))
	for i, cert := range v.Certificates {
		certificates[i] = certificateModel{
			CertificatePEM:     types.StringValue(cert.CertificatePEM),
			Description:        types.StringValue(cert.Description),
			CreatedBy:          types.StringValue(cert.CreatedBy),
			CreatedDate:        types.StringValue(cert.CreatedDate.String()),
			StartDate:          types.StringValue(cert.StartDate.String()),
			EndDate:            types.StringValue(cert.EndDate.String()),
			Fingerprint:        types.StringValue(cert.Fingerprint),
			Issuer:             types.StringValue(cert.Issuer),
			SerialNumber:       types.StringValue(cert.SerialNumber),
			SignatureAlgorithm: types.StringValue(cert.SignatureAlgorithm),
			Subject:            types.StringValue(cert.Subject),
		}
	}
	m.Certificates = certificates
}

func (m *caSetDataSourceModel) setData(data caSetDataSourceModel) {
	m.ID = data.ID
	m.Name = data.Name
	m.Description = data.Description
	m.AccountID = data.AccountID
	m.CreatedBy = data.CreatedBy
	m.CreatedDate = data.CreatedDate
	m.VersionCreatedBy = data.VersionCreatedBy
	m.VersionCreatedDate = data.VersionCreatedDate
	m.VersionModifiedBy = data.VersionModifiedBy
	m.VersionModifiedDate = data.VersionModifiedDate
	m.AllowInsecureSHA1 = data.AllowInsecureSHA1
	m.VersionDescription = data.VersionDescription
	m.StagingVersion = data.StagingVersion
	m.ProductionVersion = data.ProductionVersion
	m.Certificates = data.Certificates
	m.DeletedBy = data.DeletedBy
	m.DeletedDate = data.DeletedDate
	// Set the version only if it wasn't provided.
	if m.Version.IsNull() || m.Version.IsUnknown() {
		m.Version = data.Version
	}
}
