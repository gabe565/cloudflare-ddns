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
	sources, err := conf.Sources()
	if err != nil {
		return "", err
	}

	var errs []error //nolint:prealloc
	for _, source := range sources {
		slogx.Trace("Querying source", "name", source)
		var content string
		var err error
		switch source {
		case config.SourceCloudflareTLS:
			content, err = Cloudflare(ctx, true, conf.DNSUseTCP)
		case config.SourceCloudflare:
			content, err = Cloudflare(ctx, false, conf.DNSUseTCP)
		case config.SourceOpenDNSTLS:
			content, err = OpenDNS(ctx, true, conf.DNSUseTCP)
		case config.SourceOpenDNS:
			content, err = OpenDNS(ctx, false, conf.DNSUseTCP)
		case config.SourceIPInfo:
			content, err = IPInfo(ctx)
		case config.SourceIPify:
			content, err = IPify(ctx)
		}
		if err == nil {
			slogx.Trace("Got response", "source", source, "content", content)
			return content, nil
		}
		errs = append(errs, fmt.Errorf("%s: %w", source, err))
		slog.Debug("Source failed", "source", source, "error", err)
	}
	return "", fmt.Errorf("%w: %w", ErrAllSourcesFailed, errors.Join(errs...))
}
