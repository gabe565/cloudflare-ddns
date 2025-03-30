package lookup

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/errsgroup"
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
		switch source {
		case config.SourceCloudflareTLS:
			response, err = Cloudflare(ctx, true, conf.DNSUseTCP, conf.UseV4, conf.UseV6)
		case config.SourceCloudflare:
			response, err = Cloudflare(ctx, false, conf.DNSUseTCP, conf.UseV4, conf.UseV6)
		case config.SourceOpenDNSTLS:
			response, err = OpenDNS(ctx, true, conf.DNSUseTCP, conf.UseV4, conf.UseV6)
		case config.SourceOpenDNS:
			response, err = OpenDNS(ctx, false, conf.DNSUseTCP, conf.UseV4, conf.UseV6)
		case config.SourceICanHazIP:
			response, err = ICanHazIP(ctx, conf.UseV4, conf.UseV6)
		case config.SourceIPInfo:
			response, err = IPInfo(ctx, conf.UseV4, conf.UseV6)
		case config.SourceIPify:
			response, err = IPify(ctx, conf.UseV4, conf.UseV6)
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

func lookupV4V6(v4, v6 bool, v4func, v6func func() (string, error)) (Response, error) {
	var response Response
	var group errsgroup.Group

	if v4 {
		group.Go(func() error {
			var err error
			response.IPV4, err = v4func()
			return err
		})
	}

	if v6 {
		group.Go(func() error {
			var err error
			response.IPV6, err = v6func()
			return err
		})
	}

	err := group.Wait()
	return response, err
}
