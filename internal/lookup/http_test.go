package lookup

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHTTPServer(t *testing.T, network string) *httptest.Server {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		switch network {
		case tcp4:
			_, err := w.Write([]byte(cfV4))
			assert.NoError(t, err)
		case tcp6:
			_, err := w.Write([]byte(cfV6))
			assert.NoError(t, err)
		}
	}))

	addr := "127.0.0.1:0"
	if network == tcp6 {
		addr = "[::1]:0"
	}

	var err error
	server.Listener, err = net.Listen(network, addr)
	require.NoError(t, err)

	server.Start()
	t.Cleanup(server.Close)

	return server
}

func Test_HTTPPlain(t *testing.T) {
	t.Run("v4", func(t *testing.T) {
		server := newHTTPServer(t, tcp4)
		got, err := HTTP(t.Context(), tcp4, server.URL)
		require.NoError(t, err)
		assert.Equal(t, cfV4, got)
	})

	t.Run("v6", func(t *testing.T) {
		server := newHTTPServer(t, tcp6)
		got, err := HTTP(t.Context(), tcp6, server.URL)
		require.NoError(t, err)
		assert.Equal(t, cfV6, got)
	})
}

func TestHTTPv4v6(t *testing.T) {
	t.Run("both", func(t *testing.T) {
		c := Client{v4: true, v6: true}
		got, err := c.HTTPv4v6(t.Context(), HTTPv4v6{
			URLv4: newHTTPServer(t, tcp4).URL,
			URLv6: newHTTPServer(t, tcp6).URL,
		})
		require.NoError(t, err)

		expect := Response{IPV4: cfV4, IPV6: cfV6}
		assert.Equal(t, expect, got)
	})

	t.Run("only v4", func(t *testing.T) {
		c := Client{v4: true}
		got, err := c.HTTPv4v6(t.Context(), HTTPv4v6{
			URLv4: newHTTPServer(t, tcp4).URL,
		})
		require.NoError(t, err)

		expect := Response{IPV4: cfV4}
		assert.Equal(t, expect, got)
	})

	t.Run("only v6", func(t *testing.T) {
		c := Client{v6: true}
		got, err := c.HTTPv4v6(t.Context(), HTTPv4v6{
			URLv6: newHTTPServer(t, tcp6).URL,
		})
		require.NoError(t, err)

		expect := Response{IPV6: cfV6}
		assert.Equal(t, expect, got)
	})
}
