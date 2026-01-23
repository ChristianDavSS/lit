package domain

import (
	"os"

	tree "github.com/tree-sitter/go-tree-sitter"
)

/*
 * language.go: Definition of all the types used in the parsing
 */

// LanguageData - > struct made to register the language an all it's complements (used by the parser)
type LanguageData struct {
	Language *tree.Language
	Queries  string
}

// FunctionData - > struct made to register all the data returned by the parser and save it
type FunctionData struct {
	Name                     string
	TotalParams, Complexity  int
	StartPosition            tree.Point
	StartByte, EndByte, Size uint
	Feedback                 string
}

// Directory saves up the dir name with it's content
type Directory struct {
	DirName string
	Content []os.DirEntry
}

// Feedback is used to send messages to the user when a something is >= value
type Feedback struct {
	Value   int
	Message string
}

type Config struct {
	NamingConventionIndex int8
}

// NodeManagement defines the functions every language struct uses
type NodeManagement interface {
	ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *FunctionData)
	GetLanguage() *tree.Language
	GetQueries() string
	GetVarAppearancesQuery(name string) string
}

var Messages = map[string][]Feedback{
	"parameters": {
		Feedback{Value: 5, Message: "   WARNING: The function takes more parameters than recommended.\n"},
		Feedback{Value: 8, Message: "   WARNING: The function is dangerous and has many responsibilities.\n    You may want to make it cleaner/readable.\n"},
		Feedback{Value: 10, Message: "   POTENTIAL ERROR: The function takes too much parameters itself. You should simplify it.\n"},
	},
	"complexity": {
		Feedback{Value: 10, Message: "   WARNING: The function it's a little bit complex.\n"},
		Feedback{Value: 15, Message: "   POTENTIAL ERROR: The function it's very complex itself. You may simplify or break it down.\n"},
		Feedback{Value: 20, Message: "   POTENTIAL ERROR: The function it's too complex and unreadable, taking too much decisions.\n    You must modify it and make it cleaner!.\n"},
	},
	"size": {
		Feedback{Value: 85, Message: "   WARNING: The functions might be long. You might want to make it smaller.\n"},
		Feedback{Value: 120, Message: "   WARNING: The function is longer than recommended. You should consider a refactor.\n"},
		Feedback{Value: 150, Message: "   POTENTIAL ERROR: The function is too long. You must refactor now!.\n"},
	},
}

// AddInitialData Method to append the initial data into a FunctionData "object"
func (f *FunctionData) AddInitialData(name string, totalParams int, startByte, endByte, size uint, startPos tree.Point) {
	f.Name = name
	f.TotalParams = totalParams
	f.StartByte = startByte
	f.EndByte = endByte
	f.Size = size
	f.StartPosition = startPos
}

// IsTargetInRange validates the range given by another function to validate it's written on the same byte range
func (f *FunctionData) IsTargetInRange(startByte, endByte uint) bool {
	return f.StartByte <= startByte && f.EndByte >= endByte
}

func (f *FunctionData) SetFunctionFeedback() {
	for key, value := range Messages {
		var msg string
		for _, item := range value {
			if f.getValue(key) >= item.Value {
				msg = item.Message
			}
		}
		f.Feedback += msg
	}
}

func (f *FunctionData) getValue(key string) int {
	dict := map[string]int{
		"parameters": f.TotalParams,
		"complexity": f.Complexity,
		"size":       int(f.Size),
	}

	return dict[key]
}
