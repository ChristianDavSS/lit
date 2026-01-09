package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte(" const FUNCION1 = (x, y) => { [].forEach(e => { [].map(a => { const a = 1; }); }) }"), "js")
}
