package targets

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type CollectConfig struct {
	Targets     []string
	TargetLists []string

	ReadStdin  bool
	Stdin      *os.File
	TrimSpaces bool
	SkipEmpty  bool
	Dedupe     bool

	VerboseLogFn func(format string, args ...any)
}

func CollectTargets(cfg CollectConfig) ([]string, error) {
	var out []string
	seen := map[string]bool{}

	add := func(v string) {
		if cfg.TrimSpaces {
			v = strings.TrimSpace(v)
		}
		if cfg.SkipEmpty && v == "" {
			return
		}
		if cfg.Dedupe {
			key := strings.ToLower(v)
			if seen[key] {
				return
			}
			seen[key] = true
		}
		out = append(out, v)
	}

	// stdin (pipeline)
	if cfg.ReadStdin && cfg.Stdin != nil {
		if cfg.VerboseLogFn != nil {
			cfg.VerboseLogFn("reading targets from stdin\n")
		}
		sc := bufio.NewScanner(cfg.Stdin)
		for sc.Scan() {
			line := sc.Text()
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Support comma-separated even on stdin lines
			parts := strings.Split(line, ",")
			for _, p := range parts {
				add(strings.TrimSpace(p))
			}
		}
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("failed reading stdin: %w", err)
		}
	}

	// -t / --target
	for _, t := range cfg.Targets {
		parts := strings.Split(t, ",")
		for _, p := range parts {
			add(strings.TrimSpace(p))
		}
	}

	// -tL / --target-list
	for _, pathItem := range cfg.TargetLists {
		paths := strings.Split(pathItem, ",")
		for _, p := range paths {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if cfg.VerboseLogFn != nil {
				cfg.VerboseLogFn("reading targets from file=%s\n", p)
			}
			f, err := os.Open(p)
			if err != nil {
				return nil, fmt.Errorf("failed to open target list file '%s': %w", p, err)
			}
			if err := readTargetsFromReader(f, add); err != nil {
				f.Close()
				return nil, fmt.Errorf("failed reading target list file '%s': %w", p, err)
			}
			f.Close()
		}
	}

	return out, nil
}

func readTargetsFromReader(r io.Reader, add func(string)) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				add(p)
			}
		}
	}
	return sc.Err()
}
