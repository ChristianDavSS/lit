package grammar

import (
	"CLI_App/src/config"
	"CLI_App/src/internals/ast/languages"
	"CLI_App/src/internals/utils"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	goGrammar "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

var GoLanguage = languages.LanguageInformation{
	Language: tree.NewLanguage(goGrammar.Language()),
	Queries:
	// Functions
	"(function_declaration name: (_) @function.name " +
		"parameters: (_) @function.parameters " +
		"body: (_) @function.body ) @function " +
		// Variable names
		"(expression_list (identifier) @variable.name" +
		"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.GetActiveNamingConvention() + "\"))" +
		"(var_declaration (var_spec name: (identifier) @variable.name)" +
		"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.GetActiveNamingConvention() + "\"))" +
		// Keywords
		"[" +
		"(if_statement) (for_statement) (expression_case)" +
		// Binary expressions
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"] @keyword",
	ManageNode: func(captureNames []string, code []byte, node tree.QueryCapture, nodeInfo *languages.FunctionData) {
		// Search the 'alternative' node in the children
		alternative := node.Node.ChildByFieldName("alternative")
		switch {
		case captureNames[node.Index] == "variable.name":
			nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
				string(code[node.Node.StartByte():node.Node.EndByte()]),
				node.Node.StartPosition().Row, node.Node.StartPosition().Column,
			)
			return
		case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "expression_list":
			return
		case alternative != nil && alternative.GrammarName() == "block":
			nodeInfo.Complexity++
		}
		nodeInfo.Complexity++
	},
}
