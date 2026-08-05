package cloudcertificates

import (
	"context"
	"math"
	"slices"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &lineagesDataSource{}
	_ datasource.DataSourceWithConfigure = &lineagesDataSource{}
)

type (
	lineagesDataSource struct {
		meta.DataSource
	}

	lineagesDataSourceModel struct {
		// Filter inputs
		ContractID        types.String   `tfsdk:"contract_id"`
		LineageName       types.String   `tfsdk:"lineage_name"`
		LineageIDs        []types.Int64  `tfsdk:"lineage_ids"`
		SecureNetwork     types.String   `tfsdk:"secure_network"`
		StackMode         types.String   `tfsdk:"stack_mode"`
		LineageType       types.String   `tfsdk:"lineage_type"`
		Domain            types.String   `tfsdk:"domain"`
		GenerationStatus  []types.String `tfsdk:"generation_status"`
		ExpiringInDays    types.Int64    `tfsdk:"expiring_in_days"`
		KeyType           types.String   `tfsdk:"key_type"`
		Issuer            types.String   `tfsdk:"issuer"`
		ExpandGenerations types.Bool     `tfsdk:"expand_generations"`
		Sort              types.String   `tfsdk:"sort"`
		Limit             types.Int64    `tfsdk:"limit"`

		// Computed output
		Lineages []lineageItemModel `tfsdk:"lineages"`
	}

	lineageItemModel struct {
		LineageID types.Int64 `tfsdk:"lineage_id"`
		commonLineageModel
	}
)

// NewLineagesDataSource returns a new lineages data source.
func NewLineagesDataSource() datasource.DataSource {
	return &lineagesDataSource{}
}

// Metadata configures the data source's type name.
func (d *lineagesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_lineages"
}

// Schema defines the Terraform schema for the lineages data source.
func (d *lineagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of Cloud Certificate Manager lineages matching the specified search and filter criteria.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter lineages by Akamai contract ID.",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"lineage_name": schema.StringAttribute{
				Optional:    true,
				Description: "Filter lineages by substring match on lineage name.",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"lineage_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "Filter lineages by one or more lineage IDs.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.NoNullValues(),
					setvalidator.ValueInt64sAre(int64validator.AtLeast(1)),
				},
			},
			"secure_network": schema.StringAttribute{
				Optional:    true,
				Description: "Filter lineages by secure network type, e.g. ENHANCED_TLS or STANDARD_TLS.",
				Validators: []validator.String{stringvalidator.OneOf(
					string(cloudcertificates.SecureNetworkEnhancedTLS),
					string(cloudcertificates.SecureNetworkStandardTLS),
				)},
			},
			"stack_mode": schema.StringAttribute{
				Optional:    true,
				Description: "Filter lineages by stack mode, e.g. SINGLE_STACK or MULTIPLE_STACK.",
				Validators: []validator.String{stringvalidator.OneOf(
					string(cloudcertificates.StackModeSingleStack),
					string(cloudcertificates.StackModeMultipleStack),
				)},
			},
			"lineage_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter lineages by type, e.g. MULTIPLE_GENERATION or SINGLE_GENERATION.",
				Validators: []validator.String{stringvalidator.OneOf(
					string(cloudcertificates.LineageTypeMultipleGeneration),
					string(cloudcertificates.LineageTypeSingleGeneration),
				)},
			},
			"domain": schema.StringAttribute{
				Optional:    true,
				Description: "Filter lineages by case-insensitive exact or wildcard match against SANs or common name.",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"generation_status": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Filter lineages by active generation statuses (CSR_READY, READY_FOR_USE, ACTIVE, ARCHIVED, ABANDONED).",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.NoNullValues(),
					setvalidator.ValueStringsAre(stringvalidator.OneOf(
						string(cloudcertificates.GenerationStatusCSRReady),
						string(cloudcertificates.GenerationStatusReadyForUse),
						string(cloudcertificates.GenerationStatusActive),
						string(cloudcertificates.GenerationStatusArchived),
						string(cloudcertificates.GenerationStatusAbandoned),
					)),
				},
			},
			"expiring_in_days": schema.Int64Attribute{
				Optional:    true,
				Description: "Filter lineages with a READY_FOR_USE signed certificate expiring within N days.",
				Validators:  []validator.Int64{int64validator.Between(0, math.MaxInt32)},
			},
			"key_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter lineages declaring the given cryptographic key algorithm, e.g. RSA or ECDSA.",
				Validators: []validator.String{stringvalidator.OneOf(
					string(cloudcertificates.CryptographicAlgorithmRSA),
					string(cloudcertificates.CryptographicAlgorithmECDSA),
				)},
			},
			"issuer": schema.StringAttribute{
				Optional:    true,
				Description: "Filter lineages by substring match on the signed certificate issuer.",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"expand_generations": schema.BoolAttribute{
				Optional:    true,
				Description: "When true, expands all active generation pointers (HEAD, CURRENT_PRODUCTION, PREVIOUS_PRODUCTION, CURRENT_STAGING) to include full algorithm instance details across returned lineages.",
			},
			"sort": schema.StringAttribute{
				Optional:    true,
				Description: "Sort criterion for returned lineages, e.g. -modifiedDate.",
				Validators: []validator.String{stringvalidator.OneOf(
					"-modifiedDate",
					"+modifiedDate",
					"-createdDate",
					"+createdDate",
					"-lineageName",
					"+lineageName",
					"-expirationDate",
					"+expirationDate",
				)},
			},
			"limit": schema.Int64Attribute{
				Optional: true,
				Description: "The maximum number of lineages to return. " +
					"When not specified, all lineages matching the filter criteria are returned.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"lineages": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of lineages matching the given filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: lineageItemNestedAttributes(),
				},
			},
		},
	}
}

