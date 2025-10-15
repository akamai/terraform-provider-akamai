// Package slicesets provides utility functions for manipulation on sets represented as slices.
package slicesets

// Subtract returns a new slice that contains all elements of 'a' that are not in 'b'.
func Subtract[E comparable](a []E, b []E) []E {
	toRemove := make(map[E]struct{}, len(b))
	for _, vb := range b {
		toRemove[vb] = struct{}{}
	}
	var result []E
	for _, v := range a {
		if _, found := toRemove[v]; !found {
			result = append(result, v)
		}
	}
	return result
}
