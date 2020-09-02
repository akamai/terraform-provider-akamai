package property

import "errors"

var (
	ErrPapiLookingUpGroupByName = errors.New("looking up group with name")
	ErrPapiNoGroupsFound        = errors.New("no groups found")
	ErrPapiGroupNotFound        = errors.New("could not find group for given group ID")
	ErrPapiFindingGroupsByName  = errors.New("could not find groups for given name")
	ErrPapiNoContractProvided   = errors.New("contract ID is required for non-default name")
	ErrPapiGroupNotInContract   = errors.New("group does not belong to contract")
	ErrPAPICPCodeModify         = errors.New("CP Code with provided name already exists for provided group and contract IDs and it cannot be managed through this API - please contact Customer Support")
	ErrInvalidPropertyType      = errors.New("property must be of the specified type")
)
