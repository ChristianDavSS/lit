package languages

import (
	"CLI_App/internal/adapter/analysis"
	"CLI_App/internal/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	goGrammar "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

type golang struct {
	data            analysis.LanguageData
	shouldFix       bool
	conventionIndex int8
}

// NewGoLanguage creates a new Go language adapter.
func NewGoLanguage(pattern string, shouldFix bool, conventionIndex int8) analysis.NodeManagement {
	return golang{
		shouldFix:       shouldFix,
		conventionIndex: conventionIndex,
		data: analysis.LanguageData{
			Language: tree.NewLanguage(goGrammar.Language()),
			Queries:  buildGoQueries(pattern),
		},
	}
}

func buildGoQueries(pattern string) string {
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

func (g golang) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
			string(code[node.Node.StartByte():node.Node.EndByte()]),
			node.Node.StartPosition().Row, node.Node.StartPosition().Column,
		)
		analysis.ModifyVariableName(g, node.Node, code, filepath, g.shouldFix, g.conventionIndex)
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
