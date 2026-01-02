package analysis

import (
	"CLI_App/src/internals"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// String validScriptPattern for the regex that validates scripts
var validScriptPattern = "^[a-zA-Z0-9._-]+\\.(py|go|java|js|jsx|dart|c|cpp|css|html|ts|md)$"

// String notValidDirPattern for the regex that validates you won´t visit unwanted sites
var notValidDirPattern = "^(node_modules|.*\\.exe|target|venv|__pycache__|" +
	"\\.(git|idea|mvn|cmd))$"

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
	fmt.Println("Files command")
}

// Test loc flag development: get the lines of code of every language
func loc(files []os.DirEntry) {
	languagesMap := make(map[string]int)
	traverseFiles(languagesMap, files, "")
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

// Navigate through the file system with a DFS algorithm.
func traverseFiles(languages map[string]int, files []os.DirEntry, dirName string) {
	for _, v := range files {
		// Check out if the current position contains a file or a directory
		if v.IsDir() {
			// If we should ignore a directory based on our regex, we do.
			if r, _ := regexp.Match(notValidDirPattern, []byte(v.Name())); r {
				continue
			}
			currentDirName := dirName + v.Name()
			fmt.Println("Reading", currentDirName)
			dir, err := os.ReadDir(currentDirName)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			traverseFiles(languages, dir, currentDirName+"/")
		} else {
			// Check if the current file is a programming language script
			if r, _ := regexp.Match(validScriptPattern, []byte(v.Name())); !r {
				continue
			}
			file, err := os.ReadFile(dirName + v.Name())
			if err != nil {
				return
			}
			totalLines := len(strings.Split(string(file), "\n"))
			nameSplit := strings.Split(v.Name(), ".")
			nameLanguage := nameSplit[len(nameSplit)-1]
			language, ok := languages[nameLanguage]
			if !ok {
				languages[nameLanguage] = totalLines
			} else {
				languages[nameLanguage] = language + totalLines
			}
		}
	}
}
