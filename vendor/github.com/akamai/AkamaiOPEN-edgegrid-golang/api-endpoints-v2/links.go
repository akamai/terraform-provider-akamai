package apiendpoints

type Links []struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
