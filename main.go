package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte("a = True and False\nif True or False:\n    x = 1 if True or False else 12"), "py")
}
