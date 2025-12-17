package cli

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Targets     []string
	TargetLists []string

	ModeAutomatic bool
	ModeDomain    bool
	ModeEmail     bool

	Silent  bool
	Quiet   bool
	Verbose bool
	Debug   bool

	NoColor bool
}

type multiStringFlag []string

func (m *multiStringFlag) String() string { return strings.Join(*m, ",") }

// Accepts repeated usage and comma-separated values.
func (m *multiStringFlag) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			*m = append(*m, p)
		}
	}
	return nil
}

func ParseFlags(args []string) (Config, error) {
	var cfg Config
	var targets multiStringFlag
	var lists multiStringFlag

	fs := flag.NewFlagSet("search-leaks", flag.ContinueOnError)
	fs.SetOutput(nil) // avoid default printing

	fs.Var(&targets, "target", "Define the query value (repeatable, accepts comma-separated values)")
	fs.Var(&targets, "t", "Alias for --target")

	fs.Var(&lists, "target-list", "Define a file containing query values (repeatable, accepts comma-separated paths)")
	fs.Var(&lists, "tL", "Alias for --target-list")

	fs.BoolVar(&cfg.ModeAutomatic, "automatic", true, "Force automatic mode (default)")
	fs.BoolVar(&cfg.ModeAutomatic, "a", true, "Alias for --automatic")

	fs.BoolVar(&cfg.ModeDomain, "domain", false, "Force domain endpoint for all items")
	fs.BoolVar(&cfg.ModeDomain, "d", false, "Alias for --domain")

	fs.BoolVar(&cfg.ModeEmail, "email", false, "Force email endpoint for all items")
	fs.BoolVar(&cfg.ModeEmail, "e", false, "Alias for --email")

	fs.BoolVar(&cfg.Silent, "silent", false, "Silent mode (results only)")
	fs.BoolVar(&cfg.Silent, "s", false, "Alias for --silent")
	fs.BoolVar(&cfg.Quiet, "quiet", false, "Quiet mode (results only)")
	fs.BoolVar(&cfg.Quiet, "q", false, "Alias for --quiet")

	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose/debug logs")
	fs.BoolVar(&cfg.Verbose, "v", false, "Alias for --verbose")
	fs.BoolVar(&cfg.Debug, "debug", false, "Enable verbose/debug logs")

	fs.BoolVar(&cfg.NoColor, "no-color", false, "Disable ANSI colors")
	fs.BoolVar(&cfg.NoColor, "nc", false, "Alias for --no-color")

	// Custom usage
	fs.Usage = func() { PrintUsage() }

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			PrintUsage()
			return cfg, fmt.Errorf("help requested")
		}
		PrintUsage()
		return cfg, err
	}

	cfg.Targets = []string(targets)
	cfg.TargetLists = []string(lists)

	// If user explicitly sets --domain/--email, automatic must be considered disabled.
	// But we still enforce only one mode below in ResolveMode().
	return cfg, nil
}

type Mode string

const (
	ModeAutomatic Mode = "automatic"
	ModeDomain    Mode = "domain"
	ModeEmail     Mode = "email"
)

func ResolveMode(auto bool, domain bool, email bool) (Mode, error) {
	// If user set domain or email, automatic is effectively off.
	if domain || email {
		auto = false
	}

	count := 0
	if auto {
		count++
	}
	if domain {
		count++
	}
	if email {
		count++
	}

	if count == 0 {
		// default to automatic
		return ModeAutomatic, nil
	}
	if count > 1 {
		return "", fmt.Errorf("only one mode is allowed: use exactly one of --automatic/-a, --domain/-d, --email/-e")
	}

	if domain {
		return ModeDomain, nil
	}
	if email {
		return ModeEmail, nil
	}
	return ModeAutomatic, nil
}
