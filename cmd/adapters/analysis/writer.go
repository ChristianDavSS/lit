package analysis

import (
	"CLI_App/cmd/adapters/config"
	"CLI_App/cmd/domain"
	"fmt"
	"os"
	"strings"

	tree "github.com/tree-sitter/go-tree-sitter"
)

/*
 * writer.go - > file to write some data on a file using some ast queries.
 */

// FileModifier is -
type FileModifier struct {
	management          domain.NodeManagement
	varAppearancesQuery string
}

func NewFileModifier(management domain.NodeManagement, varName string) FileModifier {
	return FileModifier{
		management:          management,
		varAppearancesQuery: management.GetVarAppearancesQuery(varName),
	}
}

// ModifyVariableName - > this function modifies the variable written the wrong way in the code, rewriting it for you.
// Takes the node captured (the one we want to modify) and the path.
// Only converts from one convention to another (safety conditions)
func (f FileModifier) ModifyVariableName(node tree.Node, code []byte, filePath string) {
	// currentVarName is the current name of the variable on the code
	currentVarName := string(code[node.StartByte():node.EndByte()])

	// If the variable isn't  camelCase, CamelCase or snake_case, we don't modify it (for code safety)
	if !domain.RegexMatch(domain.CamelCase+"|"+domain.SnakeCase, currentVarName) {
		return
	}

	// get the indexes where there's a separator
	upperIndexes := getSeparatorIndexes(currentVarName)
	// with the indexes, separate the line into valid tokens
	tokens := getTokens(upperIndexes, currentVarName)
	// get the new variable name (according to the current naming conventions selected)
	newVarName := refactorVarName(tokens)

	// get the query, cursor and captures (applying the query to fetch them)
	root := GetAST(code, f.management.GetLanguageData().Language)
	defer root.Close()
	query, cursor, captures := GetCapturesByQueries(f.management.GetLanguageData().Language,
		f.varAppearancesQuery, code, root.RootNode())
	defer query.Close()
	defer cursor.Close()

	// cache of the sum of the difference between lengths of the variables
	diff := make(map[uint]uint)
	// slice the code into lines (just as the script
	slicedCode := strings.Split(string(code), "\n")

	// loop through the node captures
	for {
		// get the next match
		match := captures.Next()
		if match == nil {
			break
		}
		// get the node from the captures (just one capture per match)
		node = match.Captures[0].Node
		// get the current line of code (the one that'll be modified
		str := slicedCode[node.StartPosition().Row]
		// check if there's a value of this row in the cache
		_, ok := diff[node.StartPosition().Row]
		// if there's not, we initialize it to 0
		if !ok {
			diff[node.StartPosition().Row] = 0
		}
		// get the value from the cache of that position (at this point, there key will always have a value)
		value := diff[node.StartPosition().Row]

		// modify the line of code using the cache and slicing
		slicedCode[node.StartPosition().Row] = str[:node.StartPosition().Column+value] + newVarName + str[node.EndPosition().Column+value:]

		// update the value of the row, adding up the difference of lengths
		diff[node.StartPosition().Row] += uint(len(newVarName) - len(currentVarName))
	}

	// Write the modified code into the file (with the new variable names where they belong)
	err := os.WriteFile(filePath, []byte(strings.Join(slicedCode, "\n")), 0644)
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
func refactorVarName(tokens []string) string {
	var newName = tokens[0]

	jsonAdapter := config.NewJSONAdapter()
	switch jsonAdapter.GetConfig().NamingConventionIndex {
	case 2:
		newName = ""
		camelCases(&newName, tokens)
	case 1, 3:
		camelCases(&newName, tokens[1:])
	case 4:
		for _, token := range tokens[1:] {
			newName += "_" + token
		}
	}
	return newName
}

// function with the logics for the camelCase and CamelCase conversions
func camelCases(target *string, tokens []string) {
	for _, token := range tokens {
		*target += string(token[0]-32) + token[1:]
	}
}
