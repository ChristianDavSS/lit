package utils

import (
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

// NodeManagement defines the functions every language struct uses
type NodeManagement interface {
	ManageNode(captureNames []string, code []byte, node tree.QueryCapture, nodeInfo *FunctionData)
	GetLanguage() *tree.Language
	GetQueries() string
	GetVarAppearancesQuery(name string) string
}

// FunctionData - > struct made to register all the data returned by the parser and save it
type FunctionData struct {
	Name                     string
	TotalParams, Complexity  int
	StartPosition            tree.Point
	StartByte, EndByte, Size uint
	Feedback                 string
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
