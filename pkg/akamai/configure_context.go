package akamai

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/google/uuid"
	"github.com/mgwoj/go-retryablehttp"
	"github.com/spf13/cast"
)

type contextConfig struct {
	edgegridConfig *edgegrid.Config
	userAgent      string
	ctx            context.Context
	requestLimit   int
	enableCache    bool
}

func configureContext(cfg contextConfig) (*meta.OperationMeta, error) {
	operationID := uuid.NewString()
	log := logger.FromContext(cfg.ctx, "OperationID", operationID)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.RetryWaitMin = 10 * time.Second
	retryClient.RetryWaitMax = 5 * time.Minute

	sess, err := session.New(
		session.WithSigner(cfg.edgegridConfig),
		session.WithUserAgent(cfg.userAgent),
		session.WithLog(log),
		session.WithHTTPTracing(cast.ToBool(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED"))),
		session.WithRequestLimit(cfg.requestLimit),
		session.WithClient(retryClient.StandardClient()),
	)
	if err != nil {
		return nil, err
	}

	retryClient.Postprocess = func(r *http.Request) error {
		return sess.Sign(r)
	}
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		var urlErr *url.Error
		if resp != nil && resp.Request.Method == "GET" ||
			resp == nil && errors.As(err, &urlErr) && strings.ToUpper(urlErr.Op) == "GET" {
			return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
		}
		return false, nil
	}

	cache.Enable(cfg.enableCache)

	return meta.New(sess, log.HCLog(), operationID)
}
