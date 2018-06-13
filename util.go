package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

/* make sure we work with absolute paths, the symlinks are resolved and
   there is no duplicates */
func sanitize(dir string) (cleandir string, err error) {
	var absolutePath, resolvedSymlink string
	if absolutePath, err = filepath.Abs(dir); err != nil {
		return
	}
	if resolvedSymlink, err = filepath.EvalSymlinks(absolutePath); err != nil {
		return
	}
	if cleandir, err = filepath.Abs(resolvedSymlink); err != nil {
		return
	}
	var finfo os.FileInfo
	if finfo, err = os.Stat(cleandir); os.IsNotExist(err) {
		return
	}
	if !finfo.IsDir() {
		return cleandir, fmt.Errorf("%q: not a directory", cleandir)
	}
	return
}

func checkBinaries(binaries ...string) {
	for _, binary := range binaries {
		_, e := exec.LookPath(binary)
		if e != nil {
			log.Fatalf("%s executable not found in path! Aborting...", binary)
		}
	}
}

func ternary(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}
