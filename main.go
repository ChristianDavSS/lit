package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte("const F=()=>{if(a){const F=(a)=>{if(b){}}}}\n"), "js")
}
