package property

import (
	"context"

	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePropertyVariables() *schema.Resource {
	return &schema.Resource{
		CreateContext: resPropVarsOperation,
		ReadContext:   resPropVarsOperation,
		DeleteContext: resPropVarsOperation,
		UpdateContext: resPropVarsOperation,
		Schema: map[string]*schema.Schema{
			"json": {Type: schema.TypeString, Computed: true, Description: "JSON variables representation"},
			"variables": {
				Type:       schema.TypeSet,
				Optional:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("akamai_property_variables"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variable": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name":        {Type: schema.TypeString, Required: true, ValidateDiagFunc: tf.IsNotBlank},
									"hidden":      {Type: schema.TypeBool, Required: true},
									"sensitive":   {Type: schema.TypeBool, Required: true},
									"description": {Type: schema.TypeString, Optional: true},
									"value":       {Type: schema.TypeString, Optional: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resPropVarsOperation(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Always fail for new resources and changed values
	if d.IsNewResource() || d.HasChange("variables") {
		return diag.Errorf(`resource "akamai_property_variables" is no longer supported - See Akamai Terraform Upgrade Guide`)
	}

	// No changes and resource already exists, must be from previous version of plugin
	return nil
}
