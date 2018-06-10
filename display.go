package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
Possible display colors (ANSI colors)
*/
const (
	BLACK   = 30
	RED     = 31
	GREEN   = 32
	YELLOW  = 33
	BLUE    = 34
	MAGENTA = 35
	CYAN    = 36
	WHITE   = 37
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

func color(color int, message string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, message)
}

func red(message string) string {
	return color(RED, message)
}

func green(message string) string {
	return color(GREEN, message)
}

func blue(message string) string {
	return color(BLUE, message)
}

func yellow(message string) string {
	return color(YELLOW, message)
}

func cyan(message string) string {
	return color(CYAN, message)
}

func magenta(message string) string {
	return color(MAGENTA, message)
}
