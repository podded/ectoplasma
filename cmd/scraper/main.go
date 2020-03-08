package main

import (
	"github.com/podded/ectoplasma/podgoo"
	"time"
)

func main() {

	// TODO Make all of this configurable!

	goop := podgoo.NewPodGoo("http://localhost:13270", 10*time.Second, "ecto_ingest")

	goop.StartScraper()

}
