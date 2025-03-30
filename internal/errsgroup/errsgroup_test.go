package errsgroup

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		var group Group
		group.Go(func() error {
			return nil
		})
		require.NoError(t, group.Wait())
	})

	t.Run("1 error", func(t *testing.T) {
		var group Group
		group.Go(func() error {
			return errors.New("some error")
		})
		err := group.Wait()
		require.Error(t, err)
		assert.Equal(t, "some error", err.Error())
	})

	t.Run("2 errors", func(t *testing.T) {
		var group Group
		group.Go(func() error {
			return errors.New("some error")
		})
		group.Go(func() error {
			return errors.New("another error")
		})
		err := group.Wait()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "some error")
		assert.Contains(t, err.Error(), "another error")
	})

	t.Run("waits", func(t *testing.T) {
		var group Group
		var n int
		group.Go(func() error {
			time.Sleep(time.Millisecond)
			n++
			return nil
		})
		require.NoError(t, group.Wait())
		assert.Equal(t, 1, n)
	})
}
