package cli

import "fmt"

func PrintUsage() {
	fmt.Println(`
search-leaks - Hudson Rock (cavalier) leak statistics checker (domain/email)

Usage:
  search-leaks [flags]
  cat targets.txt | search-leaks [flags]

Targets input:
  --target, -t <target>               Define the query value (repeatable; supports comma-separated)
  --target-list, -tL <file>           Define a file with query values (repeatable; supports comma-separated)

Modes (only one allowed; default is automatic):
  --automatic, -a                     Auto-detect endpoint per target (default)
  --domain, -d                        Force domain endpoint (emails will be converted to their domain)
  --email, -e                         Force email endpoint (domains will expand to common mailbox aliases)

Output:
  --silent, -s                        Results only (no banner)
  --quiet, -q                         Results only (no banner)
  --verbose, -v                       Debug logs to stderr
  --debug                             Debug logs to stderr
  --no-color, -nc                     Disable ANSI colors

Examples:
  search-leaks -t google.com -t twitter.com
  search-leaks -t google.com,twitter.com
  search-leaks -tL targets1.txt -tL targets2.txt
  search-leaks -tL targets1.txt,targets2.txt
  cat targets.txt | search-leaks -a
  cat targets.txt | search-leaks -d
  search-leaks -e -t example.com
`)
}
