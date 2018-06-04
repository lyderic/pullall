package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func init() {
	log.SetFlags(log.Lshortfile)
	checkBinaries("git", "stty")
	termWidth = getTermWidth()
}

func main() {

	var err error
	var showVersion bool
	flag.BoolVar(&showVersion, "V", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s - v.%s (c) Lyderic Landry, London 2018\n",
			APPNAME, VERSION)
	}

	inputs := []string{"."}
	if len(flag.Args()) > 0 {
		inputs = flag.Args()
	}

	var basedirs []string
	if basedirs, err = sanitize(inputs); err != nil {
		log.Fatal(err)
	}

	fmt.Print("Looking for .git directories...")
	getGitDirs(basedirs)
	wipeLine()

	if len(gitdirs) == 0 {
		fmt.Println("git repository not found")
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
		var statusOut []byte
		if statusOut, err = getStatus(repodir, results); err != nil {
			log.Fatal(err)
		}
		results[repodir] = Result{pullSuccess, pullOut, statusOut}
		displayRepositoryStatus(repodir, result)
	}

}

func pull(repodir string, results map[string]Result) (err error) {
	defer os.Stdout.WriteString(".")
	defer wg.Done()
  var out []byte
	out, err = exec.Command("git", "-C", repodir, "pull").CombinedOutput()
	var statusOut []byte
	if statusOut, err = getStatus(repodir, results); err != nil {
		return err
	}
	lock.Lock()
	if err != nil {
		results[repodir] = Result{false, out, statusOut}
	} else {
		results[repodir] = Result{true, out, statusOut}
	}
	lock.Unlock()
	return
}

func getStatus(repodir string, results map[string]Result) (output []byte, err error) {
	if output, err = exec.Command("git", "-C", repodir, "status", "-sb").CombinedOutput(); err != nil {
		return
	}
	return
}

func getGitDirs(inputs []string) (err error) {
	for _, input := range inputs {
		if err = filepath.Walk(input, addGitDir); err != nil {
			return
		}
	}
	return
}

func addGitDir(item string, finfo os.FileInfo, errin error) (err error) {
	var base, abspath string
	base = path.Base(item)
	if abspath, err = filepath.Abs(item); err != nil {
		return
	}
	if base == ".git" {
		gitdirs = append(gitdirs, abspath)
	}
	return
}
