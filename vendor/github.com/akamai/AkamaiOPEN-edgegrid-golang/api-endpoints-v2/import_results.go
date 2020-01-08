package apiendpoints

type ImportResult struct {
	APIEndpointDetails Endpoint `json:"apiEndpointDetails"`
	Problems           Problems `json:"problems"`
}
