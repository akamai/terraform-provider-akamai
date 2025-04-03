package property

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var (
	// PAPI group errors

	// ErrLookingUpGroupByName is returned when fetching group from API client by groupName returned an error or no group was found
	ErrLookingUpGroupByName = errors.New("looking up group with name")
	// ErrNoGroupsFound is returned when no groups were found
	ErrNoGroupsFound = errors.New("no groups found")
	// ErrGroupNotInContract is returned when none of the groups could be associated with given contractID
	ErrGroupNotInContract = errors.New("group does not belong to contract")
	// ErrFetchingGroups represents error while fetching groups
	ErrFetchingGroups = errors.New("fetching groups")
	// ErrGroupNotFound is returned when group with provided ID is not found
	ErrGroupNotFound = errors.New("group not found")

	// PAPI Contract errors

	// ErrLookingUpContract is returned when fetching contract from API client by groupId returned an error or no contract was found
	ErrLookingUpContract = errors.New("looking up contract for provided group")
	// ErrMultipleContractsInGroup is returned when fetching contract from API client by groupId returned multiple different contracts
	ErrMultipleContractsInGroup = diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "multiple contracts found for given group",
			Detail: "Resource doesn't support groups with multiple contracts. " +
				"Make sure your group has only one contract assigned before proceeding.",
		},
	}
	// ErrNoContractProvided is retured when no contract ID was provided but "name" was
	ErrNoContractProvided = errors.New("'contractId' is required for non-default name")
	// ErrNoGroupProvided is returned when no "group" property is provided
	ErrNoGroupProvided = errors.New("'group' not provided and it is a required input")
	// ErrNoContractsFound is returned when no contracts were found
	ErrNoContractsFound = errors.New("no contracts were found")
	// ErrContractNotFound is returned when contract with provided ID does not exist
	ErrContractNotFound = errors.New("contract not found")
	// ErrFetchingContracts represents error while fetching contracts
	ErrFetchingContracts = errors.New("fetching contracts")
	// ErrMultipleContractsFound is returned when more than one contract was found
	ErrMultipleContractsFound = diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "multiple contracts found",
			Detail: "Resource cannot unambiguously identify the contract. " +
				"Please provide either a 'group_id' or 'group_name' for accurate identification.",
		},
	}

	// PAPI Product errors

	// ErrNoProductProvided is returned when no "product" property is provided
	ErrNoProductProvided = errors.New("'product' not provided and it is a required input")
	// ErrProductFetch represents error while fetching product
	ErrProductFetch = errors.New("fetching product")
	// ErrProductNotFound is returned when product with provided ID does not exist
	ErrProductNotFound = errors.New("product not found")

	// PAPI CP Code errors

	// ErrLookingUpCPCode is returned when fetching CP Code from API client by contractID returned an error or no CP Code was found
	ErrLookingUpCPCode = errors.New("looking up cp code by name")
	// ErrCPCodeNotFound is returned when cp code with provided ID does not exist
	ErrCPCodeNotFound = errors.New("cp code not found")
	// ErrMoreCPCodesFound is returned when cp code with provided ID does not exist
	ErrMoreCPCodesFound = errors.New("more cp codes found")
	// ErrCPCodeUpdateTimeout is returned when waiting for a cp code update results in timeout
	ErrCPCodeUpdateTimeout = errors.New("cp code update timeout")

	// PAPI Property errors

	// ErrPropertyNotFound is returned when no property was found for given name
	ErrPropertyNotFound = errors.New("property not found")
	// ErrRulesNotFound is returned when no rules were found
	ErrRulesNotFound = errors.New("property rules not found")

	// PAPI property version errors

	// ErrVersionCreate represents an error while creating new property version
	ErrVersionCreate = errors.New("creating property version")
	// ErrPropertyVersionNotFound is returned when no property versions were found
	ErrPropertyVersionNotFound = errors.New("property version not found")

	// PAPI rule format errors

	// ErrRuleFormatsNotFound is returned when no rule formats were found
	ErrRuleFormatsNotFound = errors.New("no rule formats found")

	// ErrEdgeHostnameNotFound is returned when no edgehostname were found
	ErrEdgeHostnameNotFound = errors.New("unable to find edge hostname")

	// Property includes errors

	// ErrNoLatestIncludeActivation is returned when there is no activation for provided include
	ErrNoLatestIncludeActivation = errors.New("no latest activation for given include")

	// ErrPropertyInclude is returned when operation on property include fails
	ErrPropertyInclude = errors.New("property include")

	// DiagErrActivationTimeout returned on activation poll timeout
	DiagErrActivationTimeout = diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Timeout waiting for activation status",
		Detail: `
The activation creation request has been started successfully, however the operation timeout was 
exceeded while waiting for the remote resource to update. You may retry the operation to continue 
to wait for the final status.

It is recommended that the timeout for activation resources be set to greater than 90 minutes.
See: https://www.terraform.io/docs/configuration/resources.html#operation-timeouts
`,
	}

	// DiagErrActivationCanceled is returned on activation poll cancel
	DiagErrActivationCanceled = diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Operation canceled while waiting for activation status",
		Detail: `
The activation creation request has been started successfully, however the a cancellation was received
while waiting for the remote resource to update. You may retry the operation to continue to wait for 
the final status.

It is recommended that the timeout for activation resources be set to greater than 90 minutes.
See: https://www.terraform.io/docs/configuration/resources.html#operation-timeouts
`,
	}
)
