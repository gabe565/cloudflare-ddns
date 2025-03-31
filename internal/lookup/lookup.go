package lookup

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/utils/slogx"
)

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

var ErrAllSourcesFailed = errors.New("all sources failed")

func GetPublicIP(ctx context.Context, conf *config.Config) (Response, error) {
	sources, err := conf.Sources()
	if err != nil {
		return Response{}, err
	}

	var errs []error //nolint:prealloc
	for _, source := range sources {
		slogx.Trace("Querying source", "name", source)
		var response Response
		var err error
		switch req := source.Request().(type) {
		case config.DNSv4v6:
			response, err = DNSv4v6(ctx, conf.UseV4, conf.UseV6, conf.DNSUseTCP, req)
		case config.HTTPv4v6:
			response, err = HTTPv4v6(ctx, conf.UseV4, conf.UseV6, req)
		default:
			panic("unknown request type")
		}
		if err == nil {
			slogx.Trace("Got response", "source", source, "ip", response)
			return response, nil
		}
		errs = append(errs, fmt.Errorf("%s: %w", source, err))
		slog.Debug("Source failed", "source", source, "error", err)
	}
	return Response{}, fmt.Errorf("%w: %w", ErrAllSourcesFailed, errors.Join(errs...))
}
