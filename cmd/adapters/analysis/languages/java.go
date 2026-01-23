package languages

import (
	"CLI_App/cmd/adapters/analysis"
	"CLI_App/cmd/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	javaGrammar "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

type java struct {
	shouldFix bool
	data      domain.LanguageData
}

func NewJavaLanguage(pattern string, shouldFix bool) domain.NodeManagement {
	return &java{
		shouldFix: shouldFix,
		data: domain.LanguageData{
			Language: tree.NewLanguage(javaGrammar.Language()),
			Queries:  buildJavaQuery(pattern),
		},
	}
}

func (j java) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
			string(code[node.Node.StartByte():node.Node.EndByte()]),
			node.Node.StartPosition().Row, node.Node.StartPosition().Column,
		)
		if j.shouldFix {
			writer := analysis.NewFileModifier(j, string(code[node.Node.StartByte():node.Node.EndByte()]))
			writer.ModifyVariableName(node.Node, code, filepath)
		}
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

func (j java) GetLanguage() *tree.Language {
	return j.data.Language
}

func (j java) GetQueries() string {
	return j.data.Queries
}

func (j java) GetVarAppearancesQuery(name string) string {
	return name
}
