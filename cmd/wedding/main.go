package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"weddingtemplates/internal/api"
	"weddingtemplates/internal/flow080"
	"weddingtemplates/internal/repository"
	"weddingtemplates/internal/service"
	"weddingtemplates/internal/storage"
)

func main() {
	path := flag.String("db", "wedding.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	s, e := storage.Open(*path)
	if e != nil {
		log.Fatal(e)
	}
	defer s.Close()
	svc := service.New(repository.New(s), flow080.DeterministicClock())
	log.Printf("wedding template service listening on %s", *addr)
	if e = http.ListenAndServe(*addr, api.New(svc).Handler()); e != nil && !os.IsTimeout(e) {
		log.Fatal(e)
	}
}
