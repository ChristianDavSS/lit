package ast

import (
	tree "github.com/tree-sitter/go-tree-sitter"
)

type Function struct {
	name   string
	status string
}

func (f *Function) SetStatus(status string) {
	f.status = status
}

var Warnings []*Function

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
