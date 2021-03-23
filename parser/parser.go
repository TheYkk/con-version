package parser

import "regexp"

type Commit struct {
	Type           string
	Scope          string
	BreakingChange bool
	Message        string
}

func breaking(message string) bool {
	reg := regexp.MustCompile(`(\n?)BREAKING( -)CHANGE:`)
	return reg.Match([]byte(message))
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
	reg := regexp.MustCompile(`(?i)(docs|fix|feat|chore|style|refactor|perf|test)(?:\((.*)\))?(!?): (.*)`)
	parsed := reg.FindStringSubmatch(message)

	return Commit{
		Type:           parsed[0],
		Scope:          parsed[1],
		BreakingChange: len(parsed[2]) >= 1,
		Message:        parsed[3],
	}, nil
}
