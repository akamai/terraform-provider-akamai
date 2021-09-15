package tools

// IntPtr returns a pointer to an integer passed as an argument
func IntPtr(i int) *int {
	return &i
}

