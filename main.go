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
	logpath     = filepath.Join(os.TempDir(), "pullall.log")
	accumulator strings.Builder
	wg          sync.WaitGroup
	counter     int
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

	dir := "."
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}
	log.Println("dir:", dir)

	var basedir string
	if basedir, err = sanitize(dir); err != nil {
		fmt.Println("input not valid:", dir)
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Print("Please wait, actualising repository ")
	hideCursor()
	if err = walkBasedir(basedir); err != nil {
		return
	}
	wipeLine()
	showCursor()
	if counter == 0 {
		fmt.Println("no git repository found in", basedir)
	} else {
		less(accumulator.String())
	}

	log.Printf("Processed %d repositor%s in %s\n",
		counter, ternary(counter > 1, "ies", "y"),
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
