package main

import (
	"flag"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/theykk/con-version/parser"
	"log"
	"os"
	"time"
)

var (
	Version = "dev"
)

func main() {
	// Include file and line number on log
	log.SetFlags(log.Lshortfile)
	log.Println("Con-version:", Version)
	// Get current executed dir
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := flag.String("dir", pwd, "Path of git repo")
	flag.Parse()

	r, err := git.PlainOpen(*path)
	if err != nil {
		log.Fatal(err)
	}

	// ... retrieving the HEAD reference
	tags, err := r.Tags()
	if err != nil {
		log.Fatal(err)
	}

	var lastTagHash plumbing.Hash
	var lastTag string
	err = tags.ForEach(func(reference *plumbing.Reference) error {
		lastTagHash = reference.Hash()
		lastTag = reference.Name().Short()
		log.Println(lastTag, reference.Name(), reference.Hash())
		return nil
	})

	if len(lastTag) == 0 {
		lastTag = "0.0.0"
	}

	if err != nil {
		log.Fatal(err)
	}
	var comTime time.Time

	if !lastTagHash.IsZero() {
		ob, err := r.CommitObject(lastTagHash)
		if err != nil {
			log.Fatal(err)
		}

		comTime = ob.Committer.When.Add(time.Millisecond)
	}

	commitsSince, err := r.Log(&git.LogOptions{Since: &comTime})
	if err != nil {
		log.Fatal(err)
	}

	var commits []*object.Commit
	_ = commitsSince.ForEach(func(commit *object.Commit) error {
		commits = append(commits, commit)
		return nil
	})

	var newVersion semver.Version

	// Parse tag with semver
	cVersion, err := semver.NewVersion(lastTag)
	if err != nil {
		log.Fatal(err)
	}

	for _, commit := range commits {
		cmt, err := parser.Parse(commit.Message)
		if err != nil {
			log.Fatal(err)
		}

		if cmt.BreakingChange {
			newVersion = cVersion.IncMajor()
		} else if cmt.Type == "feat" {
			newVersion = cVersion.IncMinor()
		} else if cmt.Type == "fix" {
			newVersion = cVersion.IncPatch()
		}

		// Update current version
		cVersion = &(newVersion)

		log.Printf("%#v\n", cmt)
	}

	log.Printf("\033[1;32m%s\033[0m", newVersion.String())
}
