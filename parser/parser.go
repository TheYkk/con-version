package parser

import (
	"regexp"
)

type Commit struct {
	Type           string
	Scope          string
	BreakingChange bool
	Message        string
}

var (
	commitRegex   = regexp.MustCompile(`(?i)(docs|fix|feat|chore|style|refactor|perf|test)(?:\((.*)\))?(!?): (.*)`)
	breakingRegex = regexp.MustCompile(`(\n?)BREAKING( -)CHANGE:`)
)

func breaking(message string) bool {
	return breakingRegex.Match([]byte(message))
}

func IsBreaking(message string) (bool, error) {
	brk := breaking(message)
	if brk {
		return true, nil
	}
	// Parse commit to find breaking change character `!`
	cmt, err := Parse(message)
	if err != nil {
		return false, err
	}

	if cmt.BreakingChange {
		return true, nil
	}

	return false, nil

}

func Parse(message string) (Commit, error) {
	parsed := commitRegex.FindStringSubmatch(message)

	return Commit{
		Type:           parsed[1],
		Scope:          parsed[2],
		BreakingChange: len(parsed[3]) >= 1,
		Message:        parsed[4],
	}, nil
}
