package main

import(
  "fmt"
  "log"
  "os"
  "os/exec"
  "path/filepath"
)

func init() {
  log.SetFlags(log.Lshortfile)
}

func sanitize(input string) string {
  // first, we make sure we work with absolute paths
  absoluteInput,err := filepath.Abs(input)
  if err != nil { log.Fatal(err) }
  resolvedsymlink,err := filepath.EvalSymlinks(absoluteInput)
  if err != nil { log.Fatal(err) }
  basedir,err := filepath.Abs(resolvedsymlink)
  if err != nil { log.Fatal(err) }
  basedirinfo,err := os.Stat(input)
  if err != nil { log.Fatal(err) }
  if !basedirinfo.IsDir() { log.Fatalf("Not a dir!") }
  return basedir
}

func printRed(message string) {
  fmt.Printf("\033[31m%s\033[0m\n", message)
}

func getTermWidth() int {
  _,w := getTermDim()
  return w
}

func getTermDim() (int, int) {
  cmd := exec.Command("stty", "size")
  cmd.Stdin = os.Stdin
  termDim,err := cmd.Output()
  if err != nil { log.Fatal(err) }
  var h, w int
  fmt.Sscan(string(termDim), &h, &w)
  return h, w
}

func wipeLine() {
  fmt.Print("\r")
  for i := 0 ; i < termWidth ; i++ {
    fmt.Print(" ")
  }
  fmt.Print("\r")
}

func checkBinaries(binaries ...string) {
  for _,binary := range(binaries) {
    _,e := exec.LookPath(binary)
    if e != nil {
      log.Fatalf("%s executable not found in path! Aborting...", binary)
    }
  }
}
