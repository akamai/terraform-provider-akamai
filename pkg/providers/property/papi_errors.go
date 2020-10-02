package property

import "errors"

var (
	// PAPI group errors

	// ErrLookingUpGroupByName is returned when fetching group from API client by groupName returned an error or no group was found
	ErrLookingUpGroupByName = errors.New("looking up group with name")
	// ErrNoGroupsFound is returned when no groups were found
	ErrNoGroupsFound = errors.New("no groups found")
	// ErrGroupNotInContract is returned when none of the groups could be associated with given contractID
	ErrGroupNotInContract = errors.New("group does not belong to contract")
	ErrFetchingGroups     = errors.New("fetching groups")
	ErrGroupNotFound      = errors.New("group not found")

	// PAPI Contract errors

	// ErrLookingUpContract is returned when fetching contract from API client by contractID returned an error or no contract was found
	ErrLookingUpContract = errors.New("looking up contract for provided group")
	// ErrNoContractProvided is retured when no contract ID was provided but "name" was
	ErrNoContractProvided = errors.New("contract ID is required for non-default name")
	//
	ErrNoGroupProvided = errors.New("group not provided and it is a required input")
	// ErrNoContractsFound is returned when no contracts were found
	ErrNoContractsFound  = errors.New("no contracts were found")
	ErrContractNotFound  = errors.New("contract not found")
	ErrFetchingContracts = errors.New("fetching contracts")

	// PAPI Product errors
	ErrNoProductProvided = errors.New("product not provided and it is a required input")
	ErrProductFetch      = errors.New("fetching product")
	ErrProductNotFound   = errors.New("product not found")
	// PAPI CP Code errors

	// ErrLookingUpCPCode is returned when fetching CP Code from API client by contractID returned an error or no CP Code was found
	ErrLookingUpCPCode = errors.New("looking up CP Code by name")
	// ErrCPCodeModify is returned while attempting to modify existing CP code
	ErrCPCodeModify = errors.New("CP Code with provided name already exists for provided group and contract IDs and it cannot be managed through this API - please contact Customer Support")

	// PAPI Property errors

	// ErrPropertyNotFound is returned when no property was found for given name
	ErrPropertyNotFound = errors.New("property not found")
	// ErrRulesNotFound is returned when no rules were found
	ErrRulesNotFound = errors.New("property rules not found")

	// PAPI property version errors
	ErrVersionNotFound = errors.New("property version not found")
	ErrVersionCreate   = errors.New("creating property version")
)
