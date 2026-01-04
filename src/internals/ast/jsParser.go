package ast

import (
	"fmt"
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
	jsGrammar "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

// query list of strings with all the queries to analyse
var queries = []string{
	// Error
	"(ERROR) @error",
	// Normal Function
	"(function_declaration name: (identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (statement_block) @function.body ) @function",
	// Arrow function
	"(variable_declarator name: (identifier) @function.name " +
		"value: (arrow_function parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body )) @function",
	// Class methods
	"(method_definition name: (property_identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (statement_block) @function.body ) @function",
}

func Test(code []byte) {
	// Get our ast bases in our code and grammar
	ast, err := GetAST(code, tree.NewLanguage(jsGrammar.Language()))
	if err != nil {
		return
	}
	// Get the root node (program) from the AST generated
	root := ast.RootNode()
	defer ast.Close()

	// Get the query, cursor and captures from the helper function
	query, cursor, captures := GetCapturesByQueries(tree.NewLanguage(jsGrammar.Language()),
		strings.Join(queries, " "), code, root)
	// Defer the closing (because the iterative process uses this two variables)
	defer query.Close()
	defer cursor.Close()
	fmt.Printf("%s\n\n", root.ToSexp())

	// Executes 'till there's not a next capture
	for {
		// Get the next capture
		value := captures.Next()
		// If there´s no next, we break the loop
		if value == nil {
			break
		}
		// Look up to the first node type and decide the function to run.
		switch query.CaptureNames()[value.Captures[0].Index] {
		case "function":
			funcDecl(value.Captures[1:], code)
		}
	}
}

// funcDecl manages the function_declaration nodes (children).
func funcDecl(captures []tree.QueryCapture, code []byte) {
	fmt.Println("Function detected:")
	for _, capture := range captures {
		node := capture.Node
		fmt.Printf("%s: %s\n", node.Kind(), code[node.StartByte():node.EndByte()])
		if node.Kind() == "formal_parameters" {
			function := &Function{
				name: string(code[captures[0].Node.StartByte():captures[0].Node.EndByte()]),
			}
			switch {
			case node.NamedChildCount() > 8:
				function.SetStatus("DANGEROUS")
				Warnings = append(Warnings, function)
			}
		}
	}
	fmt.Println()
}
