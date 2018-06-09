package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
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

func (r Result) process() {
	addln(r.repodir)
	if !r.pullSuccess {
		reportFail("pull", r)
		return
	}
	if !r.statusSuccess {
		reportFail("status", r)
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

func (r Result) String() string {
	return fmt.Sprintf("Result{\n  • repodir: %s\n  • reponame: %s\n  • pullSuccess: %t\n  • statusSuccess: %t\n  • pullOutput: %q\n  • statusOutput: %q\n}",
		r.repodir,
		r.reponame,
		r.pullSuccess,
		r.statusSuccess,
		string(r.pullOutput),
		string(r.statusOutput),
	)
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

func reportFail(action string, r Result) {
	addln(red(fmt.Sprintf("--> git %s failed: see log file: %q", action, logpath)))
	log.Printf("git %s failed for:\n%s", action, r)
}
