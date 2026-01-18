package internals

import (
	"fmt"
	"os"
)

// GetWorkingDirectory Method to return the current working directory
func GetWorkingDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error with the path: ", err)
		os.Exit(1)
	}

	return path
}
