package grammar

import (
	"CLI_App/src/config"
	"CLI_App/src/internals/analysis/utils"
	"CLI_App/src/internals/ast/languages"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	pyGrammar "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

var PyLanguage = languages.LanguageInformation{
	Language: tree.NewLanguage(pyGrammar.Language()),
	Queries:
	// Function definition
	"(function_definition name: (identifier) @function.name " +
		"parameters: (parameters) @function.parameters " +
		"body: (block) @function.body) @function " +
		// Variable names
		"(assignment left: ([(identifier) @variable.name (pattern_list (identifier) @variable.name)])" +
		"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.GetActiveNamingConvention() + "\"))" +
		"(for_statement left: ([(identifier) @variable.name (pattern_list (identifier) @variable.name)])" +
		"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.GetActiveNamingConvention() + "\"))" +
		// Keywords
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
	ManageNode: func(captureNames []string, code []byte, node tree.QueryCapture, nodeInfo *languages.FunctionData) {
		switch {
		case captureNames[node.Index] == "variable.name":
			nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
				string(code[node.Node.StartByte():node.Node.EndByte()]),
				node.Node.StartPosition().Row, node.Node.StartPosition().Column,
			)
			return
		case node.Node.GrammarName() == "boolean_operator" && node.Node.Parent().GrammarName() == "assignment":
			return
		}
		nodeInfo.Complexity++
	},
}
