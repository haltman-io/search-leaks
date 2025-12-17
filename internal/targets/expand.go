package targets

import (
	"fmt"
	"net/url"
	"strings"

	"search-leaks/internal/api"
	"search-leaks/internal/cli"
)

type Request struct {
	OriginalTarget string // used in the output prefix
	Endpoint       string // "domain" or "email"
	URL            string // full request URL
}

type RequestPlan struct {
	Requests []Request
}

type PlanConfig struct {
	RawTarget string
	Mode      cli.Mode
}

// Common mailbox aliases (often present, not guaranteed).
var DefaultMailboxUsernames = []string{
	"postmaster",
	"abuse",
	"hostmaster",
	"webmaster",
	"admin",
	"administrator",
	"root",
}

func BuildRequestPlan(cfg PlanConfig) (RequestPlan, error) {
	t := strings.TrimSpace(cfg.RawTarget)
	if t == "" {
		return RequestPlan{}, fmt.Errorf("empty target")
	}

	switch cfg.Mode {
	case cli.ModeAutomatic:
		if IsEmail(t) {
			return RequestPlan{Requests: []Request{makeEmailReq(t)}}, nil
		}
		return RequestPlan{Requests: []Request{makeDomainReq(t)}}, nil

	case cli.ModeDomain:
		// If it's an email, extract domain and query domain endpoint.
		if IsEmail(t) {
			d, ok := ExtractDomainFromEmail(t)
			if !ok {
				return RequestPlan{}, fmt.Errorf("invalid email, cannot extract domain: %s", t)
			}
			return RequestPlan{Requests: []Request{makeDomainReqWithOriginal(t, d)}}, nil
		}
		return RequestPlan{Requests: []Request{makeDomainReq(t)}}, nil

	case cli.ModeEmail:
		// If it's a domain, expand into default mailbox aliases.
		if !IsEmail(t) {
			reqs := make([]Request, 0, len(DefaultMailboxUsernames))
			for _, u := range DefaultMailboxUsernames {
				email := fmt.Sprintf("%s@%s", u, t)
				reqs = append(reqs, makeEmailReqWithOriginal(t, email))
			}
			return RequestPlan{Requests: reqs}, nil
		}
		return RequestPlan{Requests: []Request{makeEmailReq(t)}}, nil

	default:
		return RequestPlan{}, fmt.Errorf("unknown mode: %s", cfg.Mode)
	}
}

func makeDomainReq(domain string) Request {
	u := api.SearchByDomain + url.QueryEscape(domain)
	return Request{
		OriginalTarget: domain,
		Endpoint:       "domain",
		URL:            u,
	}
}

func makeDomainReqWithOriginal(original string, domain string) Request {
	u := api.SearchByDomain + url.QueryEscape(domain)
	return Request{
		OriginalTarget: original,
		Endpoint:       "domain",
		URL:            u,
	}
}

func makeEmailReq(email string) Request {
	u := api.SearchByEmail + url.QueryEscape(email)
	return Request{
		OriginalTarget: email,
		Endpoint:       "email",
		URL:            u,
	}
}

func makeEmailReqWithOriginal(original string, email string) Request {
	u := api.SearchByEmail + url.QueryEscape(email)
	return Request{
		OriginalTarget: original,
		Endpoint:       "email",
		URL:            u,
	}
}
