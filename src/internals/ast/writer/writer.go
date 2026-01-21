package writer

import (
	parser "CLI_App/src/internals/ast/tree"
	languages "CLI_App/src/internals/ast/utils"
	"CLI_App/src/internals/utils"
	"fmt"
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
)

/*
 * writer.go - > file to write some data on a file using some ast queries.
 */

// ModifyVariableName - > this function modifies the variable written the wrong way in the code, rewriting it for you.
// Takes the node captured (the one we want to modify) and the path.
// Only converts from one convention to another (safety conditions)
func ModifyVariableName(management languages.NodeManagement, node tree.Node, code []byte, filePath string) {
	// currentVarName is the current name of the variable on the code
	currentVarName := string(code[node.StartByte():node.EndByte()])

	// If the variable isn't  camelCase, CamelCase or snake_case, we don't modify it (for code safety)
	if !utils.RegexMatch(utils.CamelCase+"|"+utils.SnakeCase, currentVarName) {
		return
	}

	// get the indexes where there's a separator
	upperIndexes := getSeparatorIndexes(currentVarName)
	// with the indexes, separate the line into valid tokens
	tokens := getTokens(upperIndexes, currentVarName)
	newVarName := refactorVarName(tokens)
	fmt.Println(string(newVarName), currentVarName)

	root := parser.GetAST(code, management.GetLanguage())
	defer root.Close()
	query, cursor, captures := parser.GetCapturesByQueries(management.GetLanguage(),
		management.GetVarAppearancesQuery(currentVarName), code, root.RootNode())
	defer query.Close()
	defer cursor.Close()

	index := make(map[uint]string)
	splitedCode := strings.Split(string(code), "\n")

	for {
		match := captures.Next()
		if match == nil {
			break
		}
		copyOf := *match
		node = copyOf.Captures[0].Node
		str := splitedCode[node.StartPosition().Row]

		s, ok := index[node.StartPosition().Row]
		if !ok {
			index[node.StartPosition().Row] = str[:node.StartPosition().Column] + string(newVarName) + str[node.EndPosition().Column:]
			continue
		}
		index[node.StartPosition().Row] = s[:node.EndPosition().Column] + string(newVarName) + str[node.EndPosition().Column:]
	}
	fmt.Println("CODE:", index)

	// Write the new code into the file (just with the modified lines)
	/*err := os.WriteFile("main.py", code, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to write the variable name into the file...")
		os.Exit(1)
	}*/
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
	var selected int8 = 4
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
