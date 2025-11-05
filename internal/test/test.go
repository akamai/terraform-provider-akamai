// Package test contains utility code used in tests
package test

import (
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewTimeFromString returns a time value parsed from a string
// in the RFC3339Nano format
func NewTimeFromString(t *testing.T, s string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339Nano, s)
	require.NoError(t, err)
	return parsedTime
}

// NewTimeFromStringMust returns a time value parsed from a string
// in the RFC3339Nano format. It assumes that the input string is always valid
func NewTimeFromStringMust(s string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}
	return parsedTime
}

// RateLimitHTTPHandler first returns status 429 with the Akamai-RateLimit-Next or X-RateLimit-Next headers set to
// time.Now() plus a random value between 1 and 5 milliseconds. It keeps sending 429 until the Akamai-RateLimit-Next
// or X-RateLimit-Next point in time. Then it starts to return SuccessCode and SuccessBody
// indefinitely.
type RateLimitHTTPHandler struct {
	T           *testing.T
	SuccessCode int
	SuccessBody string

	mutex         sync.Mutex
	availableAt   time.Time
	returnedCodes []int
	returnTimes   []time.Time
}

func (h *RateLimitHTTPHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request, header string) {
	av := h.AvailableAt()

	if av.IsZero() {
		busyInterval := time.Duration(1+rand.Intn(4)) * time.Millisecond
		h.setAvailableAt(time.Now().Add(busyInterval))
		h.setTooManyRequests(w, header)
		return
	}

	now := time.Now()
	if now.Before(av) {
		h.setTooManyRequests(w, header)
	} else {
		h.setStatusCode(w, h.SuccessCode)
		_, err := w.Write([]byte(h.SuccessBody))
		assert.NoError(h.T, err)
	}
}

// AvailableAt returns the point in time at which the handler stops returning status code 429
func (h *RateLimitHTTPHandler) AvailableAt() time.Time {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.availableAt
}

// ReturnedCodes returns a list of status codes from subsequent handler responses
func (h *RateLimitHTTPHandler) ReturnedCodes() []int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	res := make([]int, len(h.returnedCodes))
	copy(res, h.returnedCodes)
	return res
}

// ReturnTimes returns a list of times at which subsequent responses were written
func (h *RateLimitHTTPHandler) ReturnTimes() []time.Time {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	res := make([]time.Time, len(h.returnTimes))
	copy(res, h.returnTimes)
	return res
}

func (h *RateLimitHTTPHandler) setTooManyRequests(w http.ResponseWriter, header string) {
	// Do not use Add() to avoid canonicalization to Akamai-RateLimit-Next or X-RateLimit-Next
	nextStr := h.availableAt.Format(time.RFC3339Nano)
	w.Header()[header] = []string{nextStr}
	h.setStatusCode(w, http.StatusTooManyRequests)
	body := "Your request did not succeed as this operation has reached the limit " +
		"for your account. Please try after " + nextStr
	_, err := w.Write([]byte(body))
	assert.NoError(h.T, err)
}

func (h *RateLimitHTTPHandler) setStatusCode(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.returnedCodes = append(h.returnedCodes, statusCode)
	h.returnTimes = append(h.returnTimes, time.Now())
}

func (h *RateLimitHTTPHandler) setAvailableAt(availableAt time.Time) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.availableAt = availableAt
}
