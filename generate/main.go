package main

import (
	"flag"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"html/template"
	"log"
	"os"

	"github.com/facebookgo/structtag"
)

var (
	typeName = flag.String("type", "", "type name; must be set")
	output   = flag.String("output", "", "output file name; default srcdir/<type>_factory.go")
)

func main() {
	flag.Parse()

	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	buildPkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}

	fs := token.NewFileSet()
	var files []*ast.File
	for _, goFile := range buildPkg.GoFiles {
		parsedFile, err := parser.ParseFile(fs, goFile, nil, 0)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, parsedFile)
	}

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(buildPkg.Name, fs, files, nil)
	if err != nil {
		log.Fatal(err)
	}

	target := pkg.Scope().Lookup(*typeName)
	structType := target.Type().Underlying().(*types.Struct)
	var templateData struct {
		Type   string
		Fields []TemplateField
	}

	templateData.Type = *typeName
	templateData.Fields = make([]TemplateField, 0, 0)
	for i := 0; i < structType.NumFields(); i++ {
		tag := structType.Tag(i)
		if ok, _, err := structtag.Extract("inject", tag); err != nil {
			log.Fatal(err)
		} else if !ok {
			continue
		}
		field := structType.Field(i)
		templateData.Fields = append(templateData.Fields, TemplateField{
			ArgName:   string(97 + i),
			FieldName: field.Name(),
			Type:      field.Type().String(),
		})
	}

	tmpl, err := template.New("module").Parse(`func Provide{{.Type}}({{range $i, $e := .Fields}}{{if ne $i 0 }}, {{end}}{{ .ArgName }} {{ .Type }}{{end}}) *{{.Type}} {
		return &{{.Type}}{
			{{range .Fields}}{{ .FieldName }}: {{ .ArgName }}, {{end}}
		}
}`)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(os.Stdout, templateData)
	if err != nil {
		log.Fatal(err)
	}
}

type TemplateField struct {
	ArgName   string
	FieldName string
	Type      string
}
