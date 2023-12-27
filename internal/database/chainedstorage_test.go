package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/landru29/dump1090/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChainedStorage(t *testing.T) {
	t.Parallel()

	t.Run("new database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		storage := database.NewChainedStorage[string, float64](
			ctx,
			database.ChainedWithLifetime[string, float64](time.Millisecond*100),
			database.ChainedWithCleanCycle[string, float64](time.Millisecond*30),
		)

		t.Cleanup(func() {
			require.NoError(t, storage.Close())
		})

		time.Sleep(time.Millisecond * 500)

		assert.Empty(t, storage.Keys())
	})

	t.Run("chain elements", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		storage := database.NewChainedStorage[string, float64](
			ctx,
			database.ChainedWithLifetime[string, float64](time.Second*3000),
			database.ChainedWithCleanCycle[string, float64](time.Second*3000),
		)

		t.Cleanup(func() {
			require.NoError(t, storage.Close())
		})

		storage.Add("42", 42.0)
		storage.Add("42", 24.0)
		storage.Add("24", 42.0)

		assert.Len(t, storage.Keys(), 2)

		assert.Len(t, storage.Elements("42"), 2)
		assert.Len(t, storage.Elements("24"), 1)
	})

	t.Run("add elements", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		storage := database.NewChainedStorage[string, float64](
			ctx,
			database.ChainedWithLifetime[string, float64](time.Millisecond*100),
			database.ChainedWithCleanCycle[string, float64](time.Millisecond*130),
		)

		t.Cleanup(func() {
			require.NoError(t, storage.Close())
		})

		storage.Add("42", 42.0)
		storage.Add("42", 24.0)
		storage.Add("24", 42.0)

		assert.Len(t, storage.Keys(), 2)

		time.Sleep(time.Millisecond * 500)

		assert.Empty(t, storage.Keys())
	})

	t.Run("Keep some elements", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		storage := database.NewChainedStorage[string, float64](
			ctx,
			database.ChainedWithLifetime[string, float64](time.Millisecond*100),
			database.ChainedWithCleanCycle[string, float64](time.Millisecond*30),
		)

		t.Cleanup(func() {
			require.NoError(t, storage.Close())
		})

		storage.Add("42", 42.0)
		storage.Add("24", 42.0)
		time.Sleep(time.Millisecond * 50)
		storage.Add("42", 24.0)

		assert.Len(t, storage.Keys(), 2)

		time.Sleep(time.Millisecond * 80)

		assert.Len(t, storage.Keys(), 1)
		require.Len(t, storage.Elements("42"), 1)
		assert.Equal(t, storage.Elements("42")[0], 24.0) //nolint: testifylint

		assert.Nil(t, storage.Elements("24"))
	})

	t.Run("context cancelled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())

		cancel()

		storage := database.NewChainedStorage[string, float64](
			ctx,
			database.ChainedWithLifetime[string, float64](time.Millisecond*100),
			database.ChainedWithCleanCycle[string, float64](time.Millisecond*30),
		)

		t.Cleanup(func() {
			require.NoError(t, storage.Close())
		})

		storage.Add("42", 42.0)
		storage.Add("42", 24.0)
		storage.Add("24", 42.0)

		assert.Len(t, storage.Keys(), 2)

		time.Sleep(time.Millisecond * 500)

		assert.Len(t, storage.Keys(), 2)
	})
}
