package languages

import tree "github.com/tree-sitter/go-tree-sitter"

/*
 * language.go: Definition of all the types used in the parsing
 */

// LanguageInformation - > struct made to register the language an all it's complements (used by the parser)
type LanguageInformation struct {
	Language        *tree.Language
	Queries         string
	RegexComplexity *RegexComplexity
}

// FunctionData - > struct made to register all the data returned by the parser and save it
type FunctionData struct {
	Name, Parameters        string
	TotalParams, Complexity int
	StartByte, EndByte      uint
}

// RegexComplexity - > struct made to give the parser enough data to parse our source code
type RegexComplexity struct {
	Code []byte
	// Function made to manage each node of the nodes found of every function. Called in a for loop.
	ManageNode func(captureNames []string, code []byte, node tree.QueryCapture, complexity *int)
	MainFunc   *FunctionData
}

// AddInitialData Method to append the initial data into a FunctionData "object"
func (f *FunctionData) AddInitialData(name, parameters string, totalParams int, startByte, endByte uint) {
	f.Name = name
	f.Parameters = parameters
	f.TotalParams = totalParams
	f.StartByte = startByte
	f.EndByte = endByte
}

// IsTargetInRange validates the range given by another function to validate it's written on the same byte range
func (f *FunctionData) IsTargetInRange(startByte, endByte uint) bool {
	return f.StartByte <= startByte && f.EndByte >= endByte
}
