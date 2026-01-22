package languages

import (
	"CLI_App/src/internals/ast/config"
	language "CLI_App/src/internals/ast/utils"
	"CLI_App/src/internals/ast/writer"
	"CLI_App/src/internals/utils"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	goGrammar "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

type golang struct {
	data language.LanguageData
}

func (g golang) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *language.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
			string(code[node.Node.StartByte():node.Node.EndByte()]),
			node.Node.StartPosition().Row, node.Node.StartPosition().Column,
		)
		writer.ModifyVariableName(g, node.Node, code, filepath)
		return
	case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "expression_list":
		return
	case alternative != nil && alternative.GrammarName() == "block":
		nodeInfo.Complexity++
	}
	nodeInfo.Complexity++
}

func (g golang) GetLanguage() *tree.Language {
	return g.data.Language
}

func (g golang) GetQueries() string {
	return g.data.Queries
}

func (g golang) GetVarAppearancesQuery(name string) string {
	return name
}

var GoLanguage = golang{
	language.LanguageData{
		Language: tree.NewLanguage(goGrammar.Language()),
		Queries:
		// Functions
		"(function_declaration name: (_) @function.name " +
			"parameters: (_) @function.parameters " +
			"body: (_) @function.body ) @function " +
			// Variable names
			"(expression_list (identifier) @variable.name" +
			"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.ActivePattern + "\"))" +
			"(var_declaration (var_spec name: (identifier) @variable.name)" +
			"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.ActivePattern + "\"))" +
			// Keywords
			"[" +
			"(if_statement) (for_statement) (expression_case)" +
			// Binary expressions
			"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
			"] @keyword",
	},
}
