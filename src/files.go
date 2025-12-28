package src

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type Error struct {
	message  string
	position token.Pos
}

type BasicLiteral struct {
	kind     token.Token
	value    string
	position token.Pos
}

type FunctionCall struct {
	args     []ast.Expr
	position token.Pos
}

type FunctionDecl struct {
	parameters *ast.FieldList
	position   token.Pos
}

type CodeInfo interface {
	GetPosition() token.Pos
}

// Interface impl
func (e *Error) GetPosition() token.Pos {
	return e.position
}

func (b *BasicLiteral) GetPosition() token.Pos {
	return b.position
}

func (f *FunctionCall) GetPosition() token.Pos {
	return f.position
}

func (f *FunctionDecl) GetPosition() token.Pos {
	return f.position
}

func Init() {
	tokens := getFileTokens("src/root.go")
	manageToken(tokens)
}

// Functionality
func getFileTokens(filename string) []CodeInfo {
	var Parsed []CodeInfo
	// Get a new fileset for our AST
	fileset := token.NewFileSet()
	// Parse a certain script to scan it
	node, err := parser.ParseFile(fileset, filename, nil, 0)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	// Use the Inspect() function to traverse the AST
	ast.Inspect(node, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.BasicLit:
			Parsed = append(Parsed, &BasicLiteral{
				kind:     x.Kind,
				value:    x.Value,
				position: x.Pos(),
			})
		case *ast.CallExpr:
			Parsed = append(Parsed, &FunctionCall{
				args:     x.Args,
				position: x.Pos(),
			})
		case *ast.FuncDecl:
			Parsed = append(Parsed, &FunctionDecl{
				parameters: x.Type.Params,
				position:   x.Pos(),
			})
			/* Errors cases */
		case *ast.BadDecl:
			Parsed = append(Parsed, &Error{
				message:  "Bad declaration (syntax error)",
				position: x.Pos(),
			})
		case *ast.BadStmt:
			Parsed = append(Parsed, &Error{
				message:  "Bas statement.",
				position: x.Pos(),
			})
		case *ast.BadExpr:
			Parsed = append(Parsed, &Error{
				message:  "Bad expresion.",
				position: x.Pos(),
			})
		}
		return true
	})
	return Parsed
}

func manageToken(tokens []CodeInfo) {
	for _, v := range tokens {
		fmt.Println(v, v.GetPosition())
	}
}
