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
	name, parameters        string
	totalParams, complexity int
}

// RegexComplexity - > struct made to give the parser enough data to parse our source code
type RegexComplexity struct {
	Keyword, BodyStatements string
	Code                    []byte
	KeywordMatchFunc        func(node *tree.Node, complexity *int)
	NoMatchFunc             func(node *tree.Node, complexity *int, code []byte)
}
