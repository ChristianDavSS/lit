package ast

import (
	"CLI_App/src/internals/ast/languages"
	"CLI_App/src/internals/ast/tree"
	"CLI_App/src/internals/ast/utils"
	"fmt"
	"os"
	"path/filepath"
)

/*
 * registry.go - > File for configuration of the language queries. This contains each programming language
 * config for their parsing and analysis. Also, contains methods to manage the Functions...
 */

var languageQueries = map[string]utils.NodeManagement{
	"js":   languages.JsLanguage,
	"jsx":  languages.JsLanguage,
	"py":   languages.PythonData,
	"java": languages.JavaLanguage,
	"go":   languages.GoLanguage,
}

func RunParser(code []byte, filename string) []*utils.FunctionData {
	languageInfo, ok := languageQueries[filepath.Ext(filename)[1:]]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error with the language: not supported yet.")
		os.Exit(1)
	}

	// Find the cyclical complexity of the functions
	functions := tree.CyclicalComplexity(languageInfo, code, filename)

	// Remove the 'normal' functions from the list, keeping the dangerous ones.
	i := 0
	for i < len(functions) {
		function := functions[i]
		// If the functions isn't dangerous, we delete it
		if function.TotalParams <= utils.Messages["parameters"][0].Value &&
			function.Complexity <= utils.Messages["complexity"][0].Value &&
			int(function.Size) <= utils.Messages["size"][0].Value &&
			function.Feedback == "" {
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
