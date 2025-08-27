package mtlstruststore

import (
	"context"
	"fmt"
	"regexp"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &caSetAssociationsDataSource{}
	_ datasource.DataSourceWithConfigure = &caSetAssociationsDataSource{}
)

type (
	caSetAssociationsDataSource struct {
		meta meta.Meta
	}

	caSetAssociationsModel struct {
		ID          types.String                 `tfsdk:"id"`
		Name        types.String                 `tfsdk:"name"`
		Properties  []propertyAssociationModel   `tfsdk:"properties"`
		Enrollments []enrollmentAssociationModel `tfsdk:"enrollments"`
	}

	propertyAssociationModel struct {
		PropertyID   types.String    `tfsdk:"property_id"`
		PropertyName types.String    `tfsdk:"property_name"`
		AssetID      types.Int64     `tfsdk:"asset_id"`
		GroupID      types.Int64     `tfsdk:"group_id"`
		Hostnames    []hostnameModel `tfsdk:"hostnames"`
	}

	hostnameModel struct {
		Hostname types.String `tfsdk:"hostname"`
		Network  types.String `tfsdk:"network"`
		Status   types.String `tfsdk:"status"`
	}

	enrollmentAssociationModel struct {
		EnrollmentID    types.Int64  `tfsdk:"enrollment_id"`
		StagingSlots    types.List   `tfsdk:"staging_slots"`
		ProductionSlots types.List   `tfsdk:"production_slots"`
		CN              types.String `tfsdk:"cn"`
	}
)

// NewCASetAssociationsDataSource returns a new CA Set Associations data source.
func NewCASetAssociationsDataSource() datasource.DataSource {
	return &caSetAssociationsDataSource{}
}

// Metadata configures data source's meta information.
func (d *caSetAssociationsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set_associations"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *caSetAssociationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *caSetAssociationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of the properties and/or enrollments where a given ca set is used.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID is a unique identifier representing the CA set. Either 'id' or 'name' must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
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
			"properties": schema.ListNestedAttribute{
				Description: "Properties associated with given CASet.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"property_id": schema.StringAttribute{
							Description: "A unique identifier for the property.",
							Computed:    true,
						},
						"property_name": schema.StringAttribute{
							Description: "A unique, descriptive name for the property.",
							Computed:    true,
						},
						"asset_id": schema.Int64Attribute{
							Description: "An alternative identifier for the property.",
							Computed:    true,
						},
						"group_id": schema.Int64Attribute{
							Description: "Identifies the group to which the property is assigned.",
							Computed:    true,
						},
						"hostnames": schema.ListNestedAttribute{
							Description: "Contains details about associated hostnames.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"hostname": schema.StringAttribute{
										Description: "The name of the device.",
										Computed:    true,
									},
									"network": schema.StringAttribute{
										Description: "The network on which CA set to hostname association is formed/removed/in progress. The values for this are 'STAGING', 'PRODUCTION'.",
										Computed:    true,
									},
									"status": schema.StringAttribute{
										Description: "The status of CA set to hostname association. The values for it are - 'ATTACHING', 'DETACHING', 'ATTACHED'.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"enrollments": schema.ListNestedAttribute{
				Description: "Enrollments associated with a given CA Set.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enrollment_id": schema.Int64Attribute{
							Description: "A unique identifier for the enrollment.",
							Computed:    true,
						},
						"staging_slots": schema.ListAttribute{
							Description: "Slots where the certificate is deployed on the staging network.",
							Computed:    true,
							ElementType: types.Int64Type,
						},
						"production_slots": schema.ListAttribute{
							Description: "Slots where the certificate is deployed on the production network.",
							Computed:    true,
							ElementType: types.Int64Type,
						},
						"cn": schema.StringAttribute{
							Description: "The domain name to use for the certificate, also known as the common name.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *caSetAssociationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "mTLS TrustStore CA Set Associations Data Source")

	var data caSetAssociationsModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client = Client(d.meta)

	if !data.Name.IsNull() {
		setID, err := findCASetID(ctx, client, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Could not fetch CA set ID for provided name", err.Error())
			return
		}
		data.ID = types.StringValue(setID)
	} else {
		caSet, err := client.GetCASet(ctx, mtlstruststore.GetCASetRequest{CASetID: data.ID.ValueString()})
		if err != nil {
			resp.Diagnostics.AddError("Could not fetch CA set", err.Error())
			return
		}
		data.Name = types.StringValue(caSet.CASetName)
	}

	associations, err := client.ListCASetAssociations(ctx, mtlstruststore.ListCASetAssociationsRequest{CASetID: data.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error fetching CA set associations", err.Error())
		return
	}

	diags := data.assignAssociations(ctx, associations)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *caSetAssociationsModel) assignAssociations(ctx context.Context, associations *mtlstruststore.ListCASetAssociationsResponse) diag.Diagnostics {
	enrollments := make([]enrollmentAssociationModel, 0, len(associations.Associations.Enrollments))
	for _, enrollment := range associations.Associations.Enrollments {
		stagingSlots, diags := types.ListValueFrom(ctx, types.Int64Type, enrollment.StagingSlots)
		if diags.HasError() {
			return diags
		}
		productionSlots, diags := types.ListValueFrom(ctx, types.Int64Type, enrollment.ProductionSlots)
		if diags.HasError() {
			return diags
		}

		enrollments = append(enrollments, enrollmentAssociationModel{
			EnrollmentID:    types.Int64Value(enrollment.EnrollmentID),
			StagingSlots:    stagingSlots,
			ProductionSlots: productionSlots,
			CN:              types.StringValue(enrollment.CN),
		})
	}
	m.Enrollments = enrollments

	props := make([]propertyAssociationModel, 0, len(associations.Associations.Properties))
	for _, property := range associations.Associations.Properties {
		props = append(props, propertyAssociationModel{
			PropertyID:   types.StringValue(property.PropertyID),
			PropertyName: types.StringPointerValue(property.PropertyName),
			AssetID:      types.Int64PointerValue(property.AssetID),
			GroupID:      types.Int64PointerValue(property.GroupID),
			Hostnames:    m.getHostnameAssociations(property.Hostnames),
		})
	}
	m.Properties = props
	return nil
}

func (m *caSetAssociationsModel) getHostnameAssociations(hostnames []mtlstruststore.AssociationHostname) []hostnameModel {
	hostnamesModel := make([]hostnameModel, 0, len(hostnames))
	for _, hostname := range hostnames {
		hostnamesModel = append(hostnamesModel, hostnameModel{
			Hostname: types.StringValue(hostname.Hostname),
			Network:  types.StringValue(hostname.Network),
			Status:   types.StringValue(hostname.Status),
		})
	}
	return hostnamesModel
}
