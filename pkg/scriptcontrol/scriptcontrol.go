package scriptcontrol

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"unicode"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
	"github.com/traefik/yaegi/interp"
)

type trackfunc func(tc traincontrol.TrainControl)

type ScriptControl struct {
	funcs map[string]Func
}

type Func struct {
	Name        string
	Description string
	Func        interface{}
}

// LoadFromDir parses and interprets all train control functions from the directory
func LoadFromDir(interp *interp.Interpreter, dir string) (map[string]Func, error) {
	allFuncs := make(map[string]Func)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to read files in directory: %w", err))
	}

	for _, file := range files {
		if file.IsDir() {
			log.Printf("only the root folder '%s' gets loaded and interpreted, '%s' does not.", dir, filepath.Join(dir, file.Name()))
			continue
		}

		funcs, err := LoadFromFile(interp, filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		allFuncs, err = join(allFuncs, funcs)
		if err != nil {
			return nil, fmt.Errorf("unable to append functions from '%s' to all functions in folder: %w", file.Name(), err)
		}
	}

	return allFuncs, nil
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

		trFunc, ok := v.Interface().(func(*traincontrol.TrainControl))
		if !ok {
			trFunc, ok := v.Interface().(func(*traincontrol.TrainControl, int))
			if !ok {
				trFunc, ok := v.Interface().(func(*traincontrol.TrainControl, string))
				if !ok {
					continue
				}

				out[funcName] = Func{
					Name:        funcName,
					Description: fn.Doc.Text(),
					Func:        trFunc,
				}
				continue
			}

			out[funcName] = Func{
				Name:        funcName,
				Description: fn.Doc.Text(),
				Func:        trFunc,
			}
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

func join(a, b map[string]Func) (map[string]Func, error) {
	for name, func_ := range b {
		if _, ok := a[name]; ok {
			return nil, fmt.Errorf("the function '%s' was already declared", name)
		}

		a[name] = func_
	}

	return a, nil
}
