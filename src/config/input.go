package config

import (
	patterns "CLI_App/src/internals/analysis/utils"
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

// getNamingConvention - > Asks the user to select their favorite naming convention for the variables.
func getNamingConvention() string {
	// Variable to save the value the user selects.
	var selectedConvention int8
	// Variable to save up the key-value pairs with the index-regex for the naming conventions.
	conventions := map[int8]string{
		1: patterns.LowerCamelCase,
		2: patterns.UpperCamelCase,
		3: patterns.CamelCase,
		4: patterns.SnakeCase,
	}

	// Ask the user for a valid input n times
	for selectedConvention > 4 || selectedConvention < 1 {
		clearScreen(osClear[runtime.GOOS][0], osClear[runtime.GOOS][1:]...)
		fmt.Println("Select a valid naming convention.")
		color.Red("[1] camelCase")
		color.Magenta("[2] CamelCase")
		color.Blue("[3] CamelCase/camelCase (for languages like Go)")
		color.Cyan("[4] snake_case")
		fmt.Scanf("%d", &selectedConvention)
	}

	// Return the pattern to use
	return conventions[selectedConvention]
}
