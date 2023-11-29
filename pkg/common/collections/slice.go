package collections

// ForEachInSlice applies the provided function for each element of the slice.
func ForEachInSlice[S ~[]E, E any](s S, fn func(a E) E) {
	for i, v := range s {
		s[i] = fn(v)
	}
}

// StringInSlice determines if the searched string appears in the array.
func StringInSlice(s []string, searchTerm string) bool {
	for _, v := range s {
		if v == searchTerm {
			return true
		}
	}
	return false
}
