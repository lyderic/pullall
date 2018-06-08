package main

const APPNAME = "pullall"
const VERSION = "0.1.4"

type Result struct {
	pullSuccess  bool
	pullOutput   []byte
	statusOutput []byte
}
