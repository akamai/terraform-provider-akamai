package apiendpoints

type Methods []Method

type Method struct {
	APIResourceMethodID      int          `json:"apiResourceMethodId"`
	APIResourceMethod        MethodValue  `json:"apiResourceMethod"`
	APIResourceMethodLogicID int          `json:"apiResourceMethodLogicId"`
	APIParameters            []Parameters `json:"apiParameters"`
}

type MethodValue string

const (
	MethodGet     MethodValue = "get"
	MethodPost    MethodValue = "post"
	MethodPut     MethodValue = "put"
	MethodDelete  MethodValue = "delete"
	MethodHead    MethodValue = "head"
	MethodPatch   MethodValue = "patch"
	MethodOptions MethodValue = "options"
)
