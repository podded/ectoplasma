package main

import (
	"time"

	"github.com/podded/ectoplasma/podgoo"
)

func main() {
	// Just start up the minimal server to put things onto the ingest queue
	// TODO Make all of this configurable!
	goop := podgoo.NewPodGoo("http://localhost:13271", 10*time.Second, "ecto_server")
	goop.BoundHost = "0.0.0.0"
	goop.BoundPort = 13271

	goop.ListenAndServe()
}
