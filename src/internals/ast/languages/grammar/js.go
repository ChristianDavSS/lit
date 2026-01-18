package grammar

import (
	"CLI_App/src/config"
	"CLI_App/src/internals/analysis/utils"
	"CLI_App/src/internals/ast/languages"

	tree "github.com/tree-sitter/go-tree-sitter"
	jsGrammar "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

var JsLanguage = languages.LanguageInformation{
	Language: tree.NewLanguage(jsGrammar.Language()),
	Queries:
	// Normal Function
	"(function_declaration name: (identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body ) @function" +
		// Arrow Function
		"(variable_declarator name: (identifier) @function.name " +
		"value: (arrow_function parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body )) @function" +
		// Class Methods
		"(method_definition name: (property_identifier) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (_) @function.body ) @function" +
		// Variable names
		"(variable_declarator name: ([(identifier) @variable.name (array_pattern (identifier) @variable.name)])" +
		"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.ActiveNamingConvention + "\"))" +
		"(for_in_statement left: ([(identifier) @variable.name (array_pattern (identifier) @variable.name)])" +
		"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.ActiveNamingConvention + "\"))" +
		// Functions body information (keywords)
		"[" +
		// if, else-if, else
		"(if_statement condition: (_) consequence: (_) alternative: (else_clause)?)" +
		"(else_clause (statement_block))" +
		// Statements
		"(for_statement) (for_in_statement) (while_statement) (switch_case) (catch_clause)" +
		// Expressions
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"(ternary_expression)" +
		// ForEach, Map, etc.
		"(call_expression function: (member_expression object: (_) property: (_) @call.name)" +
		"arguments: (arguments (arrow_function)) (#match? @call.name \"^(forEach)$\"))" +
		"] @keyword",
	ManageNode: func(captureNames []string, code []byte, node tree.QueryCapture, nodeInfo *languages.FunctionData) {
		switch {
		case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "variable_declarator":
			return
		}
		nodeInfo.Complexity++
	},
}
