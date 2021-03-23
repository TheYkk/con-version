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

	if lastTagHash.IsZero() {
		comTime = time.Time{}
	} else {
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

	for _, commit := range commits {
		cmt, err := parser.Parse(commit.Message)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%#v\n", cmt)
	}
}
