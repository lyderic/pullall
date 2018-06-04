package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

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
		if basedirinfo, err = os.Stat(input); err != nil {
			return
		}
		if !basedirinfo.IsDir() {
			return output, fmt.Errorf("not a dir")
		}
		output = append(output, basedir)
	}
	return
}

func printRed(message string) {
	fmt.Printf("\033[31m%s\033[0m\n", message)
}

func getTermWidth() int {
	_, w := getTermDim()
	return w
}

func getTermDim() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	termDim, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	var h, w int
	fmt.Sscan(string(termDim), &h, &w)
	return h, w
}

func wipeLine() {
	fmt.Printf("\r%s\r", strings.Repeat(" ", termWidth))
}

func checkBinaries(binaries ...string) {
	for _, binary := range binaries {
		_, e := exec.LookPath(binary)
		if e != nil {
			log.Fatalf("%s executable not found in path! Aborting...", binary)
		}
	}
}
