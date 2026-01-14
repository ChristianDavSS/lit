package analysis

import (
	"CLI_App/src/internals"
	"CLI_App/src/internals/analysis/utils"
	"CLI_App/src/internals/ast"
	"CLI_App/src/internals/ast/languages"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// languagesMap - > map where we save all the data from the loc flag
var languagesMap = make(map[string]int)

// DangerousFunctions map - > Map to save up the dangerous functions per script
var DangerousFunctions = make(map[string][]*languages.FunctionData)

// Files - > Entry point for the command line with the flags
func Files(locFlag bool) {
	files, err := os.ReadDir(internals.GetWorkingDirectory())
	if err != nil {
		fmt.Println("Error reading the structure: ", err)
		os.Exit(1)
	}
	if locFlag {
		loc(files)
		return
	}
	traverseFiles(files, "", fileScanner, utils.ScanValidScriptPattern)
	printDangerousFunctions()
}

// printDangerousFunctions function - > Prints out the result of the file scanner.
func printDangerousFunctions() {
	fmt.Printf("\nDangerous functions found in the project: %d\n", len(DangerousFunctions))
	for key, value := range DangerousFunctions {
		fmt.Printf("- %s:\n", key)
		for _, item := range value {
			fmt.Printf(" * Function %s (at %d:%d)\n", item.Name, item.StartPosition.Row, item.StartPosition.Column)
			fmt.Printf("   Parameters: %d\n   Total lines of code: %d\n", item.TotalParams, item.Size)
			fmt.Println(item.Feedback)
		}
		fmt.Println()
	}
}

// Test loc flag development: get the lines of code of every language
func loc(files []os.DirEntry) {
	traverseFiles(files, "", addToLanguagesMap, utils.LocValidScriptPattern)
	fmt.Println()
	fmt.Println("Results (language -> total lines of code):")
	total := 0.0
	// Get the total lines of code (so I can show the percentages per language)
	for _, v := range languagesMap {
		total += float64(v)
	}

	// Get the results
	for k, v := range languagesMap {
		fmt.Printf("%s %d (%.1f%%)\n", k, v, (float64(v)*100)/total)
	}
	fmt.Println("Total lines of code:", total)
}

// addToLanguagesMap - > Used as an argument when the loc flag is true
func addToLanguagesMap(filename string, code []byte) {
	totalLines := len(strings.Split(string(code), "\n"))
	nameSplit := strings.Split(filename, ".")
	nameLanguage := nameSplit[len(nameSplit)-1]
	language, ok := languagesMap[nameLanguage]
	if !ok {
		languagesMap[nameLanguage] = totalLines
	} else {
		languagesMap[nameLanguage] = language + totalLines
	}
}

// ---------------------------------------------------------------------

// fileScanner get the full name of the file and the code, calling the parser on the code
func fileScanner(filename string, code []byte) {
	language := strings.Split(filename, ".")
	functions := ast.RunParser(code, language[len(language)-1])
	// If there's any function returned, we save it up
	if len(functions) > 0 {
		DangerousFunctions[filename] = append(DangerousFunctions[filename], functions...)
	}
}

// Navigate through the file system with a DFS algorithm.
func traverseFiles(files []os.DirEntry, dirName string, fileFunction func(filename string, code []byte), validScriptPattern string) {
	for _, v := range files {
		// Check out if the current position contains a file or a directory
		if v.IsDir() {
			// If we should ignore a directory based on our regex, we do.
			if r, _ := regexp.Match(utils.NotValidDirPattern, []byte(v.Name())); r {
				continue
			}
			currentDirName := dirName + v.Name()
			fmt.Println("Reading", currentDirName)
			dir, err := os.ReadDir(currentDirName)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			traverseFiles(dir, currentDirName+"/", fileFunction, validScriptPattern)
		} else {
			// Check if the current file is a programming language script
			if r, _ := regexp.Match(validScriptPattern, []byte(v.Name())); !r {
				continue
			}
			file, err := os.ReadFile(dirName + v.Name())
			if err != nil {
				return
			}
			fileFunction(dirName+v.Name(), file)
		}
	}
}
