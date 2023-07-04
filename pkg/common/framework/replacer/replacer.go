// Package replacer ...
package replacer

// Replacer represents a string replacement operation.
type Replacer struct {
	Source       []string
	Replacements []string
	EqFunc       func(string, string) bool
}

// Replace replaces strings in the Source slice with corresponding strings from the Replacements slice.
// It returns a new slice with the replaced strings, while leaving the original Source slice unchanged.
// The replacement is performed by iterating over the Source slice and checking if each element has a match in the Replacements slice.
// If a match is found, the corresponding element from the Replacements slice is used as the replacement.
func (r Replacer) Replace() []string {
	if r.EqFunc == nil {
		r.EqFunc = func(_, _ string) bool { return false }
	}
	newslice := make([]string, len(r.Source))
	copy(newslice, r.Source)
	for i, val := range newslice {
		if v, ok := r.findPreferred(val); ok {
			newslice[i] = v
		}
	}
	return newslice
}

func (r Replacer) findPreferred(s string) (string, bool) {
	for _, v := range r.Replacements {
		if r.EqFunc(v, s) {
			return v, true
		}
	}
	return "", false
}
