package apiendpoints

type CacheSettings struct {
	Enabled           bool                  `json:"enabled"`
	Option            string                `json:"option"`
	MaxAge            *MaxAge               `json:"maxAge"`
	ServeStale        bool                  `json:"serveStale"`
	DownstreamCaching DownstreamCaching     `json:"downstreamCaching"`
	ErrorCaching      ErrorCaching          `json:"errorCaching"`
	Resources         map[int]CacheResource `json:"resources"`
}

type DownstreamCaching struct {
	Option        string  `json:"option"`
	Lifetime      string  `json:"lifetime"`
	MaxAge        *MaxAge `json:"maxAge"`
	Headers       string  `json:"headers"`
	MarkAsPrivate bool    `json:"markAsPrivate"`
}

type ErrorCaching struct {
	Enabled       bool    `json:"enabled"`
	MaxAge        *MaxAge `json:"maxAge"`
	PreserveStale bool    `json:"preserveStale"`
}

type MaxAge struct {
	Duration int    `json:"duration"`
	Unit     string `json:"unit"`
}

type CacheResource struct {
	ResourceSettings
	Option     CacheResourceOptionValue `json:"option"`
	MaxAge     *MaxAge                  `json:"maxAge"`
	ServeStale bool                     `json:"serveStale"`
}

type MaxAgeUnitValue string
type CacheResourceOptionValue string

const (
	MaxAgeUnitSeconds MaxAgeUnitValue = "SECONDS"
	MaxAgeUnitMinutes MaxAgeUnitValue = "MINUTES"
	MaxAgeUnitHours   MaxAgeUnitValue = "HOURS"
	MaxAgeUnitDays    MaxAgeUnitValue = "DAYS"

	CacheResourceOptionCache                             CacheResourceOptionValue = "CACHE"
	CacheResourceOptionBypassCache                       CacheResourceOptionValue = "BYPASS_CACHE"
	CacheResourceOptionNoStore                           CacheResourceOptionValue = "NO_STORE"
	CacheResourceOptionHonorOriginCacheControl           CacheResourceOptionValue = "HONOR_ORIGIN_CACHE_CONTROL"
	CacheResourceOptionHonorOriginExpires                CacheResourceOptionValue = "HONOR_ORIGIN_EXPIRES"
	CacheResourceOptionHonorOriginCacheControlAndExpires CacheResourceOptionValue = "HONOR_ORIGIN_CACHE_CONTROL_AND_EXPIRES"
)
