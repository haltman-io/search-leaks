package output

import (
	"fmt"
	"io"
)

type PrinterConfig struct {
	NoColor bool
	Silent  bool
	Verbose bool
	Out     io.Writer
	Err     io.Writer
}

type Printer struct {
	noColor bool
	silent  bool
	verbose bool
	out     io.Writer
	err     io.Writer
}

func NewPrinter(cfg PrinterConfig) *Printer {
	return &Printer{
		noColor: cfg.NoColor,
		silent:  cfg.Silent,
		verbose: cfg.Verbose,
		out:     cfg.Out,
		err:     cfg.Err,
	}
}

func (p *Printer) Silent() bool  { return p.silent }
func (p *Printer) NoColor() bool { return p.noColor }

func (p *Printer) Printf(format string, args ...any) {
	fmt.Fprintf(p.out, format, args...)
}

func (p *Printer) Errorf(format string, args ...any) {
	fmt.Fprintf(p.err, format, args...)
}

func (p *Printer) Debugf(format string, args ...any) {
	if !p.verbose {
		return
	}
	fmt.Fprintf(p.err, Wrap(p.noColor, cDim, "[debug] ")+format, args...)
}

func (p *Printer) PrintHeader(target string, url string) {
	// Style similar to: [target] [url]
	left := Wrap(p.noColor, cCyan, target)
	right := Wrap(p.noColor, cBlue, url)
	fmt.Fprintf(p.out, "[%s] [%s]\n", left, right)
}

func (p *Printer) PrintKV(target string, brackets []string, key string, value string) {
	t := Wrap(p.noColor, cCyan, target)

	// bracket segments (e.g., stealer(1), stats(2))
	seg := ""
	for _, b := range brackets {
		seg += fmt.Sprintf(" [%s]", Wrap(p.noColor, cMag, b))
	}

	k := Wrap(p.noColor, cYellow, key)
	v := Wrap(p.noColor, cGreen, value)

	fmt.Fprintf(p.out, "[%s]%s [%s: %s]\n", t, seg, k, v)
}
