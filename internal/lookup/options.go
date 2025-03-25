package lookup

type Option func(c *Client)

func WithDNSUseTCP(dnsUseTCP bool) Option {
	return func(c *Client) {
		c.DNSUseTCP = dnsUseTCP
	}
}
