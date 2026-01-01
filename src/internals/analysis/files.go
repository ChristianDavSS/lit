package analysis

import (
	"CLI_App/src/internals"
	"fmt"
	"os"
	"strings"
)

// Test loc flag development: get the lines of code of every language
func Test() {
	files, err := os.ReadDir(internals.GetWorkingDirectory())
	if err != nil {
		fmt.Println("Error reading the structure: ", err)
		os.Exit(1)
	}

	languagesMap := make(map[string]int)
	traverseFiles(languagesMap, files, "")
	fmt.Println(languagesMap)
}

// Read the .gitignore content
func readGitIgnore() string {
	gitign, err := os.ReadFile(".gitignore")
	if err != nil {
		fmt.Println("There´s not a .gitignore defined")
		os.Exit(1)
	}
	return string(gitign)
}

// Read the files
func traverseFiles(languages map[string]int, files []os.DirEntry, dirName string) {
	for _, v := range files {
		// Use gitignore content instead
		if strings.Contains(readGitIgnore(), v.Name()) || v.Name() == ".gitignore" {
			continue
		}
		if v.IsDir() {
			currentDirName := dirName + v.Name()
			fmt.Println("Reading ", currentDirName)
			dir, err := os.ReadDir(currentDirName)
			if err != nil {
				fmt.Println("ERROR -> ", err)
				os.Exit(1)
			}
			traverseFiles(languages, dir, currentDirName+"/")
		} else {
			file, err := os.ReadFile(dirName + v.Name())
			if err != nil {
				return
			}
			totalLines := len(strings.Split(string(file), "\n"))
			nameLanguage := strings.Split(v.Name(), ".")[1]
			language, ok := languages[nameLanguage]
			if !ok {
				languages[nameLanguage] = totalLines
			} else {
				languages[nameLanguage] = language + 1
			}
		}
	}
}
