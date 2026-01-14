package utils

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// helper.go - > File that keeps all reusable functions and types of analysis/*.

// Contributor type: Struct to save up all the Author data.
type Contributor struct {
	Name, Email string
	Commits     []*Commit
}

// Commit struct: Struct to save up all the commit information.
type Commit struct {
	Hash          plumbing.Hash
	When, Message string
	Stats         object.FileStats
}

// Directory struct: Saves up the dir data (for the DFS)
type Directory struct {
	DirName string
	Content []os.DirEntry
}

// ValidateDate Function to validate dates and show feedback if needed. (since, until flags)
func ValidateDate(date string) *time.Time {
	if len(date) != 10 {
		if len(date) > 0 {
			fmt.Printf("error with the date format")
		}
		return nil
	}
	d, err := time.Parse("02/01/2006", date)
	if err != nil {
		if len(date) > 0 {
			fmt.Printf("error parsing the date. %s", err)
		}
		return nil
	}
	return &d
}

// MatchRegex function: Takes a target string and compares it to a regex pattern, returning the match result
func MatchRegex(target, pattern string) bool {
	match, err := regexp.MatchString(pattern, target)
	if err != nil {
		os.Exit(1)
	}
	return match
}
