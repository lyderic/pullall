package main

import "fmt"

type Result struct {
	pullSuccess  bool
	pullOutput   []byte
	statusOutput []byte
}

func (r Result) String() string {
	return fmt.Sprintf(" PullSuccess: %t\n pullOutput: %s\n statusOutput: %s",
		r.pullSuccess,
		string(r.pullOutput),
		string(r.statusOutput),
	)
}
