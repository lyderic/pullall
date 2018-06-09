package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func pull(repodir string, results map[string]Result) (err error) {
	start := time.Now()
	defer os.Stdout.WriteString(".")
	defer wg.Done()
	var pullOut []byte
	pullArgs := []string{"pull"}
	if pullOut, err = git(repodir, pullArgs...); err != nil {
		return
	}
	var statusOut []byte
	if statusOut, err = getStatus(repodir, results); err != nil {
		return
	}
	lock.Lock() // this lock is a problem and I think it can be resolved with a pointer to this map
	if err != nil {
		results[repodir] = Result{false, pullOut, statusOut}
	} else {
		results[repodir] = Result{true, pullOut, statusOut}
	}
	lock.Unlock()
	message := fmt.Sprintf("%q pulled in %s", path.Base(repodir), time.Now().Sub(start))
	//fmt.Print(message)
	log.Println(message)
	return
}

func getStatus(repodir string, results map[string]Result) (output []byte, err error) {
	args := []string{"status", "-sb"}
	if output, err = git(repodir, args...); err != nil {
		return
	}
	return
}
