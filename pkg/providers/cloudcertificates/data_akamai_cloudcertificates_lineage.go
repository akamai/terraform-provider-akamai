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
	_ datasource.DataSource              = &lineageDataSource{}
	_ datasource.DataSourceWithConfigure = &lineageDataSource{}
)

type (
	lineageDataSource struct {
		meta.DataSource
	}

	lineageDataSourceModel struct {
		LineageID         types.Int64 `tfsdk:"lineage_id"`
		ExpandGenerations types.Bool  `tfsdk:"expand_generations"`
		commonLineageModel
	}
)

// NewLineageDataSource returns a new lineage data source.
func NewLineageDataSource() datasource.DataSource {
	return &lineageDataSource{}
}

// Metadata configures the data source's type name.
func (d *lineageDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_lineage"
}

// Schema defines the Terraform schema for the lineage data source.
func (d *lineageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a single Cloud Certificate Manager lineage by its ID. " +
			"Set expand_generations to true to include full algorithm instance details for all active generation pointers " +
			"(HEAD, CURRENT_PRODUCTION, PREVIOUS_PRODUCTION, CURRENT_STAGING).",
		Attributes: mergeAttributes(lineageCommonAttributes(), map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the lineage to retrieve.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"expand_generations": schema.BoolAttribute{
				Optional: true,
				Description: "When true, expands all active generation pointers (HEAD, CURRENT_PRODUCTION, " +
					"PREVIOUS_PRODUCTION, CURRENT_STAGING) to include full algorithm instance details.",
			},
		}),
	}
}

// Read fetches the lineage state from the API and populates Terraform state.
func (d *lineageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineage DataSource Read")

	var data lineageDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	lineageID := data.LineageID.ValueInt64()

	var expandGenerations []cloudcertificates.ExpandGenerations
	if data.ExpandGenerations.ValueBool() {
		expandGenerations = []cloudcertificates.ExpandGenerations{
			cloudcertificates.ExpandGenerationsHead,
			cloudcertificates.ExpandGenerationsCurrentProduction,
			cloudcertificates.ExpandGenerationsPreviousProduction,
			cloudcertificates.ExpandGenerationsCurrentStaging,
		}
	}

	client := d.Client.GetCloudCertificates()
	lineage, err := client.GetLineage(ctx, cloudcertificates.GetLineageRequest{
		LineageID:         lineageID,
		ExpandGenerations: expandGenerations,
	})
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrLineageNotFound) {
			resp.Diagnostics.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d", lineageID))
		} else {
			resp.Diagnostics.AddError("Read Cloud Certificates Lineage failed", err.Error())
		}
		return
	}

	data = mapLineageResponseToModel(data.LineageID, data.ExpandGenerations, lineage)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapLineageResponseToModel(lineageID types.Int64, expandGenerations types.Bool, lineage *cloudcertificates.GetLineageResponse) lineageDataSourceModel {
	return lineageDataSourceModel{
		LineageID:          lineageID,
		ExpandGenerations:  expandGenerations,
		commonLineageModel: mapCommonLineage(cloudcertificates.Lineage(*lineage)),
	}
}
