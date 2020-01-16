package apiendpoints

type AcgPair []struct {
	DisplayName string `json:"displayName"`
	AcgID       string `json:"acgId"`
	GroupID     int    `json:"groupId"`
	ContractID  string `json:"contractId"`
}
