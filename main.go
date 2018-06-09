/* TODO
--- leaner main file (shorter functions, split into more files)
--- don't pass big result map on each goroutine! don't lock it!
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
	"time"
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
	checkBinaries("git", "stty", "less")
	if termWidth, err = getTermWidth(); err != nil {
		termWidth = 80 // *very* conservative
	}
}

func main() {
	start := time.Now()
	var err error
	var showVersion bool
	logpath := filepath.Join(os.TempDir(), "pullall.log")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&logpath, "log", logpath, "log file")
	flag.Parse()

	var logfile *os.File
	if logfile, err = os.Create(logpath); err != nil {
		fmt.Printf("cannot log to %q, please choose another file with --log\n", logpath)
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

	fmt.Print("Looking for .git directories..")
	getGitDirs(basedirs)
	wipeLine()
	if len(gitdirs) == 0 {
		fmt.Println("no git repository found in", inputs)
		os.Exit(1)
	}

	results := make(map[string]Result)

	fmt.Printf("Pulling %d repositor%s..",
		len(gitdirs),
		ternary(len(gitdirs) > 1, "ies", "y"),
	)
	for _, gitdir := range gitdirs {
		wg.Add(1)
		repodir := filepath.Dir(gitdir)
		go pull(repodir, results)
	}
	wg.Wait()
	// we retry the pulls that failed, sequentially this time:
	for repodir, result := range results {
		if result.pullSuccess == false {
			wg.Add(1)
			log.Println("Retrying:", repodir)
			pull(repodir, results)
		}
	}
	wipeLine()
	fmt.Print("Processing...")
	for repodir, result := range results {
		pullSuccess := results[repodir].pullSuccess
		pullOut := results[repodir].pullOutput
		var statusOut []byte
		if statusOut, err = getStatus(repodir, results); err != nil {
			log.Fatal(err)
		}
		results[repodir] = Result{pullSuccess, pullOut, statusOut}
		processRepositoryStatus(repodir, result)
	}
	wipeLine()
	less(accumulator.String())

	log.Printf("Processed %d repositor%s in %s\n",
		len(gitdirs),
		ternary(len(gitdirs) > 1, "ies", "y"),
		time.Now().Sub(start))
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
	start := time.Now()
	for _, input := range inputs {
		if err = filepath.Walk(input, addGitDir); err != nil {
			return
		}
	}
	log.Printf("Got %d git dir%s in %s\n", len(gitdirs),
		ternary(len(gitdirs) > 1, "s", ""),
		time.Now().Sub(start))
	return
}

func addGitDir(item string, finfo os.FileInfo, errin error) (err error) {
	var base, abspath string
	base = path.Base(item)
	if abspath, err = filepath.Abs(item); err != nil {
		return
	}
	if base == ".git" {
		os.Stdout.Write([]byte{'.'})
		gitdirs = append(gitdirs, abspath)
	}
	return
}

func version() {
	fmt.Printf("%s - v.%s (c) Lyderic Landry, London 2018\n",
		APPNAME, VERSION)
}
