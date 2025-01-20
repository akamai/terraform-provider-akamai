package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &domainsDataSource{}
var _ datasource.DataSourceWithConfigure = &domainsDataSource{}

// NewGTMDomainsDataSource returns a new GTM domains data source
func NewGTMDomainsDataSource() datasource.DataSource {
	return &domainsDataSource{}
}

var (
	domainsBlock = map[string]schema.Block{
		"domains": schema.SetNestedBlock{
			Description: "List of domains under given contract.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "A unique domain name.",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "The current status of the domain.",
					},
					"acg_id": schema.StringAttribute{
						Computed:    true,
						Description: "The contract's identifier, with which the domain is associated.",
					},
					"last_modified": schema.StringAttribute{
						Computed:    true,
						Description: "An ISO 8601 timestamp that indicates the time of the last domain change.",
					},
					"last_modified_by": schema.StringAttribute{
						Computed:    true,
						Description: "The email address of the administrator who made the last change to the domain.",
					},
					"change_id": schema.StringAttribute{
						Computed:    true,
						Description: "UUID that identifies a version of the domain configuration.",
					},
					"activation_state": schema.StringAttribute{
						Computed:    true,
						Description: "'PENDING' when a change has been made but not yet propagated; 'COMPLETE' when the last configuration change has propagated successfully; 'DENIED' if the domain configuration failed validation; 'DELETED' if the domain has been deleted.",
					},
					"modification_comments": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive note about changes to the domain.",
					},
					"sign_and_serve": schema.BoolAttribute{
						Computed:    true,
						Description: "If set (true) we will sign the domain's resource records so that they can be validated by a validating resolver.",
					},
					"sign_and_serve_algorithm": schema.StringAttribute{
						Computed:    true,
						Description: "The signing algorithm to use for signAndServe. One of the following values: RSA_SHA1, RSA_SHA256, RSA_SHA512, ECDSA_P256_SHA256, ECDSA_P384_SHA384, ED25519, ED448.",
					},
					"delete_request_id": schema.StringAttribute{
						Computed:    true,
						Description: "UUID for delete request during domain deletion. Null if the domain is not being deleted.",
					},
				},
				Blocks: map[string]schema.Block{
					"links": &schema.SetNestedBlock{
						Description: "Provides a URL path that allows direct navigation to the domain.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Computed:    true,
									Description: "Indicates the link relationship of the object.",
								},
								"href": schema.StringAttribute{
									Computed:    true,
									Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
								},
							},
						},
					},
				},
			},
		},
	}
)

type domainsDataSource struct {
	meta meta.Meta
}

type (
	domainsDataSourceModel struct {
		Domains []domain `tfsdk:"domains"`
	}

	domain struct {
		Name                  types.String `tfsdk:"name"`
		Status                types.String `tfsdk:"status"`
		AcgID                 types.String `tfsdk:"acg_id"`
		LastModified          types.String `tfsdk:"last_modified"`
		LastModifiedBy        types.String `tfsdk:"last_modified_by"`
		ChangeID              types.String `tfsdk:"change_id"`
		ActivationState       types.String `tfsdk:"activation_state"`
		ModificationComments  types.String `tfsdk:"modification_comments"`
		SignAndServe          types.Bool   `tfsdk:"sign_and_serve"`
		SignAndServeAlgorithm types.String `tfsdk:"sign_and_serve_algorithm"`
		DeleteRequestID       types.String `tfsdk:"delete_request_id"`
		Links                 []link       `tfsdk:"links"`
	}
)

// Schema is used to define data source's terraform schema
func (d *domainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "List of domains under given contract.",
		Blocks:              domainsBlock,
	}
}

// Configure configures data source at the beginning of the lifecycle
func (d *domainsDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta, ok := request.ProviderData.(meta.Meta)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
	d.meta = meta
}

// Metadata configures data source's meta information
func (d *domainsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_gtm_domains"
}

// Read is called when the provider must read data source values in order to update state
func (d *domainsDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM Domains DataSource Read")
	var data domainsDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	domains, err := client.ListDomains(ctx)
	if err != nil {
		response.Diagnostics.AddError("fetching domains failed", err.Error())
		return
	}

	data.Domains = getDomains(domains)
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func getDomains(domainList []gtm.DomainItem) []domain {
	var result []domain
	for _, dom := range domainList {
		domain := domain{
			Name:                  types.StringValue(dom.Name),
			LastModified:          types.StringValue(dom.LastModified),
			Status:                types.StringValue(dom.Status),
			AcgID:                 types.StringValue(dom.AcgID),
			LastModifiedBy:        types.StringValue(dom.LastModifiedBy),
			ChangeID:              types.StringValue(dom.ChangeID),
			ActivationState:       types.StringValue(dom.ActivationState),
			ModificationComments:  types.StringValue(dom.ModificationComments),
			SignAndServe:          types.BoolValue(dom.SignAndServe),
			SignAndServeAlgorithm: types.StringValue(dom.SignAndServeAlgorithm),
			DeleteRequestID:       types.StringValue(dom.DeleteRequestID),
			Links:                 getLinks(dom.Links),
		}

		result = append(result, domain)
	}
	return result
}
