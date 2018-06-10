/* TODO
--- get rid of the global variables, whenever possible: use pointers
--- when one cannot pull a repository because the ssh key is missing or the passphrase needs
    to be provided, the process should skip, not wait forever until Ctrl-C
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Globals
var (
	logpath     = filepath.Join(os.TempDir(), "pullall.log") // default, can be set with --log flag
	accumulator strings.Builder
	wg          sync.WaitGroup
	lock        = sync.RWMutex{}
)

func init() {
	checkBinaries("git", "stty", "less")
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func main() {

	start := time.Now()
	var err error
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&logpath, "log", logpath, "log file")
	flag.Parse()

	if showVersion {
		version()
		return
	}

	var logfile *os.File
	defer logfile.Close()
	if err = initlog(logfile); err != nil {
		log.Fatal(err)
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
	var gitdirs []string
	if gitdirs, err = getGitDirs(basedirs); err != nil {
		return
	}
	wipeLine()
	if len(gitdirs) == 0 {
		fmt.Println("no git repository found in", inputs)
		os.Exit(1)
	}

	fmt.Printf("Actualising %d repositor%s..",
		len(gitdirs),
		ternary(len(gitdirs) > 1, "ies", "y"),
	)
	for _, gitdir := range gitdirs {
		wg.Add(1)
		repodir := filepath.Dir(gitdir)
		go actualise(repodir)
	}
	wg.Wait()

	wipeLine()
	less(accumulator.String())

	log.Printf("Processed %d repositor%s in %s\n",
		len(gitdirs),
		ternary(len(gitdirs) > 1, "ies", "y"),
		time.Now().Sub(start))
	log.Println("=== END OF MAIN ===")
	log.Println()
}

func version() {
	fmt.Printf("%s - v.%s (c) Lyderic Landry, London 2018\n",
		APPNAME, VERSION)
}

func initlog(logfile *os.File) (err error) {
	if logfile, err = os.Create(logpath); err != nil {
		fmt.Printf("cannot log to %q, please choose another file with --log\n", logpath)
		return
	}
	log.SetOutput(logfile)
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	return
}
