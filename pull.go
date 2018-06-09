package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func pull(repodir string) (err error) {
	start := time.Now()
	defer os.Stdout.WriteString(".")
	defer wg.Done()
	var result Result
	result.pullSuccess = false
	var pullOut []byte
	pullArgs := []string{"pull"}
	if pullOut, err = git(repodir, pullArgs...); err != nil {
		return
	}
	result.pullOutput = pullOut
	var statusOut []byte
	statusArgs := []string{"status", "-sb"}
	if statusOut, err = git(repodir, statusArgs...); err != nil {
		return
	}
	result.statusOutput = statusOut
	result.pullSuccess = true
	processResult(repodir, result)
	message := fmt.Sprintf("%q pulled in %s", path.Base(repodir), time.Now().Sub(start))
	//fmt.Println(message)
	log.Println(message)
	return
}
