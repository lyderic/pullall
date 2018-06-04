package main

import (
	"sync"
)

const appname = "pullall"
const appversion = "0.1.1"

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
