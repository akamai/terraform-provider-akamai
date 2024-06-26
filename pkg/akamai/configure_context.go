package akamai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/retryablehttp"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type contextConfig struct {
	edgegridConfig *edgegrid.Config
	userAgent      string
	ctx            context.Context
	requestLimit   int
	enableCache    bool
	retryMax       int
	retryWaitMin   time.Duration
	retryWaitMax   time.Duration
	retryDisabled  bool
}

func configureContext(cfg contextConfig) (*meta.OperationMeta, error) {
	operationID := uuid.NewString()
	log := logger.FromContext(cfg.ctx, "OperationID", operationID)

	opts := []session.Option{
		session.WithSigner(cfg.edgegridConfig),
		session.WithUserAgent(cfg.userAgent),
		session.WithLog(log),
		session.WithHTTPTracing(cast.ToBool(os.Getenv("AKAMAI_HTTP_TRACE_ENABLED"))),
		session.WithRequestLimit(cfg.requestLimit),
	}
	var sess session.Session
	var err error
	if cfg.retryDisabled {
		sess, err = sessionWithoutRetry(opts)
	} else {
		sess, err = sessionWithRetry(cfg, opts)
	}
	if err != nil {
		return nil, err
	}
	cache.Enable(cfg.enableCache)

	return meta.New(sess, log.HCLog(), operationID)
}

func sessionWithoutRetry(opts []session.Option) (session.Session, error) {
	return session.New(opts...)
}

func sessionWithRetry(cfg contextConfig, opts []session.Option) (session.Session, error) {
	if cfg.retryMax == 0 {
		cfg.retryMax = 10
	}
	if cfg.retryWaitMin == 0 {
		cfg.retryWaitMin = time.Duration(1) * time.Second
	}
	if cfg.retryWaitMax == 0 {
		cfg.retryWaitMax = time.Duration(30) * time.Second
	}

	err := validateRetryConfiguration(cfg)
	if err != nil {
		return nil, err
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = cfg.retryMax
	retryClient.RetryWaitMin = cfg.retryWaitMin
	retryClient.RetryWaitMax = cfg.retryWaitMax

	opts = append(opts, session.WithClient(retryClient.StandardClient()))
	sess, err := session.New(opts...)
	if err != nil {
		return nil, err
	}

	retryClient.PrepareRetry = func(r *http.Request) error {
		return sess.Sign(r)
	}

	retryClient.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return sess.Sign(req)
	}

	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		var urlErr *url.Error
		if (resp != nil && resp.Request.Method == http.MethodGet) ||
			(resp == nil && errors.As(err, &urlErr) && strings.ToUpper(urlErr.Op) == http.MethodGet) {
			if ctx.Err() != nil {
				return false, ctx.Err()
			}
			if resp != nil && resp.StatusCode == http.StatusConflict {
				return true, nil
			}
			return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
		}
		return false, nil
	}

	return sess, nil
}

func validateRetryConfiguration(cfg contextConfig) error {
	maxRetries := 50
	maxWaitTime := time.Hour * 24

	if cfg.retryMax < 0 || cfg.retryWaitMin < 0 || cfg.retryWaitMax < 0 {
		return fmt.Errorf("wrong retry values: maximum number of retries (%d), minimum retry wait time (%v), maximum retry wait time (%v) cannot be negative", cfg.retryMax, cfg.retryWaitMin, cfg.retryWaitMax)
	}

	if cfg.retryWaitMax < cfg.retryWaitMin {
		return fmt.Errorf("wrong retry values: maximum retry wait time (%v) cannot be lower than minimum retry wait time (%v)", cfg.retryWaitMax, cfg.retryWaitMin)
	}

	if cfg.retryMax > maxRetries {
		return fmt.Errorf("wrong retry values: too many retries, maximum number of retries (%d) cannot be higher than %d", cfg.retryMax, maxRetries)
	}

	if cfg.retryWaitMin > maxWaitTime || cfg.retryWaitMax > maxWaitTime {
		return fmt.Errorf("wrong retry values: retry wait time too long, minimum retry wait time (%v) cannot be higher than %v or maximum retry wait time (%v) cannot be higher than %v", cfg.retryWaitMin, maxWaitTime, cfg.retryWaitMax, maxWaitTime)

	}
	return nil
}
