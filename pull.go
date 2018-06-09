package main

import (
	"log"
	"os"
	"time"
)

func pull(repodir string, results map[string]Result) (err error) {
	start := time.Now()
	defer os.Stdout.WriteString(".")
	defer wg.Done()
	var pullOut []byte
	args := []string{"pull"}
	if pullOut, err = git(repodir, args...); err != nil {
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
	log.Printf("%s pulled in %s\n", repodir, time.Now().Sub(start))
	return
}
