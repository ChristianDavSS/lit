package languages

import (
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
	javaGrammar "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

var JavaLanguage = LanguageInformation{
	Language: tree.NewLanguage(javaGrammar.Language()),
	Queries:
	// Method definition
	"(method_declaration type: (_) name: (_) @function.name " +
		"parameters: (formal_parameters) @function.parameters " +
		"body: (block) @function.body ) @function " +
		"[" +
		// Loops
		"(for_statement) (while_statement) (do_statement)" +
		// If, else-if, else
		"(if_statement condition: (_) consequence: (_) alternative: (_)?) (ternary_expression)" +
		// Expressions
		"(binary_expression left: ([(true) (false)]) right: ([(true) (false)])) (switch_block_statement_group) (catch_clause)" +
		"] @keyword",
	RegexComplexity: &RegexComplexity{
		ManageNode: func(captureNames []string, code []byte, node tree.QueryCapture, complexity *int) {
			if node.Node.GrammarName() == "binary_expression" {
				line := string(code[node.Node.StartByte():node.Node.EndByte()])
				totalBooleans := strings.Count(line, "&&") + strings.Count(line, "||")
				if totalBooleans < 1 {
					return
				}
			}
			alternative := node.Node.ChildByFieldName("alternative")
			if alternative != nil && alternative.GrammarName() == "block" {
				*complexity++
			}
			*complexity++
		},
		MainFunc: &FunctionData{Name: "GlobalData", TotalParams: 0, Complexity: 1},
	},
}
