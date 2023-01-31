package main

import (
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":5000",
		Handler:      RunServer(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	var g errgroup.Group

	g.Go(func() error {
		return server.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
