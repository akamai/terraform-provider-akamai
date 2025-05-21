package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
)

func dataSourceProperty() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of property.",
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The current version of the property.",
			},
			"contract_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contract ID assigned to the property.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Group ID assigned to the property.",
			},
			"latest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Property's current latest version.",
			},
			"note": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The client property notes.",
			},
			"production_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Property's version currently activated in production (zero when not active in production).",
			},
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product ID assigned to the property.",
			},
			"property_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the property.",
			},
			"rules": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property rules as JSON.",
			},
			"rule_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule format version.",
			},
			"staging_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Property's version currently activated in staging (zero when not active in staging).",
			},
			"asset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the property in the Identity and Access Management API.",
			},
			"property_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the type of the property.",
			},
		},
	}
}

func dataPropertyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	log := meta.Log("PAPI", "dataPropertyRead")
	log.Debug("Reading Property")

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	prop, err := findProperty(ctx, name, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tf.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		prop.LatestVersion = version
	}

	rules, err := getRulesForProperty(ctx, prop, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	body, err := json.Marshal(rules)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rules", string(body)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %q", tf.ErrValueSet, err.Error()))
	}

	propVersion, err := getPropertyVersion(ctx, meta, prop)
	if err != nil {
		return diag.FromErr(err)
	}

	propertyAttr := getPropertyAttributes(prop, propVersion)
	err = tf.SetAttrs(d, propertyAttr)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %q", tf.ErrValueSet, err.Error()))
	}

	d.SetId(prop.PropertyID)
	return nil
}

func getPropertyVersion(ctx context.Context, meta meta.Meta, property *papi.Property) (*papi.GetPropertyVersionsResponse, error) {
	client := Client(meta)
	req := papi.GetPropertyVersionRequest{
		PropertyID:      property.PropertyID,
		PropertyVersion: property.LatestVersion,
		ContractID:      property.ContractID,
		GroupID:         property.GroupID,
	}
	resp, err := client.GetPropertyVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrPropertyVersionNotFound, err.Error())
	}
	return resp, nil
}

func getRulesForProperty(ctx context.Context, property *papi.Property, meta meta.Meta) (*papi.GetRuleTreeResponse, error) {
	client := Client(meta)
	req := papi.GetRuleTreeRequest{
		PropertyID:      property.PropertyID,
		PropertyVersion: property.LatestVersion,
		ContractID:      property.ContractID,
		GroupID:         property.GroupID,
	}
	rules, err := client.GetRuleTree(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRulesNotFound, err.Error())
	}
	return rules, nil
}

func getPropertyAttributes(propertyResponse *papi.Property, propertyVersionResponse *papi.GetPropertyVersionsResponse) map[string]interface{} {
	propertyVersion := propertyVersionResponse.Version
	property := map[string]interface{}{
		"asset_id":           propertyResponse.AssetID,
		"contract_id":        propertyResponse.ContractID,
		"group_id":           propertyResponse.GroupID,
		"latest_version":     propertyResponse.LatestVersion,
		"note":               propertyVersion.Note,
		"product_id":         propertyVersion.ProductID,
		"production_version": decodeVersion(propertyResponse.ProductionVersion),
		"property_id":        propertyResponse.PropertyID,
		"rule_format":        propertyVersion.RuleFormat,
		"staging_version":    decodeVersion(propertyResponse.StagingVersion),
	}

	if propertyResponse.PropertyType != nil {
		property["property_type"] = *propertyResponse.PropertyType
	} else {
		property["property_type"] = ""
	}
	return property
}
