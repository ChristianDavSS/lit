package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte(" const FUNCION1 = () => { const a = true ? 1 : 0; } "), "js")
}
