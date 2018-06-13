package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Repository struct {
	repodir       string
	reponame      string
	pullSuccess   bool
	statusSuccess bool
	pullOutput    []byte
	statusOutput  []byte
}

func (r Repository) process() {
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
			addln(green(line))
		}
	}
	for statusScanner.Scan() {
		line := statusScanner.Text()
		if strings.HasPrefix(line, "##") {
			if strings.Contains(line, "[") {
				addln(yellow(strings.ToUpper(line[26:])))
			} else {
				continue
			}
		} else {
			addln(blue(line))
		}
	}
}

func (r Repository) String() string {
	return fmt.Sprintf("Repository{\n  • repodir: %s\n  • reponame: %s\n  • pullSuccess: %t\n  • statusSuccess: %t\n  • pullOutput: %q\n  • statusOutput: %q\n}",
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

func reportFail(action string, r Repository) {
	addln(red(fmt.Sprintf("⯁ git %s failed:", action)))
	scanner := bufio.NewScanner(bytes.NewReader(r.pullOutput))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			addln(red(line))
		}
	}
	log.Printf("git %s failed for:\n%s", action, r)
}
