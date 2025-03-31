package config

import (
	"strings"

	"gabe565.com/cloudflare-ddns/internal/output"
	"github.com/fatih/color"
	"github.com/miekg/dns"
)

//go:generate go tool enumer -type Source -trimprefix Source -transform snake -linecomment -output source_string.go

type Source uint8

const (
	SourceCloudflareTLS Source = iota
	SourceCloudflare
	SourceOpenDNSTLS // opendns_tls
	SourceOpenDNS    // opendns
	SourceICanHazIP  // icanhazip
	SourceIPInfo     // ipinfo
	SourceIPify      // ipify
)

func (s Source) Description(format output.Format) string {
	req := s.Request()
	switch req := req.(type) {
	case DNSv4v6, HTTPv4v6:
		return req.Description(format)
	default:
		panic("invalid request type")
	}
}

type Requestv4v6 interface {
	Description(format output.Format) string
}

type HTTPv4v6 struct {
	URLv4, URLv6 string
}

func (d HTTPv4v6) Description(format output.Format) string {
	switch format {
	case output.FormatANSI:
		bold := color.New(color.Bold).Sprint
		return "Makes HTTPS requests to " + bold(d.URLv4) + " and " + bold(d.URLv6) + "."
	case output.FormatMarkdown:
		return "Makes HTTPS requests to `" + d.URLv4 + "` and `" + d.URLv6 + "`."
	default:
		panic("unimplemented format: " + format)
	}
}

type DNSv4v6 struct {
	ServerV4   string
	QuestionV4 dns.Question
	ServerV6   string
	QuestionV6 dns.Question
	TLS        bool
}

func (d DNSv4v6) Description(format output.Format) string {
	proto := "DNS"
	if d.TLS {
		proto = "DNS-over-TLS"
	}
	switch format {
	case output.FormatANSI:
		bold := color.New(color.Bold).Sprint
		return "Queries " + bold(
			strings.TrimSuffix(d.QuestionV4.Name, "."),
		) + " using " + proto + " via " + bold(
			d.ServerV4,
		) + "."
	case output.FormatMarkdown:
		return "Queries `" + strings.TrimSuffix(
			d.QuestionV4.Name,
			".",
		) + "` using " + proto + " via `" + d.ServerV4 + "`."
	default:
		panic("unimplemented format: " + format)
	}
}

func (s Source) Request() Requestv4v6 { //nolint:ireturn
	var server string
	var tls bool
	switch s {
	case SourceCloudflareTLS:
		server = "one.one.one.one:853"
		tls = true
		fallthrough
	case SourceCloudflare:
		if server == "" {
			server = "one.one.one.one:53"
		}
		question := dns.Question{
			Name:   "whoami.cloudflare.",
			Qtype:  dns.TypeTXT,
			Qclass: dns.ClassCHAOS,
		}
		return DNSv4v6{
			ServerV4:   server,
			QuestionV4: question,
			ServerV6:   server,
			QuestionV6: question,
			TLS:        tls,
		}
	case SourceOpenDNSTLS:
		server = "dns.opendns.com:853"
		tls = true
		fallthrough
	case SourceOpenDNS:
		if server == "" {
			server = "dns.opendns.com:53"
		}
		return DNSv4v6{
			ServerV4: server,
			QuestionV4: dns.Question{
				Name:   "myip.opendns.com.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			},
			ServerV6: server,
			QuestionV6: dns.Question{
				Name:   "myip.opendns.com.",
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassINET,
			},
			TLS: tls,
		}
	case SourceICanHazIP:
		return HTTPv4v6{
			URLv4: "https://ipv4.icanhazip.com",
			URLv6: "https://ipv6.icanhazip.com",
		}
	case SourceIPInfo:
		return HTTPv4v6{
			URLv4: "https://ipinfo.io/ip",
			URLv6: "https://v6.ipinfo.io/ip",
		}
	case SourceIPify:
		return HTTPv4v6{
			URLv4: "https://api.ipify.org",
			URLv6: "https://api6.ipify.org",
		}
	default:
		panic("source request unimplemented: " + s.String())
	}
}