// Read queries the API for lineages, iterating through paginated results, and populates the state.
func (d *lineagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Lineages DataSource Read")

	var data lineagesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	listReq := cloudcertificates.ListLineagesRequest{
		ContractID:    data.ContractID.ValueString(),
		LineageName:   data.LineageName.ValueString(),
		SecureNetwork: cloudcertificates.SecureNetwork(data.SecureNetwork.ValueString()),
		StackMode:     cloudcertificates.StackMode(data.StackMode.ValueString()),
		LineageType:   cloudcertificates.LineageType(data.LineageType.ValueString()),
		Domain:        data.Domain.ValueString(),
		KeyType:       cloudcertificates.CryptographicAlgorithm(data.KeyType.ValueString()),
		Issuer:        data.Issuer.ValueString(),
		PageSize:      cloudcertificates.MaxListLineagesPageSize,
	}

	if tf.IsKnown(data.ExpiringInDays) {
		expiringInDays := int(data.ExpiringInDays.ValueInt64())
		listReq.ExpiringInDays = ptr.To(expiringInDays)
	}

	if len(data.LineageIDs) > 0 {
		lineageIDs := make([]int64, len(data.LineageIDs))
		for i, lineageID := range data.LineageIDs {
			lineageIDs[i] = lineageID.ValueInt64()
		}
		slices.Sort(lineageIDs)
		listReq.LineageIDs = lineageIDs
	}

	if len(data.GenerationStatus) > 0 {
		statuses := make([]cloudcertificates.GenerationStatus, len(data.GenerationStatus))
		for i, s := range data.GenerationStatus {
			statuses[i] = cloudcertificates.GenerationStatus(s.ValueString())
		}
		slices.Sort(statuses)
		listReq.GenerationStatus = statuses
	}

	if tf.IsKnown(data.Sort) {
		listReq.Sort = data.Sort.ValueString()
	}

	if data.ExpandGenerations.ValueBool() {
		listReq.ExpandGenerations = []cloudcertificates.ExpandGenerations{
			cloudcertificates.ExpandGenerationsHead,
			cloudcertificates.ExpandGenerationsCurrentProduction,
			cloudcertificates.ExpandGenerationsPreviousProduction,
			cloudcertificates.ExpandGenerationsCurrentStaging,
		}
	}

	allLineages, err := d.listAllLineages(ctx, listReq, data.Limit.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError("Read Cloud Certificates Lineages failed", err.Error())
		return
	}

	lineageModels := make([]lineageItemModel, len(allLineages))
	for i, l := range allLineages {
		lineageModels[i] = mapLineageItemToModel(l)
	}
	data.Lineages = lineageModels

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// listAllLineages follows the API's keyset pagination cursor until the lineage list is exhausted or the
// requested limit is collected. A nil limit collects all lineages.
func (d *lineagesDataSource) listAllLineages(ctx context.Context, listReq cloudcertificates.ListLineagesRequest, limit *int64) ([]cloudcertificates.Lineage, error) {
	client := d.Client.GetCloudCertificates()

	var allLineages []cloudcertificates.Lineage
	for {
		if limit != nil {
			remaining := *limit - int64(len(allLineages))
			if remaining <= 0 {
				break
			}
			listReq.PageSize = int(min(remaining, int64(cloudcertificates.MaxListLineagesPageSize)))
		}

		res, err := client.ListLineages(ctx, listReq)
		if err != nil {
			return nil, err
		}

		pageLineages := res.Lineages
		if limit != nil {
			remaining := *limit - int64(len(allLineages))
			if int64(len(pageLineages)) > remaining {
				pageLineages = pageLineages[:int(remaining)]
			}
		}
		allLineages = append(allLineages, pageLineages...)

		if len(pageLineages) == 0 || res.NextCursor == nil || *res.NextCursor == "" {
			break
		}

		nextCursor := *res.NextCursor
		listReq.After = nextCursor
	}

	return allLineages, nil
}

func mapLineageItemToModel(lineage cloudcertificates.Lineage) lineageItemModel {
	return lineageItemModel{
		LineageID:          types.Int64Value(lineage.LineageID),
		commonLineageModel: mapCommonLineage(lineage),
	}
}

func lineageItemNestedAttributes() map[string]schema.Attribute {
	return mergeAttributes(lineageCommonAttributes(), map[string]schema.Attribute{
		"lineage_id": schema.Int64Attribute{
			Computed:    true,
			Description: "Unique identifier of the lineage.",
		},
	})
}
