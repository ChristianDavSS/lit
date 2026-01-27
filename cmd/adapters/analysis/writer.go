package analysis

import (
	"CLI_App/cmd/adapters/analysis/types"
	"CLI_App/cmd/adapters/config"
	"strings"
)

/*
 * writer.go - > file to write some data on a file using some ast queries.
 */

// FileModifier is an adapter for modifying data into the code script. This way, we can manage the file modification safely
type FileModifier struct {
	management    types.NodeManagement
	activePattern string
}

func NewFileModifier(management types.NodeManagement, activePattern string) FileModifier {
	return FileModifier{
		management:    management,
		activePattern: activePattern,
	}
}

// ModifyVariableName - > this function modifies the variable written the wrong way in the code, rewriting it for you.
// Takes the initial variable name and converts it
// Only converts from one convention to another (safety conditions)
func (f FileModifier) ModifyVariableName(code *[]string) {
	// get the query, cursor and captures (applying the query to fetch them)
	root := GetAST(code, f.management.GetLanguageData().Language)
	defer root.Close()
	query, cursor, captures := GetCapturesByQueries(f.management.GetLanguageData().Language,
		f.management.GetVarAppearancesQuery(f.activePattern), code, root.RootNode())
	defer query.Close()
	defer cursor.Close()

	localCache := map[uint]struct {
		names   []string
		lastIdx uint
	}{}

	// loop through the node captures
	for {
		// get the next match
		match := captures.Next()
		if match == nil {
			break
		}
		copyOf := *match
		// get the node from the captures (just one capture per match)
		node := copyOf.Captures[0].Node
		newName := refactorVarName(getTokens((*code)[node.StartPosition().Row][node.StartPosition().Column:node.EndPosition().Column]))

		value, ok := localCache[node.StartPosition().Row]
		if !ok {
			value = struct {
				names   []string
				lastIdx uint
			}{
				names:   make([]string, 0),
				lastIdx: 0,
			}
			localCache[node.StartPosition().Row] = value
		}
		value.names = append(value.names, newName)
		value.lastIdx = node.EndPosition().Column

		localCache[node.StartPosition().Row] = value
	}

	for key, value := range localCache {
		(*code)[key] = strings.Join(value.names, ", ") + (*code)[key][value.lastIdx:]
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
