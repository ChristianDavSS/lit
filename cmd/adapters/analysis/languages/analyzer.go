package languages

import (
	"CLI_App/cmd/adapters/analysis"
	"CLI_App/cmd/domain"
	"path/filepath"
)

type FileAnalyzer struct {
	shouldFix     bool
	activePattern string
}

func NewFileAnalyzer(shouldFix bool, activePattern string) *FileAnalyzer {
	return &FileAnalyzer{
		shouldFix:     shouldFix,
		activePattern: activePattern,
	}
}

// AnalyzeFile analyses the file via DFS (is executed in the scanner)
func (analyzer *FileAnalyzer) AnalyzeFile(filePath string, code []byte) []*domain.FunctionData {
	// Set the variable to save up the language of the current script
	var activeLanguage domain.NodeManagement

	// Save the language for the complexity
	switch filepath.Ext(filePath)[1:] {
	case "js", "jsx":
		activeLanguage = NewJSLanguage(analyzer.activePattern, analyzer.shouldFix)
	case "go":
		activeLanguage = NewGolangLanguage(analyzer.activePattern, analyzer.shouldFix)
	case "java":
		activeLanguage = NewJavaLanguage(analyzer.activePattern, analyzer.shouldFix)
	case "py":
		activeLanguage = NewPythonLanguage(analyzer.activePattern, analyzer.shouldFix)
	}

	// Calculate the cyclical complexity and get the functions returned
	functions := analysis.CyclicalComplexity(activeLanguage, code, filePath)

	i := 0
	for i < len(functions) {
		if functions[i].TotalParams < domain.Messages["parameters"][0].Value &&
			int(functions[i].Size) < domain.Messages["size"][0].Value &&
			functions[i].Complexity < domain.Messages["complexity"][0].Value &&
			functions[i].Feedback == "" {
			functions = append(functions[:i], functions[i+1:]...)
			continue
		}
		functions[i].SetFunctionFeedback()
		i++
	}

	return functions
}
