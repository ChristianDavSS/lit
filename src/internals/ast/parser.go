package ast

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
)

type FunctionData struct {
	name, parameters        string
	totalParams, complexity int
}

type RegexComplexity struct {
	keyword, bodyStatements, andKw, orKw string
	code                                 []byte
}

var FunctionsData map[string][]*FunctionData

func GetAST(code []byte, language *tree.Language) (*tree.Tree, error) {
	// Create a parser for the code
	parser := tree.NewParser()
	defer parser.Close()
	// Set the grammar of the language to parse
	err := parser.SetLanguage(language)
	if err != nil {
		return nil, err
	}
	// Parse the source code with the configured parser
	treeParser := parser.Parse(code, nil)

	// Get the root (program node) node from the parser.
	return treeParser, nil
}

func GetCapturesByQueries(language *tree.Language, queries string, code []byte, root *tree.Node) (*tree.Query,
	*tree.QueryCursor, tree.QueryMatches) {
	// Create a query to extract the data we need
	query, _ := tree.NewQuery(language, queries)

	// Create a query cursor for our custom query (just to keep the necessary data)
	cursor := tree.NewQueryCursor()

	// Execute the query and generate the captures
	return query, cursor, cursor.Matches(query, root, code)
}

// CyclicalComplexity Function that calculates the cyclical complexity of the code. Useful for the user feedback
func CyclicalComplexity(language *tree.Language, queries string, root *tree.Node, config *RegexComplexity) {
	query, cursor, captures := GetCapturesByQueries(language, queries, config.code, root)
	defer query.Close()
	defer cursor.Close()
	var nodes []*tree.QueryMatch

	for {
		match := captures.Next()
		if match == nil {
			break
		}
		nodes = append(nodes, match)
	}

	for i := range len(nodes) {
		complexity := 1
		functionBody(&nodes[i].Captures[2].Node, &complexity, config)
		fmt.Printf("Complexity of %s - > %d\n", config.code[nodes[i].Captures[0].Node.StartByte():nodes[i].Captures[0].Node.EndByte()], complexity)
	}
}

func functionBody(node *tree.Node, complexity *int, config *RegexComplexity) {
	for i := range node.NamedChildCount() {
		child := node.NamedChild(i)
		if regexMatching(config.bodyStatements, child.GrammarName()) {
			functionBody(child, complexity, config)
			continue
		}

		if regexMatching(config.keyword, child.GrammarName()) {
			if child.GrammarName() == "else_clause" {
				if node.Child(node.ChildCount()-1).GrammarName() == "if_statement" {
					*complexity++
				}
			} else {
				*complexity++
			}
			functionBody(child, complexity, config)
			continue
		}

		line := string(config.code[child.StartByte():child.EndByte()])
		*complexity += strings.Count(line, config.andKw) + strings.Count(line, config.orKw)
	}
}

func regexMatching(pattern, s string) bool {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		os.Exit(1)
	}
	return matched
}
