package iam

import (
	"testing"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	t.Run("Validate provider schema", func(t *testing.T) {
		t.Parallel()

		p := provider{}

		if err := p.ProviderSchema().InternalValidate(); err != nil {
			t.Fatal(err)
		}
	})
}
