package languages

import (
	"fmt"

	tree "github.com/tree-sitter/go-tree-sitter"
	pyGrammar "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

var PyLanguage = LanguageInformation{
	Language: tree.NewLanguage(pyGrammar.Language()),
	Queries: []string{
		"(function_definition name: (identifier) @function.name " +
			"parameters: (parameters) @function.parameters " +
			"body: (block) @function.body )",
	},
	RegexComplexity: &RegexComplexity{
		BodyStatements: "lol",
		Keyword:        "lol",
		KeywordMatchFunc: func(node *tree.Node, complexity *int) {
			return
		},
		NoMatchFunc: func(node *tree.Node, complexity *int, code []byte) {
			fmt.Printf("%s\n", code[node.StartByte():node.EndByte()])
		},
	},
}
