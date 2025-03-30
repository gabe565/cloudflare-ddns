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

func lookupDNS(ctx context.Context, host, port string, tcp, tls bool, question dns.Question) (string, error) {
	start := time.Now()
	c := &dns.Client{}
	switch {
	case tls:
		c.Net = "tcp-tls"
	case tcp:
		c.Net = "tcp"
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
		func() (string, error) { return lookupDNS(ctx, "1.1.1.1", port, tcp, tls, question) },
		func() (string, error) { return lookupDNS(ctx, "2606:4700:4700::1111", port, tcp, tls, question) },
	)
}

func OpenDNS(ctx context.Context, tls, tcp, v4, v6 bool) (Response, error) {
	port := "53"
	if tls {
		port = "853"
	}

	return lookupV4V6(v4, v6,
		func() (string, error) {
			return lookupDNS(ctx, "208.67.222.222", port, tcp, tls, dns.Question{
				Name:   "myip.opendns.com.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassANY,
			})
		},
		func() (string, error) {
			return lookupDNS(ctx, "2620:119:35::35", port, tcp, tls, dns.Question{
				Name:   "myip.opendns.com.",
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassANY,
			})
		},
	)
}
