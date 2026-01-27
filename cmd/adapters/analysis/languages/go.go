package languages

import (
	"CLI_App/cmd/adapters/analysis/types"
	"CLI_App/cmd/domain"

	tree "github.com/tree-sitter/go-tree-sitter"
	goGrammar "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

type golang struct {
	data types.LanguageData
}

func NewGolangLanguage(pattern string) types.NodeManagement {
	return &golang{
		data: types.LanguageData{
			Language: tree.NewLanguage(goGrammar.Language()),
			Queries:  buildGolangQuery(pattern),
		},
	}
}

func (g golang) ManageNode(captureNames []string, code *[]string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		g.variableManagement(node, nodeInfo, code)
		return
	case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "expression_list":
		return
	case alternative != nil && alternative.GrammarName() == "block":
		nodeInfo.Complexity++
	}
	nodeInfo.Complexity++
}

func (g golang) variableManagement(varNode tree.QueryCapture, functionData *domain.FunctionData, code *[]string) {
	// Set the initial feedback
	functionData.SetVariableFeedback(
		(*code)[varNode.Node.StartPosition().Row][varNode.Node.StartPosition().Column:varNode.Node.EndPosition().Column],
		domain.Point(varNode.Node.StartPosition()),
	)
}

func buildGolangQuery(pattern string) string {
	return "(function_declaration name: (_) @function.name " +
		"parameters: (_) @function.parameters " +
		"body: (_) @function.body ) @function " +
		// Variable names
		"(expression_list (identifier) @variable.name" +
		"(#not-match? @variable.name \"" + domain.AllowNonNamedVar + "|" + pattern + "\"))" +
		"(var_declaration (var_spec name: (identifier) @variable.name)" +
		"(#not-match? @variable.name \"" + domain.AllowNonNamedVar + "|" + pattern + "\"))" +
		// Keywords
		"[" +
		"(if_statement) (for_statement) (expression_case)" +
		// Binary expressions
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"] @keyword"
}

func (g golang) GetLanguageData() types.LanguageData {
	return g.data
}

func (g golang) GetVarAppearancesQuery(name string) string {
	return name
}
