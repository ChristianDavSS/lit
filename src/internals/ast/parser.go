package ast

import (
	"CLI_App/src/internals/ast/languages"
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
)

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

// CyclicalComplexity Function that calculates the cyclical complexity of the code. Useful for the user feedback.
func CyclicalComplexity(language *tree.Language, queries string, root *tree.Node, config *languages.RegexComplexity) []*languages.FunctionData {
	// Get the basics to iterate through the captures and keep the data
	query, cursor, captures := GetCapturesByQueries(language, queries, config.Code, root)
	defer query.Close()
	defer cursor.Close()
	// List used as a stack to get the subfunctions and it's complexity right
	var Stack []*languages.FunctionData
	// Slice to save up the functions we take out the stack (with their final complexity)
	var Functions []*languages.FunctionData
	if config.MainFunc != nil {
		Functions = append(Functions, config.MainFunc)
	}

	// Get the Functions data
	for {
		// Iterator of all the captures of the source code.
		match := captures.Next()
		if match == nil {
			break
		}
		// Get a deep copy of the match (because the memory of it will be overwritten)
		copyOf := *match

		switch {
		// While we iterate through the captures, we save up the copies on a slice of objects.
		case query.CaptureNames()[copyOf.Captures[0].Index] == "function":
			Stack = append(Stack, &languages.FunctionData{Complexity: 1})
			// Add the initial data to the object reference in the stack
			Stack[len(Stack)-1].AddInitialData(
				string(config.Code[copyOf.Captures[1].Node.StartByte():copyOf.Captures[1].Node.EndByte()]),
				string(config.Code[copyOf.Captures[2].Node.StartByte():copyOf.Captures[2].Node.EndByte()]),
				int(copyOf.Captures[2].Node.NamedChildCount()),
				copyOf.Captures[3].Node.StartByte(), copyOf.Captures[3].Node.EndByte(),
			)

		// If there´s code without a function (JS, Python) before a function definition, we count it as main.
		case len(Stack) <= 0:
			config.ManageNode(query.CaptureNames(), config.Code, copyOf.Captures[0], &Functions[0].Complexity)

		// Validate node ranges with the function body. This is the logic for functions inside functions or the main one
		case copyOf.Captures[0].Node.StartByte() >= Stack[len(Stack)-1].StartByte && copyOf.Captures[0].Node.EndByte() <= Stack[len(Stack)-1].EndByte:
			// + 1 in complexity in the function.
			config.ManageNode(query.CaptureNames(), config.Code, copyOf.Captures[0], &Stack[len(Stack)-1].Complexity)

		// The code only gets here when there's a line out of the scope of a function
		// Remove the most recent element from the stack and add it to the Functions list
		default:
			// Flag to detect if the node is a main function node or not
			isMain := false
			// While the last element on the stack doesn't satisfy the range of the current node, we add it up
			// to the final Functions slice and remove it from the stack to keep going until it finds the right Function.
			for !(copyOf.Captures[0].Node.StartByte() >= Stack[len(Stack)-1].StartByte && copyOf.Captures[0].Node.EndByte() <= Stack[len(Stack)-1].EndByte) {
				if len(Stack) <= 1 && config.MainFunc != nil {
					isMain = true
					config.ManageNode(query.CaptureNames(), config.Code, copyOf.Captures[0], &Functions[0].Complexity)
					break
				}
				Functions = append(Functions, Stack[len(Stack)-1])
				Stack = Stack[:len(Stack)-1]
			}
			// At the end, in the verified node, we sum +1 to the complexity.
			if !isMain {
				config.ManageNode(query.CaptureNames(), config.Code, copyOf.Captures[0], &Stack[len(Stack)-1].Complexity)
			}
		}
	}
	// If there's any function still on the stack, we copy it into the Functions slice.
	return append(Functions, Stack...)
}
