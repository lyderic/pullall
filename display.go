package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var termWidth int // needed for wiping the whole line

func init() {
	var err error
	if termWidth, err = getTermWidth(); err != nil {
		termWidth = 80 // *very* conservative
	}
}

func getTermWidth() (w int, err error) {
	if _, w, err = getTermDim(); err != nil {
		return
	}
	return
}

func getTermDim() (h, w int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	var termDim []byte
	if termDim, err = cmd.Output(); err != nil {
		return
	}
	fmt.Sscan(string(termDim), &h, &w)
	return
}

func wipeLine() {
	fmt.Printf("\r%s\r", strings.Repeat(" ", termWidth))
}

func less(s string) (err error) {
	less := exec.Command("less", "-FRIX")
	less.Stdin = strings.NewReader(s)
	less.Stdout, less.Stderr = os.Stdout, os.Stdout
	return less.Run()
}
