package lookup

import (
	"context"
	"errors"
	"net"
	"time"

	"gabe565.com/utils/slogx"
	"github.com/miekg/dns"
)

var ErrNoDNSAnswer = errors.New("no DNS answer")

func lookupDNS(ctx context.Context, host, port string, tcp, ipv6, tls bool, question dns.Question) (string, error) {
	start := time.Now()
	c := &dns.Client{}
	if ipv6 {
		switch {
		case tls:
			c.Net = "tcp6-tls"
		case tcp:
			c.Net = "tcp6"
		default:
			c.Net = "udp6"
		}
	} else {
		switch {
		case tls:
			c.Net = "tcp4-tls"
		case tcp:
			c.Net = "tcp4"
		default:
			c.Net = "udp4"
		}
	}
	m := &dns.Msg{Question: []dns.Question{question}}

	server := net.JoinHostPort(host, port)
	slogx.Trace("DNS query",
		"server", server,
		"net", c.Net,
		"name", question.Name,
		"type", dns.TypeToString[question.Qtype],
		"class", dns.ClassToString[question.Qclass],
	)

	res, _, err := c.ExchangeContext(ctx, m, server)
	if err != nil {
		return "", err
	}

	slogx.Trace("DNS response", "took", time.Since(start), "server", server, "response", res)

	if len(res.Answer) == 0 {
		return "", ErrNoDNSAnswer
	}

	var val string
	switch answer := res.Answer[0].(type) {
	case *dns.A:
		val = answer.A.String()
	case *dns.AAAA:
		val = answer.AAAA.String()
	case *dns.TXT:
		if len(answer.Txt) == 0 {
			return "", ErrNoDNSAnswer
		}
		val = answer.Txt[0]
	}
	return val, nil
}

func Cloudflare(ctx context.Context, tls, tcp, v4, v6 bool) (Response, error) {
	port := "53"
	if tls {
		port = "853"
	}

	question := dns.Question{
		Name:   "whoami.cloudflare.",
		Qtype:  dns.TypeTXT,
		Qclass: dns.ClassCHAOS,
	}

	return lookupV4V6(v4, v6,
		func() (string, error) { return lookupDNS(ctx, "one.one.one.one", port, tcp, false, tls, question) },
		func() (string, error) { return lookupDNS(ctx, "one.one.one.one", port, tcp, true, tls, question) },
	)
}

func OpenDNS(ctx context.Context, tls, tcp, v4, v6 bool) (Response, error) {
	port := "53"
	if tls {
		port = "853"
	}

	return lookupV4V6(v4, v6,
		func() (string, error) {
			return lookupDNS(ctx, "dns.opendns.com", port, tcp, false, tls, dns.Question{
				Name:   "myip.opendns.com.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassANY,
			})
		},
		func() (string, error) {
			return lookupDNS(ctx, "dns.opendns.com", port, tcp, true, tls, dns.Question{
				Name:   "myip.opendns.com.",
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassANY,
			})
		},
	)
}
