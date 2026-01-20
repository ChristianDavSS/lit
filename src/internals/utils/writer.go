package utils

import (
	"fmt"
	"os"
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
)

// ModifyVariableName - > this function modifies the variable written the wrong way in the code, rewriting it for you.
// Takes the node captured (the one we want to modify) and the path.
// Only converts from one convention to another (safety conditions)
func ModifyVariableName(node tree.Node, code []byte, filePath string) {
	// currentVarName is the current name of the variable on the code
	currentVarName := string(code[node.StartByte():node.EndByte()])

	// If the variable isn't  camelCase, CamelCase or snake_case, we don't modify it (for code safety)
	if !RegexMatch(CamelCase+"|"+SnakeCase, currentVarName) {
		return
	}

	// get the indexes where there's a separator
	upperIndexes := getSeparatorIndexes(currentVarName)
	// with the indexes, separate the line into valid tokens
	tokens := getTokens(upperIndexes, currentVarName)
	newVarName := refactorVarName(tokens)

	// modify just the line of code we need to.
	code = append(append(code[:node.StartByte()], newVarName...), code[node.EndByte():]...)

	// Write the new code into the file (just with the modified lines)
	err := os.WriteFile("main.py", code, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to write the variable name into the file...")
		os.Exit(1)
	}
}

// ---- Writing on files and renaming ----
// getSeparatorIndexes: get the indexes on the line of code where there's a separator (uppercase or underscore)
func getSeparatorIndexes(line string) []int16 {
	// slice to save up the indexes
	var indexes []int16

	// iterate through the line of code
	for i, ch := range line {
		// if the character is an uppercase letter or an underscore, we save the index of that
		if ch >= 65 && ch <= 90 && i > 0 || ch == 95 {
			indexes = append(indexes, int16(i))
		}
	}

	return indexes
}

// getTokens: gets the tokens (valid substrings) from a line of code with a determined naming convention
func getTokens(upperIndexes []int16, line string) []string {
	// slice to save up every cleaned token
	var tokens []string
	// variable to keep track of the previous index of the list
	var prevIdx int16

	// traverse our indexes
	for i, currIdx := range upperIndexes {
		// variable that is just 0 or 1 depending on the previous index value in the line of code.
		var sum int16
		// if the line of code in the position of the previous index is an underscore, we sum it up to one
		if line[prevIdx] == 95 {
			sum++
		}
		// we add up the clean token into the slice of tokens
		tokens = append(tokens, strings.ToLower(line[prevIdx+sum:currIdx]))
		// set up the previous token
		prevIdx = currIdx

		// if it's the last iteration, we add the values without a limit
		if i >= len(upperIndexes)-1 {
			tokens = append(tokens, strings.ToLower(line[prevIdx+sum:]))
		}
	}

	// return the clean slice of tokens
	return tokens
}

// refactorVarName: with the strings split in tokens, returns a []byte of the new line of code.
func refactorVarName(tokens []string) []byte {
	var selected int8 = 3
	var newName string = tokens[0]

	switch selected {
	case 3:
		for _, token := range tokens[1:] {
			newName += string(token[0]-32) + token[1:]
		}
	case 4:
		for _, token := range tokens[1:] {
			newName += "_" + token
		}
	}
	return []byte(newName)
}
