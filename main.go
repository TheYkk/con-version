package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/theykk/con-version/parser"
	"log"
	"os"
	"time"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.Lshortfile)
	// ... retrieving the HEAD reference
	tags, err := r.Tags()
	if err != nil {
		log.Fatal(err)
	}

	var lastTag plumbing.Hash
	err = tags.ForEach(func(reference *plumbing.Reference) error {
		lastTag = reference.Hash()
		log.Println(reference.Name(), reference.Hash())
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	var comTime time.Time

	if lastTag.IsZero() {
		comTime = time.Time{}
	} else {
		ob, err := r.CommitObject(lastTag)
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

	for _, commit := range commits {
		//println(commit.Message)
		//_ = commit
		log.Println(parser.Parse(commit.Message))
	}
}
