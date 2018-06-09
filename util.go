package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

/* make sure we work with absolute paths and the symlinks are resolved */
func sanitize(inputs []string) (output []string, err error) {
	for _, input := range inputs {
		var absoluteInput, resolvedsymlink, basedir string
		if absoluteInput, err = filepath.Abs(input); err != nil {
			return
		}
		if resolvedsymlink, err = filepath.EvalSymlinks(absoluteInput); err != nil {
			return
		}
		if basedir, err = filepath.Abs(resolvedsymlink); err != nil {
			return
		}
		var basedirinfo os.FileInfo
		if basedirinfo, err = os.Stat(input); os.IsNotExist(err) {
			return
		}
		if !basedirinfo.IsDir() {
			return output, fmt.Errorf("not a dir")
		}
		output = append(output, basedir)
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
