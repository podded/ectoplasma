package main

import (
	"time"

	"github.com/podded/ectoplasma/podgoo"
)

func main() {
	// TODO Make all of this configurable!
	goop := podgoo.NewPodGoo("http://localhost:13271", 10*time.Second, "ecto_ingest")
	goop.StartScraper()
}
