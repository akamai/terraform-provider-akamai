package cloudlets

import (
	"errors"
)

var (
	// ErrPolicyActivation is returned when policy activation fails
	ErrPolicyActivation = errors.New("policy activation")
	// ErrPolicyActivationTimeout is returned when policy activation fails due to a timeout
	ErrPolicyActivationTimeout = errors.New("policy activation timeout")

	// ErrPolicyActivationCanceled is returned on activation poll cancel
	ErrPolicyActivationCanceled = errors.New("operation canceled while waiting for policy activation status")
	// ErrPolicyActivationContextTerminated is returned on activation context termination
	ErrPolicyActivationContextTerminated = errors.New("policy activation context terminated")

	// ErrApplicationLoadBalancerActivation is returned when application load balancer activation fails
	ErrApplicationLoadBalancerActivation = errors.New("application load balancer activation")
	// ErrApplicationLoadBalancerActivationTimeout is returned when application load balancer activation fails due to a timeout
	ErrApplicationLoadBalancerActivationTimeout = errors.New("application load balancer activation timeout")

	// ErrApplicationLoadBalancerActivationCanceled is returned on activation poll cancel
	ErrApplicationLoadBalancerActivationCanceled = errors.New("operation canceled while waiting for application load balancer activation status")
	// ErrApplicationLoadBalancerActivationContextTerminated is returned on activation context termination
	ErrApplicationLoadBalancerActivationContextTerminated = errors.New("application load balancer activation context terminated")
)
