package tools

import (
	"context"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Set many attributes of a schema.ResourceData in one call
func SetBatch(ctx context.Context, d *schema.ResourceData, AttributeValues map[string]interface{}) error {
	logger := log.FromContext(ctx)

	for attr, value := range AttributeValues {
		if err := d.Set(attr, value); err != nil {
			logger.WithError(err).Errorf("could not set %q", attr)
			return err
		}
	}

	return nil
}
