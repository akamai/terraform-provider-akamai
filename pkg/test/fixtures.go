package test

import (
	"fmt"
	"io/ioutil"
)

// FixtureBytes returns the entire contents of the given file as a byte slice. Path can be given as Sprintf format and
// args. Panics on error.
func FixtureBytes(path string, args ...interface{}) []byte {
	contents, err := ioutil.ReadFile(fmt.Sprintf(path, args...))
	if err != nil {
		panic(err)
	}
	return contents
}

// Fixture returns the entire contents of the given file as a string. Path be given as Sprintf format and args. Panics
// on error
func Fixture(path string, args ...interface{}) string {
	return string(FixtureBytes(path, args...))
}
