/* TODO
--- inputs should be a set i.e. if a directory is passed twice (or a symlink...), it
    should be pulled only once.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Globals
var (
	termWidth   int
	gitdirs     []string
	accumulator strings.Builder
	wg          sync.WaitGroup
	lock        = sync.RWMutex{}
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
	defer logfile.Close()
	if err = initlog(logfile, logpath); err != nil {
		log.Fatal(err)
	}

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

	fmt.Printf("Pulling %d repositor%s..",
		len(gitdirs),
		ternary(len(gitdirs) > 1, "ies", "y"),
	)
	for _, gitdir := range gitdirs {
		wg.Add(1)
		repodir := filepath.Dir(gitdir)
		go pull(repodir)
	}
	wg.Wait()

	less(accumulator.String())

	log.Printf("Processed %d repositor%s in %s\n",
		len(gitdirs),
		ternary(len(gitdirs) > 1, "ies", "y"),
		time.Now().Sub(start))
	log.Println("=== END OF MAIN ===\n")
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

func initlog(logfile *os.File, logpath string) (err error) {
	if logfile, err = os.Create(logpath); err != nil {
		fmt.Printf("cannot log to %q, please choose another file with --log\n", logpath)
		return
	}
	log.SetOutput(logfile)
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	return
}
