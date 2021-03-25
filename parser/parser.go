// Copyright 2021 Kaan Karakaya
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
