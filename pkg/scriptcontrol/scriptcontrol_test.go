package scriptcontrol

import (
	"fmt"
	"testing"

	"github.com/codepuree/tilo-railway-company/pkg/trclib"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func TestLoadFromFile(t *testing.T) {
	interp := interp.New(interp.Options{})
	interp.Use(stdlib.Symbols)
	interp.Use(trclib.Symbols)

	funcs, err := LoadFromFile(interp, "E:/Users/LC/Documents/Projects/tilo-railway-company/testdata/track1.go")
	if err != nil {
		t.Error(fmt.Errorf("unable to load from file: %w", err))
	}

	t.Errorf("%+v", funcs)
}
