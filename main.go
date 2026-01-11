package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte("public class Main {\n    boolean b = true && false;\n    public static void main(String[] args) {\n        int a = true ? 1 : false ? 1 : 0;\n    }\n    public static class Pene {\n        boolean c = true && false;\n        \n        public static void semen() {\n            if (true) {}\n        }\n    }\n}\n"), "java")
}
