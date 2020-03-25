package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gobuffalo/envy"

	"github.com/podded/ectoplasma/podgoo"
)

func main() {

	envy.Load()

	bouncerAddress := envy.Get("BOUNCER_ADDRESS", "http://localhost:13271")
	timeoutEnv := envy.Get("ECTO_SCRAPER_TIMEOUT", "10")
	descriptor := envy.Get("ECTO_SCRAPER_DESC", "ecto_ingest")

	numWorkersEnv := envy.Get("ECTO_SCRAPER_WORKERS", "1")

	timeout := 10
	i ,err := strconv.Atoi(timeoutEnv)
	if err != nil {
		timeout = i
	}

	numWorkers := 1
	i ,err = strconv.Atoi(numWorkersEnv)
	if err != nil {
		numWorkers = i
	}

	ctx := NewCtx()

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			goop := podgoo.NewPodGoo(bouncerAddress, time.Duration(timeout)*time.Second, descriptor)
			goop.StartScraper(ctx)
		}()
	}

	wg.Wait()


}

func NewCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		cancel()
	}()
	return ctx
}

