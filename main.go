package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.Test([]byte("function suma(a, b) { return a + b; } console.log(suma(1, 2))"))
}
