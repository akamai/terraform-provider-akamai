package tf

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

// ResetAttrs resets (sets to nil) the provided set of attributes
func ResetAttrs(d *schema.ResourceData, attributes ...string) error {
	for _, attr := range attributes {
		if err := d.Set(attr, nil); err != nil {
			return err
		}
	}

	return nil
}
