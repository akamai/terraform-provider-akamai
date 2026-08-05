package cloudcertificates

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &certificatesActivityDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesActivityDataSource{}
)

type (
	certificatesActivityDataSource struct {
		meta.DataSource
	}

	certificatesActivityDataSourceModel struct {
		LineageID  types.Int64     `tfsdk:"lineage_id"`
		Limit      types.Int64     `tfsdk:"limit"`
		Activities []activityModel `tfsdk:"activities"`
	}

	activityModel struct {
		ActivityID   types.Int64  `tfsdk:"activity_id"`
		EventType    types.String `tfsdk:"event_type"`
		LineageID    types.Int64  `tfsdk:"lineage_id"`
		GenerationID types.Int64  `tfsdk:"generation_id"`
		CreatedBy    types.String `tfsdk:"created_by"`
		Network      types.String `tfsdk:"network"`
		Outcome      types.String `tfsdk:"outcome"`
		CreatedTime  types.String `tfsdk:"created_time"`
	}
)

// NewCertificatesActivityDataSource returns a new Cloud Certificates Activity data source.
func NewCertificatesActivityDataSource() datasource.DataSource {
	return &certificatesActivityDataSource{}
}

// Metadata configures data source's meta information.
func (d *certificatesActivityDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_activity"
}

// Schema is used to define data source's terraform schema.
func (d *certificatesActivityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for retrieving the activity history for a certificate lineage.",
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
				Description: "The maximum number of most recent activity events to return. " +
					"When not specified, all activity events for the lineage are returned.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"activities": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Activity events for the certificate lineage, ordered newest first.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"activity_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The unique identifier of the activity event.",
						},
						"event_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of operation recorded by the activity event.",
						},
						"lineage_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The unique identifier of the certificate lineage associated with the activity event.",
						},
						"generation_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The generation identifier associated with the activity event, or null for lineage-level events.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The user or system that triggered the activity event.",
						},
						"network": schema.StringAttribute{
							Computed:    true,
							Description: "The target network associated with the activity event, or null if not applicable.",
						},
						"outcome": schema.StringAttribute{
							Computed:    true,
							Description: "The execution result of the activity event, or null if not applicable.",
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "The time the activity event was recorded in UTC.",
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *certificatesActivityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Cloud Certificates Activity DataSource Read")

	var data certificatesActivityDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	events, err := d.listAllActivity(ctx, data.LineageID.ValueInt64(), data.Limit.ValueInt64Pointer())
	if err != nil {
		if errors.Is(err, cloudcertificates.ErrLineageNotFound) {
			resp.Diagnostics.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d", data.LineageID.ValueInt64()))
		} else {
			resp.Diagnostics.AddError("Failed to retrieve lineage activity", err.Error())
		}
		return
	}

	data.Activities = activitiesToModels(events)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// listAllActivity follows the API's keyset pagination cursor until the activity history is exhausted or the
// requested limit is collected. A nil limit collects all events.
func (d *certificatesActivityDataSource) listAllActivity(ctx context.Context, lineageID int64, limit *int64) ([]cloudcertificates.LineageActivityEvent, error) {
	client := d.Client.GetCloudCertificates()

	request := cloudcertificates.ListLineageActivityRequest{
		LineageID: lineageID,
		PageSize:  cloudcertificates.MaxListLineageActivityPageSize,
	}

	var events []cloudcertificates.LineageActivityEvent
	for {
		if limit != nil {
			remaining := int(*limit) - len(events)
			if remaining <= 0 {
				break
			}
			request.PageSize = min(remaining, cloudcertificates.MaxListLineageActivityPageSize)
		}

		resp, err := client.ListLineageActivity(ctx, request)
		if err != nil {
			return nil, err
		}

		pageEvents := resp.Events
		if limit != nil {
			remaining := int(*limit) - len(events)
			if len(pageEvents) > remaining {
				pageEvents = pageEvents[:remaining]
			}
		}
		events = append(events, pageEvents...)

		if resp.NextCursor == nil || len(resp.Events) == 0 {
			break
		}
		request.Cursor = *resp.NextCursor
	}

	return events, nil
}

func activitiesToModels(events []cloudcertificates.LineageActivityEvent) []activityModel {
	activities := make([]activityModel, 0, len(events))
	for _, e := range events {
		activities = append(activities, activityModel{
			ActivityID:   types.Int64Value(e.ActivityID),
			EventType:    types.StringValue(e.EventType),
			LineageID:    types.Int64Value(e.LineageID),
			GenerationID: types.Int64PointerValue(e.GenerationID),
			CreatedBy:    types.StringValue(e.CreatedBy),
			Network:      types.StringPointerValue(e.Network),
			Outcome:      types.StringPointerValue(e.Outcome),
			CreatedTime:  date.TimeRFC3339NanoValue(e.CreatedTime),
		})
	}
	return activities
}
