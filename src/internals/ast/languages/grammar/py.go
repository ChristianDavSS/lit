package grammar

import (
	"CLI_App/src/internals/ast/languages"

	tree "github.com/tree-sitter/go-tree-sitter"
	pyGrammar "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

var PyLanguage = languages.LanguageInformation{
	Language: tree.NewLanguage(pyGrammar.Language()),
	Queries: "(function_definition name: (identifier) @function.name " +
		"parameters: (parameters) @function.parameters " +
		"body: (block) @function.body) @function " +
		"[" +
		// If, else-if, else
		"(if_statement condition: (_)) (elif_clause condition: (_)) (else_clause body: (_))" +
		// Loops
		"(for_statement) (while_statement condition: (_) body: (_))" +
		// Operators
		"(boolean_operator left: (_) right: (_))" +
		// Clauses
		"(except_clause value: (_)) (conditional_expression) (case_clause (_))" +
		// List comprehension
		"(list_comprehension body: (_) (for_in_clause left: (_) right: (_))) (if_clause (_))" +
		"] @keyword",
	RegexComplexity: &languages.RegexComplexity{
		ManageNode: func(captureNames []string, code []byte, node tree.QueryCapture, complexity *int) {
			if node.Node.GrammarName() == "boolean_operator" && node.Node.Parent().GrammarName() == "assignment" {
				return
			}
			*complexity++
		},
	},
}
