package main

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Result struct {
	repodir       string
	reponame      string
	pullSuccess   bool
	statusSuccess bool
	pullOutput    []byte
	statusOutput  []byte
}

func (r Result) String() string {
	return fmt.Sprintf("Result{pullSuccess: %t - statusSuccess: \t\n pullOutput: %s\n statusOutput: %s",
		r.statusSuccess,
		r.pullSuccess,
		string(r.pullOutput),
		string(r.statusOutput),
	)
}

func (r Result) process() {
	addln(r.repodir)
	if !r.pullSuccess {
		addln(red("--> incorrectly pulled!"))
		return
	}
	pullScanner := bufio.NewScanner(bytes.NewReader(r.pullOutput))
	statusScanner := bufio.NewScanner(bytes.NewReader(r.statusOutput))
	for pullScanner.Scan() {
		line := pullScanner.Text()
		match, _ := regexp.MatchString("(?i)already up.*to.*date.*", line)
		if match {
			continue
		} else {
			addln(red(line))
		}
	}
	for statusScanner.Scan() {
		line := statusScanner.Text()
		if strings.HasPrefix(line, "##") {
			if strings.Contains(line, "[") {
				addln(red(strings.ToUpper(line[26:])))
			} else {
				continue
			}
		} else {
			addln(red(line))
		}
	}
}

func add(message string) {
	accumulator.WriteString(message)
}

func addln(message string) {
	accumulator.WriteString(message)
	accumulator.WriteString("\n")
}

func red(message string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", message)
}
