package description_test

import (
	"fmt"
	"os"

	"github.com/Konstantin8105/description"
)

func Example() {
	descr, err := description.New(".")
	if err != nil {
		panic(err)
	}
	rep, err := descr.Report()
	fmt.Fprintf(os.Stdout, "Report:\n%s\n", rep)
	if err != nil {
		fmt.Fprintf(os.Stdout, "\n%v\n", err)
	}

	tmpl := description.Template(".")
	fmt.Fprintf(os.Stdout, "Golang template:\n%s\n", tmpl)

	// Output:
	// Report:
	//
	// buf        bytes buffer for templorary strings
	// err        typical error
	// filename   name of file
	// filenames  names of files
	// gofiles    prepared list of Go files, except Go test files
	// lines      list of report strings
	// ok         true if acceptable
	// rootFolder folder name for start seaching
	// unsort     unsorted list of variable names
	// warnings   list of warnings, you can ignore them
	//
	// report errors:
	// ├──Name `descr` have not description
	// ├──Name `isTypeOk` have not description
	// └──Name `list` have not description
	//
	// Golang template:
	// var Description map[string]string = map[string]string{
	//     "ListName":   "",
	//     "a":          "",
	//     "add":        "",
	//     "bl":         "",
	//     "buf":        "",
	//     "cl":         "",
	//     "d":          "",
	//     "descr":      "",
	//     "ds":         "",
	//     "err":        "",
	//     "et":         "",
	//     "f":          "",
	//     "fd":         "",
	//     "filename":   "",
	//     "filenames":  "",
	//     "fset":       "",
	//     "gocode":     "",
	//     "gofiles":    "",
	//     "isTypeOk":   "",
	//     "key":        "",
	//     "kv":         "",
	//     "lines":      "",
	//     "list":       "",
	//     "mt":         "",
	//     "n":          "",
	//     "name":       "",
	//     "names":      "",
	//     "ns":         "",
	//     "ok":         "",
	//     "rep":        "",
	//     "rootFolder": "",
	//     "unsort":     "",
	//     "v":          "",
	//     "val":        "",
	//     "value":      "",
	//     "vl":         "",
	//     "vs":         "",
	//     "w":          "",
	//     "warnings":   "",
	// }
}
