package tools

// IntPtr returns a pointer to an integer passed as an argument
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr returns a pointer to an int64 passed as an argument
func Int64Ptr(i int64) *int64 {
	return &i
}

// Float64Ptr returns a pointer to a float64 passed as an argument
func Float64Ptr(f float64) *float64 {
	return &f
}

// StringPtr returns a pointer to a string passed as an argument
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to a bool passed as an argument
func BoolPtr(b bool) *bool {
	return &b
}
