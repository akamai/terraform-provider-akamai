package property

import (
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func resourceProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyCreate,
		ReadContext:   resourcePropertyRead,
		UpdateContext: resourcePropertyUpdate,
		DeleteContext: resourcePropertyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePropertyImport,
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

			"latest_version":     {Type: schema.TypeInt, Computed: true},
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
	logger := akamai.Meta(m).Log("PAPI", "resourcePropertyCreate")
	client := inst.Client(akamai.Meta(m))
	ctx = log.NewContext(ctx, logger)

	if err := resPropAssertNoForbiddenAttr(d); err != nil {
		// User will also see messages saying which attributes are not supported
		return diag.Errorf("unsupported attributes given. See the Akamai Terraform Upgrade Guide")
	}

	// Schema guarantees name, contract_id, group_id, and product_id are strings
	PropertyName := d.Get("name").(string)

	GroupID := d.Get("group_id").(string)
	if GroupID == "" {
		GroupID = d.Get("group").(string)
	}
	GroupID = tools.AddPrefix(GroupID, "grp_")

	ContractID := d.Get("contract_id").(string)
	if ContractID == "" {
		ContractID = d.Get("contract").(string)
	}
	ContractID = tools.AddPrefix(ContractID, "ctr_")

	ProductID := d.Get("product_id").(string)
	if ProductID == "" {
		ProductID = d.Get("product").(string)
	}
	ProductID = tools.AddPrefix(ProductID, "prd_")

	PropertyID, err := createProperty(ctx, client, PropertyName, GroupID, ContractID, ProductID)
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := map[string]interface{}{
		"group_id":    GroupID,
		"group":       GroupID,
		"contract_id": ContractID,
		"contract":    ContractID,
		"product_id":  ProductID,
		"product":     ProductID,
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(PropertyID)
	return resourcePropertyRead(ctx, d, m)
}

func resourcePropertyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resPropHostnameRead"))
	client := inst.Client(akamai.Meta(m))

	if err := resPropAssertNoForbiddenAttr(d); err != nil {
		// User will also see messages saying which attributes are not supported
		return diag.Errorf("unsupported attributes given. See the Akamai Terraform Upgrade Guide")
	}

	// Schema guarantees property_id, group_id, and contract_id are strings
	PropertyID := d.Id()
	ContractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")
	GroupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")

	Property, err := fetchProperty(ctx, client, PropertyID, GroupID, ContractID)
	if err != nil {
		return diag.FromErr(err)
	}

	var StagingVersion int
	if Property.StagingVersion == nil {
		StagingVersion = 0
	} else {
		StagingVersion = *Property.StagingVersion
	}

	var ProductionVersion int
	if Property.ProductionVersion == nil {
		ProductionVersion = 0
	} else {
		ProductionVersion = *Property.ProductionVersion
	}

	attrs := map[string]interface{}{
		"name":               Property.PropertyName,
		"group_id":           Property.GroupID,
		"group":              Property.GroupID,
		"contract_id":        Property.ContractID,
		"contract":           Property.ContractID,
		"product_id":         Property.ProductID,
		"product":            Property.ProductID,
		"latest_version":     Property.LatestVersion,
		"staging_version":    StagingVersion,
		"production_version": ProductionVersion,
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePropertyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger := akamai.Meta(m).Log("PAPI", "resourcePropertyUpdate")

	if err := resPropAssertNoForbiddenAttr(d); err != nil {
		d.Partial(true)
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
	if diags.HasError() {
		d.Partial(true)
		return diags
	}

	// There are no attributes that can be updated by this resource

	diags = append(diags, resourcePropertyRead(ctx, d, m)...)
	if diags.HasError() {
		d.Partial(true)
	}
	return diags
}

func resourcePropertyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resourcePropertyDelete"))
	client := inst.Client(akamai.Meta(m))

	PropertyID := d.Id()
	ContractID := d.Get("contract_id").(string)
	GroupID := d.Get("group_id").(string)

	if err := removeProperty(ctx, client, PropertyID, GroupID, ContractID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePropertyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resourcePropertyImport"))

	// User-supplied import ID is a comma-separated list of PropertyID[,GroupID[,ContractID]]
	// ContractID and GroupID are optional as long as the PropertyID is sufficient to fetch the property
	var PropertyID, GroupID, ContractID string
	parts := strings.Split(d.Id(), ",")

	switch len(parts) {
	case 1:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
	case 2:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
		GroupID = tools.AddPrefix(parts[1], "grp_")
	case 3:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
		GroupID = tools.AddPrefix(parts[1], "grp_")
		ContractID = tools.AddPrefix(parts[2], "ctr_")
	default:
		return nil, fmt.Errorf("invalid property identifier: %q", d.Id())
	}

	// Import only needs to set the resource ID and enough attributes that the read opertaion will function, so there's
	// no need to fetch anything if the user gave both GroupID and ContractID
	if GroupID != "" && ContractID != "" {
		attrs := map[string]interface{}{
			"group_id":    GroupID,
			"contract_id": ContractID,
		}
		if err := rdSetAttrs(ctx, d, attrs); err != nil {
			return nil, err
		}

		d.SetId(PropertyID)
		return []*schema.ResourceData{d}, nil
	}

	// Missing GroupID, ContractID, or both -- Attempt to fetch them. If the PropertyID is not sufficient, PAPI
	// will return an error.
	Property, err := fetchProperty(ctx, inst.Client(akamai.Meta(m)), PropertyID, GroupID, ContractID)
	if err != nil {
		return nil, err
	}

	attrs := map[string]interface{}{
		"group_id":    Property.GroupID,
		"contract_id": Property.ContractID,
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return nil, err
	}

	d.SetId(Property.PropertyID)
	return []*schema.ResourceData{d}, nil
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

func createProperty(ctx context.Context, client papi.PAPI, PropertyName, GroupID, ContractID, ProductID string) (PropertyID string, err error) {
	req := papi.CreatePropertyRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
		Property: papi.PropertyCreate{
			ProductID:    ProductID,
			PropertyName: PropertyName,
		},
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("creating property")
	res, err := client.CreateProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not create property")
		return
	}
	PropertyID = res.PropertyID

	logger.WithFields(logFields(*res)).Info("property created")
	return
}

func removeProperty(ctx context.Context, client papi.PAPI, PropertyID, GroupID, ContractID string) error {
	req := papi.RemovePropertyRequest{
		PropertyID: PropertyID,
		GroupID:    GroupID,
		ContractID: ContractID,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))
	logger.Debug("removing property")

	_, err := client.RemoveProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not remove property")
		return err
	}

	logger.Info("property removed")

	return nil
}
