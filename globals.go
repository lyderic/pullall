package main

import (
	"sync"
)

const appname = "pullall"
const appversion = "0.1.0"

var termWidth int
var gitdirs []string
var wg sync.WaitGroup
var lock = sync.RWMutex{}
