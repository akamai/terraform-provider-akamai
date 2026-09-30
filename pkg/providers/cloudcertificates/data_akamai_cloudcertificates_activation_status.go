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
	_ datasource.DataSource              = &activationStatusDataSource{}
	_ datasource.DataSourceWithConfigure = &activationStatusDataSource{}
)

type (
	activationStatusDataSource struct {
		meta.DataSource
	}

	activationStatusDataSourceModel struct {
		LineageID    types.Int64 `tfsdk:"lineage_id"`
		ActivationID types.Int64 `tfsdk:"activation_id"`
		activationDetailModel
	}
)

// NewActivationStatusDataSource returns a new CloudCertificates Activation Status data source.
func NewActivationStatusDataSource() datasource.DataSource {
	return &activationStatusDataSource{}
}

func (d *activationStatusDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_activation_status"
}

// Schema is used to define data source's terraform schema.
func (d *activationStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the status of a certificate lineage activation request, as created by a promote, rollback, or replace-staging operation.",
		Attributes: mergeAttributes(activationDetailAttributes(), map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the certificate lineage.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
			"activation_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the activation request.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
			},
		}),
	}
}

func (d *activationStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CCM Activation Status DataSource Read")
	var data activationStatusDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	activation, err := d.Client.GetCloudCertificates().GetActivationStatus(ctx, cloudcertificates.GetActivationStatusRequest{
		LineageID:    data.LineageID.ValueInt64(),
		ActivationID: data.ActivationID.ValueInt64(),
	})
	if err != nil {
		switch {
		case errors.Is(err, cloudcertificates.ErrActivationNotFound):
			resp.Diagnostics.AddError("Activation Not Found", fmt.Sprintf("No activation found with ID %d for lineage %d", data.ActivationID.ValueInt64(), data.LineageID.ValueInt64()))
		case errors.Is(err, cloudcertificates.ErrLineageNotFound):
			resp.Diagnostics.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d", data.LineageID.ValueInt64()))
		default:
			resp.Diagnostics.AddError("Failed to retrieve activation status", err.Error())
		}
		return
	}

	data.activationDetailModel = mapActivationDetail(*activation)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
