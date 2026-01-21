package languages

import (
	"CLI_App/src/config"
	language "CLI_App/src/internals/ast/utils"
	"CLI_App/src/internals/ast/writer"
	"CLI_App/src/internals/utils"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	pyGrammar "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

// python: struct with LanguageInformation embedded
type python struct {
	data language.LanguageData
}

// ManageNode - > Function to implement the NodeManagement interface
func (p python) ManageNode(captureNames []string, code []byte, node tree.QueryCapture, nodeInfo *language.FunctionData) {
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
			string(code[node.Node.StartByte():node.Node.EndByte()]),
			node.Node.StartPosition().Row, node.Node.StartPosition().Column,
		)
		writer.ModifyVariableName(p, node.Node, code, "main.py")
		return
	case node.Node.GrammarName() == "boolean_operator" && node.Node.Parent().GrammarName() == "assignment":
		return
	}
	nodeInfo.Complexity++
}

func (p python) GetLanguage() *tree.Language {
	return p.data.Language
}

func (p python) GetQueries() string {
	return p.data.Queries
}

func (p python) GetVarAppearancesQuery(name string) string {
	return fmt.Sprintf("((identifier) @variable.name (#match? @variable.name \"^%s$\"))", name)
}

// PythonData initialize the needed data for the Python language usage
var PythonData = python{
	language.LanguageData{
		Language: tree.NewLanguage(pyGrammar.Language()),
		Queries:
		// Function definition
		"(function_definition name: (identifier) @function.name " +
			"parameters: (parameters) @function.parameters " +
			"body: (block) @function.body) @function " +
			// Variable names
			"(assignment left: ([(identifier) @variable.name (pattern_list (identifier) @variable.name)])" +
			"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.GetActiveNamingConvention() + "\"))" +
			"(for_statement left: ([(identifier) @variable.name (pattern_list (identifier) @variable.name)])" +
			"(#not-match? @variable.name \"" + utils.AllowNonNamedVar + "|" + config.GetActiveNamingConvention() + "\"))" +
			// Keywords
			"[" +
			// If, else-if, else
			"(if_statement condition: (_)) (elif_clause condition: (_)) (else_clause body: (_))" +
			// Loops
			"(for_statement) (while_statement condition: (_) body: (_))" +
			// Operators
			"(boolean_operator left: (_) right: (_))" +
			// Clauses
			"(except_clause value: (_)) (conditional_expression) (case_clause (_))" +
			// List comprehension
			"(list_comprehension body: (_) (for_in_clause left: (_) right: (_))) (if_clause (_))" +
			"] @keyword",
	},
}
