package ast

import (
	"fmt"
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
	jsGrammar "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

// queries query list of strings with all the queries to analyse functions
var queries = []string{
	// Error
	"(ERROR) @error",
	// Normal Function
	"(function_declaration name: (identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (statement_block) @function.body )",
	// Arrow function
	"(variable_declarator name: (identifier) @function.name " +
		"value: (arrow_function parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body ))",
	// Class methods
	"(method_definition name: (property_identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (statement_block) @function.body )",
	/*"[(if_statement) (else_clause) (binary_expression) (switch_statement) (catch_clause)" +
	"(while_statement) (do_statement) (expression_statement) (for_statement)] @keyword",*/
}

func Test(code []byte) {
	language := tree.NewLanguage(jsGrammar.Language())
	// Get our ast bases in our code and grammar
	ast, err := GetAST(code, language)
	if err != nil {
		return
	}
	// Get the root node (program) from the AST generated
	root := ast.RootNode()
	defer ast.Close()

	// Get the query, cursor and captures from the helper function
	query, cursor, _ := GetCapturesByQueries(language, strings.Join(queries, " "), code, root)
	// Defer the closing (because the iterative process uses this two variables)
	defer query.Close()
	defer cursor.Close()
	fmt.Printf("%s\n\n", root.ToSexp())

	CyclicalComplexity(language, strings.Join(queries, " "), root, &RegexComplexity{
		// Keywords: Their body will be visited, and they will sum up + 1
		keyword: "(if_statement|else_clause|binary_expression|catch_clause|" +
			"while_statement|do_statement|expression_statement|for_statement|switch_case|ternary_expression)",
		// Body statements: Their body will be visited, but they don't sum up by themselves
		bodyStatements: "(body|statement_block|switch_statement|lexical_declaration|variable_declarator)",
		andKw:          "&&",
		orKw:           "||",
		code:           code,
	})
}
