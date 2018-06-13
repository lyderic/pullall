package main

import (
	"log"
	"os"
	"path/filepath"
)

func walkBasedir(basedir string) (err error) {
	if err = filepath.Walk(basedir, actualiseGitDir); err != nil {
		return
	}
	wg.Wait()
	return
}

func actualiseGitDir(item string, finfo os.FileInfo, errin error) (err error) {
	if errin != nil {
		log.Printf("cannot access %q: %v", item, errin)
		return nil // simply skip, don't kill the whole thing
	}
	if finfo.Name() == ".git" {
		counter++
		var abspath string
		if abspath, err = filepath.Abs(item); err != nil {
			return
		}
		wg.Add(1)
		repodir := filepath.Dir(abspath)
		go actualise(repodir)
	}
	return
}
