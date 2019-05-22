package main

import "fmt"

type endpoint struct {
	File         string
	Pattern      string
	StripPattern bool
}

func (e endpoint) String() string {
	return fmt.Sprintf("file: %s, pattern: %s, strip pattern: %v", e.File, e.Pattern, e.StripPattern)
}

type cfg struct {
	Endpoints []endpoint
}
