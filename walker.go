package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var unique map[string]bool

// We have to make sure that each git dir is added only once
func getGitDirs(inputs []string) (gitdirs []string, err error) {
	start := time.Now()
	if err = walkInputs(inputs); err != nil {
		return
	}
	for gitdir := range unique {
		gitdirs = append(gitdirs, gitdir)
	}
	log.Printf("found %d git dir%s in %s\n", len(gitdirs),
		ternary(len(gitdirs) > 1, "s", ""),
		time.Now().Sub(start))
	return
}

func walkInputs(inputs []string) (err error) {
	unique = make(map[string]bool, len(inputs))
	for _, input := range inputs {
		if err = filepath.Walk(input, addGitDir); err != nil {
			return
		}
	}
	return
}

func addGitDir(item string, finfo os.FileInfo, errin error) (err error) {
	var base, abspath string
	base = path.Base(item)
	if abspath, err = filepath.Abs(item); err != nil {
		return
	}
	if base == ".git" {
		os.Stdout.Write([]byte{'.'})
		unique[abspath] = true
	}
	return
}
