package ui

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
}

// GetNamingConvention - > Asks the user to select their favorite naming convention for the variables.
func GetNamingConvention() int8 {
	// Variable to save the value the user selects.
	var selectedConvention int8

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
	return selectedConvention
}

// Function used to clear the screen in Windows or Linux...
func clearScreen(name string, args ...string) {
	// Clear the screen
	cmd := exec.Command(name, args...)
	// Set the stdout we want to clear
	cmd.Stdout = os.Stdout
	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running the command to clear your screen.")
		os.Exit(1)
	}
}
