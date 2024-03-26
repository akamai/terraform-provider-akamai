package edgeworkers

import (
	"errors"
)

var (
	// ErrEdgeworkerActivation is returned when edgeworker activation fails
	ErrEdgeworkerActivation = errors.New("edgeworker activation")
	// ErrEdgeworkerNoCurrentActivation is returned when edgeworker activation fails
	ErrEdgeworkerNoCurrentActivation = errors.New("edgeworker is not active")
	// ErrEdgeworkerDeactivation is returned when edgeworker deactivation fails
	ErrEdgeworkerDeactivation = errors.New("edgeworker deactivation")
	// ErrEdgeworkerNoLatestDeactivation is returned when edgeworker deactivation fails
	ErrEdgeworkerNoLatestDeactivation = errors.New("edgeworker does not have any deactivations")
	// ErrEdgeworkerActivationFailure is returned when edgeworker activation fails due to a timeout
	ErrEdgeworkerActivationFailure = errors.New("edgeworker activation failure")
	// ErrEdgeworkerDeactivationFailure is returned when edgeworker deactivation fails due to a timeout
	ErrEdgeworkerDeactivationFailure = errors.New("edgeworker deactivation failure")
	// ErrEdgeworkerActivationTimeout is returned when edgeworker activation fails due to a timeout
	ErrEdgeworkerActivationTimeout = errors.New("edgeworker activation timeout")
	// ErrEdgeworkerDeactivationTimeout is returned when edgeworker deactivation fails due to a timeout
	ErrEdgeworkerDeactivationTimeout = errors.New("edgeworker deactivation timeout")
	// ErrEdgeworkerActivationCancelled is returned on activation poll cancel
	ErrEdgeworkerActivationCancelled = errors.New("operation cancelled while waiting for edgeworker activation status")
	// ErrEdgeworkerDeactivationCancelled is returned on deactivation poll cancel
	ErrEdgeworkerDeactivationCancelled = errors.New("operation cancelled while waiting for edgeworker deactivation status")
	// ErrEdgeworkerActivationContextTerminated is returned on activation context termination
	ErrEdgeworkerActivationContextTerminated = errors.New("edgeworker activation context terminated")
	// ErrEdgeworkerDeactivationContextTerminated is returned on deactivation context termination
	ErrEdgeworkerDeactivationContextTerminated = errors.New("edgeworker deactivation context terminated")
)
