package iam

import (
	"testing"
)

func TestProviderOld(t *testing.T) {
	t.Parallel()

	t.Run("Validate provider schema", func(t *testing.T) {
		t.Parallel()

		p := providerOld{}

		if err := p.ProviderSchema().InternalValidate(); err != nil {
			t.Fatal(err)
		}
	})
}
