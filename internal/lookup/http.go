package lookup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"gabe565.com/utils/slogx"
)

var ErrUpstreamStatus = errors.New("upstream error")

func httpPlain(ctx context.Context, network, url string) (string, error) {
	start := time.Now()
	slogx.Trace("HTTP request", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	dialer := net.Dialer{}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, _, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}
	client := http.Client{Transport: transport}

	res, err := client.Do(req)
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
		func() (string, error) { return httpPlain(ctx, "tcp4", "https://ipv4.icanhazip.com") },
		func() (string, error) { return httpPlain(ctx, "tcp6", "https://ipv6.icanhazip.com") },
	)
}

func IPInfo(ctx context.Context, v4, v6 bool) (Response, error) {
	return lookupV4V6(v4, v6,
		func() (string, error) { return httpPlain(ctx, "tcp4", "https://ipinfo.io/ip") },
		func() (string, error) { return httpPlain(ctx, "tcp6", "https://v6.ipinfo.io/ip") },
	)
}

func IPify(ctx context.Context, v4, v6 bool) (Response, error) {
	return lookupV4V6(v4, v6,
		func() (string, error) { return httpPlain(ctx, "tcp4", "https://api.ipify.org") },
		func() (string, error) { return httpPlain(ctx, "tcp6", "https://api6.ipify.org") },
	)
}
