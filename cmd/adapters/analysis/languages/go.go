package languages

import (
	"CLI_App/cmd/adapters/analysis"
	"CLI_App/cmd/adapters/analysis/types"
	"CLI_App/cmd/domain"

	tree "github.com/tree-sitter/go-tree-sitter"
	goGrammar "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

type golang struct {
	shouldFix    bool
	fileModifier analysis.FileModifier
	data         types.LanguageData
}

func NewGolangLanguage(pattern string, shouldFix bool) types.NodeManagement {
	g := &golang{
		shouldFix: shouldFix,
		data: types.LanguageData{
			Language: tree.NewLanguage(goGrammar.Language()),
			Queries:  buildGolangQuery(pattern),
		},
	}
	g.fileModifier = analysis.NewFileModifier(g)
	return g
}

func (g golang) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		g.fileModifier.ModifyVariableName(code, filepath, string(code[node.Node.StartByte():node.Node.EndByte()]), nodeInfo, g.shouldFix)
		return
	case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "expression_list":
		return
	case alternative != nil && alternative.GrammarName() == "block":
		nodeInfo.Complexity++
	}
	nodeInfo.Complexity++
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
