package property

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
)

func resourcePropertyBootstrap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyBootstrapCreate,
		ReadContext:   resourcePropertyBootstrapRead,
		UpdateContext: resourcePropertyBootstrapUpdate,
		DeleteContext: resourcePropertyBootstrapDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePropertyBootstrapImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validatePropertyName,
				Description:      "Name to give to the Property (must be unique)",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				StateFunc:   addPrefixToState("grp_"),
				Description: "Group ID to be assigned to the Property",
			},
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				StateFunc:   addPrefixToState("ctr_"),
				Description: "Contract ID to be assigned to the Property",
			},
			"product_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product ID to be assigned to the Property",
				StateFunc:   addPrefixToState("prd_"),
			},
		},
	}
}

func resourcePropertyBootstrapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyBootstrapCreate")
	client := Client(meta)
	ctx = log.NewContext(ctx, logger)

	// Schema guarantees these types
	propertyName := d.Get("name").(string)
	contractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")
	groupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")
	productID := tools.AddPrefix(d.Get("product_id").(string), "prd_")

	// we use default rule format (probably v2022-
	propertyID, err := createProperty(ctx, client, propertyName, groupID, contractID, productID, "")
	if err != nil {
		return interpretCreatePropertyError(ctx, err, meta, groupID, contractID, productID)
	}
	d.SetId(propertyID)

	return resourcePropertyBootstrapRead(ctx, d, m)
}

func resourcePropertyBootstrapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, meta.Must(m).Log("PAPI", "resourcePropertyBootstrapRead"))
	logger := log.FromContext(ctx)
	client := Client(meta.Must(m))

	// Schema guarantees group_id, and contract_id are strings
	propertyID := d.Id()
	contractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")
	groupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")

	_, err := fetchLatestProperty(ctx, client, propertyID, groupID, contractID)
	if errors.Is(err, tf.ErrNotFound) {
		logger.Warnf("property %q removed on server. Removing from local state", propertyID)

		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePropertyBootstrapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, meta.Must(m).Log("PAPI", "resourcePropertyBootstrapUpdate"))
	logger := log.FromContext(ctx)

	diags := diag.Diagnostics{}

	immutable := []string{
		"group_id",
		"contract_id",
		"product_id",
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

	return resourcePropertyBootstrapRead(ctx, d, m)
}

func resourcePropertyBootstrapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, meta.Must(m).Log("PAPI", "resourcePropertyBootstrapDelete"))
	client := Client(meta.Must(m))

	propertyID := d.Id()
	contractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")
	groupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")

	if err := removeProperty(ctx, client, propertyID, groupID, contractID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePropertyBootstrapImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ctx = log.NewContext(ctx, meta.Must(m).Log("PAPI", "resourcePropertyImport"))

	// User-supplied import ID is a comma-separated list of propertyID[,groupID,contractID]
	// contractID and groupID are optional as long as the propertyID is sufficient to fetch the property
	var propertyID, groupID, contractID string
	parts := strings.Split(d.Id(), ",")
	switch len(parts) {
	case 3:
		propertyID = tools.AddPrefix(parts[0], "prp_")
		contractID = tools.AddPrefix(parts[1], "ctr_")
		groupID = tools.AddPrefix(parts[2], "grp_")
	case 2:
		return nil, fmt.Errorf("missing group id or contract id")
	case 1:
		propertyID = tools.AddPrefix(parts[0], "prp_")

	default:
		return nil, fmt.Errorf("invalid property identifier: %q", d.Id())
	}

	client := Client(meta.Must(m))
	property, err := fetchLatestProperty(ctx, client, propertyID, groupID, contractID)
	if err != nil {
		return nil, err
	}

	res, err := fetchPropertyVersion(ctx, client, property.PropertyID, property.GroupID, property.ContractID, property.LatestVersion)
	if err != nil {
		return nil, err
	}
	property.ProductID = res.Version.ProductID

	attrs := map[string]interface{}{
		"name":        property.PropertyName,
		"group_id":    property.GroupID,
		"contract_id": property.ContractID,
		"product_id":  property.ProductID,
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return nil, err
	}

	d.SetId(property.PropertyID)
	return []*schema.ResourceData{d}, nil
}
