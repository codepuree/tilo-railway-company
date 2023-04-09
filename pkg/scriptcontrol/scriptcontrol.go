package scriptcontrol

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"unicode"

	"github.com/codepuree/tilo-railway-company/internal/communication"
	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
	"github.com/codepuree/tilo-railway-company/pkg/trclib"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type trackfunc func(tc traincontrol.TrainControl)

type ScriptControl struct {
	Funcs         map[string]Func
	directoryPath string
	interp        *interp.Interpreter
	// reader for the interpreter
	reader *io.PipeReader
}

type Func struct {
	Name        string
	Description string
	Func        interface{}
}

func NewScriptControl(directoryPath string) *ScriptControl {
	return &ScriptControl{
		directoryPath: directoryPath,
	}
}

func (sc *ScriptControl) Bind(pub communication.Publisher) {
	var writer *io.PipeWriter
	sc.reader, writer = io.Pipe()
	// defer reader.Close() // TODO: add `Close` method to scriptcontrol and close reader there
	rbuf := bufio.NewReader(sc.reader)
	go func() {
		lb, _, err := rbuf.ReadLine()
		if err != nil {
			log.Fatal(fmt.Errorf("unable to read form stdout reader: %w", err))
		}

		pub.SendToAll(communication.Message{
			From: "scriptcontrol",
			To:   "all",
			Data: string(lb),
		})
	}()
	sc.interp = interp.New(interp.Options{
		Stdout: writer,
	})
	sc.interp.Use(stdlib.Symbols)
	sc.interp.Use(trclib.Symbols)
}

func (sc *ScriptControl) Close() error {
	return sc.reader.Close()
}

func (sc *ScriptControl) Load() error {
	err := sc.LoadFromDir(sc.directoryPath)
	if err != nil {
		return fmt.Errorf("unable to load and interpret functions in directory '%s': %w", sc.directoryPath, err)
	}

	return nil
}

// LoadFromDir parses and interprets all train control functions from the directory
func (sc *ScriptControl) LoadFromDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("unable to read files in directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			log.Printf("only the root folder '%s' gets loaded and interpreted, '%s' does not.", dir, filepath.Join(dir, file.Name()))
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		err := sc.LoadFromFile(filePath)
		if err != nil {
			return fmt.Errorf("unable to load functions from file '%s': %w", filePath, err)
		}
	}

	return nil
}

// LoadFromFile loads all train control functions from the file and interprets them
func (sc *ScriptControl) LoadFromFile(path string) error {
	// Parse file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("unable to parse file '%s': %w", path, err)
	}

	packageName := node.Name.Name

	if sc.interp == nil {
		return fmt.Errorf("The interpreter for the script control that should load %s was not initialized!", path)
	}

	// Interpret file
	_, err = sc.interp.EvalPath(path)
	if err != nil {
		return fmt.Errorf("unable to interpret file '%s': %w", path, err)
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
		v, err := sc.interp.Eval(fmt.Sprintf("%s.%s", packageName, funcName))
		if err != nil {
			return err
		}

		trFunc, ok := v.Interface().(func(*traincontrol.TrainControl))
		if !ok {
			trFunc, ok := v.Interface().(func(*traincontrol.TrainControl, int))
			if !ok {
				trFunc, ok := v.Interface().(func(*traincontrol.TrainControl, string))
				if !ok {
					continue
				}

				sc.Funcs[funcName] = Func{
					Name:        funcName,
					Description: fn.Doc.Text(),
					Func:        trFunc,
				}
				continue
			}

			sc.Funcs[funcName] = Func{
				Name:        funcName,
				Description: fn.Doc.Text(),
				Func:        trFunc,
			}
			continue
		}

		sc.Funcs[funcName] = Func{
			Name:        funcName,
			Description: fn.Doc.Text(),
			Func:        trFunc,
		}
	}

	return nil
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
