package main

const APPNAME = "pullall"
const VERSION = "0.1.7"

type Result struct {
	pullSuccess  bool
	pullOutput   []byte
	statusOutput []byte
}
