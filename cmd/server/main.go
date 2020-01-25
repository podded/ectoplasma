package main

import "github.com/podded/ectoplasma"

func main() {
	// Just start up the minimal server to put things onto the ingest queue

	// TODO Make all of this configurable!

	goop := ectoplasma.PodGoo{
		BoundHost: "0.0.0.0",
		BoundPort: 13270,
	}

	goop.ListenAndServe()

}
