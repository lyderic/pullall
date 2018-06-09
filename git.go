package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func git(repodir string, a ...string) (output []byte, err error) {
	return simpleGit(repodir, a...)
	//return gitTimeOut(10, repodir, a...)
}

/* execute git command into a repo dir and with timeout */
func gitTimeOut(timeout int, repodir string, a ...string) (output []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	args := []string{"-C", repodir}
	args = append(args, a...)
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdin = os.Stdin
	output, err = cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return output, fmt.Errorf("git %s t timed out after %d seconds", timeout)
	}
	return
}

func simpleGit(repodir string, a ...string) (output []byte, err error) {
	args := []string{"-C", repodir}
	args = append(args, a...)
	cmd := exec.Command("git", args...)
	return cmd.CombinedOutput()
}
