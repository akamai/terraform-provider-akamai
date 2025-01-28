package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &accountSwitchKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &accountSwitchKeysDataSource{}
)

type (
	accountSwitchKeysDataSource struct {
		meta meta.Meta
	}

	accountSwitchKeysModel struct {
		ClientID          types.String       `tfsdk:"client_id"`
		Filter            types.String       `tfsdk:"filter"`
		AccountSwitchKeys []accountSwitchKey `tfsdk:"account_switch_keys"`
	}

	accountSwitchKey struct {
		AccountName      types.String `tfsdk:"account_name"`
		AccountSwitchKey types.String `tfsdk:"account_switch_key"`
	}
)

// NewAccountSwitchKeysDataSource returns a new account switch keys data source.
func NewAccountSwitchKeysDataSource() datasource.DataSource {
	return &accountSwitchKeysDataSource{}
}

func (d *accountSwitchKeysDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_account_switch_keys"
}

func (d *accountSwitchKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountSwitchKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Identity and Access Management account switch keys",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "Unique identifier for each API client. If not provided it assumes your client id.",
			},
			"filter": schema.StringAttribute{
				Optional:    true,
				Description: "Filters results by accountId or accountName. Enter at least three characters to filter the results by substring.",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(3)},
			},
			"account_switch_keys": schema.ListNestedAttribute{
				Description: "List of account switch keys and account names available to the client.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_name": schema.StringAttribute{
							Computed:    true,
							Description: "Descriptive label for the account.",
						},
						"account_switch_key": schema.StringAttribute{
							Computed:    true,
							Description: "The identifier for an account other than your API client's default.",
						},
					},
				},
			},
		},
	}
}

func (d *accountSwitchKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Account Switch Keys DataSource Read")

	var data accountSwitchKeysModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := inst.Client(d.meta)
	accountSwitchKeysResponse, err := client.ListAccountSwitchKeys(ctx, iam.ListAccountSwitchKeysRequest{
		ClientID: data.ClientID.ValueString(),
		Search:   data.Filter.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Reading IAM Account Switch Keys Failed", err.Error())
		return
	}

	keys := make([]accountSwitchKey, 0, len(accountSwitchKeysResponse))
	for _, key := range accountSwitchKeysResponse {
		switchKey := accountSwitchKey{
			AccountName:      types.StringValue(key.AccountName),
			AccountSwitchKey: types.StringValue(key.AccountSwitchKey),
		}
		keys = append(keys, switchKey)
	}
	data.AccountSwitchKeys = keys
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
