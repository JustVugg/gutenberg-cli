package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gutenberg.local/dev-pulse/internal/aggr"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	args := os.Args[2:]
	switch os.Args[1] {
	case "sources":
		runSources(args)
	case "fan-out", "search":
		runFanOut(args)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println("Dev Pulse aggregator - aggregator over 2 source(s)")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  sources             List configured sources")
	fmt.Println("  fan-out [--param k=v ...] [--parallel N] [--json]   Call all sources in parallel")
	fmt.Println("  search [...]         Alias for fan-out")
}

func runSources(_ []string) {
	data, _ := json.MarshalIndent(aggr.Sources, "", "  ")
	fmt.Println(string(data))
}

func runFanOut(args []string) {
	params := map[string]string{}
	parallel := 4
	asJSON := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--param", "-p":
			i++
			if i >= len(args) {
				continue
			}
			equals := strings.Index(args[i], "=")
			if equals == -1 {
				params[args[i]] = "true"
			} else {
				params[args[i][:equals]] = args[i][equals+1:]
			}
		case "--parallel":
			i++
			if i < len(args) {
				if n, err := strconv.Atoi(args[i]); err == nil {
					parallel = n
				}
			}
		case "--json":
			asJSON = true
		}
	}
	items := aggr.FanOut(context.Background(), params, parallel)
	if asJSON {
		data, _ := json.MarshalIndent(items, "", "  ")
		fmt.Println(string(data))
		return
	}
	for _, item := range items {
		if item.Error != "" {
			fmt.Printf("- %s ERROR: %s\n", item.Source, item.Error)
			continue
		}
		fmt.Printf("- %s status=%d url=%s\n", item.Source, item.Status, item.URL)
	}
}
