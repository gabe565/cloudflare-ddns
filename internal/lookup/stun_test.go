package lookup

import (
	"net"
	"testing"

	"github.com/pion/stun/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSTUNServer(t *testing.T, network string, build func(txID [stun.TransactionIDSize]byte) *stun.Message) string {
	t.Helper()

	addr := loopbackV4
	if network == udp6 {
		addr = loopbackV6
	}

	conn, err := (&net.ListenConfig{}).ListenPacket(t.Context(), network, addr)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	go func() {
		buf := make([]byte, 1280)
		for {
			n, from, err := conn.ReadFrom(buf)
			if err != nil {
				return
			}

			req := &stun.Message{Raw: buf[:n]}
			if err := req.Decode(); err != nil {
				continue
			}

			if _, err := conn.WriteTo(build(req.TransactionID).Raw, from); err != nil {
				return
			}
		}
	}()

	return conn.LocalAddr().String()
}

func respond(
	t *testing.T,
	msgType stun.MessageType,
	attrs ...stun.Setter,
) func([stun.TransactionIDSize]byte) *stun.Message {
	t.Helper()
	return func(txID [stun.TransactionIDSize]byte) *stun.Message {
		setters := append([]stun.Setter{stun.NewTransactionIDSetter(txID), msgType}, attrs...)
		msg, err := stun.Build(setters...)
		require.NoError(t, err)
		return msg
	}
}

func xorMapped(t *testing.T, ip string) func([stun.TransactionIDSize]byte) *stun.Message {
	t.Helper()
	return respond(t, stun.BindingSuccess, &stun.XORMappedAddress{IP: net.ParseIP(ip), Port: 1234})
}

func Test_STUN(t *testing.T) {
	t.Run("v4", func(t *testing.T) {
		server := newSTUNServer(t, udp4, xorMapped(t, cfV4))
		got, err := STUN(t.Context(), udp4, server)
		require.NoError(t, err)
		assert.Equal(t, cfV4, got)
	})

	t.Run("v6", func(t *testing.T) {
		server := newSTUNServer(t, udp6, xorMapped(t, cfV6))
		got, err := STUN(t.Context(), udp6, server)
		require.NoError(t, err)
		assert.Equal(t, cfV6, got)
	})

	t.Run("mapped address fallback", func(t *testing.T) {
		server := newSTUNServer(t, udp4,
			respond(t, stun.BindingSuccess, &stun.MappedAddress{IP: net.ParseIP(cfV4), Port: 1234}),
		)
		got, err := STUN(t.Context(), udp4, server)
		require.NoError(t, err)
		assert.Equal(t, cfV4, got)
	})

	t.Run("skips unknown attributes", func(t *testing.T) {
		server := newSTUNServer(t, udp4, respond(t, stun.BindingSuccess,
			stun.NewSoftware("test software"),
			&stun.XORMappedAddress{IP: net.ParseIP(cfV4), Port: 1234},
		))
		got, err := STUN(t.Context(), udp4, server)
		require.NoError(t, err)
		assert.Equal(t, cfV4, got)
	})

	t.Run("error response", func(t *testing.T) {
		server := newSTUNServer(t, udp4, respond(t, stun.BindingError, stun.CodeBadRequest))
		_, err := STUN(t.Context(), udp4, server)
		require.ErrorIs(t, err, ErrSTUNError)
		assert.Contains(t, err.Error(), "400 Bad Request")
	})

	t.Run("no address", func(t *testing.T) {
		server := newSTUNServer(t, udp4, respond(t, stun.BindingSuccess))
		_, err := STUN(t.Context(), udp4, server)
		require.ErrorIs(t, err, ErrNoSTUNAddr)
	})

	t.Run("transaction id mismatch", func(t *testing.T) {
		var zero [stun.TransactionIDSize]byte
		server := newSTUNServer(t, udp4, func([stun.TransactionIDSize]byte) *stun.Message {
			return xorMapped(t, cfV4)(zero)
		})
		_, err := STUN(t.Context(), udp4, server)
		require.ErrorIs(t, err, ErrSTUNResponse)
	})
}

func TestSTUNv4v6(t *testing.T) {
	t.Run("both", func(t *testing.T) {
		c := Client{v4: true, v6: true}
		got, err := c.STUNv4v6(t.Context(), STUNv4v6{
			ServerV4: newSTUNServer(t, udp4, xorMapped(t, cfV4)),
			ServerV6: newSTUNServer(t, udp6, xorMapped(t, cfV6)),
		})
		require.NoError(t, err)

		expect := Response{IPV4: cfV4, IPV6: cfV6}
		assert.Equal(t, expect, got)
	})

	t.Run("only v4", func(t *testing.T) {
		c := Client{v4: true}
		got, err := c.STUNv4v6(t.Context(), STUNv4v6{
			ServerV4: newSTUNServer(t, udp4, xorMapped(t, cfV4)),
		})
		require.NoError(t, err)

		expect := Response{IPV4: cfV4}
		assert.Equal(t, expect, got)
	})

	t.Run("only v6", func(t *testing.T) {
		c := Client{v6: true}
		got, err := c.STUNv4v6(t.Context(), STUNv4v6{
			ServerV6: newSTUNServer(t, udp6, xorMapped(t, cfV6)),
		})
		require.NoError(t, err)

		expect := Response{IPV6: cfV6}
		assert.Equal(t, expect, got)
	})
}
