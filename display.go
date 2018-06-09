package main

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var accumulator strings.Builder

func processRepositoryStatus(repodir string, result Result) {
	addln(repodir)
	if !result.pullSuccess {
		addln(red("--> incorrectly pulled!"))
		return
	}
	pullScanner := bufio.NewScanner(bytes.NewReader(result.pullOutput))
	statusScanner := bufio.NewScanner(bytes.NewReader(result.statusOutput))
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
			add(red(line))
		}
	}
}

func addln(message string) {
	accumulator.WriteString(message)
	accumulator.WriteString("\n")
}

func add(message string) {
	accumulator.WriteString(message)
}

func red(message string) string {
	return fmt.Sprintf("\033[31m%s\033[0m\n", message)
}
