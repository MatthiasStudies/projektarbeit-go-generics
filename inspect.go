package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

const inspectCode = `
package main

type MyInt int

func isEven(n MyInt) bool {
	return n%2 == 0
}

type MyStruct struct {
	Field1 string
	Field2 int
}

func main() {
	x := MyInt(42)
	x = MyInt(43)
	s := MyStruct{Field1: "hello", Field2: 10}
	// inspect: MyStruct, 1, s, s.Field1
	_ = x
	_ = s
}
`

const inspectPrefix = "inspect:"

func findLookupNames(commentText string) []string {
	if !strings.HasPrefix(commentText, inspectPrefix) {
		return nil
	}
	text := strings.TrimPrefix(commentText, inspectPrefix)
	parts := strings.Split(text, ",")
	var names []string
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}

func formatObj(fset *token.FileSet, obj types.Object) string {
	if obj == nil {
		return "\t<not found>\n"
	}
	pos := fset.Position(obj.Pos())

	buff := &strings.Builder{}
	fmt.Fprintf(buff, "\tKind: %T\n", obj)
	fmt.Fprintf(buff, "\tType: %s\n", obj.Type().String())
	fmt.Fprintf(buff, "\tPkg: %v\n", obj.Pkg())
	fmt.Fprintf(buff, "\tPos: %v\n", pos)
	if v, ok := obj.(*types.Var); ok {
		fmt.Fprintf(buff, "\tVar isExported: %v\n", v.Exported())
	}
	if f, ok := obj.(*types.Func); ok {
		sig := f.Type().(*types.Signature)
		fmt.Fprintf(buff, "\tFunc Params: %s\n", sig.Params().String())
		fmt.Fprintf(buff, "\tFunc Results: %s\n", sig.Results().String())
	}
	underlying := obj.Type().Underlying()
	fmt.Fprintf(buff, "\tUnderlying Type: %T %s\n", underlying, underlying.String())
	return buff.String()
}

func printObj(fset *token.FileSet, pos token.Pos, name string, obj types.Object) {
	fmt.Printf("%s,\t%q\n", fset.Position(pos), name)
	fmt.Println(formatObj(fset, obj))
}

func main() {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "test.go", inspectCode, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	conf := types.Config{
		Importer: importer.Default(),
	}

	pkg, err := conf.Check("main", fset, []*ast.File{f}, nil)
	if err != nil {
		panic(err)
	}

	for _, comment := range f.Comments {
		names := findLookupNames(comment.Text())
		if names == nil {
			continue
		}

		pos := comment.Pos()
		scope := pkg.Scope().Innermost(pos) // Find the scope closest to the comment position

		for _, name := range names {
			_, obj := scope.LookupParent(name, pos)
			printObj(fset, pos, name, obj)
		}
	}
}
