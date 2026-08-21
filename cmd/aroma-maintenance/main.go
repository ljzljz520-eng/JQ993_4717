package main

import (
	"aroma-maintenance/internal/api"
	"aroma-maintenance/internal/catalog"
	"aroma-maintenance/internal/importer"
	"aroma-maintenance/internal/report"
	"aroma-maintenance/internal/review"
	"aroma-maintenance/internal/search"
	"aroma-maintenance/internal/store"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	path := flag.String("db", "aroma-maintenance.db", "bbolt database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	if err := run(*path, *addr); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
func run(path, addr string) error {
	s, err := store.Open(path)
	if err != nil {
		return err
	}
	defer s.Close()
	h := api.NewHandler(catalog.New(s), review.New(s), search.New(s), importer.New(), report.New(s))
	return http.ListenAndServe(addr, h)
}
