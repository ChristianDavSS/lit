package analysis

import (
	"CLI_App/internal/domain"

	tree "github.com/tree-sitter/go-tree-sitter"
)

// LanguageData struct made to register the language an all it's complements (used by the parser)
type LanguageData struct {
	Language *tree.Language
	Queries  string
}

// NodeManagement defines the functions every language struct uses
type NodeManagement interface {
	ManageNode(captureNames []string, code []byte, filepath string, node tree.QueryCapture, nodeInfo *domain.FunctionData)
	GetLanguage() *tree.Language
	GetQueries() string
	GetVarAppearancesQuery(name string) string
}
