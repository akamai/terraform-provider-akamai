package cloudcertificates

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &activationsDataSource{}
	_ datasource.DataSourceWithConfigure = &activationsDataSource{}
)

type (
	activationsDataSource struct {
		meta.DataSource
	}

	activationsDataSourceModel struct {
		LineageID   types.Int64       `tfsdk:"lineage_id"`
		Limit       types.Int64       `tfsdk:"limit"`
		Activations []activationModel `tfsdk:"activations"`
	}
)

// NewActivationsDataSource returns a new Cloud Certificates Activations data source.
func NewActivationsDataSource() datasource.DataSource {
	return &activationsDataSource{}
}

// Metadata configures the data source's metadata.
func (d *activationsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_activations"
}

// Schema defines the data source's Terraform schema.
func (d *activationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the history of activation requests for a certificate lineage.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique numeric identifier of the certificate lineage.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"limit": schema.Int64Attribute{
				Optional: true,
				Description: "The maximum number of most recent activation records to return. " +
					"When not specified, all activation records for the lineage are returned.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"activations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Activation records for the certificate lineage, ordered newest first.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: mergeAttributes(activationDetailAttributes(), map[string]schema.Attribute{
						"activation_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The unique identifier of the activation request.",
						},
						"lineage_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The unique identifier of the certificate lineage associated with the activation.",
						},
					}),
				},
			},
		},
	}
}

func (d *activationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Activations DataSource Read")

	var data activationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var limit *int64
	if tf.IsKnown(data.Limit) {
		limit = data.Limit.ValueInt64Pointer()
	}

	activations, err := d.listAllActivations(ctx, data.LineageID.ValueInt64(), limit)
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrLineageNotFound) {
			resp.Diagnostics.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d", data.LineageID.ValueInt64()))
		} else {
			resp.Diagnostics.AddError("Failed to retrieve lineage activations", err.Error())
		}
		return
	}

	data.Activations = activationsToModels(activations)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// listAllActivations follows the API's keyset pagination cursor until the activation history is exhausted or
// the requested number of records has been collected. Requests use at most 100 records, and the final request
// is reduced to the remaining limit. A nil limit collects all records.
func (d *activationsDataSource) listAllActivations(ctx context.Context, lineageID int64, limit *int64) ([]cloudcertificates.GetActivationStatusResponse, error) {
	client := d.Client.GetCloudCertificates()

	req := cloudcertificates.ListActivationsRequest{
		LineageID: lineageID,
		PageSize:  cloudcertificates.MaxListActivationsPageSize,
	}

	var activations []cloudcertificates.GetActivationStatusResponse
	for {
		var remaining int64
		if limit != nil {
			remaining = *limit - int64(len(activations))
			if remaining <= 0 {
				break
			}
			req.PageSize = int(min(remaining, int64(cloudcertificates.MaxListActivationsPageSize)))
		}

		resp, err := client.ListActivations(ctx, req)
		if err != nil {
			return nil, err
		}

		items := resp.Items
		if limit != nil && int64(len(items)) > remaining {
			items = items[:remaining]
		}
		activations = append(activations, items...)

		if resp.NextCursor == nil || *resp.NextCursor == "" || *resp.NextCursor == req.Cursor {
			break
		}
		req.Cursor = *resp.NextCursor
	}

	return activations, nil
}

func activationsToModels(items []cloudcertificates.GetActivationStatusResponse) []activationModel {
	activations := make([]activationModel, len(items))
	for i, activation := range items {
		activations[i] = mapActivation(activation)
	}
	return activations
}
