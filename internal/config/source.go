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
	V4URL, V6URL string
}

func (d HTTPv4v6) Description(format output.Format) string {
	switch format {
	case output.FormatANSI:
		bold := color.New(color.Bold).Sprint
		return "Makes HTTPS requests to " + bold(d.V4URL) + " and " + bold(d.V6URL) + "."
	case output.FormatMarkdown:
		return "Makes HTTPS requests to `" + d.V4URL + "` and `" + d.V6URL + "`."
	default:
		panic("unimplemented format: " + format)
	}
}

type DNSv4v6 struct {
	Server                 string
	TLS                    bool
	V4Question, V6Question dns.Question
}

func (d DNSv4v6) Description(format output.Format) string {
	proto := "DNS"
	if d.TLS {
		proto = "DNS-over-TLS"
	}
	switch format {
	case output.FormatANSI:
		bold := color.New(color.Bold).Sprint
		return "Queries " + bold(strings.TrimSuffix(d.V4Question.Name, ".")) + " using " + proto + " via " + bold(d.Server) + "."
	case output.FormatMarkdown:
		return "Queries `" + strings.TrimSuffix(d.V4Question.Name, ".") + "` using " + proto + " via `" + d.Server + "`."
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
			Server:     server,
			TLS:        tls,
			V4Question: question,
			V6Question: question,
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
			Server: server,
			TLS:    tls,
			V4Question: dns.Question{
				Name:   "myip.opendns.com.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassANY,
			},
			V6Question: dns.Question{
				Name:   "myip.opendns.com.",
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassANY,
			},
		}
	case SourceICanHazIP:
		return HTTPv4v6{
			V4URL: "https://ipv4.icanhazip.com",
			V6URL: "https://ipv6.icanhazip.com",
		}
	case SourceIPInfo:
		return HTTPv4v6{
			V4URL: "https://ipinfo.io/ip",
			V6URL: "https://v6.ipinfo.io/ip",
		}
	case SourceIPify:
		return HTTPv4v6{
			V4URL: "https://api.ipify.org",
			V6URL: "https://api6.ipify.org",
		}
	default:
		panic("source request unimplemented: " + s.String())
	}
}
