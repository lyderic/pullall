package main

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func displayRepositoryStatus(repodir string, result Result) {
	fmt.Println(repodir)
	if !result.pullSuccess {
		printRed("--> incorrectly pulled!")
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
			printRed(line)
		}
	}
	for statusScanner.Scan() {
		line := statusScanner.Text()
		if strings.HasPrefix(line, "##") {
			if strings.Contains(line, "[") {
				printRed(strings.ToUpper(line[26:]))
			} else {
				continue
			}
		} else {
			printRed(line)
		}
	}
}
