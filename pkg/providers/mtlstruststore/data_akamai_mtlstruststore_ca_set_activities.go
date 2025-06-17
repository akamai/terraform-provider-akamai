package mtlstruststore

import (
	"context"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &caSetActivitiesDataSource{}
	_ datasource.DataSourceWithConfigure = &caSetActivitiesDataSource{}
)

type (
	caSetActivitiesDataSource struct {
		meta meta.Meta
	}

	caSetActivitiesDataSourceModel struct {
		ID          types.String    `tfsdk:"id"`
		Name        types.String    `tfsdk:"name"`
		Start       types.String    `tfsdk:"start"`
		End         types.String    `tfsdk:"end"`
		CreatedDate types.String    `tfsdk:"created_date"`
		CreatedBy   types.String    `tfsdk:"created_by"`
		Status      types.String    `tfsdk:"status"`
		DeletedDate types.String    `tfsdk:"deleted_date"`
		DeletedBy   types.String    `tfsdk:"deleted_by"`
		Activities  []activityModel `tfsdk:"activities"`
	}

	activityModel struct {
		Type         types.String `tfsdk:"type"`
		Network      types.String `tfsdk:"network"`
		Version      types.Int64  `tfsdk:"version"`
		ActivityDate types.String `tfsdk:"activity_date"`
		ActivityBy   types.String `tfsdk:"activity_by"`
	}
)

// NewCASetActivitiesDataSource returns a new mtls truststore ca set activities data source.
func NewCASetActivitiesDataSource() datasource.DataSource {
	return &caSetActivitiesDataSource{}
}

