package appsec

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

//suppressJsonProvided to handle when json supplied vs HCL values
func suppressJsonProvided(_, old, new string, d *schema.ResourceData) bool {

	json := d.Get("json").(string)
	if json != "" {
		if old == "" && new == "" {
			return true
		}
		return true
	}

	return false
}
