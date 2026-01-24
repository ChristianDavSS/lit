package analysis

import (
	"CLI_App/cmd/adapters/analysis/types"
	"CLI_App/cmd/adapters/config"
	"CLI_App/cmd/domain"
	"fmt"
	"os"
	"strings"
)

/*
 * writer.go - > file to write some data on a file using some ast queries.
 */

// FileModifier is an adapter for modifying data into the code script. This way, we can manage the file modification safely
type FileModifier struct {
	management types.NodeManagement
}

func NewFileModifier(management types.NodeManagement) FileModifier {
	return FileModifier{
		management: management,
	}
}

// ModifyVariableName - > this function modifies the variable written the wrong way in the code, rewriting it for you.
// Takes the initial variable name and converts it
// Only converts from one convention to another (safety conditions)
func (f FileModifier) ModifyVariableName(code []byte, filePath, varName string, shouldFix bool) {
	// If the variable isn't  camelCase, CamelCase or snake_case, we don't modify it (for code safety)
	if !domain.RegexMatch(domain.CamelCase+"|"+domain.SnakeCase, varName) || !shouldFix {
		return
	}

	// with the indexes, separate the line into valid tokens
	tokens := getTokens(varName)
	fmt.Println(tokens)
	// get the new variable name (according to the current naming conventions selected)
	newVarName := refactorVarName(tokens)

	// get the query, cursor and captures (applying the query to fetch them)
	root := GetAST(code, f.management.GetLanguageData().Language)
	defer root.Close()
	query, cursor, captures := GetCapturesByQueries(f.management.GetLanguageData().Language,
		f.management.GetVarAppearancesQuery(varName), code, root.RootNode())
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
		node := match.Captures[0].Node
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
		diff[node.StartPosition().Row] += uint(len(newVarName) - len(varName))
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
func getTokens(line string) []string {
	// slice to save up the indexes
	var tokens []string
	var i int

	// iterate through the line of code
	for j, ch := range line {
		if ch == 95 || ch >= 65 && ch <= 90 {
			tokens = append(tokens, strings.ToLower(line[i:j]))
			if ch == 95 {
				i = j + 1
			} else {
				i = j
			}
		}
		j++
	}

	return append(tokens, strings.ToLower(line[i:]))
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
