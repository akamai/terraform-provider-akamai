package collections

// ForEachInSlice applies the provided function for each element of the slice.
func ForEachInSlice[S ~[]E, E any](s S, fn func(a E) E) {
	for i, v := range s {
		s[i] = fn(v)
	}
}
