package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

var filePath string

func init() {
	flag.StringVar(&filePath, "path", "", "Path to Go file.")
	flag.Parse()
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

	fmt.Println("| Constant | Description |")
	fmt.Println("| --- | --- |")
	for constName, comment := range table {
		fmt.Printf("| %s | %s |\n", constName, comment)
	}
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
	value = strings.Join(parts, " ")
	if _, ok := (*table)[key]; ok {
		(*table)[key] += fmt.Sprintf(" %s", value)
	} else {
		(*table)[key] = value
	}
}
