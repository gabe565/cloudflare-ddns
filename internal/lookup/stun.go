package lookup

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"gabe565.com/cloudflare-ddns/internal/errsgroup"
	"gabe565.com/utils/slogx"
	"github.com/pion/stun/v4"
)

const (
	udp4 = "udp4"
	udp6 = "udp6"
)

const stunTimeout = 2 * time.Second

var (
	ErrSTUNResponse = errors.New("invalid STUN response")
	ErrSTUNError    = errors.New("STUN error response")
	ErrNoSTUNAddr   = errors.New("no address in STUN response")
)

func STUN(ctx context.Context, network, server string) (string, error) {
	start := time.Now()

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, network, server)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = conn.Close()
	}()

	slogx.Trace("STUN request", "server", server, "net", network, "addr", conn.RemoteAddr())

	deadline := time.Now().Add(stunTimeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err := conn.SetDeadline(deadline); err != nil {
		return "", err
	}

	req := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	if _, err := conn.Write(req.Raw); err != nil {
		return "", err
	}

	buf := make([]byte, 1280)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	ip, err := parseSTUNResponse(buf[:n], req.TransactionID)
	if err != nil {
		return "", err
	}

	slogx.Trace("STUN response", "took", time.Since(start), "server", server, "ip", ip)
	return ip, nil
}

func parseSTUNResponse(b []byte, txID [stun.TransactionIDSize]byte) (string, error) {
	res := &stun.Message{Raw: b}
	if err := res.Decode(); err != nil {
		return "", fmt.Errorf("%w: %w", ErrSTUNResponse, err)
	}

	if res.TransactionID != txID {
		return "", fmt.Errorf("%w: transaction ID mismatch", ErrSTUNResponse)
	}

	if res.Type != stun.BindingSuccess {
		var code stun.ErrorCodeAttribute
		if err := code.GetFrom(res); err == nil {
			return "", fmt.Errorf("%w: %d %s", ErrSTUNError, code.Code, code.Reason)
		}
		return "", fmt.Errorf("%w: unexpected message type %s", ErrSTUNResponse, res.Type)
	}

	var addr stun.XORMappedAddress
	switch err := addr.GetFrom(res); {
	case err == nil:
		return addr.IP.String(), nil
	case !errors.Is(err, stun.ErrAttributeNotFound):
		return "", fmt.Errorf("%w: %w", ErrSTUNResponse, err)
	}

	var mapped stun.MappedAddress
	switch err := mapped.GetFrom(res); {
	case err == nil:
		return mapped.IP.String(), nil
	case !errors.Is(err, stun.ErrAttributeNotFound):
		return "", fmt.Errorf("%w: %w", ErrSTUNResponse, err)
	}

	return "", ErrNoSTUNAddr
}

func (c *Client) STUNv4v6(ctx context.Context, req STUNv4v6) (Response, error) {
	var response Response
	var group errsgroup.Group

	if c.v4 {
		group.Go(func() error {
			var err error
			response.IPV4, err = STUN(ctx, udp4, req.ServerV4)
			return err
		})
	}

	if c.v6 {
		group.Go(func() error {
			var err error
			response.IPV6, err = STUN(ctx, udp6, req.ServerV6)
			return err
		})
	}

	err := group.Wait()
	return response, err
}
