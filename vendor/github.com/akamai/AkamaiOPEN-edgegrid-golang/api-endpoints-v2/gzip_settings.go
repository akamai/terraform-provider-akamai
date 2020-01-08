package apiendpoints

type GzipSettings struct {
	CompressResponse CompressResponseValue `json:"compressResponse"`
}

type CompressResponseValue string

const (
	CompressResponseAlways       CompressResponseValue = "ALWAYS"
	CompressResponseNever        CompressResponseValue = "NEVER"
	CompressResponseSameAsOrigin CompressResponseValue = "SAME_AS_ORIGIN"
)
