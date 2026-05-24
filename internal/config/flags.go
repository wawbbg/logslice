package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/logslice/internal/parser"
)

// ParseFlags parses os.Args and returns a populated Config.
// On error or --help it writes to stderr and exits.
func ParseFlags(args []string) (*Config, error) {
	cfg := DefaultConfig()

	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	var fromStr, toStr string

	fs.StringVar(&fromStr, "from", "", "start of time range (e.g. 2024-01-02T15:04:05)")
	fs.StringVar(&toStr, "to", "", "end of time range (e.g. 2024-01-02T16:04:05)")
	fs.StringVar(&cfg.Pattern, "pattern", "", "filter lines containing this substring or pattern")
	fs.StringVar(&cfg.Format, "format", cfg.Format, "output format: plain, json, numbered")
	fs.StringVar(&cfg.OutputFile, "output", "", "write results to file instead of stdout")
	fs.BoolVar(&cfg.ShowSummary, "summary", false, "print a summary line after output")
	fs.IntVar(&cfg.MaxLines, "max-lines", cfg.MaxLines, "maximum number of lines to output (0 = unlimited)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Fprintln(os.Stderr, "usage: logslice [flags] <logfile>")
		fs.PrintDefaults()
		return nil, fmt.Errorf("no log file specified")
	}
	cfg.FilePath = remaining[0]

	if fromStr != "" {
		t, err := parser.ParseTimestamp(fromStr)
		if err != nil {
			return nil, fmt.Errorf("invalid --from value: %w", err)
		}
		cfg.From = t
	}

	if toStr != "" {
		t, err := parser.ParseTimestamp(toStr)
		if err != nil {
			return nil, fmt.Errorf("invalid --to value: %w", err)
		}
		cfg.To = t
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
