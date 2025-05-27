package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

var filePath string
var tmplFilePath string

func init() {
	flag.StringVar(&filePath, "path", "", "Path to Go file.")
	flag.StringVar(&tmplFilePath, "template-path", "", "Path to a Go template file.")
	flag.Parse()
}

type Template struct {
	GoConstTable string
}

func main() {
	if filePath == "" {
		fmt.Println("File path not provided.")
		return
	}

	src, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	table := make(map[string]string)

	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		if genDecl.Tok != token.CONST {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			constName := valueSpec.Names[0].Name
			addTableEntry(constName, genDecl.Doc.Text(), &table)

			if valueSpec.Doc != nil {
				for _, comment := range valueSpec.Doc.List {
					constName := valueSpec.Names[0].Name
					addTableEntry(constName, comment.Text, &table)
				}
			}
		}

		return true
	})

	tmplStr := DefaultTemplate
	if tmplFilePath != "" {
		tmplBytes, err := os.ReadFile(tmplFilePath)
		if err != nil {
			panic(err)
		}
		tmplStr = string(tmplBytes)
	}

	tmpl, err := template.New("goconsttable").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	tableStr := strings.Builder{}
	tableStr.WriteString("| Constant | Description |\n")
	tableStr.WriteString("| --- | --- |\n")

	for constName, comment := range table {
		tableStr.WriteString(fmt.Sprintf("| %s | %s |\n", constName, comment))
	}

	tmpl.Execute(os.Stdout, Template{ GoConstTable: tableStr.String() })
}

func addTableEntry(key, value string, table *map[string]string) {
	parts := strings.Split(value, "\n")
	for i := range parts {
		text := strings.Trim(parts[i], "/")
		text = strings.TrimSpace(text)
		if text != "" {
			parts[i] = text
		}
	}
	comment := strings.Join(parts, " ")
	current := (*table)[key]
	if current == "" {
		(*table)[key] += fmt.Sprintf("%s", comment)
	} else {
		(*table)[key] += fmt.Sprintf(" %s", comment)
	}
}
