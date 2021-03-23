package main

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os"
	"time"
)

func main() {
	path, err := os.Getwd()
	CheckIfError(err)
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// ... retrieving the HEAD reference
	tags, err := r.Tags()
	CheckIfError(err)

	var lastTag plumbing.Hash
	err = tags.ForEach(func(reference *plumbing.Reference) error {
		lastTag = reference.Hash()
		log.Println(reference.Name(), reference.Hash())
		return nil
	})
	log.Println(lastTag)
	CheckIfError(err)

	ob, err := r.CommitObject(lastTag)
	CheckIfError(err)
	comTime := ob.Committer.When.Add(time.Millisecond)

	commitsSince, err := r.Log(&git.LogOptions{Since: &comTime})
	CheckIfError(err)
	var commits []*object.Commit
	_ = commitsSince.ForEach(func(commit *object.Commit) error {
		commits = append(commits, commit)
		return nil
	})

	for _, commit := range commits {
		println(commit.Message)
		//_ = commit
	}
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
