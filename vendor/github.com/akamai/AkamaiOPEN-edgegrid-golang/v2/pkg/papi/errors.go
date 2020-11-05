package papi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// Error is a papi error interface
	Error struct {
		Type          string          `json:"type"`
		Title         string          `json:"title"`
		Detail        string          `json:"detail"`
		Instance      string          `json:"instance,omitempty"`
		BehaviorName  string          `json:"behaviorName,omitempty"`
		ErrorLocation string          `json:"errorLocation,omitempty"`
		StatusCode    int             `json:"statusCode,omitempty"`
		Errors        json.RawMessage `json:"errors,omitempty"`
		Warnings      json.RawMessage `json:"warnings,omitempty"`
	}
)

// Error parses an error from the response
func (p *papi) Error(r *http.Response) error {
	var e Error

	var body []byte

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Log(r.Request.Context()).Errorf("reading error response body: %s", err)
		e.StatusCode = r.StatusCode
		e.Title = fmt.Sprintf("Failed to read error body")
		e.Detail = err.Error()
		return &e
	}

	if err := json.Unmarshal(body, &e); err != nil {
		p.Log(r.Request.Context()).Errorf("could not unmarshal API error: %s", err)
		e.Title = fmt.Sprintf("Failed to unmarshal error body")
		e.Detail = err.Error()
	}

	e.StatusCode = r.StatusCode

	return &e
}

func (e *Error) Error() string {
	msg, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Sprintf("error marshaling API error: %s", err)
	}
	return fmt.Sprintf("API error: \n%s", msg)
}

// Is handles error comparisons
func (e *Error) Is(target error) bool {
	var t *Error
	if !errors.As(target, &t) {
		return false
	}

	if e == t {
		return true
	}

	if e.StatusCode != t.StatusCode {
		return false
	}

	return e.Error() == t.Error()
}
