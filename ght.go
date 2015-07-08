package main

import (
	//"ght/cmd"
	"ght/mod"
)

func main() {
	paths := []string{"/data/go/src/ght/mod"}
	mod.Watch(paths)
}
