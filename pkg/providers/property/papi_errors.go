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

	// ErrLookingUpContract is returned when fetching contract from API client by contractID returned an error or no contract was found
	ErrLookingUpContract = errors.New("looking up contract for provided group")
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

	// PAPI Product errors

	// ErrNoProductProvided is returned when no "product" property is provided
	ErrNoProductProvided = errors.New("'product' not provided and it is a required input")
	// ErrProductFetch represents error while fetching product
	ErrProductFetch = errors.New("fetching product")
	// ErrProductNotFound is returned when product with provided ID does not exist
	ErrProductNotFound = errors.New("product not found")

	// PAPI CP Code errors

	// ErrLookingUpCPCode is returned when fetching CP Code from API client by contractID returned an error or no CP Code was found
	ErrLookingUpCPCode = errors.New("looking up CP Code by name")
	// ErrCPCodeModify is returned while attempting to modify existing CP code
	ErrCPCodeModify = errors.New("CP Code with provided name already exists for provided group and contract IDs and it cannot be managed through this API - please contact Customer Support")
	// ErrCpCodeNotFound is returned when cp code with provided ID does not exist
	ErrCpCodeNotFound = errors.New("cp code not found")

	// PAPI Property errors

	// ErrPropertyNotFound is returned when no property was found for given name
	ErrPropertyNotFound = errors.New("property not found")
	// ErrRulesNotFound is returned when no rules were found
	ErrRulesNotFound = errors.New("property rules not found")

	// PAPI property version errors

	// ErrVersionCreate represents an error while creating new property version
	ErrVersionCreate = errors.New("creating property version")

	// PAPI rule format errors

	// ErrRuleFormatsNotFound is returned when no rule formats were found
	ErrRuleFormatsNotFound = errors.New("no rule formats found")

	// ErrEdgeHostnameNotFound is returned when no edgehostname were found
	ErrEdgeHostnameNotFound = errors.New("unable to find edge hostname")

	// DiagWarnActivationTimeout returned on activation poll timeout
	DiagWarnActivationTimeout = diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Timeout waiting for activation status",
		Detail: `
The activation creation request has been started successfully, however the operation timeout was 
exceeded while waiting for the remote resource to update. You may retry the operation to continue 
to wait for the final status.

It is recommended that the timeout for activation resources be set to greater than 90 minutes.
See: https://www.terraform.io/docs/configuration/resources.html#operation-timeouts
`,
	}

	// DiagWarnActivationCanceled is returned on activation poll cancel
	DiagWarnActivationCanceled = diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Operation canceled while waiting for activation status",
		Detail: `
The activation creation request has been started successfully, however the a cancellation was recived
while waiting for the remote resource to update. You may retry the operation to continue to wait for 
the final status.

It is recommended that the timeout for activation resources be set to greater than 90 minutes.
See: https://www.terraform.io/docs/configuration/resources.html#operation-timeouts
`,
	}
)
