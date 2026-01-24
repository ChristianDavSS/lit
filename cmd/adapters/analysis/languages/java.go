package languages

import (
	"CLI_App/cmd/adapters/analysis"
	"CLI_App/cmd/adapters/analysis/types"
	"CLI_App/cmd/domain"

	tree "github.com/tree-sitter/go-tree-sitter"
	javaGrammar "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

type java struct {
	shouldFix    bool
	fileModifier analysis.FileModifier
	data         types.LanguageData
}

func NewJavaLanguage(pattern string, shouldFix bool) types.NodeManagement {
	j := &java{
		shouldFix: shouldFix,
		data: types.LanguageData{
			Language: tree.NewLanguage(javaGrammar.Language()),
			Queries:  buildJavaQuery(pattern),
		},
	}
	j.fileModifier = analysis.NewFileModifier(j)
	return j
}

func (j java) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		// Set the initial feedback
		nodeInfo.SetVariableFeedback(string(code[node.Node.StartByte():node.Node.EndByte()]), domain.Point(node.Node.StartPosition()))
		j.fileModifier.ModifyVariableName(filepath, string(code[node.Node.StartByte():node.Node.EndByte()]), j.shouldFix)
		return
	case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "variable_declarator":
		return
	case alternative != nil && alternative.GrammarName() == "block":
		nodeInfo.Complexity++
	}
	nodeInfo.Complexity++
}

func buildJavaQuery(pattern string) string {
	return "(method_declaration type: (_) name: (_) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (block) @function.body ) @function " +
		// Variable names
		"(variable_declarator name: (identifier) @variable.name " +
		"(#not-match? @variable.name \"" + domain.AllowNonNamedVar + "|" + pattern + "\"))" +
		"(enhanced_for_statement name: (identifier) @variable.name " +
		"(#not-match? @variable.name \"" + domain.AllowNonNamedVar + "|" + pattern + "\"))" +
		// Keywords (+1 complexity)
		"[" +
		// Loops
		"(for_statement) (while_statement) (do_statement) (enhanced_for_statement)" +
		// If, else-if, else
		"(if_statement condition: (_) consequence: (_) alternative: (_)?) (ternary_expression)" +
		// Expressions
		"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
		"(switch_block_statement_group) (catch_clause)" +
		"] @keyword"
}

func (j java) GetLanguageData() types.LanguageData {
	return j.data
}

func (j java) GetVarAppearancesQuery(name string) string {
	return name
}
