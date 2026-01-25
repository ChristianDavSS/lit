package types

import (
	"CLI_App/cmd/domain"

	tree "github.com/tree-sitter/go-tree-sitter"
)

// LanguageData - > struct made to register the language an all it's complements (used by the parser)
type LanguageData struct {
	Language *tree.Language
	Queries  string
}

// NodeManagement defines the functions every language struct uses
type NodeManagement interface {
	ManageNode(captureNames []string, code *[]string, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData)
	GetLanguageData() LanguageData
	GetVarAppearancesQuery(name string) string
}
