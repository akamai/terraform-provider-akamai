package apiendpoints

type APIPrivacySettings struct {
	Resources map[int]APIPrivacyResource `json:"resources"`
	Public    bool                       `json:"public"`
}

type APIPrivacyResource struct {
	ResourceSettings
	Notes  string `json:"notes"`
	Public bool   `json:"public"`
}
