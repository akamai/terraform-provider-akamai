package tools

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetAttrs allows you to set many attributes of a schema.ResourceData in one call
func SetAttrs(d *schema.ResourceData, AttributeValues map[string]interface{}) error {
	for attr, value := range AttributeValues {
		if err := d.Set(attr, value); err != nil {
			return err
		}
	}

	return nil
}
