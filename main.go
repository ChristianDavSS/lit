package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.Test([]byte("function testCase(x,y){  const a = true && false; if(x>0&&y>0) { if (x < y || y < 1000) { switch (x) { case true && false || false && true: if (true && false) { return false || false; } return false || true; } } return 1; } else if(x>0||y>0) { return 2; } const v = 1; return 0; }"))
}
