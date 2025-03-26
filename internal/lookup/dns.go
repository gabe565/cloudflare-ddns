package lookup

import (
	"context"
	"errors"
	"time"

	"gabe565.com/utils/slogx"
	"github.com/miekg/dns"
)

type lookupDNSOptions struct {
	server string
	useTCP bool
	useTLS bool
	req    dns.Question
}

func (l lookupDNSOptions) net() string {
	switch {
	case l.useTLS:
		return "tcp-tls"
	case l.useTCP:
		return "tcp"
	default:
		return ""
	}
}

var ErrNoDNSAnswer = errors.New("no DNS answer")

func lookupDNS(ctx context.Context, opts lookupDNSOptions) (string, error) {
	start := time.Now()
	c := &dns.Client{Net: opts.net()}
	m := &dns.Msg{Question: []dns.Question{opts.req}}

	slogx.Trace("DNS query", "server", opts.server, "net", c.Net, "name", opts.req.Name, "type", dns.TypeToString[opts.req.Qtype], "class", dns.ClassToString[opts.req.Qclass])

	res, _, err := c.ExchangeContext(ctx, m, opts.server)
	if err != nil {
		return "", err
	}

	slogx.Trace("DNS response", "took", time.Since(start), "server", opts.server, "response", res)

	if len(res.Answer) == 0 {
		return "", ErrNoDNSAnswer
	}

	var val string
	switch answer := res.Answer[0].(type) {
	case *dns.A:
		val = answer.A.String()
	case *dns.TXT:
		if len(answer.Txt) == 0 {
			return "", ErrNoDNSAnswer
		}
		val = answer.Txt[0]
	}
	return val, nil
}

func Cloudflare(ctx context.Context, tls, tcp bool) (string, error) {
	server := "1.1.1.1:53"
	if tls {
		server = "1.1.1.1:853"
	}
	return lookupDNS(ctx, lookupDNSOptions{
		server: server,
		useTLS: tls,
		useTCP: tcp,
		req: dns.Question{
			Name:   "whoami.cloudflare.",
			Qtype:  dns.TypeTXT,
			Qclass: dns.ClassCHAOS,
		},
	})
}

func OpenDNS(ctx context.Context, tls, tcp bool) (string, error) {
	server := "dns.opendns.com:53"
	if tls {
		server = "dns.opendns.com:853"
	}
	return lookupDNS(ctx, lookupDNSOptions{
		server: server,
		useTLS: tls,
		useTCP: tcp,
		req: dns.Question{
			Name:   "myip.opendns.com.",
			Qtype:  dns.TypeA,
			Qclass: dns.ClassANY,
		},
	})
}
