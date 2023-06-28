package akamai

import (
	"context"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type contextConfig struct {
	edgercPath    string
	edgercSection string
	edgercConfig  *edgegrid.Config
	userAgent     string
	ctx           context.Context
	requestLimit  int
	enableCache   bool
}

func configureContext(cfg contextConfig) (*meta.OperationMeta, error) {
	operationID := uuid.NewString()
	log := logger.FromContext(cfg.ctx, "OperationID", operationID)

	edgerc, err := newEdgegridConfig(cfg.edgercPath, cfg.edgercSection, cfg.edgercConfig)
	if err != nil {
		return nil, err
	}

	sess, err := session.New(
		session.WithSigner(edgerc),
		session.WithUserAgent(cfg.userAgent),
		session.WithLog(log),
		session.WithHTTPTracing(cast.ToBool(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED"))),
		session.WithRequestLimit(cfg.requestLimit),
	)
	if err != nil {
		return nil, err
	}

	cache.Enable(cfg.enableCache)

	return meta.New(sess, log.HCLog(), operationID)
}
