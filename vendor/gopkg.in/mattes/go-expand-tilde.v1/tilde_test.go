package tilde

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestHome(t *testing.T) {
	if _, err := Home(); err != nil {
		t.Fatal(err)
	}
}

func TestExpand(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		// TODO

	default:
		home, _ := Home()
		h, _ := Expand("~")
		if h != home {
			t.Errorf("Expected '%v', got '%v'", home, h)
		}

		h2, _ := Expand("~/")
		if h2 != home+"/" {
			t.Errorf("Expected '%v', got '%v'", home, h2)
		}

		h3, _ := Expand("~/path")
		if h3 != filepath.Join(home, "path") {
			t.Errorf("Expected '%v', got '%v'", home, h3)
		}
	}
}
