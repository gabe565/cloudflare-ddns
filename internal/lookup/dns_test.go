package lookup

import (
	"net"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	domain = "test."
	cfV4   = "1.1.1.1"
	cfV6   = "2606:4700:4700::1111"
)

func newDNSServer(t *testing.T, network string) string {
	mux := dns.NewServeMux()
	mux.HandleFunc(domain, func(w dns.ResponseWriter, r *dns.Msg) {
		m := &dns.Msg{}
		m.SetReply(r)
		switch r.Question[0].Qtype {
		case dns.TypeA:
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET},
				A:   net.ParseIP(cfV4),
			})
		case dns.TypeAAAA:
			m.Answer = append(m.Answer, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET},
				AAAA: net.ParseIP(cfV6),
			})
		}
		err := w.WriteMsg(m)
		require.NoError(t, err)
	})

	addr := "127.0.0.1:0"
	if network == tcp6 {
		addr = "[::1]:0"
	}

	ready := make(chan struct{})
	server := dns.Server{
		Addr:              addr,
		Net:               network,
		Handler:           mux,
		NotifyStartedFunc: func() { close(ready) },
	}
	t.Cleanup(func() {
		_ = server.Shutdown()
	})

	go func() {
		assert.NoError(t, server.ListenAndServe())
	}()

	select {
	case <-t.Context().Done():
		return ""
	case <-ready:
		return server.Listener.Addr().String()
	}
}

func Test_DNS(t *testing.T) {
	t.Run("v4", func(t *testing.T) {
		got, err := DNS(t.Context(), newDNSServer(t, tcp4), true, false, false, dns.Question{
			Name:   domain,
			Qtype:  dns.TypeA,
			Qclass: dns.ClassINET,
		})
		require.NoError(t, err)
		assert.Equal(t, cfV4, got)
	})

	t.Run("v6", func(t *testing.T) {
		got, err := DNS(t.Context(), newDNSServer(t, tcp6), true, true, false, dns.Question{
			Name:   domain,
			Qtype:  dns.TypeAAAA,
			Qclass: dns.ClassINET,
		})
		require.NoError(t, err)
		assert.Equal(t, cfV6, got)
	})
}

func TestDNSv4v6(t *testing.T) {
	t.Run("both", func(t *testing.T) {
		c := Client{v4: true, v6: true, tcp: true}
		got, err := c.DNSv4v6(t.Context(), DNSv4v6{
			ServerV4: newDNSServer(t, tcp4),
			QuestionV4: dns.Question{
				Name:   domain,
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			},
			ServerV6: newDNSServer(t, tcp6),
			QuestionV6: dns.Question{
				Name:   domain,
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassINET,
			},
		})
		require.NoError(t, err)

		expect := Response{IPV4: cfV4, IPV6: cfV6}
		assert.Equal(t, expect, got)
	})

	t.Run("only v4", func(t *testing.T) {
		c := Client{v4: true, tcp: true}
		got, err := c.DNSv4v6(t.Context(), DNSv4v6{
			ServerV4: newDNSServer(t, tcp4),
			QuestionV4: dns.Question{
				Name:   domain,
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			},
		})
		require.NoError(t, err)

		expect := Response{IPV4: cfV4}
		assert.Equal(t, expect, got)
	})

	t.Run("only v6", func(t *testing.T) {
		c := Client{v6: true, tcp: true}
		got, err := c.DNSv4v6(t.Context(), DNSv4v6{
			ServerV6: newDNSServer(t, tcp6),
			QuestionV6: dns.Question{
				Name:   domain,
				Qtype:  dns.TypeAAAA,
				Qclass: dns.ClassINET,
			},
		})
		require.NoError(t, err)

		expect := Response{IPV6: cfV6}
		assert.Equal(t, expect, got)
	})
}
