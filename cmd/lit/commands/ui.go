package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

// osClear - > Takes the name of the os and returns the command lines to clear the stdout.
var osClear = map[string][]string{
	"windows": {"cmd", "/c", "cls"},
	"linux":   {"clear"},
	"darwin":  {"clear"},
}

// GetNamingConvention - > Asks the user to select their favorite naming convention for the variables.
func GetNamingConvention() int8 {
	// Variable to save the value the user selects.
	var selectedConvention int8

	// Ask the user for a valid input n times
	for selectedConvention > 4 || selectedConvention < 1 {
		ClearScreen()
		fmt.Println("Select a valid naming convention.")
		color.Red("[1] camelCase")
		color.Magenta("[2] CamelCase")
		color.Blue("[3] CamelCase/camelCase (for languages like Go)")
		color.Cyan("[4] snake_case")
		fmt.Scanf("%d", &selectedConvention)
	}

	// Return the pattern to use
	return selectedConvention
}

// ClearScreen clears the terminal screen.
func ClearScreen() {
	name := "linux"
	if runtime.GOOS == "windows" {
		name = "windows"
	} else if runtime.GOOS == "darwin" {
		name = "darwin"
	}

	cmdArgs, ok := osClear[name]
	if !ok {
		// Fallback or ignore
		return
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running the command to clear your screen.")
	}
}
