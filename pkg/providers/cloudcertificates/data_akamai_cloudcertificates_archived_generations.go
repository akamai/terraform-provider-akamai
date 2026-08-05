package cloudcertificates

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &archivedGenerationsDataSource{}
	_ datasource.DataSourceWithConfigure = &archivedGenerationsDataSource{}
)

type (
	archivedGenerationsDataSource struct {
		meta.DataSource
	}

	archivedGenerationsDataSourceModel struct {
		LineageID         types.Int64               `tfsdk:"lineage_id"`
		IncludeAlgorithms types.Bool                `tfsdk:"include_algorithms"`
		Generations       []archivedGenerationModel `tfsdk:"generations"`
	}

	archivedGenerationModel struct {
		GenerationID     types.Int64  `tfsdk:"generation_id"`
		GenerationStatus types.String `tfsdk:"generation_status"`
		commonGenerationModel
	}
)

// NewArchivedGenerationsDataSource returns a new Cloud Certificates Archived Generations data source.
func NewArchivedGenerationsDataSource() datasource.DataSource {
	return &archivedGenerationsDataSource{}
}

// Metadata configures the data source's type name.
func (d *archivedGenerationsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_archived_generations"
}

// Schema defines the Terraform schema for the archived generations data source.
func (d *archivedGenerationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves archived and abandoned generations for a Cloud Certificate Manager lineage, oldest first.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the certificate lineage.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"include_algorithms": schema.BoolAttribute{
				Optional:    true,
				Description: "When set to true, complete algorithm details, including certificate and CSR PEM material, are included. By default, only sparse algorithm details are returned.",
			},
			"generations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Archived and abandoned certificate generations, ordered oldest first.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: mergeAttributes(archivedGenerationCommonAttributes(),
						map[string]schema.Attribute{
							"generation_id": schema.Int64Attribute{
								Computed:    true,
								Description: "The unique identifier of the archived generation.",
							},
							"generation_status": schema.StringAttribute{
								Computed:    true,
								Description: "The status of the archived generation.",
							},
						}),
				},
			},
		},
	}
}

func archivedGenerationCommonAttributes() map[string]schema.Attribute {
	attributes := generationCommonAttributes()
	algorithmsAttribute := attributes["algorithms"].(schema.MapNestedAttribute)
	algorithmsAttribute.Description = "Per key-type certificate details for this generation. Sparse details are returned by default; complete details are returned when include_algorithms is true."
	attributes["algorithms"] = algorithmsAttribute

	descriptions := map[string]string{
		"first_promoted_to_production_time": "Time the generation was first promoted to production, in RFC3339 format. Null if never promoted.",
		"generation_created_by":             "Username of the person who created this generation.",
		"generation_created_time":           "Time the generation was created, in RFC3339 format.",
		"generation_modified_by":            "Username of the person who last modified this generation. Null if never modified.",
		"generation_modified_time":          "Time the generation was last modified, in RFC3339 format. Null if the generation has never been modified since creation.",
	}
	for name, description := range descriptions {
		attribute := attributes[name].(schema.StringAttribute)
		attribute.Description = description
		attributes[name] = attribute
	}

	return attributes
}

// Read fetches archived generations from the API and populates Terraform state.
func (d *archivedGenerationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Archived Generations DataSource Read")

	var data archivedGenerationsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	lineageID := data.LineageID.ValueInt64()
	generations, err := d.Client.GetCloudCertificates().ListArchivedGenerations(ctx, cloudcertificates.ListArchivedGenerationsRequest{
		LineageID:         lineageID,
		IncludeAlgorithms: data.IncludeAlgorithms.ValueBool(),
	})
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrLineageNotFound) {
			resp.Diagnostics.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d", lineageID))
			return
		}
		resp.Diagnostics.AddError("Failed to retrieve archived generations", err.Error())
		return
	}

	data.Generations = archivedGenerationsToModels(generations.Items)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func archivedGenerationsToModels(generations []cloudcertificates.ArchivedGeneration) []archivedGenerationModel {
	result := make([]archivedGenerationModel, len(generations))
	for i, generation := range generations {
		result[i] = archivedGenerationModel{
			GenerationID:          types.Int64Value(generation.GenerationID),
			GenerationStatus:      types.StringValue(generation.GenerationStatus),
			commonGenerationModel: mapCommonGeneration(generation.Generation),
		}
	}
	return result
}
