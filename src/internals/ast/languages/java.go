package languages

import (
	"CLI_App/src/internals/ast/config"
	language "CLI_App/src/internals/ast/utils"
	"CLI_App/src/internals/ast/writer"
	"CLI_App/src/internals/utils"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	javaGrammar "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

type java struct {
	data language.LanguageData
}

func (j java) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *language.FunctionData) {
	// Search the 'alternative' node in the children
	alternative := node.Node.ChildByFieldName("alternative")
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
			string(code[node.Node.StartByte():node.Node.EndByte()]),
			node.Node.StartPosition().Row, node.Node.StartPosition().Column,
		)
		writer.ModifyVariableName(j, node.Node, code, filepath)
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

func (j java) GetVarAppearancesQuery(name string) string {
	return name
}

var JavaLanguage = java{
	language.LanguageData{
		Language: tree.NewLanguage(javaGrammar.Language()),
		Queries:
		// Method definition
		"(method_declaration type: (_) name: (_) @function.name " +
			"parameters: (formal_parameters) @function.parameters " +
			"body: (block) @function.body ) @function " +
			// Variable names
			"(variable_declarator name: (identifier) @variable.name " +
			"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.ActivePattern + "\"))" +
			"(enhanced_for_statement name: (identifier) @variable.name " +
			"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.ActivePattern + "\"))" +
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
