package ast

import (
	"CLI_App/src/internals/ast/languages"
	"fmt"
	"strings"
)

/*
 * registry.go - > File for configuration of the language queries. This contains each programming language
 * config for their parsing and analysis.
 */

var languageQueries = map[string]*languages.LanguageInformation{
	"js": &languages.JsLanguage,
	"py": nil,
}

func RunParser(code []byte, language string) {
	languageInfo, ok := languageQueries[language]
	if !ok {
		return
	}
	// Get our ast bases in our code and grammar
	ast, err := GetAST(code, languageInfo.Language)
	if err != nil {
		return
	}
	// Get the root node (program) from the AST generated
	root := ast.RootNode()
	defer ast.Close()

	// Get the query, cursor and captures from the helper function
	query, cursor, _ := GetCapturesByQueries(languageInfo.Language, strings.Join(languageInfo.Queries, " "), code, root)
	// Defer the closing (because the iterative process uses this two variables)
	defer query.Close()
	defer cursor.Close()
	fmt.Printf("%s\n\n", root.ToSexp())

	languageInfo.RegexComplexity.Code = code
	CyclicalComplexity(languageInfo.Language, strings.Join(languageInfo.Queries, " "), root, languageInfo.RegexComplexity)
}
