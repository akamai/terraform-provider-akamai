package retry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/stretchr/testify/assert"
)

var alwaysRetry = func(_ error) bool { return true }

func contextWithTimeout(timeout time.Duration) func() (context.Context, context.CancelFunc) {
	return func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.Background(), timeout)
	}
}

func returnCallsTimes100FailingOnFirst(failOnFirstCalls int) func(*int) Fn[int] {
	return func(calls *int) Fn[int] {
		return func(_ context.Context) (*int, error) {
			*calls++
			if failOnFirstCalls != 0 && *calls <= failOnFirstCalls {
				return nil, fmt.Errorf("polling error")
			}
			return ptr.To(*calls * 100), nil
		}
	}
}

func TestPoll(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		setupContext     func() (context.Context, context.CancelFunc)
		setupFn          func(*int) Fn[int]
		shouldRetryData  func(int) bool
		shouldRetryError func(error) bool
		interval         time.Duration
		deadline         time.Duration
		expectedResult   int
		expectedError    string
		expectedCalls    int
		minDuration      time.Duration
	}{
		"success on first attempt": {
			expectedResult: 100,
			expectedCalls:  1,
		},
		"retry on error and succeed": {
			shouldRetryError: alwaysRetry,
			setupFn:          returnCallsTimes100FailingOnFirst(2),
			interval:         10 * time.Millisecond,
			expectedResult:   300,
			expectedCalls:    3,
			minDuration:      20 * time.Millisecond,
		},
		"retry on data condition and succeed": {
			shouldRetryData: func(v int) bool { return v < 500 },
			interval:        10 * time.Millisecond,
			expectedResult:  500,
			expectedCalls:   5,
			minDuration:     40 * time.Millisecond,
		},
		"non-retryable error returns immediately": {
			setupFn:       returnCallsTimes100FailingOnFirst(1),
			expectedError: "not retrying due to non-retriable error: polling error",
			expectedCalls: 1,
		},
		"context timeout during retry loop": {
			setupContext:     contextWithTimeout(50 * time.Millisecond),
			setupFn:          returnCallsTimes100FailingOnFirst(2),
			shouldRetryError: alwaysRetry,
			interval:         30 * time.Millisecond,
			expectedError:    "context terminated while waiting to retry: context deadline exceeded",
			expectedCalls:    2,
		},
		"context timeout due to deadline option": {
			setupContext: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			setupFn:          returnCallsTimes100FailingOnFirst(2),
			shouldRetryError: alwaysRetry,
			interval:         30 * time.Millisecond,
			deadline:         50 * time.Millisecond,
			expectedError:    "context terminated while waiting to retry: context deadline exceeded",
			expectedCalls:    2,
		},
		"uses default interval when not set": {
			setupContext:     contextWithTimeout(15 * time.Second),
			setupFn:          returnCallsTimes100FailingOnFirst(1),
			shouldRetryError: alwaysRetry,
			interval:         0, // reset to default
			expectedResult:   200,
			expectedCalls:    2,
			minDuration:      10 * time.Second,
		},
		"error nil function": {
			setupFn:       func(_ *int) Fn[int] { return nil },
			expectedError: "retry function cannot be nil",
			expectedCalls: 0,
		},
		"error nil context": {
			setupContext: func() (context.Context, context.CancelFunc) {
				return nil, func() {}
			},
			expectedError: "context cannot be nil",
			expectedCalls: 0,
		},
		"error context without deadline and no deadline option": {
			setupContext: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			expectedError: "either provide a context with deadline or set Deadline in polling options",
			expectedCalls: 0,
		},
		"error context cancelled during interval wait": {
			setupContext: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					time.Sleep(150 * time.Millisecond)
					cancel()
				}()
				return ctx, nil
			},
			shouldRetryError: alwaysRetry,
			interval:         100 * time.Millisecond,
			deadline:         5 * time.Second,
			setupFn:          returnCallsTimes100FailingOnFirst(2),
			expectedError:    "context terminated while waiting to retry: context canceled",
			expectedCalls:    2,
		},
		"function not executed when context is already done": {
			setupContext: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				// Wait to ensure context expires before Poll is called
				time.Sleep(5 * time.Millisecond)
				return ctx, cancel
			},
			expectedError: "context terminated before function execution: context deadline exceeded",
			expectedCalls: 0, // Function should NOT be called
		},
		"error negative deadline": {
			deadline:      -5,
			expectedError: "deadline -5ns cannot be negative",
		},
		"error negative interval": {
			interval:      -10 * time.Millisecond,
			expectedError: "interval -10ms cannot be negative",
		},
		"error nil response without error": {
			setupFn: func(_ *int) Fn[int] {
				return func(_ context.Context) (*int, error) {
					return nil, nil
				}
			},
			expectedError: "received nil response without error",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			calls := 0
			setupCtx := tc.setupContext
			if setupCtx == nil {
				setupCtx = contextWithTimeout(5 * time.Second)
			}
			ctx, cancel := setupCtx()
			if cancel != nil {
				defer cancel()
			}

			opts := PollingOpts[int]{
				Fn: returnCallsTimes100FailingOnFirst(0)(&calls),
			}
			if tc.setupFn != nil {
				opts.Fn = tc.setupFn(&calls)
			}
			if tc.shouldRetryData != nil {
				opts.ShouldRetryData = tc.shouldRetryData
			}
			if tc.shouldRetryError != nil {
				opts.ShouldRetryError = tc.shouldRetryError
			}
			if tc.interval != 0 {
				opts.Interval = tc.interval
			}
			if tc.deadline != 0 {
				opts.Deadline = tc.deadline
			}

			start := time.Now()
			resp, err := Poll(ctx, opts)
			duration := time.Since(start)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, *resp)
			}
			assert.Equal(t, tc.expectedCalls, calls)
			if tc.minDuration != 0 {
				assert.GreaterOrEqual(t, duration, tc.minDuration)
			}
		})
	}
}
