package collections

import "errors"

// ErrDuplicateKey is returned when the key already exists in the map
var ErrDuplicateKey = errors.New("duplicate key")

// AddMap populates 'to' with elements found in 'from'
func AddMap[K comparable, V any](to, from map[K]V) error {
	for k, v := range from {
		if _, ok := to[k]; ok {
			return ErrDuplicateKey
		}
		to[k] = v
	}
	return nil
}
