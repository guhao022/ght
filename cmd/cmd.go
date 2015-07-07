package cmd

import (
	"fmt"
	"os"
	"ght/mod"
)

func help(helpcode string) {
	switch helpcode {
	case "help":
		fmt.Println(Helptext)
	case "run":
		fmt.Println(Runhelp)
	}
}
func Anget() {
	var ang mod.Anget

	commands := os.Args
	if len(commands) < 3 {
		help("help")
		os.Exit(0)
	}
	//anget.dev = true
	switch commands[1] {
	case "help":
		fmt.Println(Helptext)
	case "run":
		switch commands[2] {
		case "--help", "-h":
			help("help")
			os.Exit(0)
		default:
			ang.Server = commands[2]
		}
		ang.Run()
	default:
		help("help")
	}
}
