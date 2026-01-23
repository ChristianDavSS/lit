package languages

import (
	"CLI_App/cmd/adapters/analysis"
	"CLI_App/cmd/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	goGrammar "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

type golang struct {
	shouldFix bool
	data      domain.LanguageData
}

func NewGolangLanguage(pattern string, shouldFix bool) domain.NodeManagement {
	return &golang{
		shouldFix: shouldFix,
		data: domain.LanguageData{
			Language: tree.NewLanguage(goGrammar.Language()),
			Queries:  buildGolangQuery(pattern),
		},
	}
}

func (g golang) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
			string(code[node.Node.StartByte():node.Node.EndByte()]),
			node.Node.StartPosition().Row, node.Node.StartPosition().Column,
		)
		if g.shouldFix {
			writer := analysis.NewFileModifier(g, string(code[node.Node.StartByte():node.Node.EndByte()]))
			writer.ModifyVariableName(node.Node, code, filepath)
		}
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

func (g golang) GetLanguage() *tree.Language {
	return g.data.Language
}

func (g golang) GetQueries() string {
	return g.data.Queries
}

func (g golang) GetVarAppearancesQuery(name string) string {
	return name
}
