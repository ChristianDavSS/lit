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
		"(binary_expression left: (_) right: (_)) (switch_block_statement_group) (catch_clause)" +
		"] @keyword",
	RegexComplexity: &RegexComplexity{
		ManageNode: func(captureNames []string, code []byte, node tree.QueryCapture, complexity *int) {
			// Search the 'alternative' node in the children
			alternative := node.Node.ChildByFieldName("alternative")
			switch {
			case node.Node.GrammarName() == "binary_expression":
				// Get the line of code of the current node
				line := string(code[node.Node.StartByte():node.Node.EndByte()])
				totalBooleans := strings.Count(line, "&&") + strings.Count(line, "||")
				if totalBooleans < 1 || node.Node.Parent().GrammarName() == "variable_declarator" {
					return
				}
			case alternative != nil && alternative.GrammarName() == "block":
				*complexity++
			}
			*complexity++
		},
		MainFunc: &FunctionData{Name: "GlobalData", TotalParams: 0, Complexity: 1},
	},
}
