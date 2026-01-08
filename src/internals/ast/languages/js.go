package languages

import (
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
	jsGrammar "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

var JsLanguage = LanguageInformation{
	Language: tree.NewLanguage(jsGrammar.Language()),
	Queries: []string{
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
	},
	RegexComplexity: &RegexComplexity{
		// Keywords: Their body will be visited, and they will sum up + 1
		Keyword: "(if_statement|else_clause|catch_clause|for_in_statement|" +
			"while_statement|do_statement|for_statement|switch_case|ternary_expression)",
		// Body statements: Their body will be visited, but they don't sum up by themselves
		BodyStatements: "(body|statement_block|switch_statement|lexical_declaration|variable_declarator)",
		KeywordMatchFunc: func(node *tree.Node, complexity *int) {
			if node.GrammarName() == "else_clause" {
				if node.Child(node.ChildCount()-1).GrammarName() == "if_statement" {
					*complexity++
				}
				return
			}
			*complexity++
		},
		NoMatchFunc: func(node *tree.Node, complexity *int, code []byte) {
			line := string(code[node.StartByte():node.EndByte()])
			*complexity += strings.Count(line, "&&") + strings.Count(line, "||")
		},
	},
}
