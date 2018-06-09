package main

import (
	"os/exec"
)

func git(repodir string, a ...string) (output []byte, err error) {
	args := []string{"-C", repodir}
	args = append(args, a...)
	cmd := exec.Command("git", args...)
	return cmd.CombinedOutput()
}