// Metadata configures data source's meta information.
func (d *caSetActivitiesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set_activities"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *caSetActivitiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema is used to define data source's terraform schema.
func (d *caSetActivitiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve activities for a specific MTLS Truststore CA Set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifies each CA set. Either `id` or `name` must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the CA set. Either `id` or `name` must be provided.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"start": schema.StringAttribute{
				Description: "Filters out any activities after this time, expressed as an ISO 8601 timestamp. To specify a fixed time range, pair this with an 'end' parameter.",
				Optional:    true,
			},
			"end": schema.StringAttribute{
				Description: "Filters out any activities before this time, expressed as an ISO 8601 timestamp. To specify a fixed time range, pair this with a 'start' parameter.",
				Optional:    true,
			},
			"created_date": schema.StringAttribute{
				Description: "When the CA set was created.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "The user who created the CA set.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Indicates the status of the CA set.",
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
			"activities": schema.ListNestedAttribute{
				Description: "Activities performed on the CA set.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "The type of CA set activity. 'CREATE_CA_SET' indicates creating a CA set, or 'CREATE_CA_SET_VERSION' for creating a version. " +
								"'ACTIVATE_CA_SET_VERSION' indicates activating a CA set version, while 'DEACTIVATE_CA_SET_VERSION' indicates deactivation. " +
								"'DELETE_CA_SET' indicates deleting a CA set.",
							Computed: true,
						},
						"network": schema.StringAttribute{
							Description: "Indicates the network for any activation-related activities, either 'STAGING' or 'PRODUCTION'.",
							Computed:    true,
						},
						"version": schema.Int64Attribute{
							Description: "The CA set's incremental version number.",
							Computed:    true,
						},
						"activity_date": schema.StringAttribute{
							Description: "When this CA set activity occurred.",
							Computed:    true,
						},
						"activity_by": schema.StringAttribute{
							Description: "The user who initiated this CA set activity.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *caSetActivitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "MTLS TrustStore CA Set Activities DataSource Read")

	var data caSetActivitiesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client = Client(d.meta)

	if !data.Name.IsNull() {
		tflog.Debug(ctx, "'name' provided, attempting to find CA set ID")
		setID, err := findCASetID(ctx, client, data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Read CA set activities failed", err.Error())
			return
		}

		data.ID = types.StringValue(setID)
	}

	activities, err := data.getActivities(ctx, client)
	if err != nil {
		resp.Diagnostics.AddError("Read CA set activities failed", err.Error())
		return
	}

	modelData := convertDataToModel(*activities)
	data.setData(modelData)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *caSetActivitiesDataSourceModel) getActivities(ctx context.Context, client mtlstruststore.MTLSTruststore) (*mtlstruststore.ListCASetActivitiesResponse, error) {
	var start, end time.Time
	var err error
	if m.Start.ValueString() != "" {
		if start, err = time.Parse(time.RFC3339, m.Start.ValueString()); err != nil {
			return nil, fmt.Errorf("invalid start time format: %w", err)
		}
	}
	if m.End.ValueString() != "" {
		if end, err = time.Parse(time.RFC3339, m.End.ValueString()); err != nil {
			return nil, fmt.Errorf("invalid end time format: %w", err)
		}
	}

	activities, err := client.ListCASetActivities(ctx, mtlstruststore.ListCASetActivitiesRequest{
		CASetID: m.ID.ValueString(),
		Start:   start,
		End:     end,
	})
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func convertDataToModel(activities mtlstruststore.ListCASetActivitiesResponse) caSetActivitiesDataSourceModel {
	data := caSetActivitiesDataSourceModel{
		ID:          types.StringValue(activities.CASetID),
		Name:        types.StringValue(activities.CASetName),
		CreatedDate: types.StringValue(activities.CreatedDate.String()),
		CreatedBy:   types.StringValue(activities.CreatedBy),
		Status:      types.StringValue(activities.CASetStatus),
		DeletedBy:   types.StringPointerValue(activities.DeletedBy),
	}

	if activities.DeletedDate != nil {
		data.DeletedDate = types.StringValue(activities.DeletedDate.String())
	}

	activitiesModel := make([]activityModel, len(activities.Activities))
	for i, activity := range activities.Activities {
		am := activityModel{
			Type:         types.StringValue(activity.Type),
			ActivityDate: types.StringValue(activity.ActivityDate.String()),
			ActivityBy:   types.StringValue(activity.ActivityBy),
			Network:      types.StringPointerValue(activity.Network),
			Version:      types.Int64PointerValue(activity.Version),
		}
		activitiesModel[i] = am
	}

	data.Activities = activitiesModel

	return data
}

func (m *caSetActivitiesDataSourceModel) setData(activities caSetActivitiesDataSourceModel) {
	m.ID = activities.ID
	m.Name = activities.Name
	m.CreatedDate = activities.CreatedDate
	m.CreatedBy = activities.CreatedBy
	m.Status = activities.Status
	m.DeletedDate = activities.DeletedDate
	m.DeletedBy = activities.DeletedBy
	m.Activities = activities.Activities
}

func findCASetID(ctx context.Context, client mtlstruststore.MTLSTruststore, caSetName string) (string, error) {
	caSets, err := client.ListCASets(ctx, mtlstruststore.ListCASetsRequest{
		CASetNamePrefix: caSetName,
	})
	if err != nil {
		return "", fmt.Errorf("could not find CA Set ID for the given CA Set Name '%s', API error: %w", caSetName, err)
	}

	var matchingSets []mtlstruststore.CASetResponse
	for _, caSet := range caSets.CASets {
		if caSet.CASetName == caSetName {
			matchingSets = append(matchingSets, caSet)
		}
	}

	switch len(matchingSets) {
	case 0:
		return "", fmt.Errorf("no CA set found with name '%s'", caSetName)
	case 1:
		return matchingSets[0].CASetID, nil
	default:
		return "", fmt.Errorf("multiple CA sets IDs found with name '%s': %v. Use the ID to fetch a specific CA set", caSetName, buildSetsMap(matchingSets))
	}
}

func findNotDeletedCASetID(ctx context.Context, client mtlstruststore.MTLSTruststore, caSetName string) (string, error) {
	caSets, err := client.ListCASets(ctx, mtlstruststore.ListCASetsRequest{
		CASetNamePrefix: caSetName,
	})
	if err != nil {
		return "", fmt.Errorf("could not find CA Set ID for the given CA Set Name '%s', API error: %w", caSetName, err)
	}

	var matchingSets []mtlstruststore.CASetResponse
	for _, caSet := range caSets.CASets {
		if caSet.CASetName == caSetName && caSet.CASetStatus == "NOT_DELETED" {
			matchingSets = append(matchingSets, caSet)
		}
	}

	switch len(matchingSets) {
	case 0:
		return "", fmt.Errorf("no CA set found with name '%s' and status 'NOT_DELETED'", caSetName)
	case 1:
		return matchingSets[0].CASetID, nil
	default:
		return "", fmt.Errorf("multiple CA sets IDs found with name '%s' and status 'NOT_DELETED': %v. Use the ID to fetch a specific CA set", caSetName, buildSetsMap(matchingSets))
	}
}

func buildSetsMap(sets []mtlstruststore.CASetResponse) map[string]string {
	m := make(map[string]string)
	for _, set := range sets {
		m[set.CASetID] = set.CASetStatus
	}
	return m
}
