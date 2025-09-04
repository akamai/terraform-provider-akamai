package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestEmailValidator(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		email     string
		expectErr bool
	}{
		{email: "user@example.com", expectErr: false},
		{email: "first.last@example.com", expectErr: false},
		{email: "user+tag@example.co.uk", expectErr: false},
		{email: "plainaddress", expectErr: true},
		{email: "@missinglocalpart.com", expectErr: true},
		{email: "missingatsign.com", expectErr: true},
		{email: "user@.com", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			emailValidator := EmailValidator{}

			// Prepare the request object with the test email
			req := validator.StringRequest{
				ConfigValue: types.StringValue(tt.email),
			}

			// Prepare the response object
			resp := validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}

			// Run the validator
			emailValidator.ValidateString(ctx, req, &resp)

			// Check the expected outcome
			if tt.expectErr {
				assert.True(t, resp.Diagnostics.HasError(), "Expected an error but got none")
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Expected no error but got one")
			}
		})
	}
}
