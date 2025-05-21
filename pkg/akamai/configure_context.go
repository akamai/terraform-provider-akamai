package akamai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgegrid"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/retryablehttp"
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
	log := log.FromContext(cfg.ctx, "OperationID", operationID)

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

func overrideRetryPolicy(basePolicy retryablehttp.CheckRetry) retryablehttp.CheckRetry {
	return func(ctx context.Context, resp *http.Response, err error) (bool, error) {

		// do not retry on context.Canceled or context.DeadlineExceeded
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			isPatchHostnameURL, e := regexp.MatchString(
				regexp.MustCompile(`^/papi/v1/properties/[^/]+/hostnames$`).String(),
				resp.Request.URL.Path)
			if e != nil {
				return false, e
			}
			// If the request is PATCH hostname bucket resulting in 429 (default cert limit exceeded), do not retry it
			if isPatchHostnameURL && resp.Request.Method == http.MethodPatch {
				return false, nil
			}
			// Retry all PAPI requests resulting in status code 429
			// The backoff time is calculated in getXRateLimitBackoff
			if strings.HasPrefix(resp.Request.URL.Path, "/papi/") {
				return true, nil
			}
		}

		var urlErr *url.Error
		if (resp != nil && resp.Request.Method == http.MethodGet) ||
			(resp == nil && errors.As(err, &urlErr) && strings.ToUpper(urlErr.Op) == http.MethodGet) {

			if resp != nil && resp.StatusCode == http.StatusConflict {
				return true, nil
			}
			return basePolicy(ctx, resp, err)
		}
		return false, nil
	}
}

// Note that Date's resolution is seconds (e.g. Mon, 01 Jul 2024 14:32:14 GMT),
// while X-RateLimit-Next's resolution is milliseconds (2024-07-01T14:32:28.645Z).
// This may cause the wait time to be inflated by at most one second, like for the
// actual server response time around 2024-07-01T14:32:14.999Z. This is acceptable behavior
// as retry does not occur earlier than expected.
func getXRateLimitBackoff(resp *http.Response, logger akalog.Interface) (time.Duration, bool) {
	nextHeader := resp.Header.Get("X-RateLimit-Next")
	if nextHeader == "" {
		return 0, false
	}
	next, err := time.Parse(time.RFC3339Nano, nextHeader)
	if err != nil {
		if logger != nil {

			logger.Error("Could not parse X-RateLimit-Next header", "error", err)
		}
		return 0, false
	}

	dateHeader := resp.Header.Get("Date")
	if dateHeader == "" {
		if logger != nil {
			logger.Warnf("No Date header for X-RateLimit-Next: %s", nextHeader)
		}
		return 0, false
	}
	date, err := time.Parse(time.RFC1123, dateHeader)
	if err != nil {
		if logger != nil {
			logger.Error("Could not parse Date header", "error", err)
		}
		return 0, false
	}

	// Next in the past does not make sense
	if next.Before(date) {
		if logger != nil {
			logger.Warnf("X-RateLimit-Next: %s before Date: %s", nextHeader, dateHeader)
		}
		return 0, false
	}
	return next.Sub(date), true
}

func overrideBackoff(baseBackoff retryablehttp.Backoff, logger akalog.Interface) retryablehttp.Backoff {
	return func(minT, maxT time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests {
				if wait, ok := getXRateLimitBackoff(resp, logger); ok {
					return wait
				}
			}
		}
		return baseBackoff(minT, maxT, attemptNum, resp)
	}
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

	retryClient.HTTPClient.CheckRedirect = func(req *http.Request, _ []*http.Request) error {
		return sess.Sign(req)
	}

	retryClient.CheckRetry = overrideRetryPolicy(retryablehttp.DefaultRetryPolicy)
	l := sess.Log(cfg.ctx)
	retryClient.Backoff = overrideBackoff(retryablehttp.DefaultBackoff, l)
	retryClient.Logger = session.GetRetryableLogger(l)
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
