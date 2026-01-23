package languages

import (
	"CLI_App/cmd/adapters/analysis"
	"CLI_App/cmd/domain"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	pyGrammar "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

// python: struct with LanguageInformation embedded
type python struct {
	shouldFix bool
	data      domain.LanguageData
}

func NewPythonLanguage(pattern string, shouldFix bool) domain.NodeManagement {
	return &python{
		shouldFix: shouldFix,
		data: domain.LanguageData{
			Language: tree.NewLanguage(pyGrammar.Language()),
			Queries:  buildPythonQuery(pattern),
		},
	}
}

// ManageNode - > Function to implement the NodeManagement interface
func (p python) ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData) {
	switch {
	case captureNames[node.Index] == "variable.name":
		nodeInfo.Feedback += fmt.Sprintf("   Error: The variable '%s' is not using the valid naming convention. (%d:%d).\n",
			string(code[node.Node.StartByte():node.Node.EndByte()]),
			node.Node.StartPosition().Row, node.Node.StartPosition().Column,
		)
		if p.shouldFix {
			writer := analysis.NewFileModifier(p, string(code[node.Node.StartByte():node.Node.EndByte()]))
			writer.ModifyVariableName(node.Node, code, filepath)
		}
		return
	case node.Node.GrammarName() == "boolean_operator" && node.Node.Parent().GrammarName() == "assignment":
		return
	}
	nodeInfo.Complexity++
}

func buildPythonQuery(pattern string) string {
	fmt.Println("PATTERN:", pattern)
	return "(function_definition name: (identifier) @function.name " +
		"parameters: (parameters) @function.parameters " +
		"body: (block) @function.body) @function " +
		// Variable names
		"(assignment left: ([(identifier) @variable.name (pattern_list (identifier) @variable.name)])" +
		"(#not-match? @variable.name \"" + domain.AllowNonNamedVar + "|" + pattern + "\"))" +
		"(for_statement left: ([(identifier) @variable.name (pattern_list (identifier) @variable.name)])" +
		"(#not-match? @variable.name \"" + domain.AllowNonNamedVar + "|" + pattern + "\"))" +
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
		"] @keyword"
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
