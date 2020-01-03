package configgtm

import (
	"fmt"
)

type ConfigGTMError interface {
	error
	Network() bool
	NotFound() bool
	FailedToSave() bool
	ValidationFailed() bool
}

func IsConfigGTMError(e error) bool {
	_, ok := e.(ConfigGTMError)
	return ok
}

type CommonError struct {
	entityName       string
	name             string
	httpErrorMessage string
	apiErrorMessage  string
	err              error
}

func (e CommonError) SetItem(itemName string, itemVal interface{}) {
	switch itemName {
	case "entityName":
		e.entityName = itemVal.(string)
	case "name":
		e.name = itemVal.(string)
	case "httpErrorMessage":
		e.httpErrorMessage = itemVal.(string)
	case "apiErrorMessage":
		e.apiErrorMessage = itemVal.(string)
	case "err":
		e.err = itemVal.(error)
	}
}

func (e CommonError) GetItem(itemName string) interface{} {
	switch itemName {
	case "entityName":
		return e.entityName
	case "name":
		return e.name
	case "httpErrorMessage":
		return e.httpErrorMessage
	case "apiErrorMessage":
		return e.apiErrorMessage
	case "err":
		return e.err
	}

	return nil
}

func (e CommonError) Network() bool {
	if e.httpErrorMessage != "" {
		return true
	}
	return false
}

func (e CommonError) NotFound() bool {
	if e.err == nil && e.httpErrorMessage == "" && e.apiErrorMessage == "" {
		return true
	}
	return false
}

func (CommonError) FailedToSave() bool {
	return false
}

func (e CommonError) ValidationFailed() bool {
	if e.apiErrorMessage != "" {
		return true
	}
	return false
}

func (e CommonError) Error() string {

	if e.Network() {
		return fmt.Sprintf("%s \"%s\" network error: [%s]", e.entityName, e.name, e.httpErrorMessage)
	}

	if e.NotFound() {
		return fmt.Sprintf("%s \"%s\" not found.", e.entityName, e.name)
	}

	if e.FailedToSave() {
		return fmt.Sprintf("%s \"%s\" failed to save: [%s]", e.entityName, e.name, e.err.Error())
	}

	if e.ValidationFailed() {
		return fmt.Sprintf("%s \"%s\" validation failed: [%s]", e.entityName, e.name, e.apiErrorMessage)
	}

	if e.err != nil {
		return e.err.Error()
	}

	return "<nil>"
}
