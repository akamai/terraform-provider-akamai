package iam

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &passwordPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordPolicyDataSource{}
)

type (
	passwordPolicyDataSource struct {
		meta meta.Meta
	}

	passwordPolicyModel struct {
		PwClass         types.String `tfsdk:"pw_class"`
		CaseDif         types.Int64  `tfsdk:"case_dif"`
		MaxRepeating    types.Int64  `tfsdk:"max_repeating"`
		MinDigits       types.Int64  `tfsdk:"min_digits"`
		MinLength       types.Int64  `tfsdk:"min_length"`
		MinLetters      types.Int64  `tfsdk:"min_letters"`
		MinNonAlpha     types.Int64  `tfsdk:"min_non_alpha"`
		MinReuse        types.Int64  `tfsdk:"min_reuse"`
		RotateFrequency types.Int64  `tfsdk:"rotate_frequency"`
	}
)

// NewPasswordPolicyDataSource returns a new password policy data source.
func NewPasswordPolicyDataSource() datasource.DataSource {
	return &passwordPolicyDataSource{}
}

func (d *passwordPolicyDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_password_policy"
}

func (d *passwordPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *passwordPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Identity and Access Management password policy.",
		Attributes: map[string]schema.Attribute{
			"pw_class": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for a password policy.",
			},
			"min_length": schema.Int64Attribute{
				Computed:    true,
				Description: "The minimum length of a password.",
			},
			"min_letters": schema.Int64Attribute{
				Computed:    true,
				Description: "The minimum number of letters in a password.",
			},
			"min_digits": schema.Int64Attribute{
				Computed:    true,
				Description: "The minimum number of digits in a password.",
			},
			"case_dif": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of characters that, at minimum, need to be in a different case.",
			},
			"min_non_alpha": schema.Int64Attribute{
				Computed:    true,
				Description: "The minimum number of non-alphabetic characters in a password.",
			},
			"max_repeating": schema.Int64Attribute{
				Computed:    true,
				Description: "The maximum allowed number of repeating characters.",
			},
			"min_reuse": schema.Int64Attribute{
				Computed:    true,
				Description: "The minimum number of previous passwords to retain to prevent password reuse.",
			},
			"rotate_frequency": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of days a password is valid.",
			},
		},
	}
}

func (d *passwordPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Password Policy DataSource Read")

	var data passwordPolicyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client := inst.Client(d.meta)
	passwordPolicyResponse, err := client.GetPasswordPolicy(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Reading IAM Password Policy Failed", err.Error())
		return
	}

	data.PwClass = types.StringValue(passwordPolicyResponse.PwClass)
	data.CaseDif = types.Int64Value(passwordPolicyResponse.CaseDiff)
	data.MaxRepeating = types.Int64Value(passwordPolicyResponse.MaxRepeating)
	data.MinDigits = types.Int64Value(passwordPolicyResponse.MinDigits)
	data.MinLength = types.Int64Value(passwordPolicyResponse.MinLength)
	data.MinLetters = types.Int64Value(passwordPolicyResponse.MinLetters)
	data.MinNonAlpha = types.Int64Value(passwordPolicyResponse.MinNonAlpha)
	data.MinReuse = types.Int64Value(passwordPolicyResponse.MinReuse)
	data.RotateFrequency = types.Int64Value(passwordPolicyResponse.RotateFrequency)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
