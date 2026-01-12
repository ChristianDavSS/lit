package analysis

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
	name, email string
	commits     []*Commit
}

// Commit struct: Struct to save up all the commit information.
type Commit struct {
	hash          plumbing.Hash
	when, message string
	stats         object.FileStats
}

// Function to validate dates and show feedback if needed. (since, until flags)
func validateDate(date string) *time.Time {
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

// matchRegex function: Takes a target string and compares it to a regex pattern, returning the match result
func matchRegex(target, pattern string) bool {
	match, err := regexp.MatchString(pattern, target)
	if err != nil {
		os.Exit(1)
	}
	return match
}
