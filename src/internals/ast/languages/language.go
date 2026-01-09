package languages

import tree "github.com/tree-sitter/go-tree-sitter"

/*
 * language.go: Definition of all the types used in the parsing
 */

// LanguageInformation - > struct made to register the language an all it's complements (used by the parser)
type LanguageInformation struct {
	Language        *tree.Language
	Queries         []string
	RegexComplexity *RegexComplexity
}

// FunctionData - > struct made to register all the data returned by the parser and save it
type FunctionData struct {
	Name, Parameters        string
	TotalParams, Complexity int
}

// RegexComplexity - > struct made to give the parser enough data to parse our source code
type RegexComplexity struct {
	Code []byte
	// Function made to manage each node of the nodes found of every function. Called in a for loop.
	ManageNode func(captureNames []string, code []byte, node tree.QueryCapture, complexity *int)
}

// AddInitialData Method to append the initial data into a FunctionData "object"
func (f *FunctionData) AddInitialData(name, parameters string, totalParams int) {
	f.Name = name
	f.Parameters = parameters
	f.TotalParams = totalParams
}
