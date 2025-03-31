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

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/errsgroup"
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
	//nolint:errcheck
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

func HTTPv4v6(ctx context.Context, v4, v6 bool, req config.HTTPv4v6) (Response, error) {
	var response Response
	var group errsgroup.Group

	if v4 {
		group.Go(func() error {
			var err error
			response.IPV4, err = httpPlain(ctx, tcp4, req.URLv4)
			return err
		})
	}

	if v6 {
		group.Go(func() error {
			var err error
			response.IPV6, err = httpPlain(ctx, tcp6, req.URLv6)
			return err
		})
	}

	err := group.Wait()
	return response, err
}
