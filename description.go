package description

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	errorTree "github.com/Konstantin8105/errors"
)

const ListName string = "Description"

var Description map[string]string = map[string]string{
	"filename":   "name of file",
	"filenames":  "names of files",
	"buf":        "bytes buffer for templorary strings",
	"rootFolder": "folder name for start seaching",
	"err":        "typical error",
	"warnings":   "list of warnings, you can ignore them",
	"unsort":     "unsorted list of variable names",
	"gofiles":    "prepared list of Go files, except Go test files",
	"lines":      "list of report strings",
	"ok":         "true if acceptable",

	// ignore
	"a": "", "bl": "", "cl": "", "ds": "", "et": "", "f": "",
	"d": "", "fset": "", "key": "", "value": "", "mt": "",
	"add": "", "vs": "", "vl": "", "v": "", "rep": "",
	"kv": "", "ListName": "", "fd": "", "gocode": "", "n": "",
	"name": "", "names": "", "ns": "", "val": "", "w": "",
}

type List struct {
	Names []string
	Descr map[string]string
}

func (l List) Report() (rep string, warnings error) {
	et := errorTree.New("report errors:")

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.TabIndent)
	for i := range l.Names {
		name := l.Names[i]
		if name == ListName {
			continue
		}
		if v, ok := l.Descr[name]; ok {
			if 0 < len(v) {
				fmt.Fprintf(w, "%s\t%s\n", name, v)
			}
			continue
		}
		et.Add(fmt.Errorf("Name `%s` have not description", name))
	}
	if et.IsError() {
		warnings = et
	}
	w.Flush()
	lines := strings.Split(buf.String(), "\n")
	sort.Strings(lines)
	rep = strings.Join(lines, "\n")
	return
}

func New(rootFolder string) (list *List, err error) {
	list = new(List)
	list.Descr = map[string]string{} // initialization

	// find all Golang files
	var filenames []string
	filenames, err = files(rootFolder)
	if err != nil {
		return
	}

	// find all name of variables
	var unsort []string
	for _, filename := range filenames {
		var ns []string
		var ds map[string]string
		ns, ds, err = fileVars(filename)
		if err != nil {
			return
		}
		unsort = append(unsort, ns...)
		for k, v := range ds {
			list.Descr[k] = v
		}
	}

	// sort names
	sort.Strings(unsort)

	// uniq names
	list.Names = []string{}
	for i := range unsort {
		if i == 0 {
			list.Names = append(list.Names, unsort[i])
			continue
		}
		if unsort[i-1] == unsort[i] {
			continue
		}
		list.Names = append(list.Names, unsort[i])
	}

	return
}

func Template(rootFolder string) (gocode string) {
	n, _ := New(rootFolder)
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.TabIndent)
	for _, name := range n.Names {
		if name == ListName {
			continue
		}
		fmt.Fprintf(w, "    \"%s\":\t\"\",\n", name)
	}
	w.Flush()
	return "var Description map[string]string = map[string]string{\n" +
		buf.String() +
		"}"
}

func files(rootFolder string) (gofiles []string, err error) {
	filenames, err := ioutil.ReadDir("./")
	if err != nil {
		return
	}

	for _, info := range filenames {
		if info.IsDir() {
			continue
		}
		if !strings.HasSuffix(info.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(info.Name(), "_test.go") {
			continue
		}

		// Find Golang file
		gofiles = append(gofiles, filepath.Join(rootFolder, info.Name()))
	}

	return
}

func fileVars(filename string) (names []string, descr map[string]string, err error) {
	descr = map[string]string{} // initialization

	// positions are relative to fset
	fset := token.NewFileSet()

	// parse file
	var f *ast.File
	f, err = parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}

	// create new tree error
	et := errorTree.New(filename)

	add := func(n ast.Node) {
		if name, err := convert(n); err == nil {
			if name == "_" {
				return
			}
			names = append(names, name)
		} else {
			et.Add(err)
		}
	}

	// inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		if a, ok := n.(*ast.AssignStmt); ok && a.Tok == token.DEFINE {
			for i := range a.Lhs {
				add(a.Lhs[i])
			}
		}
		if d, ok := n.(*ast.ValueSpec); ok {
			for i := range d.Names {
				add(d.Names[i])
			}
		}
		if fd, ok := n.(*ast.FuncDecl); ok && fd.Type != nil {
			for _, v := range []*ast.FieldList{fd.Type.Params, fd.Type.Results} {
				if v == nil {
					continue
				}
				for _, f := range v.List {
					for _, id := range f.Names {
						add(id)
					}
				}
			}
		}
		if vs, ok := n.(*ast.ValueSpec); ok && 1 == len(vs.Names) && 1 == len(vs.Values) {
			if vs.Names[0].Name == ListName {
				isTypeOk := false
				if mt, ok := vs.Type.(*ast.MapType); ok {
					if key, ok := mt.Key.(*ast.Ident); ok && key.Name == "string" {
						if val, ok := mt.Value.(*ast.Ident); ok && val.Name == "string" {
							isTypeOk = true
						}
					}
				}
				if isTypeOk {
					if cl, ok := vs.Values[0].(*ast.CompositeLit); ok {
						for _, v := range cl.Elts {
							kv, ok := v.(*ast.KeyValueExpr)
							if !ok {
								continue
							}
							bl, ok := kv.Key.(*ast.BasicLit)
							if !ok {
								continue
							}
							if bl.Kind != token.STRING {
								continue
							}
							vl, ok := kv.Value.(*ast.BasicLit)
							if !ok {
								continue
							}
							if vl.Kind != token.STRING {
								continue
							}
							key := bl.Value
							key = strings.TrimPrefix(key, "\"")
							key = strings.TrimSuffix(key, "\"")
							value := vl.Value
							value = strings.TrimPrefix(value, "\"")
							value = strings.TrimSuffix(value, "\"")
							descr[key] = value
						}
					}
				}
			}
		}
		return true
	})

	// error handling
	if et.IsError() {
		err = et
	}

	return
}

func convert(n ast.Node) (name string, err error) {
	switch a := n.(type) {
	case *ast.Ident:
		name = a.Name
		return
	}
	err = fmt.Errorf("cannot convert ast.Node: %#v", n)
	return
}
