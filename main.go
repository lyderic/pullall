/* TODO
--- leaner main file (shorter functions, split into more files)
--- don't pass big result map on each goroutine! don't lock it!
--- don't retry every repo! only the ones that fail, in the same goroutine
--- pipe results to less
--- provide a --log option to debug and see what's going on
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var (
	termWidth int
	gitdirs   []string
	wg        sync.WaitGroup
	lock      = sync.RWMutex{}
)

func init() {
	var err error
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	checkBinaries("git", "stty")
	if termWidth, err = getTermWidth(); err != nil {
		termWidth = 80 // *very* conservative
	}
}

func main() {

	var err error
	var showVersion bool
	logpath := filepath.Join(os.TempDir(), "pullall.log")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&logpath, "log", logpath, "log file")
	flag.Parse()

	var logfile *os.File
	if logfile, err = os.Create(logpath); err != nil {
		fmt.Printf("cannot log to %q, please choose another file with --log", logpath)
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	if showVersion {
		version()
		return
	}

	inputs := []string{"."}
	if len(flag.Args()) > 0 {
		inputs = flag.Args()
	}
	log.Println("inputs:", inputs)

	var basedirs []string
	if basedirs, err = sanitize(inputs); err != nil {
		fmt.Println("input not valid:", inputs)
		fmt.Println(err)
		log.Fatal(err)
	}
	log.Println("basedirs:", basedirs)

	fmt.Print("Looking for .git directories...")
	getGitDirs(basedirs)
	wipeLine()
	log.Println("gitdirs:", gitdirs)

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
	// we retry the pulls that failed, sequentially this time:
	// this is not toooo bad as we don't expect much first pulls to have failed
	for repodir, result := range results {
		if result.pullSuccess == false {
			wg.Add(1)
			log.Println("We are retrying:", repodir)
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

	log.Println(inputs, "all pulled")
	log.Println("=== END OF MAIN ===\n")
}

func getStatus(repodir string, results map[string]Result) (output []byte, err error) {
	args := []string{"status", "-sb"}
	if output, err = git(repodir, args...); err != nil {
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

func version() {
	fmt.Printf("%s - v.%s (c) Lyderic Landry, London 2018\n",
		APPNAME, VERSION)
}
