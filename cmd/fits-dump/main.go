package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

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

	var types bool
	flag.BoolVar(&types, "types", false, "export observation types")

	var methods bool
	flag.BoolVar(&methods, "methods", false, "export observation methods")

	var sites bool
	flag.BoolVar(&sites, "sites", false, "export observation sites")

	var service string
	flag.StringVar(&service, "service", "fits.geonet.org.nz", "fits api end-point")

	var output string
	flag.StringVar(&output, "output", "", "output file")

	var typeId string
	flag.StringVar(&typeId, "typeId", "", "observation typeId")

	var methodId string
	flag.StringVar(&methodId, "methodId", "", "observation methodId")

	var siteId string
	flag.StringVar(&siteId, "siteId", "", "observation siteId")

	var days int
	flag.IntVar(&days, "days", 0, "observation number of days")

	flag.Parse()

	if days < 0 {
		days = 0
	}

	if days > 365000 {
		days = 365000
	}

	ctx := context.Background()

	client := fits.New(service, fits.Days(days))
	switch {
	case types:
		t, err := client.Types(ctx)
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case output != "":
			if fits.Types(t).WriteFile(output); err != nil {
				log.Fatal(err)
			}
		default:
			if fits.Types(t).Write(os.Stdout); err != nil {
				log.Fatal(err)
			}
		}
	case methods:
		m, err := client.Methods(ctx)
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case output != "":
			if fits.Methods(m).WriteFile(output); err != nil {
				log.Fatal(err)
			}
		default:
			if fits.Methods(m).Write(os.Stdout); err != nil {
				log.Fatal(err)
			}
		}
	case sites:
		s, err := client.Sites(ctx)
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case output != "":
			if fits.Sites(s).WriteFile(output); err != nil {
				log.Fatal(err)
			}
		default:
			if fits.Sites(s).Write(os.Stdout); err != nil {
				log.Fatal(err)
			}
		}
	case methodId == "":
		log.Fatalf("observations require a methodId option to be set")
	case typeId == "":
		log.Fatalf("observations require a typeId option to be set")
	case siteId == "":
		log.Fatalf("observations require a siteId option to be set")
	default:
		o, err := client.Observations(ctx, methodId, typeId, siteId)
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case output != "":
			if fits.Observations(o).WriteFile(output); err != nil {
				log.Fatal(err)
			}
		default:
			if fits.Observations(o).Write(os.Stdout); err != nil {
				log.Fatal(err)
			}
		}
	}
}
