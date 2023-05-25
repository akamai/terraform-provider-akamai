package property

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyIncludes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyIncludesRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the contract under which the data were requested",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the group under which the data were requested",
			},
			"parent_property": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The property's unique identifier",
						},
						"version": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The property's version for which the data is requested",
						},
					},
				},
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: tf.ValidateStringInSlice([]string{string(papi.IncludeTypeMicroServices), string(papi.IncludeTypeCommonSettings)}),
				Description:      "Specifies the type of the include, either `MICROSERVICES` or `COMMON_SETTINGS`. Use this field for filtering",
			},
			"includes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of includes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"latest_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the most recent version of the include",
						},
						"staging_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The most recent version to be activated to the staging network",
						},
						"production_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The most recent version to be activated to the production network",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The include's unique identifier",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A descriptive name for the include",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the type of the include, either `MICROSERVICES` or `COMMON_SETTINGS`",
						},
					},
				},
			},
		},
	}
}

type propertyIncludesAttrs struct {
	contractID     string
	groupID        string
	parentProperty *parentPropertyAttr
	includeType    string
}

type parentPropertyAttr struct {
	id      string
	version int
}

func dataPropertyIncludesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	log := meta.Log("PAPI", "dataPropertyIncludesRead")
	log.Debug("Reading property includes")

	attrs, err := getPropertyIncludesAttrs(d)
	if err != nil {
		return diag.Errorf("failed to read attributes: %s", err)
	}

	includes, err := sendListIncludesRequest(ctx, client, attrs)
	if err != nil {
		return diag.Errorf("sendListIncludesRequest error: %s", err)
	}

	includeAttrs := make([]interface{}, len(includes))
	for i, include := range includes {
		includeAttrs[i] = createIncludeAttrs(include)
	}

	err = d.Set("includes", includeAttrs)
	if err != nil {
		return diag.Errorf("could not set 'includes' attribute: %s", err)
	}
	d.SetId(attrs.createID())

	return nil
}

func sendListIncludesRequest(ctx context.Context, client papi.PAPI, attrs *propertyIncludesAttrs) ([]papi.Include, error) {
	if attrs.parentProperty != nil {
		availableIncludes, err := client.ListAvailableIncludes(ctx, papi.ListAvailableIncludesRequest{
			ContractID:      attrs.contractID,
			GroupID:         attrs.groupID,
			PropertyID:      attrs.parentProperty.id,
			PropertyVersion: attrs.parentProperty.version,
		})
		if err != nil {
			return nil, fmt.Errorf("could not list available includes: %s", err)
		}

		var includes []papi.Include
		for _, availableInclude := range availableIncludes.AvailableIncludes {
			include, err := client.GetInclude(ctx, papi.GetIncludeRequest{
				ContractID: attrs.contractID,
				GroupID:    attrs.groupID,
				IncludeID:  availableInclude.IncludeID})
			if err != nil {
				return nil, fmt.Errorf("could not get an include with ID: %s, %s", availableInclude.IncludeID, err)
			}

			if len(include.Includes.Items) != 0 {
				includes = append(includes, include.Includes.Items[0])
			}
		}

		return filterIncludes(includes, attrs.includeType), nil
	}

	includes, err := client.ListIncludes(ctx, papi.ListIncludesRequest{
		ContractID: attrs.contractID,
		GroupID:    attrs.groupID,
	})

	if err != nil {
		return nil, fmt.Errorf("could not list includes: %s", err)
	}

	return filterIncludes(includes.Includes.Items, attrs.includeType), nil
}

func getPropertyIncludesAttrs(d *schema.ResourceData) (*propertyIncludesAttrs, error) {
	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return nil, fmt.Errorf("could not get `contract_id` attribute: %s", err)
	}
	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return nil, fmt.Errorf("could not get `group_id` attribute: %s", err)
	}

	parentPropertyList, err := tf.GetListValue("parent_property", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, fmt.Errorf("could not get `parent_property` attribute: %s", err)
	}

	var parentPropertyValue *parentPropertyAttr
	if len(parentPropertyList) != 0 {
		parentPropertyMap, ok := parentPropertyList[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected map[string]interface{}, got: %T", parentPropertyMap)
		}
		propertyID, ok := parentPropertyMap["id"].(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got: %T", propertyID)
		}
		propertyVersion, ok := parentPropertyMap["version"].(int)
		if !ok {
			return nil, fmt.Errorf("expected int, got: %T", propertyID)
		}

		parentPropertyValue = &parentPropertyAttr{
			id:      propertyID,
			version: propertyVersion,
		}
	}

	includeType, err := tf.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, fmt.Errorf("could not get `type` attribute: %s", err)
	}

	return &propertyIncludesAttrs{
		contractID:     contractID,
		groupID:        groupID,
		parentProperty: parentPropertyValue,
		includeType:    includeType,
	}, nil
}

func createIncludeAttrs(include papi.Include) map[string]interface{} {
	attrs := map[string]interface{}{
		"latest_version": include.LatestVersion,
		"id":             include.IncludeID,
		"name":           include.IncludeName,
		"type":           include.IncludeType,
	}
	if include.StagingVersion != nil {
		attrs["staging_version"] = strconv.Itoa(*include.StagingVersion)
	}
	if include.ProductionVersion != nil {
		attrs["production_version"] = strconv.Itoa(*include.ProductionVersion)
	}

	return attrs
}

func filterIncludes(includes []papi.Include, includeType string) []papi.Include {
	if includeType == "" {
		return includes
	}

	var filteredIncludes []papi.Include
	for _, include := range includes {
		if string(include.IncludeType) == includeType {
			filteredIncludes = append(filteredIncludes, include)
		}
	}

	return filteredIncludes
}

func (attrs propertyIncludesAttrs) createID() string {
	idElements := []string{attrs.contractID, attrs.groupID}

	if attrs.includeType != "" {
		idElements = append(idElements, attrs.includeType)
	}

	if attrs.parentProperty != nil {
		idElements = append(idElements, attrs.parentProperty.id)
		idElements = append(idElements, strconv.Itoa(attrs.parentProperty.version))
	}

	return strings.Join(idElements, ":")
}
