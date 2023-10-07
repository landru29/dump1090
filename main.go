package main

import (
	"github.com/landru29/dump1090/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		panic(err)
	}
}
