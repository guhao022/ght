package main

import (
	//"ght/cmd"
	"ght/mod"
)

func main() {
	paths := []string{"/data/go/src/ght/mod"}
	mod.Watch(paths)
	go mod.Build()
	done := make(chan bool)
	<-done
}
