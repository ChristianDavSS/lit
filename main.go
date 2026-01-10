package main

import (
	"CLI_App/src/internals/ast"
)

func main() {
	//root.Execute()
	ast.RunParser([]byte("const FUNCION1 = (x, y) => {\n    const FUNCION2 = () => {\n        if (false) {\n            switch (true) {\n                case false: ;\n            };\n        };\n\n        const FUNCTION3 = () => {\n            switch (false) {\n                case 1: return 1;\n            }\n        }\n        if (false) {}\n\n        const FUNCION4 = () => {\n            if (false) {}\n        }\n    };\n    if (true) {};\n}; const prueba = () => {\n    if (true) {}\n}"), "js")
}
