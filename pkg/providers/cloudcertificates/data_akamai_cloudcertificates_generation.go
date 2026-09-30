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
	_ datasource.DataSource              = &generationDataSource{}
	_ datasource.DataSourceWithConfigure = &generationDataSource{}
)

type (
	generationDataSource struct {
		meta.DataSource
	}

	generationDataSourceModel struct {
		LineageID        types.Int64  `tfsdk:"lineage_id"`
		GenerationID     types.Int64  `tfsdk:"generation_id"`
		GenerationStatus types.String `tfsdk:"generation_status"`
		commonGenerationModel
	}
)

// NewGenerationDataSource returns a new Cloud Certificates Generation data source.
func NewGenerationDataSource() datasource.DataSource {
	return &generationDataSource{}
}

// Metadata configures the data source's type name.
func (d *generationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_generation"
}

// Schema defines the Terraform schema for the generation data source.
func (d *generationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a single Cloud Certificate Manager certificate generation of a certificate lineage.",
		Attributes: mergeAttributes(generationCommonAttributes(), map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the certificate lineage.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"generation_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the generation to retrieve.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"generation_status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the generation.",
			},
		}),
	}
}

// Read fetches the generation from the API and populates Terraform state.
func (d *generationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Generation DataSource Read")

	var data generationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	generation, err := d.Client.GetCloudCertificates().GetGeneration(ctx, cloudcertificates.GetGenerationRequest{
		LineageID:    data.LineageID.ValueInt64(),
		GenerationID: data.GenerationID.ValueInt64(),
	})
	if err != nil {
		switch {
		case errors.Is(err, cloudcertificates.ErrGenerationNotFound):
			resp.Diagnostics.AddError("Generation Not Found", fmt.Sprintf("No generation found with ID %d for lineage %d", data.GenerationID.ValueInt64(), data.LineageID.ValueInt64()))
		case errors.Is(err, cloudcertificates.ErrLineageNotFound):
			resp.Diagnostics.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d", data.LineageID.ValueInt64()))
		default:
			resp.Diagnostics.AddError("Failed to retrieve generation", err.Error())
		}
		return
	}

	data.convertGenerationToModel(*generation)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *generationDataSourceModel) convertGenerationToModel(generation cloudcertificates.GetGenerationResponse) {
	m.GenerationStatus = types.StringValue(generation.GenerationStatus)
	m.commonGenerationModel = mapCommonGeneration(generation.Generation)
}
