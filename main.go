package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte("const a = true || false;\n\n[[], []].forEach(e => { e.map(o => {}) })"), "js")
}
