package languages

import (
	"CLI_App/cmd/adapters/analysis"
	"CLI_App/cmd/adapters/analysis/types"
	"CLI_App/cmd/domain"

	tree "github.com/tree-sitter/go-tree-sitter"
	jsGrammar "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

// javascript interface with LanguageData embedded
type javascript struct {
	shouldFix    bool
	fileModifier analysis.FileModifier
	data         types.LanguageData
}

func NewJSLanguage(pattern string, shouldFix bool) types.NodeManagement {
	js := &javascript{
		shouldFix: shouldFix,
		data: types.LanguageData{
			Language: tree.NewLanguage(jsGrammar.Language()),
			Queries:  buildJSQuery(pattern),
		},
	}
	js.fileModifier = analysis.NewFileModifier(js)
	return js
}

func (js javascript) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	switch {
	case captureNames[node.Index] == "variable.name":
		// Set the initial feedback
		nodeInfo.SetVariableFeedback(string(code[node.Node.StartByte():node.Node.EndByte()]), domain.Point(node.Node.StartPosition()))
		js.fileModifier.ModifyVariableName(code, filepath, string(code[node.Node.StartByte():node.Node.EndByte()]), js.shouldFix)
		return
	case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "variable_declarator":
		return
	}
	nodeInfo.Complexity++
}

func buildJSQuery(pattern string) string {
	return "(function_declaration name: (identifier) @function.name " +
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
		"(#not-match? @variable.name \"" + domain.AllowNonNamedVar + "|" + pattern + "\"))" +
		"(for_in_statement left: ([(identifier) @variable.name (array_pattern (identifier) @variable.name)])" +
		"(#not-match? @variable.name \"" + domain.AllowNonNamedVar + "|" + pattern + "\"))" +
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
		"] @keyword"
}

func (js javascript) GetLanguageData() types.LanguageData {
	return js.data
}

func (js javascript) GetVarAppearancesQuery(name string) string {
	return name
}
