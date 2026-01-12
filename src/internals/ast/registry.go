package ast

import (
	"CLI_App/src/internals/ast/languages"
	"fmt"
)

/*
 * registry.go - > File for configuration of the language queries. This contains each programming language
 * config for their parsing and analysis. Also, contains methods to manage the Functions...
 */

var languageQueries = map[string]*languages.LanguageInformation{
	"js":   &languages.JsLanguage,
	"py":   &languages.PyLanguage,
	"java": &languages.JavaLanguage,
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

	languageInfo.RegexComplexity.Code = code
	Functions := CyclicalComplexity(languageInfo.Language, languageInfo.Queries, root, languageInfo.RegexComplexity)

	for _, v := range Functions {
		fmt.Printf("Datos de la funcion %s%s\n", v.Name, v.Parameters)
		fmt.Printf("Total de parametros: %d\n", v.TotalParams)
		fmt.Printf("Complejidad ciclomatica: %d\n", v.Complexity)
		fmt.Println()
	}
}
