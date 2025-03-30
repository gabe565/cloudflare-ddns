package lookup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gabe565.com/cloudflare-ddns/internal/errsgroup"
	"gabe565.com/utils/slogx"
)

var ErrUpstreamStatus = errors.New("upstream error")

func HTTPPlain(ctx context.Context, url string) (string, error) {
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

	return string(b), nil
}

func IPInfo(ctx context.Context, v4, v6 bool) (Response, error) {
	var response Response
	var group errsgroup.Group

	if v4 {
		group.Go(func() error {
			var err error
			response.IPV4, err = HTTPPlain(ctx, "https://ipinfo.io/ip")
			return err
		})
	}

	if v6 {
		group.Go(func() error {
			var err error
			response.IPV6, err = HTTPPlain(ctx, "https://v6.ipinfo.io/ip")
			return err
		})
	}

	err := group.Wait()
	return response, err
}

func IPify(ctx context.Context, v4, v6 bool) (Response, error) {
	var response Response
	var group errsgroup.Group

	if v4 {
		group.Go(func() error {
			var err error
			response.IPV4, err = HTTPPlain(ctx, "https://api.ipify.org")
			return err
		})
	}

	if v6 {
		group.Go(func() error {
			var err error
			response.IPV6, err = HTTPPlain(ctx, "https://api6.ipify.org")
			return err
		})
	}

	err := group.Wait()
	return response, err
}
