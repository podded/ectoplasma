package main

import (
	"github.com/podded/ectoplasma/podgoo"
	"time"
)

func main() {
	// Just start up the minimal server to put things onto the ingest queue

	// TODO Make all of this configurable!

	//goop := podgoo.PodGoo{
	//	BoundHost: "0.0.0.0",
	//	BoundPort: 13270,
	//}

	goop := podgoo.NewPodGoo("http://localhost:13270", 10*time.Second, "ecto_serv")
	goop.BoundHost = "0.0.0.0"
	goop.BoundPort = 13271

	goop.ListenAndServe()

}
