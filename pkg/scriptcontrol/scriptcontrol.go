package scriptcontrol

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"unicode"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
	"github.com/traefik/yaegi/interp"
)

type trackfunc func(*traincontrol.TrainControl)

type ScriptControl struct {
	funcs map[string]Func
}

type Func struct {
	Name        string
	Description string
	Func        trackfunc
}

func LoadFromDir(interp *interp.Interpreter, dir string) (map[string]Func, error) {
	return nil, nil
}

// LoadFromFile loads all train control functions from the file and interprets them
func LoadFromFile(interp *interp.Interpreter, path string) (map[string]Func, error) {
	out := make(map[string]Func)

	// Parse file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("unable to parse file '%s': %w", path, err)
	}

	packageName := node.Name.Name

	// Interpret file
	_, err = interp.EvalPath(path)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret file '%s': %w", path, err)
	}

	for _, f := range node.Decls {
		// Parse raw file for name and description
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		funcName := fn.Name.Name

		// Skip not exported functions
		if !unicode.IsUpper(rune(funcName[0])) {
			continue
		}

		// Interpret file for func
		v, err := interp.Eval(fmt.Sprintf("%s.%s", packageName, funcName))
		if err != nil {
			return nil, err
		}

		trFunc, ok := v.Interface().(trackfunc)
		if !ok {
			continue
		}

		out[funcName] = Func{
			Name:        funcName,
			Description: fn.Doc.Text(),
			Func:        trFunc,
		}
	}

	return out, nil
}
