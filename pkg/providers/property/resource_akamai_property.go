package property

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
)

func resourceProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyCreate,
		ReadContext:   resourcePropertyRead,
		UpdateContext: resourcePropertyUpdate,
		DeleteContext: resourcePropertyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    resourcePropertyV0().CoreConfigSchema().ImpliedType(),
			Upgrade: upgradePropV0,
		}},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"contract": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: `use "contract_id" attribute instead`,
				StateFunc:  addPrefixToState("ctr_"),
			},
			"contract_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"contract_id", "contract"},
				StateFunc:    addPrefixToState("ctr_"),
			},

			"group": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: `use "group_id" attribute instead`,
				StateFunc:  addPrefixToState("grp_"),
			},
			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group_id", "group"},
				StateFunc:    addPrefixToState("grp_"),
			},

			"product": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: `use "product_id" attribute instead`,
				StateFunc:  addPrefixToState("prd_"),
			},
			"product_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"product_id", "product"},
				StateFunc:    addPrefixToState("prd_"),
			},

			"staging_version":    {Type: schema.TypeInt, Computed: true},
			"production_version": {Type: schema.TypeInt, Computed: true},

			// Hard-deprecated attributes: These are effectively removed, but we wanted to refer users to the upgrade guide
			"rule_format": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: `"rule_format" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"cp_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: `"cp_code" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"contact": {
				Type:       schema.TypeSet,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: `"contact" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"hostnames": {
				Type:       schema.TypeMap,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: `"hostnames" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"origin": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname":              {Type: schema.TypeString, Optional: true},
						"port":                  {Type: schema.TypeInt, Optional: true},
						"forward_hostname":      {Type: schema.TypeString, Optional: true},
						"cache_key_hostname":    {Type: schema.TypeString, Optional: true},
						"compress":              {Type: schema.TypeBool, Optional: true},
						"enable_true_client_ip": {Type: schema.TypeBool, Optional: true},
					},
				},
				Deprecated: `"origin" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"is_secure": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: `"is_secure" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"rules": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: `"rules" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"variables": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: `"variables" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
		},
	}
}

func resourcePropertyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("PAPI", "resourcePropertyCreate")

	if err := resPropAssertNoForbiddenAttr(d); err != nil {
		// User will also see messages saying which attributes are not supported
		return diag.Errorf("unsupported attributes given. See the Akamai Terraform Upgrade Guide")
	}

	// Schema guarantees contract_id/contract are strings and one or the other is set
	var ContractID string
	if got, ok := d.GetOk("contract_id"); ok {
		ContractID = got.(string)
	} else {
		ContractID = d.Get("contract").(string)
	}
	if !strings.HasPrefix(ContractID, "ctr_") {
		ContractID = fmt.Sprintf("ctr_%s", ContractID)
	}

	// Schema guarantees group_id/group are strings and one or the other is set
	var GroupID string
	if got, ok := d.GetOk("group_id"); ok {
		GroupID = got.(string)
	} else {
		GroupID = d.Get("group").(string)
	}
	if !strings.HasPrefix(GroupID, "grp_") {
		GroupID = fmt.Sprintf("grp_%s", GroupID)
	}

	// Schema guarantees product_id/product are strings and one or the other is set
	var ProductID string
	if got, ok := d.GetOk("product_id"); ok {
		ProductID = got.(string)
	} else {
		ProductID = d.Get("product").(string)
	}
	if !strings.HasPrefix(ProductID, "prd_") {
		ProductID = fmt.Sprintf("prd_%s", ProductID)
	}

	// Schema guarantees name is a string and is present
	PropertyName := d.Get("name").(string)

	req := papi.CreatePropertyRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
		Property: papi.PropertyCreate{
			ProductID:    ProductID,
			PropertyName: PropertyName,
		},
	}

	logger = logger.WithFields(logFields(req))
	logger.Debug("creating property")

	res, err := client.CreateProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not create property")
		return diag.FromErr(err)
	}

	if !strings.HasPrefix(res.PropertyID, "prp_") {
		res.PropertyID = fmt.Sprintf("prp_%s", res.PropertyID)
	}

	logger.WithFields(logFields(*res)).Info("property created")

	d.SetId(res.PropertyID)

	return resourcePropertyRead(ctx, d, m)
}

func resourcePropertyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("PAPI", "resourcePropertyRead")

	if err := resPropAssertNoForbiddenAttr(d); err != nil {
		// User will also see messages saying which attributes are not supported
		return diag.Errorf("unsupported attributes given. See the Akamai Terraform Upgrade Guide")
	}

	// PropertyID could be un-prefixed in the case of imports
	PropertyID := d.Id()
	if !strings.HasPrefix(PropertyID, "prp_") {
		PropertyID = fmt.Sprintf("prp_%s", PropertyID)
	}
	d.SetId(PropertyID)

	req := papi.GetPropertyRequest{PropertyID: PropertyID}
	logger = logger.WithFields(logFields(req))

	res, err := client.GetProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not read property")
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	prop := res.Property
	if err := d.Set("name", prop.PropertyName); err != nil {
		logger.WithError(err).Error(`could not set "name" attribute`)
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("contract_id", prop.ContractID); err != nil {
		logger.WithError(err).Error(`could not set "contract_id" attribute`)
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("group_id", prop.GroupID); err != nil {
		logger.WithError(err).Error(`could not set "group_id" attribute`)
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("product_id", prop.ProductID); err != nil {
		logger.WithError(err).Error(`could not set "product_id" attribute`)
		diags = append(diags, diag.FromErr(err)...)
	}

	if prop.StagingVersion != nil && *prop.StagingVersion > 0 {
		if err := d.Set("staging_version", *prop.StagingVersion); err != nil {
			logger.WithError(err).Error(`could not set "staging_version" attribute`)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if prop.ProductionVersion != nil && *prop.ProductionVersion > 0 {
		if err := d.Set("production_version", *prop.ProductionVersion); err != nil {
			logger.WithError(err).Error(`could not set "production_version" attribute`)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

func resourcePropertyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyUpdate")

	if err := resPropAssertNoForbiddenAttr(d); err != nil {
		// User will also see messages saying which attributes are not supported
		return diag.Errorf("unsupported attributes given. See the Akamai Terraform Upgrade Guide")
	}

	var diags diag.Diagnostics

	immutable := []string{
		"contract_id", "contract",
		"group_id", "group",
		"product_id", "product",
	}

	for _, attr := range immutable {
		if d.HasChange(attr) {
			err := fmt.Errorf(`property attribute %q cannot be changed after creation (immutable)`, attr)
			logger.WithError(err).Error("could not update property")
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// There are no attributes that can be updated by this resource

	diags = append(diags, resourcePropertyRead(ctx, d, m)...)
	return diags
}

func resourcePropertyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyDelete")
	client := inst.Client(meta)

	// Schema guarantees one of contract_id/contract are strings and one or the other is set
	var ContractID string
	if got, ok := d.GetOk("contract_id"); ok {
		ContractID = got.(string)
	} else {
		ContractID = d.Get("contract").(string)
	}

	// Schema guarantees one of group_id/group will be set and that they're strings
	var GroupID string
	if got, ok := d.GetOk("group_id"); ok {
		GroupID = got.(string)
	} else {
		GroupID = d.Get("group").(string)
	}

	req := papi.RemovePropertyRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
		PropertyID: d.Id(),
	}

	logger = logger.WithFields(logFields(req))
	logger.Debug("removing property")

	_, err := client.RemoveProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not remove property")
		return diag.FromErr(err)
	}

	logger.Info("property removed")
	return nil
}

// Returns error when any hard-deprecated attributes contain non-zero values
func resPropAssertNoForbiddenAttr(d *schema.ResourceData) error {
	deprecated := []string{
		"rule_format",
		"cp_code",
		"contact",
		"hostnames",
		"origin",
		"is_secure",
		"rules",
		"variables",
	}
	for _, attr := range deprecated {
		if _, ok := d.GetOk(attr); ok {
			return fmt.Errorf("unsupported attribute: %q", attr)
		}
	}

	return nil
}
