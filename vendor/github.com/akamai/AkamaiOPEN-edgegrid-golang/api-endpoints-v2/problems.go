package apiendpoints

type Problems []*Problem

type Problem struct {
	Detail   string   `json:"detail"`
	Errors   Problems `json:"errors"`
	Instance string   `json:"instance"`
	Status   int      `json:"status"`
	Title    string   `json:"title"`
	Type     string   `json:"type"`
}
