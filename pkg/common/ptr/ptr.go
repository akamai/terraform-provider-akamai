// Package ptr helps with creating pointers to literal values of any type.
package ptr

// To returns a pointer to the given v of any type.
func To[T any](v T) *T {
	return &v
}
