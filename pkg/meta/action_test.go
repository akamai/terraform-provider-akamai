package meta

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/hashicorp/go-hclog"
	frameworkaction "github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionConfigure(t *testing.T) {
	t.Parallel()

	t.Run("configures client from provider data", func(t *testing.T) {
		t.Parallel()

		providerMeta, err := New(session.Must(session.New()), hclog.NewNullLogger(), "operation-id")
		require.NoError(t, err)

		configured := &Action{}
		resp := &frameworkaction.ConfigureResponse{}
		configured.Configure(context.Background(), frameworkaction.ConfigureRequest{ProviderData: providerMeta}, resp)

		assert.False(t, resp.Diagnostics.HasError())
		assert.Same(t, providerMeta.Client(), configured.Client)
	})

	t.Run("ignores nil provider data", func(t *testing.T) {
		t.Parallel()

		configured := &Action{}
		resp := &frameworkaction.ConfigureResponse{}
		configured.Configure(context.Background(), frameworkaction.ConfigureRequest{}, resp)

		assert.False(t, resp.Diagnostics.HasError())
		assert.Nil(t, configured.Client)
	})

	t.Run("reports invalid provider data", func(t *testing.T) {
		t.Parallel()

		configured := &Action{}
		resp := &frameworkaction.ConfigureResponse{}
		configured.Configure(context.Background(), frameworkaction.ConfigureRequest{ProviderData: "invalid"}, resp)

		assert.True(t, resp.Diagnostics.HasError())
		assert.Len(t, resp.Diagnostics, 1)
		assert.Nil(t, configured.Client)
	})
}
