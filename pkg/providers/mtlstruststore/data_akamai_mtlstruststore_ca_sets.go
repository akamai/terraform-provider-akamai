package mtlstruststore

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	caSetsDataSource struct {
		meta meta.Meta
	}

	caSetsDataSourceModel struct {
		NamePrefix  types.String `tfsdk:"name_prefix"`
		ActivatedOn types.String `tfsdk:"activated_on"`
		CASets      []caSetModel `tfsdk:"ca_sets"`
	}

	caSetModel struct {
		AccountID         types.String `tfsdk:"account_id"`
		ID                types.String `tfsdk:"id"`
		Name              types.String `tfsdk:"name"`
		CreatedBy         types.String `tfsdk:"created_by"`
		CreatedDate       types.String `tfsdk:"created_date"`
		DeletedBy         types.String `tfsdk:"deleted_by"`
		DeletedDate       types.String `tfsdk:"deleted_date"`
		Description       types.String `tfsdk:"description"`
		Status            types.String `tfsdk:"status"`
		LatestVersion     types.Int64  `tfsdk:"latest_version"`
		StagingVersion    types.Int64  `tfsdk:"staging_version"`
		ProductionVersion types.Int64  `tfsdk:"production_version"`
	}
)

// NewCASetsDataSource returns a new mtls truststore ca sets data source.
func NewCASetsDataSource() datasource.DataSource {
	return &caSetsDataSource{}
}

func (d *caSetsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_sets"
}

func (d *caSetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve list of MTLS Truststore CA Sets.",
		Attributes: map[string]schema.Attribute{
			"name_prefix": schema.StringAttribute{
				Description: "The name prefix for CA sets filtering or empty to return all CA sets.",
				Optional:    true,
			},
			"activated_on": schema.StringAttribute{
				Description: "When provided it filters where CA sets were activated 'INACTIVE', 'STAGING', 'PRODUCTION', 'STAGING+PRODUCTION', 'PRODUCTION+STAGING', 'STAGING,PRODUCTION', 'PRODUCTION,STAGING' network.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("INACTIVE", "STAGING", "PRODUCTION", "STAGING+PRODUCTION", "PRODUCTION+STAGING", "STAGING,PRODUCTION", "PRODUCTION,STAGING", "")},
			},
			"ca_sets": schema.ListNestedAttribute{
				Description: "List of CA sets.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Description: "Identifies the account the CA set belongs to.",
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "Identifies each CA set.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the CA set.",
							Computed:    true,
						},
						"created_date": schema.StringAttribute{
							Description: "When the CA set was created.",
							Computed:    true,
						},
						"created_by": schema.StringAttribute{
							Description: "The user who created the CA set.",
							Computed:    true,
						},
						"deleted_date": schema.StringAttribute{
							Description: "When the CA set was deleted, or null if there's no request.",
							Computed:    true,
						},
						"deleted_by": schema.StringAttribute{
							Description: "The user who requested the CA set be deleted, or null if there's no request.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Any additional comments you can add to the CA set.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Indicates if the CA set was deleted, either 'NOT_DELETED', 'DELETING', or 'DELETED'.",
							Computed:    true,
						},
						"latest_version": schema.Int64Attribute{
							Description: "The most recent version based on the updated version.",
							Computed:    true,
						},
						"production_version": schema.Int64Attribute{
							Description: "The CA set version activated on the 'PRODUCTION' network.",
							Computed:    true,
						},
						"staging_version": schema.Int64Attribute{
							Description: "The CA set version activated on the 'STAGING' network.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure configures data source at the beginning of the lifecycle.
func (d *caSetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *caSetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS TrustStore CA Sets DataSource Read")

	var data caSetsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client = Client(d.meta)

	caSets, err := client.ListCASets(ctx, mtlstruststore.ListCASetsRequest{
		CASetNamePrefix: data.NamePrefix.ValueString(),
		ActivatedOn:     mtlstruststore.Network(data.ActivatedOn.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read CA sets failed", err.Error())
		return
	}

	data.convertCASetsToModel(*caSets)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *caSetsDataSourceModel) convertCASetsToModel(caSets mtlstruststore.ListCASetsResponse) {
	for _, caSet := range caSets.CASets {
		m.CASets = append(m.CASets, caSetModel{
			AccountID:         types.StringValue(caSet.AccountID),
			ID:                types.StringValue(caSet.CASetID),
			Name:              types.StringValue(caSet.CASetName),
			CreatedBy:         types.StringValue(caSet.CreatedBy),
			CreatedDate:       date.TimeRFC3339NanoValue(caSet.CreatedDate),
			Description:       types.StringPointerValue(caSet.Description),
			Status:            types.StringValue(caSet.CASetStatus),
			DeletedBy:         types.StringPointerValue(caSet.DeletedBy),
			DeletedDate:       date.TimeRFC3339NanoPointerValue(caSet.DeletedDate),
			LatestVersion:     types.Int64PointerValue(caSet.LatestVersion),
			StagingVersion:    types.Int64PointerValue(caSet.StagingVersion),
			ProductionVersion: types.Int64PointerValue(caSet.ProductionVersion),
		})
	}
}
