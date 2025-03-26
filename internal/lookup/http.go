package lookup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"gabe565.com/utils/slogx"
)

var ErrUpstreamStatus = errors.New("upstream error")

func HTTPPlain(ctx context.Context, url string, unquote bool) (string, error) {
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

	if unquote {
		return strconv.Unquote(string(b))
	}
	return string(b), nil
}

func IPInfo(ctx context.Context) (string, error) {
	return HTTPPlain(ctx, "https://ipinfo.io/ip", false)
}

func IPify(ctx context.Context) (string, error) {
	return HTTPPlain(ctx, "https://api.ipify.org", false)
}
