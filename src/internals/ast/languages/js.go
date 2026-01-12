package languages

import (
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
	jsGrammar "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

var JsLanguage = LanguageInformation{
	Language: tree.NewLanguage(jsGrammar.Language()),
	Queries:
	// Normal Function
	"(function_declaration name: (identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters) @function" +
		// Arrow Function
		"(variable_declarator name: (identifier) @function.name " +
		"value: (arrow_function parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body )) @function" +
		// Class Methods
		"(method_definition name: (property_identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters) @function" +
		// Functions body information (keywords)
		"[" +
		// if, else-if, else
		"(if_statement condition: (_) consequence: (_) alternative: (else_clause)?)" +
		"(else_clause (statement_block))" +
		// Statements
		"(for_statement) (for_in_statement) (while_statement) (switch_case) (catch_clause)" +
		// Expressions
		"(binary_expression left: (_) right: (_)) (ternary_expression)" +
		// ForEach, Map, etc.
		"(call_expression function: (member_expression object: (_) property: (_) @call.name)" +
		"arguments: (arguments (arrow_function)) (#match? @call.name \"^(forEach|map)$\"))" +
		"] @keyword",
	RegexComplexity: &RegexComplexity{
		ManageNode: func(captureNames []string, code []byte, node tree.QueryCapture, complexity *int) {
			if node.Node.GrammarName() == "binary_expression" {
				line := string(code[node.Node.StartByte():node.Node.EndByte()])
				totalBooleans := strings.Count(line, "&&") + strings.Count(line, "||")
				if totalBooleans < 1 || node.Node.Parent().GrammarName() == "variable_declarator" {
					return
				}
			}
			*complexity++
		},
		MainFunc: &FunctionData{Name: "Main", Parameters: "()", TotalParams: 0, Complexity: 1},
	},
}
