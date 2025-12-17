package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/haltman-io/search-leaks/internal/api"
	"github.com/haltman-io/search-leaks/internal/cli"
	"github.com/haltman-io/search-leaks/internal/output"
	"github.com/haltman-io/search-leaks/internal/ratelimit"
	"github.com/haltman-io/search-leaks/internal/targets"
	"github.com/haltman-io/search-leaks/internal/util"
)

const (
	ToolName    = "search-leaks"
	ToolVersion = "v1.0.1-stable"
)

func main() {
	cfg, err := cli.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, output.ColorizeError(cfg.NoColor, err.Error()))
		os.Exit(2)
	}

	pr := output.NewPrinter(output.PrinterConfig{
		NoColor: cfg.NoColor,
		Silent:  cfg.Silent || cfg.Quiet,
		Verbose: cfg.Verbose || cfg.Debug,
		Out:     os.Stdout,
		Err:     os.Stderr,
	})

	if !pr.Silent() {
		output.PrintBanner(pr, ToolName, ToolVersion)
	}

	// Collect targets from: stdin + -t + -tL
	collected, err := targets.CollectTargets(targets.CollectConfig{
		Targets:      cfg.Targets,
		TargetLists:  cfg.TargetLists,
		ReadStdin:    util.HasStdinData(),
		Stdin:        os.Stdin,
		TrimSpaces:   true,
		SkipEmpty:    true,
		Dedupe:       true,
		VerboseLogFn: pr.Debugf,
	})
	if err != nil {
		pr.Errorf("%v\n", err)
		os.Exit(2)
	}

	if len(collected) == 0 {
		pr.Errorf("no targets provided (use stdin, -t/--target, or -tL/--target-list)\n")
		os.Exit(2)
	}

	mode, err := cli.ResolveMode(cfg.ModeAutomatic, cfg.ModeDomain, cfg.ModeEmail)
	if err != nil {
		pr.Errorf("%v\n", err)
		os.Exit(2)
	}
	pr.Debugf("mode=%s targets=%d\n", mode, len(collected))

	lim := ratelimit.NewTickerLimiter(200 * time.Millisecond) // 50 req / 10s => 1 req / 200ms
	defer lim.Stop()

	httpClient := api.NewHTTPClient(15 * time.Second)
	client := api.NewClient(httpClient)

	consecutiveErrors := 0

	ctx := context.Background()

	for _, t := range collected {
		reqPlan, err := targets.BuildRequestPlan(targets.PlanConfig{
			RawTarget: t,
			Mode:      mode,
		})
		if err != nil {
			pr.Errorf("[%s] %v\n", t, err)
			consecutiveErrors++
			if consecutiveErrors >= 3 {
				pr.Errorf("API returned errors 3 times in a row (aborting)\n")
				os.Exit(1)
			}
			continue
		}

		// A plan may expand a single input into multiple requests (e.g., domain -> default emails in --email mode).
		for _, r := range reqPlan.Requests {
			lim.Wait()

			pr.Debugf("request target=%s endpoint=%s url=%s\n", r.OriginalTarget, r.Endpoint, r.URL)

			resp, status, err := client.GetJSON(ctx, r.URL)
			if err != nil {
				pr.Errorf("[%s] request failed: %v\n", r.OriginalTarget, err)
				consecutiveErrors++
				if consecutiveErrors >= 3 {
					pr.Errorf("API returned errors 3 times in a row (aborting)\n")
					os.Exit(1)
				}
				continue
			}

			if status < 200 || status > 299 {
				pr.Errorf("[%s] API error (status=%d)\n", r.OriginalTarget, status)
				consecutiveErrors++
				if consecutiveErrors >= 3 {
					pr.Errorf("API returned errors 3 times in a row (aborting)\n")
					os.Exit(1)
				}
				continue
			}

			// Success resets consecutive error counter.
			consecutiveErrors = 0

			// Header line
			pr.PrintHeader(r.OriginalTarget, r.URL)

			// Flatten JSON to lines
			lines := output.FlattenJSON(resp)

			// Print lines in requested style
			for _, line := range lines {
				pr.PrintKV(r.OriginalTarget, line.Brackets, line.Key, line.Value)
			}
		}
	}
}
