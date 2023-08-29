package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ozym/fits/internal/fits"
)

var Version = "fits-dump"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Extract FITS data suitable for archiving\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Version: %s\n", Version)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var dir string
	flag.StringVar(&dir, "dir", ".", "output directory")

	var service string
	flag.StringVar(&service, "service", "fits.geonet.org.nz", "fits api end-point")

	flag.Parse()

	ctx := context.Background()

	client := fits.New(service)

	path := filepath.Join(dir, "meta", "sites.csv")
	log.Printf("building sites: %q", path)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		log.Fatalf("fatal: unable to build dir %q: %v", filepath.Dir(path), err)
	}
	sites, err := client.Sites(ctx)
	if err != nil {
		log.Fatalf("fatal: unable to recover sites: %v", err)
	}
	if err := fits.Sites(sites).WriteFile(path); err != nil {
		log.Fatalf("fatal: unable to write sites %q: %v", path, err)
	}

	path = filepath.Join(dir, "meta", "methods.csv")
	log.Printf("building methods: %q", path)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		log.Fatalf("fatal: unable to build dir %q: %v", filepath.Dir(path), err)
	}
	methods, err := client.Methods(ctx)
	if err != nil {
		log.Fatalf("fatal: unable to recover methods: %v", err)
	}
	if err := fits.Methods(methods).WriteFile(path); err != nil {
		log.Fatalf("fatal: unable to write methods %q: %v", path, err)
	}

	path = filepath.Join(dir, "meta", "types.csv")
	log.Printf("building types: %q", path)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		log.Fatalf("fatal: unable to build dir %q: %v", filepath.Dir(path), err)
	}
	types, err := client.Types(ctx)
	if err != nil {
		log.Fatalf("fatal: unable to recover types: %v", err)
	}
	if err := fits.Types(types).WriteFile(path); err != nil {
		log.Fatalf("fatal: unable to write types %q: %v", path, err)
	}

	log.Println("building observations")

	for _, t := range types {
		if strings.Contains(t.TypeId, "_") {
			continue
		}
		log.Printf("\t- %s (%q)", t.TypeId, t.Name)

		methods, err := client.Methods(ctx, t.TypeId)
		if err != nil {
			log.Fatal(err)
		}

		for _, m := range methods {
			sites, err := client.Sites(ctx, t.TypeId, m.MethodId)
			if err != nil {
				log.Fatal(err)
			}

			for _, s := range sites {
				path = filepath.Join(dir, "data", m.MethodId, t.TypeId, s.SiteId+".csv")
				log.Printf("\t\t building obs: %q\n", path)
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					log.Fatalf("fatal: unable to build dir %q: %v", filepath.Dir(path), err)
				}
				observations, err := client.Observations(ctx, m.MethodId, t.TypeId, s.SiteId)
				if err != nil {
					log.Fatalf("fatal: unable to recover observations: %v", err)
				}
				if err := fits.Observations(observations).WriteFile(path); err != nil {
					log.Fatalf("fatal: unable to write types %q: %v", path, err)
				}
			}
		}
	}

	log.Println("done")
}
