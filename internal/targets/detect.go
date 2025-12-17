package targets

import (
	"net/mail"
	"strings"
)

func IsEmail(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	_, err := mail.ParseAddress(s)
	return err == nil && strings.Contains(s, "@")
}

func ExtractDomainFromEmail(email string) (string, bool) {
	email = strings.TrimSpace(email)
	at := strings.LastIndex(email, "@")
	if at <= 0 || at >= len(email)-1 {
		return "", false
	}
	d := strings.TrimSpace(email[at+1:])
	if d == "" {
		return "", false
	}
	return d, true
}
