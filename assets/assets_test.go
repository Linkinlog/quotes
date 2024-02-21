package assets_test

import (
	"testing"

	"github.com/Linkinlog/quotes/assets"
)

func TestFiles(t *testing.T) {
	_, err := assets.Files().Open("assets.go")
	if err != nil {
		t.Errorf("expected to find assets.go; got %v", err)
	}
}
