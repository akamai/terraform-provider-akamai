package gtm

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceInstance struct {
	DataCenterID         types.Int64    `tfsdk:"datacenter_id"`
	UseDefaultLoadObject types.Bool     `tfsdk:"use_default_load_object"`
	LoadObject           types.String   `tfsdk:"load_object"`
	LoadObjectPort       types.Int64    `tfsdk:"load_object_port"`
	LoadServers          []types.String `tfsdk:"load_servers"`
}

type link struct {
	Rel  types.String `tfsdk:"rel"`
	Href types.String `tfsdk:"href"`
}
