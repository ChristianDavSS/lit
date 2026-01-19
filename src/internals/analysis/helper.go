package analysis

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

// helper.go - > File that keeps all reusable functions and types of analysis/*.

type Directory struct {
	DirName string
	Content []os.DirEntry
}

// ValidateDate Function to validate dates and show feedback if needed. (since, until flags)
func ValidateDate(date string) *time.Time {
	if len(date) != 10 {
		if len(date) > 0 {
			fmt.Fprintln(os.Stderr, "error with the date format")
		}
		return nil
	}
	d, err := time.Parse("02/01/2006", date)
	if err != nil {
		if len(date) > 0 {
			fmt.Fprintf(os.Stderr, "error parsing the date. %s\n", err)
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
