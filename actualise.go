package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

func actualise(repodir string) {
	start := time.Now()
	var err error
	defer os.Stdout.WriteString(fmt.Sprintf("#%03d\b\b\b\b", counter))
	defer wg.Done()
	var repository Repository
	repository.repodir = repodir
	repository.reponame = path.Base(repodir)
	repository.pullSuccess, repository.statusSuccess = true, true
	var pullOut []byte
	pullArgs := []string{"pull"}
	if pullOut, err = git(repodir, pullArgs...); err != nil {
		log.Printf("first pulling of %q failed....", repository.reponame)
		time.Sleep(1000 * time.Millisecond)
		if pullOut, err = git(repodir, pullArgs...); err != nil {
			log.Printf("%q didn't recover", repository.reponame)
			repository.pullSuccess = false
			repository.pullOutput = pullOut
			repository.process()
			return
		} else {
			log.Printf("%q successfully revovered", repository.reponame)
		}
	}
	repository.pullOutput = pullOut
	var statusOut []byte
	statusArgs := []string{"status", "-sb"}
	if statusOut, err = git(repodir, statusArgs...); err != nil {
		repository.statusSuccess = false
	}
	repository.statusOutput = statusOut
	repository.process()
	message := fmt.Sprintf("%q actualised in %s",
		path.Base(repodir),
		time.Now().Sub(start))
	log.Println(message)
}
