package lookup

import (
	"context"
	"errors"
	"time"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/errsgroup"
	"gabe565.com/utils/slogx"
	"github.com/miekg/dns"
)

var ErrNoDNSAnswer = errors.New("no DNS answer")

func lookupDNS(ctx context.Context, server string, tcp, ipv6, tls bool, question dns.Question) (string, error) {
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

func DNSv4v6(ctx context.Context, v4, v6, tcp bool, req config.DNSv4v6) (Response, error) {
	var response Response
	var group errsgroup.Group

	if v4 {
		group.Go(func() error {
			var err error
			response.IPV4, err = lookupDNS(ctx, req.Server, tcp, false, req.TLS, req.V4Question)
			return err
		})
	}

	if v6 {
		group.Go(func() error {
			var err error
			response.IPV6, err = lookupDNS(ctx, req.Server, tcp, true, req.TLS, req.V6Question)
			return err
		})
	}

	err := group.Wait()
	return response, err
}
