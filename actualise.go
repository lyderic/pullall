package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func actualise(repodir string) {
	var err error
	start := time.Now()
	defer os.Stdout.WriteString(".")
	defer wg.Done()
	var result Result
	result.repodir = repodir
	result.reponame = path.Base(repodir)
	result.pullSuccess, result.statusSuccess = true, true
	var pullOut []byte
	pullArgs := []string{"pull"}
	if pullOut, err = git(repodir, pullArgs...); err != nil {
		result.pullSuccess = false
	}
	result.pullOutput = pullOut
	var statusOut []byte
	statusArgs := []string{"status", "-sb"}
	if statusOut, err = git(repodir, statusArgs...); err != nil {
		result.statusSuccess = false
	}
	result.statusOutput = statusOut
	result.process()
	message := fmt.Sprintf("%q pulled in %s", path.Base(repodir), time.Now().Sub(start))
	log.Println(message)
}
