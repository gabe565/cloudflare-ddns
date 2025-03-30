package lookup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gabe565.com/utils/slogx"
)

var ErrUpstreamStatus = errors.New("upstream error")

func httpPlain(ctx context.Context, url string) (string, error) {
	start := time.Now()
	slogx.Trace("HTTP request", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	slogx.Trace("HTTP response", "took", time.Since(start), "status", res.Status, "body", string(b))

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %s", ErrUpstreamStatus, res.Status)
	}

	return string(bytes.TrimSpace(b)), nil
}

func ICanHazIP(ctx context.Context, v4, v6 bool) (Response, error) {
	return lookupV4V6(v4, v6,
		func() (string, error) { return httpPlain(ctx, "https://ipv4.icanhazip.com") },
		func() (string, error) { return httpPlain(ctx, "https://ipv6.icanhazip.com") },
	)
}

func IPInfo(ctx context.Context, v4, v6 bool) (Response, error) {
	return lookupV4V6(v4, v6,
		func() (string, error) { return httpPlain(ctx, "https://ipinfo.io/ip") },
		func() (string, error) { return httpPlain(ctx, "https://v6.ipinfo.io/ip") },
	)
}

func IPify(ctx context.Context, v4, v6 bool) (Response, error) {
	return lookupV4V6(v4, v6,
		func() (string, error) { return httpPlain(ctx, "https://api.ipify.org") },
		func() (string, error) { return httpPlain(ctx, "https://api6.ipify.org") },
	)
}
