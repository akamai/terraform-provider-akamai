package apiendpoints

type Parameters struct {
	APIChildParameters      []*Parameters             `json:"apiChildParameters"`
	APIParameterID          int                       `json:"apiParameterId"`
	APIParameterRequired    bool                      `json:"apiParameterRequired"`
	APIParameterName        string                    `json:"apiParameterName"`
	APIParameterLocation    APIParameterLocationValue `json:"apiParameterLocation"`
	APIParameterType        APIParameterTypeValue     `json:"apiParameterType"`
	APIParameterNotes       *string                   `json:"apiParameterNotes"`
	APIParamLogicID         int                       `json:"apiParamLogicId"`
	Array                   bool                      `json:"array"`
	APIParameterRestriction struct {
		RangeRestriction struct {
			RangeMin int `json:"rangeMin"`
			RangeMax int `json:"rangeMax"`
		} `json:"rangeRestriction"`
	} `json:"apiParameterRestriction"`
}

type APIParameterLocationValue string
type APIParameterTypeValue string

const (
	APIParameterLocationHeader APIParameterLocationValue = "header"
	APIParameterLocationCookie APIParameterLocationValue = "cookie"
	APIParameterLocationQuery  APIParameterLocationValue = "query"
	APIParameterLocationBody   APIParameterLocationValue = "body"

	APIParameterTypeString  APIParameterTypeValue = "string"
	APIParameterTypeInteger APIParameterTypeValue = "integer"
	APIParameterTypeNumber  APIParameterTypeValue = "number"
	APIParameterTypeBoolean APIParameterTypeValue = "boolean"
	APIParameterTypeJson    APIParameterTypeValue = "json/xml"
	APIParameterTypeXml     APIParameterTypeValue = "json/xml"
)
