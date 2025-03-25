package lookup

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/utils/slogx"
)

var ErrAllSourcesFailed = errors.New("all sources failed")

func GetPublicIP(ctx context.Context, conf *config.Config) (string, error) {
	client := New(WithDNSUseTCP(conf.DNSUseTCP))

	var errs []error //nolint:prealloc
	for _, source := range conf.Sources {
		slogx.Trace("Querying source", "name", source)
		var content string
		var err error
		switch source {
		case config.SourceCloudflareTLS:
			content, err = client.CloudflareTLS(ctx)
		case config.SourceCloudflare:
			content, err = client.Cloudflare(ctx)
		case config.SourceOpenDNSTLS:
			content, err = client.OpenDNSTLS(ctx)
		case config.SourceOpenDNS:
			content, err = client.OpenDNS(ctx)
		case config.SourceIPInfo:
			content, err = client.IPInfo(ctx)
		case config.SourceIPify:
			content, err = client.IPify(ctx)
		}
		if err == nil {
			slogx.Trace("Got response", "source", source, "content", content)
			return content, nil
		}
		errs = append(errs, fmt.Errorf("%s: %w", source, err))
		slog.Warn("Source failed", "source", source, "error", err)
	}
	return "", fmt.Errorf("%w: %w", ErrAllSourcesFailed, errors.Join(errs...))
}
