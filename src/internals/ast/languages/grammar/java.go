package grammar

import (
	"CLI_App/src/config"
	"CLI_App/src/internals/ast/languages"
	"CLI_App/src/internals/utils"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	javaGrammar "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

type java struct {
	data languages.LanguageData
}

func (j java) ManageNode(captureNames []string, code []byte, node tree.QueryCapture, nodeInfo *languages.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
			string(code[node.Node.StartByte():node.Node.EndByte()]),
			node.Node.StartPosition().Row, node.Node.StartPosition().Column,
		)
		return
	case node.Node.GrammarName() == "binary_expression" && node.Node.Parent().GrammarName() == "variable_declarator":
		return
	case alternative != nil && alternative.GrammarName() == "block":
		nodeInfo.Complexity++
	}
	nodeInfo.Complexity++
}

func (j java) GetLanguage() *tree.Language {
	return j.data.Language
}

func (j java) GetQueries() string {
	return j.data.Queries
}

var JavaLanguage = java{
	languages.LanguageData{
		Language: tree.NewLanguage(javaGrammar.Language()),
		Queries:
		// Method definition
		"(method_declaration type: (_) name: (_) @function.name " +
			"parameters: (formal_parameters) @function.parameters " +
			"body: (block) @function.body ) @function " +
			// Variable names
			"(variable_declarator name: (identifier) @variable.name " +
			"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.GetActiveNamingConvention() + "\"))" +
			"(enhanced_for_statement name: (identifier) @variable.name " +
			"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.GetActiveNamingConvention() + "\"))" +
			// Keywords (+1 complexity)
			"[" +
			// Loops
			"(for_statement) (while_statement) (do_statement) (enhanced_for_statement)" +
			// If, else-if, else
			"(if_statement condition: (_) consequence: (_) alternative: (_)?) (ternary_expression)" +
			// Expressions
			"((binary_expression left: (_) right: (_)) @bin_exp (#match? @bin_exp \".*(&&|[|]{2}).*\"))" +
			"(switch_block_statement_group) (catch_clause)" +
			"] @keyword",
	},
}
