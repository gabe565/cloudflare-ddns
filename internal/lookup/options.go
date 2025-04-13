package lookup

type Option func(*Client)

func WithV4(enabled bool) Option {
	return func(c *Client) {
		c.v4 = enabled
	}
}

func WithV6(enabled bool) Option {
	return func(c *Client) {
		c.v6 = enabled
	}
}

func WithForceTCP(enabled bool) Option {
	return func(c *Client) {
		c.tcp = enabled
	}
}

func WithSources(sources ...Source) Option {
	return func(c *Client) {
		c.sources = sources
	}
}
