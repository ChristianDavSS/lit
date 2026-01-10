package languages

import (
	tree "github.com/tree-sitter/go-tree-sitter"
	pyGrammar "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

var PyLanguage = LanguageInformation{
	Language: tree.NewLanguage(pyGrammar.Language()),
	Queries: "(function_definition name: (identifier) @function.name " +
		"parameters: (parameters) @function.parameters " +
		"body: (block) @function.body )",
	RegexComplexity: &RegexComplexity{},
}
