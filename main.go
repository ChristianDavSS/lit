package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte("function testCase(x,y){ const a = true ? 0 : 1; }"), "js")
}
