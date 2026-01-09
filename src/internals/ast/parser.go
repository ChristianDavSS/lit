package ast

import (
	"CLI_App/src/internals/ast/languages"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
)

var Functions []*languages.FunctionData

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
	query, err := tree.NewQuery(language, queries)
	if err != nil {
		fmt.Println("Error parsing - >", err)
	}

	// Create a query cursor for our custom query (just to keep the necessary data)
	cursor := tree.NewQueryCursor()

	// Execute the query and generate the captures
	return query, cursor, cursor.Matches(query, root, code)
}

// CyclicalComplexity Function that calculates the cyclical complexity of the code. Useful for the user feedback
func CyclicalComplexity(language *tree.Language, queries string, root *tree.Node, config *languages.RegexComplexity) {
	query, cursor, captures := GetCapturesByQueries(language, queries, config.Code, root)
	defer query.Close()
	defer cursor.Close()
	i := -1

	// Get the Functions data
	for {
		match := captures.Next()
		if match == nil {
			break
		}
		// Get a deep copy of the match (because the memory of it will be overwritten)
		copyOf := *match

		// While we iterate through the captures, we save up the copies on a slice of objects.
		if query.CaptureNames()[copyOf.Captures[0].Index] == "function" {
			i += 1
			Functions = append(Functions, &languages.FunctionData{Complexity: 1})
			// Add data to the
			Functions[i].AddInitialData(
				string(config.Code[copyOf.Captures[1].Node.StartByte():copyOf.Captures[1].Node.EndByte()]),
				string(config.Code[copyOf.Captures[2].Node.StartByte():copyOf.Captures[2].Node.EndByte()]),
				int(copyOf.Captures[2].Node.NamedChildCount()),
			)
		} else {
			config.ManageNode(query.CaptureNames(), config.Code, copyOf.Captures[0], &Functions[i].Complexity)
		}
	}
}
