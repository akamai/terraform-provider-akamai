package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &cidrBlockDataSource{}
	_ datasource.DataSourceWithConfigure = &cidrBlockDataSource{}

	// ErrIAMGetCIDRBlock is returned when GetCIDRBlock fails.
	ErrIAMGetCIDRBlock = errors.New("IAM get CIDR block failed")
)

type (
	cidrBlockDataSource struct {
		meta meta.Meta
	}

	cidrBlockModel struct {
		CIDRBlockID  types.Int64  `tfsdk:"cidr_block_id"`
		Actions      *actions     `tfsdk:"actions"`
		CIDRBlock    types.String `tfsdk:"cidr_block"`
		Comments     types.String `tfsdk:"comments"`
		CreatedBy    types.String `tfsdk:"created_by"`
		CreatedDate  types.String `tfsdk:"created_date"`
		Enabled      types.Bool   `tfsdk:"enabled"`
		ModifiedBy   types.String `tfsdk:"modified_by"`
		ModifiedDate types.String `tfsdk:"modified_date"`
	}

	actions struct {
		Delete types.Bool `tfsdk:"delete"`
		Edit   types.Bool `tfsdk:"edit"`
	}
)

// NewCIDRBlockDataSource returns a new iam CIDR block data source.
func NewCIDRBlockDataSource() datasource.DataSource {
	return &cidrBlockDataSource{}
}

func (d *cidrBlockDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_cidr_block"
}

func (d *cidrBlockDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cidrBlockDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Identity and Access Management CIDR block.",
		Attributes: map[string]schema.Attribute{
			"cidr_block_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier for each CIDR block.",
			},
			"actions": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Specifies activities available for the CIDR block.",
				Attributes: map[string]schema.Attribute{
					"delete": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can delete this CIDR block.",
					},
					"edit": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can edit this CIDR block.",
					},
				},
			},
			"cidr_block": schema.StringAttribute{
				Computed:    true,
				Description: "The value of an IP address or IP address range.",
			},
			"comments": schema.StringAttribute{
				Computed:    true,
				Description: "Descriptive label you provide for the CIDR block.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the CIDR block.",
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the CIDR block was created.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the CIDR block is enabled.",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who last edited the CIDR block.",
			},
			"modified_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the CIDR block was last modified.",
			},
		},
	}
}

func (d *cidrBlockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM CIDR Block DataSource Read")

	var data cidrBlockModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := inst.Client(d.meta)

	id := data.CIDRBlockID.ValueInt64()
	cidrBlock, err := client.GetCIDRBlock(ctx, iam.GetCIDRBlockRequest{
		CIDRBlockID: id,
		Actions:     true,
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("%s:", ErrIAMGetCIDRBlock), err.Error())
		return
	}

	data.read(cidrBlock)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *cidrBlockModel) read(cidrBlock *iam.GetCIDRBlockResponse) {
	d.CIDRBlock = types.StringValue(cidrBlock.CIDRBlock)
	d.CreatedBy = types.StringValue(cidrBlock.CreatedBy)
	d.CreatedDate = types.StringValue(date.FormatRFC3339Nano(cidrBlock.CreatedDate))
	d.Enabled = types.BoolValue(cidrBlock.Enabled)
	d.ModifiedBy = types.StringValue(cidrBlock.ModifiedBy)
	d.ModifiedDate = types.StringValue(date.FormatRFC3339Nano(cidrBlock.ModifiedDate))

	if cidrBlock.Actions != nil {
		d.Actions = &actions{
			Delete: types.BoolValue(cidrBlock.Actions.Delete),
			Edit:   types.BoolValue(cidrBlock.Actions.Edit),
		}
	}

	if cidrBlock.Comments != nil {
		d.Comments = types.StringValue(*cidrBlock.Comments)
	} else {
		d.Comments = types.StringNull()
	}
}
