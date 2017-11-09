package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func init() {
	log.SetFlags(log.Lshortfile)
	checkBinaries("git", "stty")
	termWidth = getTermWidth()
}

func main() {

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s - v.%s (c) Lyderic Landry, London 2017\n",
			appname, appversion)
		return
	}

	inputs := []string{"."}
	if len(os.Args) > 1 {
		inputs = os.Args[1:]
	}

	basedirs := sanitize(inputs)

	fmt.Print("Looking for .git directories...")
	getGitDirs(basedirs)
	wipeLine()

	if len(gitdirs) == 0 {
		fmt.Println("No git repository found.")
		os.Exit(1)
	}

	results := make(map[string]Result)

	fmt.Print("Pulling repositories..")
	for _, gitdir := range gitdirs {
		wg.Add(1)
		repodir := filepath.Dir(gitdir)
		go pull(repodir, results)
	}
	wg.Wait()
	for repodir, result := range results {
		if result.pullSuccess == false {
			wg.Add(1)
			pull(repodir, results)
		}
	}
	wipeLine()

	for repodir, result := range results {
		pullSuccess := results[repodir].pullSuccess
		pullOut := results[repodir].pullOutput
		statusOut := getStatus(repodir, results)
		results[repodir] = Result{pullSuccess, pullOut, statusOut}
		displayRepositoryStatus(repodir, result)
	}

}

func pull(repodir string, results map[string]Result) {
	defer os.Stdout.WriteString(".")
	defer wg.Done()
	pullOut, pullErr := exec.Command("git", "-C", repodir, "pull").CombinedOutput()
	statusOut := getStatus(repodir, results)
	lock.Lock()
	if pullErr != nil {
		results[repodir] = Result{false, pullOut, statusOut}
	} else {
		results[repodir] = Result{true, pullOut, statusOut}
	}
	lock.Unlock()
}

func getStatus(repodir string, results map[string]Result) []byte {
	statusOut, statusErr := exec.Command("git", "-C", repodir, "status", "-sb").CombinedOutput()
	if statusErr != nil {
		log.Fatalln("Error getting status", repodir, ":",
			statusErr, string(statusOut))
	}
	return statusOut
}

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
		//if line == "Already up-to-date." || line == "Already up to date." {
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

func getGitDirs(inputs []string) {
	for _, input := range inputs {
		err := filepath.Walk(input, addGitDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func addGitDir(item string, info os.FileInfo, err error) error {
	base := path.Base(item)
	abspath, _ := filepath.Abs(item)
	if base == ".git" {
		gitdirs = append(gitdirs, abspath)
	}
	return nil
}
