package main

import (
	"strconv"
	"time"

	"github.com/gobuffalo/envy"

	"github.com/podded/ectoplasma/podgoo"
)

func main() {

	envy.Load()

	bouncerAddress := envy.Get("BOUNCER_ADDRESS", "http://localhost:13271")
	timeoutEnv := envy.Get("ECTO_SCRAPER_TIMEOUT", "10")
	descriptor := envy.Get("ECTO_SCRAPER_DESC", "ecto_ingest")

	timeout := 10
	i, err := strconv.Atoi(timeoutEnv)
	if err != nil {
		timeout = i
	}

	goop := podgoo.NewPodGoo(bouncerAddress, time.Duration(timeout)*time.Second, descriptor)
	goop.StartScraper()

}
