package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte("public class Main {\n    public static void main(String[] args) {\n        int a = true ? 1 : false ? 1 : 0;\n    }\n}\n"), "java")
}
