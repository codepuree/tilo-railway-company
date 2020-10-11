package custom

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
	"github.com/codepuree/tilo-railway-company/pkg/trclib"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// func TestTrack1(t *testing.T) {
// 	// tests := []struct {
// 	// 	name string
// 	// }{
// 	// 	// TODO: Add test cases.
// 	// }
// 	// for _, tt := range tests {
// 	// 	t.Run(tt.name, func(t *testing.T) {
// 	// 		Track1()
// 	// 	})
// 	// }

// 	i := interp.New(interp.Options{})
// 	i.Use(stdlib.Symbols)

// 	// _, err := i.Use()

// 	// _, err := i.EvalPath("E:/Users/LC/Documents/Projects/tilo-railway-company/internal/custom/track1.go")
// 	_, err := i.Eval(`package custom

// 	import "github.com/codepuree/tilo-railway-company/pkg/traincontrol"
// 	import "fmt"

// 	func Track1(tc traincontrol.TrainControl) {
// 		fmt.Println("Hello tc", tc)
// 	}`)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	_, err = i.Eval("custom.Track1")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	// tc := traincontrol.TrainControl{}

// 	// Track1 := v.Interface().(func(traincontrol.TrainControl))

// 	// Track1(tc)
// }

type TrainController interface {
	SetSwitch(id string, state string)
}

func TestTrack1(t *testing.T) {
	const srcPath = "E:/Users/LC/Documents/Projects/tilo-railway-company/internal/custom/track1.go"

	srcb, err := ioutil.ReadFile(srcPath)
	if err != nil {
		t.Error(err)
	}

	src := string(srcb)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "track1.go", src, parser.ImportsOnly)
	if err != nil {
		t.Error(err)
	}

	// Get package name
	// packageName := f.Name.Name

	// Get exported functions
	fmt.Println("Functions:")
	for _, f := range f.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fmt.Println(fn.Name.Name)
	}

	// Get some comments
	t.Error("Some comments:")
	for _, c := range f.Comments {
		t.Error(c.Text())
	}

	// Print the imports from the file's AST.
	for _, s := range f.Imports {
		t.Error(s.Path.Value)
	}

	i := interp.New(interp.Options{
		
	})
	i.Use(stdlib.Symbols)
	i.Use(trclib.Symbols)

	// _, err := i.Eval(src)
	_, err = i.EvalPath(srcPath)
	if err != nil {
		t.Error(err)
	}

	v, err := i.Eval("custom.Track")
	if err != nil {
		t.Error(err)
	}

	track1 := v.Interface().(func(*traincontrol.TrainControl))
	tc := &traincontrol.TrainControl{}

	track1(tc)
}
