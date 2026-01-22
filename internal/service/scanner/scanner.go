package scanner

import (
	"CLI_App/internal/adapter/analysis"
	"CLI_App/internal/adapter/analysis/languages"
	"CLI_App/internal/domain"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// ScanService orchestrates the scanning process.
type ScanService struct {
	config          *domain.Config
	namingPattern   string
	shouldFix       bool
	wg              sync.WaitGroup
	mu              sync.Mutex

	// Results
	languagesMap       map[string]int
	dangerousFunctions map[string][]*domain.FunctionData
}

// NewScanService creates a new ScanService.
func NewScanService(cfg *domain.Config, namingPattern string, shouldFix bool) *ScanService {
	return &ScanService{
		config:             cfg,
		namingPattern:      namingPattern,
		shouldFix:          shouldFix,
		languagesMap:       make(map[string]int),
		dangerousFunctions: make(map[string][]*domain.FunctionData),
	}
}

// ScanProject scans the project for dangerous functions and naming violations.
func (s *ScanService) ScanProject(root string) (map[string][]*domain.FunctionData, error) {
	files := getDirEntries(root)
	s.traverseFiles(root, files, s.fileScanner, domain.ScanValidScriptPattern)
	return s.dangerousFunctions, nil
}

// CalculateLOC calculates the lines of code per language.
func (s *ScanService) CalculateLOC(root string) (map[string]int, error) {
	files := getDirEntries(root)
	s.traverseFiles(root, files, s.addToLanguagesMap, domain.LocValidScriptPattern)
	return s.languagesMap, nil
}

// fileScanner scans a single file.
func (s *ScanService) fileScanner(filename string, code []byte) {
	defer s.wg.Done()

	if code == nil {
		return
	}

	ext := filepath.Ext(filename)
	if len(ext) < 2 {
		return
	}
	// Remove dot
	langExt := ext[1:]

	var langAdapter analysis.NodeManagement
	switch langExt {
	case "go":
		langAdapter = languages.NewGoLanguage(s.namingPattern, s.shouldFix, s.config.NamingConventionIndex)
	case "java":
		langAdapter = languages.NewJavaLanguage(s.namingPattern, s.shouldFix, s.config.NamingConventionIndex)
	case "js", "jsx":
		langAdapter = languages.NewJSLanguage(s.namingPattern, s.shouldFix, s.config.NamingConventionIndex)
	case "py":
		langAdapter = languages.NewPythonLanguage(s.namingPattern, s.shouldFix, s.config.NamingConventionIndex)
	default:
		// Not supported or no adapter
		return
	}

	functions := analysis.CyclicalComplexity(langAdapter, code, filename)

	// Filter functions
	var filtered []*domain.FunctionData
	for _, f := range functions {
		// Logic from registry.go
		if f.TotalParams <= domain.Messages["parameters"][0].Value &&
			f.Complexity <= domain.Messages["complexity"][0].Value &&
			int(f.Size) <= domain.Messages["size"][0].Value &&
			f.Feedback == "" {
			continue
		}
		f.SetFunctionFeedback()
		filtered = append(filtered, f)
	}

	if len(filtered) > 0 {
		s.mu.Lock()
		s.dangerousFunctions[filename] = append(s.dangerousFunctions[filename], filtered...)
		s.mu.Unlock()
	}
}

// addToLanguagesMap counts LOC.
func (s *ScanService) addToLanguagesMap(filename string, code []byte) {
	defer s.wg.Done()
	totalLines := len(strings.Split(string(code), "\n"))
	nameSplit := strings.Split(filename, ".")
	nameLanguage := nameSplit[len(nameSplit)-1]

	s.mu.Lock()
	s.languagesMap[nameLanguage] += totalLines
	s.mu.Unlock()
}

type directory struct {
	DirName string
	Content []os.DirEntry
}

// traverseFiles navigates the file system.
func (s *ScanService) traverseFiles(root string, initialFiles []os.DirEntry, fileFunction func(filename string, code []byte), validScriptPattern string) {
	stack := []directory{{DirName: root, Content: initialFiles}}

	for len(stack) > 0 {
		files := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, v := range files.Content {
			if v.IsDir() {
				if r, _ := regexp.Match(domain.NotValidDirPattern, []byte(v.Name())); r {
					continue
				}
				dirPath := filepath.Join(files.DirName, v.Name())
				dirContent := getDirEntries(dirPath)
				stack = append(stack, directory{DirName: dirPath, Content: dirContent})
			} else {
				if r, _ := regexp.Match(validScriptPattern, []byte(v.Name())); !r {
					continue
				}
				filePath := filepath.Join(files.DirName, v.Name())
				file := readFile(filePath)
				s.wg.Add(1)
				go fileFunction(filePath, file)
			}
		}
	}
	s.wg.Wait()
}

// Helpers

func getDirEntries(name string) []os.DirEntry {
	files, err := os.ReadDir(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading the file structure: ", err)
		return []os.DirEntry{} // Return empty instead of exit
	}
	return files
}

func readFile(name string) []byte {
	file, err := os.ReadFile(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading the file %s. Please report the issue.\n", name)
		return nil
	}
	return file
}
