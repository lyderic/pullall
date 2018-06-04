package main

import (
	"sync"
)

const APPNAME = "pullall"
const VERSION = "0.1.3"

type Result struct {
	pullSuccess  bool
	pullOutput   []byte
	statusOutput []byte
}

var (
	termWidth int
	gitdirs   []string
	wg        sync.WaitGroup
	lock      = sync.RWMutex{}
)
