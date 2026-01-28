package languages

import (
	"CLI_App/cmd/adapters/analysis/types"
	"CLI_App/cmd/domain"
	"fmt"

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

func (g golang) ManageNode(captureNames []string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.UpdateInvalidNames()
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

func (g golang) GetVarAppearancesQuery(pattern string) string {
	return fmt.Sprintf("([(identifier) (field_identifier)] @variable.name") +
		fmt.Sprintf("(#not-match? @variable.name \"^%s|%s$\"))", pattern, domain.AllowNonNamedVar) +
		fmt.Sprintf("(type_declaration (type_spec name: (type_identifier) @variable.name) (#not-match? @variable.name \"^%s|%s$\"))", pattern, domain.AllowNonNamedVar) +
		fmt.Sprintf("(method_declaration receiver: (parameter_list (parameter_declaration type: ([(type_identifier) @variable.name (pointer_type (type_identifier) @variable.name)])))") +
		fmt.Sprintf("(#not-match? @variable.name \"^%s|%s$\"))", pattern, domain.AllowNonNamedVar)
}
