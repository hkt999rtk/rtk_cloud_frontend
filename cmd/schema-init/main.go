package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"realtek-connect/internal/analytics"
	"realtek-connect/internal/leads"
	"realtek-connect/internal/search"

	_ "modernc.org/sqlite"
)

func main() {
	analyticsPath := flag.String("analytics", "", "analytics SQLite database path")
	searchPath := flag.String("search", "", "search SQLite database path")
	leadsPath := flag.String("leads", "", "leads SQLite database path")
	flag.Parse()
	for name, path := range map[string]string{"analytics": *analyticsPath, "search": *searchPath, "leads": *leadsPath} {
		if path == "" {
			log.Fatalf("--%s is required", name)
		}
	}
	open := func(path string) *sql.DB {
		db, err := sql.Open("sqlite", path)
		if err != nil {
			log.Fatal(err)
		}
		return db
	}
	a := open(*analyticsPath)
	defer a.Close()
	if err := analytics.NewRepository(a, 30).Init(); err != nil {
		log.Fatal(err)
	}
	s := open(*searchPath)
	defer s.Close()
	if err := search.NewRepository(s).Init(context.Background()); err != nil {
		log.Fatal(err)
	}
	l := open(*leadsPath)
	defer l.Close()
	if err := leads.NewRepository(l).Init(); err != nil {
		log.Fatal(err)
	}
}
