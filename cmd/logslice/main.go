package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
)

const timeLayout = "2006-01-02T15:04:05"

func main() {
	var (
		start   = flag.String("start", "", "Start of time range (RFC3339, e.g. 2024-01-01T00:00:00)")
		end     = flag.String("end", "", "End of time range (RFC3339, e.g. 2024-01-01T23:59:59)")
		pattern = flag.String("pattern", "", "Optional substring or regex pattern to match")
		numbers = flag.Bool("n", false, "Prefix output lines with line numbers")
		jsonOut = flag.Bool("json", false, "Output results as JSON objects")
		summary = flag.Bool("summary", false, "Print a summary after results")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: logslice [options] <logfile>\n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	filepath := flag.Arg(0)

	var startTime, endTime time.Time
	var err error
	if *start != "" {
		startTime, err = time.Parse(timeLayout, *start)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid --start: %v\n", err)
			os.Exit(1)
		}
	}
	if *end != "" {
		endTime, err = time.Parse(timeLayout, *end)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid --end: %v\n", err)
			os.Exit(1)
		}
	}

	format := output.FormatPlain
	if *numbers {
		format = output.FormatNumbered
	}
	if *jsonOut {
		format = output.FormatJSON
	}

	fmt, err := output.NewFormatter(os.Stdout, format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "formatter error: %v\n", err)
		os.Exit(1)
	}

	matches, err := filter.FilterFile(filepath, startTime, endTime, *pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error filtering file: %v\n", err)
		os.Exit(1)
	}

	for _, line := range matches {
		if werr := fmt.Write(line); werr != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", werr)
			os.Exit(1)
		}
	}

	if *summary {
		fmt.WriteSummary(len(matches))
	}
}
