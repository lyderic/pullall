package main

import(
  "sync"
)

const appname = "pullall"
const appversion = "0.0.9"

var termWidth int
var gitdirs []string
var wg sync.WaitGroup
var lock = sync.RWMutex{}
