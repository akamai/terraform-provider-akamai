// Package retry defines a generic retry mechanism with customizable options for retrying API calls.
package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// PollingOpts configures the behavior of a polling operation that repeatedly executes
// a function until a success condition is met or a deadline is reached.
type PollingOpts[T any] struct {
	// Fn is the function to execute on each polling attempt. This field is required.
	Fn Fn[T]

	// ShouldRetryData is an optional predicate that examines the successful result data to
	// determine if polling should continue. If provided and returns true, polling continues.
	// If nil, polling stops on the first successful execution of Fn.
	ShouldRetryData func(T) bool

	// ShouldRetryError is an optional predicate that examines errors to determine if polling
	// should continue. If provided and returns true, polling continues despite the error.
	// If nil, any error from Fn will stop polling and be returned.
	ShouldRetryError func(error) bool

	// Interval is the duration to wait between polling attempts. It must be a non-negative integer.
	// If zero, default interval of 10 seconds is used.
	Interval time.Duration

	// Deadline is the maximum duration to continue polling before giving up. This field
	// must be provided if the context used by the Poll function does not have a deadline
	// defined. It must be a non-negative integer. If the context contains a deadline, this field can be omitted.
	// When both are present, whichever deadline is reached first will stop the polling.
	// When the deadline is exceeded, polling stops and returns a timeout error.
	Deadline time.Duration
}

// defaultInterval is the default wait time between retries.
const defaultInterval = 10 * time.Second

// Fn is a function type that represents the operation to be polled.
type Fn[T any] func(context.Context) (*T, error)

// Poll repeatedly executes a function at regular intervals until a condition is met or a deadline is reached.
//
// The polling mechanism works by calling opts.Fn in a loop, waiting opts.Interval between attempts.
// Polling continues until one of the following occurs:
//   - The function returns data that doesn't satisfy ShouldRetryData
//   - The function returns an error that doesn't satisfy ShouldRetryError
//   - The context deadline is exceeded
//   - The context is cancelled
//
// Deadline handling:
// Either provide a context with a deadline, or set opts.Deadline to create a timeout context automatically.
// If opts.Deadline is set, it will wrap the provided context with a timeout.
//
// Retry predicates:
//   - ShouldRetryData: Called when opts.Fn returns successfully (non-nil data, nil error).
//     Return true to continue polling, false to accept the result and stop.
//     If not provided, defaults to a function that always returns false (i.e., stop on first success).
//   - ShouldRetryError: Called when opts.Fn returns an error.
//     Return true to continue polling (treating the error as retriable), false to stop and return the error.
//     If not provided, defaults to a function that always returns false (i.e., stop on first error).
//
// Example usage:
//
//	result, err := Poll(ctx, PollingOpts[*MyResource]{
//	    Fn: func(ctx context.Context) (*MyResource, error) {
//	        return fetchResource(ctx, resourceID)
//	    },
//	    ShouldRetryData: func(r *MyResource) bool {
//	        return r.Status != "ACTIVE" // Keep polling until status is ACTIVE
//	    },
//	    ShouldRetryError: func(err error) bool {
//	        return errors.Is(err, ErrNotFound) // Retry on NotFound errors
//	    },
//	    Interval: 5 * time.Second,
//	    Deadline: 2 * time.Minute,
//	})
func Poll[T any](ctx context.Context, opts PollingOpts[T]) (*T, error) {
	now := time.Now()
	if opts.Fn == nil {
		return nil, fmt.Errorf("retry function cannot be nil")
	}
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if opts.ShouldRetryData == nil {
		opts.ShouldRetryData = func(_ T) bool { return false }
	}
	if opts.ShouldRetryError == nil {
		opts.ShouldRetryError = func(_ error) bool { return false }
	}
	if opts.Interval < 0 {
		return nil, fmt.Errorf("interval %v cannot be negative", opts.Interval)
	}
	if opts.Interval == 0 {
		opts.Interval = defaultInterval
	}
	if opts.Deadline < 0 {
		return nil, fmt.Errorf("deadline %v cannot be negative", opts.Deadline)
	}
	originalDeadline := getDeadlineOrNil(ctx, now)
	if opts.Deadline > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Deadline)
		defer cancel()
	} else if originalDeadline == nil {
		return nil, fmt.Errorf("either provide a context with deadline or set Deadline in polling options")
	}

	tflog.Debug(ctx, "Starting polling loop", map[string]any{
		"interval":          opts.Interval,
		"deadline":          opts.Deadline,
		"effectiveDeadline": getDeadlineOrNil(ctx, now),
		"originalDeadline":  originalDeadline,
	})

	for {
		// Check context error to handle cancellations and timeouts.
		if ctx.Err() != nil {
			tflog.Debug(ctx, "Context terminated before function execution", map[string]any{
				"error": ctx.Err(),
			})
			return nil, fmt.Errorf("context terminated before function execution: %w", ctx.Err())
		}

		tflog.Debug(ctx, "Executing the polling function")
		resp, err := opts.Fn(ctx)

		doRetry, retryErr := shouldRetry(ctx, opts, resp, err)
		// Logged inside shouldRetry
		if retryErr != nil {
			return resp, retryErr
		}
		if !doRetry {
			return resp, nil
		}

		select {
		case <-time.After(opts.Interval):
			// Continue to next iteration
			tflog.Debug(ctx, "Retrying after interval", map[string]any{
				"interval": opts.Interval,
			})
		case <-ctx.Done():
			tflog.Debug(ctx, "Context terminated while waiting to retry", map[string]any{
				"error": ctx.Err(),
			})
			return nil, fmt.Errorf("context terminated while waiting to retry: %w", ctx.Err())
		}
	}
}

// shouldRetry determines whether to retry based on the response and error from the polling function.
func shouldRetry[T any](ctx context.Context, opts PollingOpts[T], resp *T, err error) (bool, error) {
	if err != nil {
		if opts.ShouldRetryError(err) {
			tflog.Debug(ctx, "Retrying due to retriable error", map[string]any{
				"error": err,
			})
			return true, nil
		}
		tflog.Debug(ctx, "Not retrying due to non-retriable error", map[string]any{
			"error": err,
		})
		return false, fmt.Errorf("not retrying due to non-retriable error: %w", err)
	}

	if resp == nil {
		tflog.Debug(ctx, "Received nil response without error")
		return false, fmt.Errorf("received nil response without error")
	}
	if opts.ShouldRetryData(*resp) {
		tflog.Debug(ctx, "Retrying due to ShouldRetryData condition")
		return true, nil
	}

	tflog.Debug(ctx, "Not retrying as response is acceptable")
	return false, nil
}

func getDeadlineOrNil(ctx context.Context, start time.Time) *time.Duration {
	deadline, ok := ctx.Deadline()
	if ok {
		diff := deadline.Sub(start)
		return &diff
	}
	return nil
}
