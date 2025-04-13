package lookup

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gabe565.com/utils/slogx"
)

func NewClient(opts ...Option) *Client {
	c := &Client{v4: true}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Client struct {
	v4, v6, tcp bool
	sources     []Source
}

var ErrAllSourcesFailed = errors.New("all sources failed")

func (c *Client) GetPublicIP(ctx context.Context) (Response, error) {
	var errs []error //nolint:prealloc
	for _, source := range c.sources {
		slogx.Trace("Querying source", "name", source)
		var response Response
		var err error
		switch req := source.Request().(type) {
		case DNSv4v6:
			response, err = c.DNSv4v6(ctx, req)
		case HTTPv4v6:
			response, err = c.HTTPv4v6(ctx, req)
		default:
			panic("unknown request type")
		}
		if err == nil {
			slog.Debug("Got response", "source", source, "ip", response)
			return response, nil
		}
		errs = append(errs, fmt.Errorf("%s: %w", source, err))
		slog.Debug("Source failed", "source", source, "error", err)
	}
	return Response{}, fmt.Errorf("%w: %w", ErrAllSourcesFailed, errors.Join(errs...))
}

type Response struct {
	IPV4, IPV6 string
}

func (r Response) LogValue() slog.Value {
	attr := make([]slog.Attr, 0, 2)
	if r.IPV4 != "" {
		attr = append(attr, slog.String("v4", r.IPV4))
	}
	if r.IPV6 != "" {
		attr = append(attr, slog.String("v6", r.IPV6))
	}
	return slog.GroupValue(attr...)
}
