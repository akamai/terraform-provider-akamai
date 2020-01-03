package apiendpoints

type CORSSetings struct {
	Enabled          bool          `json:"enabled"`
	AllowedOrigins   []string      `json:"allowedOrigins,omitempty"`
	AllowedHeaders   []string      `json:"allowedHeaders,omitempty"`
	AllowedMethods   []MethodValue `json:"allowedMethods,omitempty"`
	AllowCredentials bool          `json:"allowCredentials,omitempty"`
	ExposedHeaders   []string      `json:"exposedHeaders,omitempty"`
	PreflightMaxAge  int           `json:"preflightMaxAge,omitempty"`
}
