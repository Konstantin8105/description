# description
Generate and check description of all variables

```
package description // import "github.com/Konstantin8105/description"


CONSTANTS

const ListName string = "Description"
    ListName is name of typical description list


VARIABLES

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

	"a": "", "bl": "", "cl": "", "ds": "", "et": "", "f": "",
	"d": "", "fset": "", "key": "", "value": "", "mt": "",
	"add": "", "vs": "", "vl": "", "v": "", "rep": "",
	"kv": "", "ListName": "", "fd": "", "gocode": "", "n": "",
	"name": "", "names": "", "ns": "", "val": "", "w": "",
}
    Description of package variables


FUNCTIONS

func Template(rootFolder string) (gocode string)
    Template return Go code for all variable names


TYPES

type List struct {
	// List of variable names
	Names []string

	// Founded description list in Go files
	Descr map[string]string
}
    List of descriptions parts

func New(rootFolder string) (list *List, err error)
    New return description data

func (l List) Report() (rep string, warnings error)
    Report return typical representation of description and warning for
    variables without description

```
