package lookup

import (
	"strings"

	"gabe565.com/cloudflare-ddns/internal/output"
	"github.com/charmbracelet/lipgloss"
	"github.com/miekg/dns"
)

//go:generate go tool enumer -type Source -transform snake -linecomment -output source_string.go

type Source uint8

const (
	CloudflareTLS Source = iota
	Cloudflare
	OpenDNSTLS // opendns_tls
	OpenDNS    // opendns
	ICanHazIP  // icanhazip
	IPInfo     // ipinfo
	IPify      // ipify
)

func (s Source) Description(format output.Format) string {
	req := s.Request()
	return req.Description(format)
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
		bold := lipgloss.NewStyle().Bold(true).Render
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
		bold := lipgloss.NewStyle().Bold(true).Render
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
	case CloudflareTLS:
		server = "one.one.one.one:853"
		tls = true
		fallthrough
	case Cloudflare:
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
	case OpenDNSTLS:
		server = "dns.opendns.com:853"
		tls = true
		fallthrough
	case OpenDNS:
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
	case ICanHazIP:
		return HTTPv4v6{
			URLv4: "https://ipv4.icanhazip.com",
			URLv6: "https://ipv6.icanhazip.com",
		}
	case IPInfo:
		return HTTPv4v6{
			URLv4: "https://ipinfo.io/ip",
			URLv6: "https://v6.ipinfo.io/ip",
		}
	case IPify:
		return HTTPv4v6{
			URLv4: "https://api.ipify.org",
			URLv6: "https://api6.ipify.org",
		}
	default:
		panic("source request unimplemented: " + s.String())
	}
}
