package tools

// IntPtr returns a pointer to an integer passed as an argument
func IntPtr(i int) *int {
	return &i
}

// StringPtr returns a pointer to a string passed as an argument
func StringPtr(s string) *string {
	return &s
}
