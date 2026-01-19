package utils

import (
	"fmt"
	"os"

	tree "github.com/tree-sitter/go-tree-sitter"
)

/*
 * files.go: file to manage the reading/writing on scripts and directories.
 */

// GetDirEntries - > Gets the entries of the directory given (the path of the dir)
func GetDirEntries(name string) []os.DirEntry {
	files, err := os.ReadDir(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading the file structure: ", err)
		os.Exit(1)
	}
	return files
}

// GetWorkingDirectory Method to return the current working directory
func GetWorkingDirectory() string {
	// Get the working directory (where we are executing the commands)
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error with the path: ", err)
		os.Exit(1)
	}

	// If everything went right, we return the path of the wd
	return path
}

// ModifyVariableName - > this function modifies the variable written the wrong way in the code, rewriting it for you.
func ModifyVariableName(node tree.Node, path string) {
	// Read the file with the given path
	file, _ := os.ReadFile(path)

	// While it finds a valid index for the substring, we'll keep iterating
	for {
		fmt.Println(string(file[node.StartByte():node.EndByte()]))
		indexes := findSeparatorIndexes(string(file[node.StartByte():node.EndByte()]))
		modifyVariableName(indexes, file[node.StartByte():node.EndByte()])
		break
	}
	// Write the new code into the file (just with the modified lines)
	/*err := os.WriteFile("main.py", []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to write the variable name into the file...")
		os.Exit(1)
	}*/
}

func findSeparatorIndexes(line string) []int16 {
	var indexes []int16

	// Traverse the string
	for i, ch := range line {
		// Find if the value is in the upper letters range or if it's an underscore
		if ch >= 65 && ch <= 90 || ch == 95 {
			indexes = append(indexes, int16(i))
		}
	}

	return indexes
}

func modifyVariableName(indexes []int16, line []byte) []byte {
	currentNamingConv := "camelCase"
	fmt.Println(currentNamingConv)
	return []byte("")
}
