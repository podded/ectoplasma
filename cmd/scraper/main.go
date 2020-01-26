package main

import (
	"fmt"
	"github.com/podded/ectoplasma"
	"os"
)

func main() {

	// TODO Make all of this configurable!

	goop := ectoplasma.PodGoo{
	}

	err := goop.ProcessIngestQueue()

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
