package ast

import (
	"CLI_App/src/internals/ast/languages"
	"CLI_App/src/internals/ast/languages/grammar"
	"fmt"
	"os"
)

/*
 * registry.go - > File for configuration of the language queries. This contains each programming language
 * config for their parsing and analysis. Also, contains methods to manage the Functions...
 */

var languageQueries = map[string]*languages.LanguageInformation{
	"js":   &grammar.JsLanguage,
	"py":   &grammar.PyLanguage,
	"java": &grammar.JavaLanguage,
}

func RunParser(code []byte, language string) []*languages.FunctionData {
	languageInfo, ok := languageQueries[language]
	if !ok {
		fmt.Println("Error with the language: not supported yet.")
		os.Exit(1)
	}
	// Get our ast bases in our code and grammar
	ast, err := GetAST(code, languageInfo.Language)
	if err != nil {
		fmt.Println("Error getting the AST for the language. Do you have a 64x C compiler installed?")
		os.Exit(1)
	}
	// Get the root node (program) from the AST generated
	root := ast.RootNode()
	defer ast.Close()

	// Set up the code on the configuration struct
	languageInfo.RegexComplexity.Code = code
	// Find the cyclical complexity of the functions
	functions := CyclicalComplexity(languageInfo.Language, languageInfo.Queries, root, languageInfo.RegexComplexity)

	// Remove the 'normal' functions from the list, keeping the dangerous ones.
	i := 0
	for i < len(functions) {
		function := functions[i]
		// If the functions isn't dangerous, we delete it
		if function.TotalParams <= languages.Messages["parameters"][0].Value &&
			function.Complexity <= languages.Messages["complexity"][0].Value &&
			int(function.Size) <= languages.Messages["size"][0].Value {
			// If the function is dangerous, we increase the counter.
			functions = append(functions[:i], functions[i+1:]...)
			continue
		}
		function.SetFunctionFeedback()
		i++
	}

	// Return the dangerous functions on the script
	return functions
}
