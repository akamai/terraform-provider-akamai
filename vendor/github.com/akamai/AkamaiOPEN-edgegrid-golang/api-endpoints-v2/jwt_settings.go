package apiendpoints

type JWTSettings struct {
	Enabled  bool `json:"enabled"`
	Settings struct {
		Location   JWTSettingsLocationValue `json:"location"`
		ParamName  string                   `json:"paramName"`
		ClockSkew  int                      `json:"clockSkew"`
		Validation *struct {
			Claims        []JWTClaim    `json:"claims"`
			RsaPublicKeyA RsaPublicKey  `json:"rsaPublicKeyA"`
			RsaPublicKeyB *RsaPublicKey `json:"rsaPublicKeyB,omitempty"`
		} `json:"validation"`
	} `json:"settings"`
	Resources map[int]JWTSettingsResource `json:"resources"`
}

type JWTSettingsResource struct {
	ResourceSettings
	Enabled bool    `json:"enabled"`
	Notes   *string `json:"notes,omitempty"`
}

type RsaPublicKey struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type JWTClaim struct {
	Name     string   `json:"name"`
	Validate bool     `json:"validate"`
	Required bool     `json:"required"`
	Value    []string `json:"value"`
	Type     string   `json:"type"`
}

type JWTSettingsLocationValue string

const (
	JWTSettingsLocationHeader JWTSettingsLocationValue = "HEADER"
	JWTSettingsLocationCookie JWTSettingsLocationValue = "COOKIE"
	JWTSettingsLocationQuery  JWTSettingsLocationValue = "QUERY"
)
