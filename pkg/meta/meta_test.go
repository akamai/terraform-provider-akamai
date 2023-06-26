package meta

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/apex/log"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeta(t *testing.T) {
	var sess = session.Must(session.New())
	var logger = hclog.New(hclog.DefaultOptions)
	operationID := "opID"

	meta, err := New(sess, logger, operationID)
	assert.NoError(t, err)

	t.Run("Session() return sess", func(t *testing.T) {
		assert.Equal(t, sess, meta.Session())
	})
	t.Run("OperationID() return operationID", func(t *testing.T) {
		assert.Equal(t, operationID, meta.OperationID())
	})
}

func TestNew_err(t *testing.T) {
	var sess = session.Must(session.New())
	var logger = hclog.New(hclog.DefaultOptions)

	t.Run("nil log", func(t *testing.T) {
		_, err := New(sess, nil, "")
		assert.Error(t, err, ErrNilLog)
	})
	t.Run("nil session", func(t *testing.T) {
		_, err := New(nil, logger, "")
		assert.Error(t, err, ErrNilSession)
	})

}

func TestMust(t *testing.T) {
	t.Run("no panic", func(t *testing.T) {
		var sess = session.Must(session.New())
		var logger = hclog.New(hclog.DefaultOptions)

		meta, err := New(sess, logger, "")
		require.NoError(t, err)

		var m interface{} = meta
		assert.NotPanics(t, func() {
			_ = Must(m)
		})
	})

	t.Run("panics for non meta", func(t *testing.T) {
		var m interface{} = "not_meta"
		assert.Panics(t, func() {
			_ = Must(m)
		})
	})
}
