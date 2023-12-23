package model_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"testing"

	"github.com/landru29/dump1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecksum(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D40621D58C382D690C8AC2863A7")
		require.NoError(t, err)

		require.NoError(t, model.ModeS(dataByte).CheckSum())
	})

	t.Run("ko", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D40621D59C382D690C8AC2863A7")
		require.NoError(t, err)

		require.Error(t, model.ModeS(dataByte).CheckSum())
	})
}

func TestSquitter(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/dump1090.txt")
	require.NoError(t, err)

	scanner := bufio.NewScanner(file)

	lineIdx := 1

	for scanner.Scan() {
		line := scanner.Text()

		message := line[1 : len(line)-1]

		dataByte, err := hex.DecodeString(message)
		require.NoError(t, err, "line #%d: %s", lineIdx, message)

		squitter, err := model.ModeS(dataByte).Squitter()
		require.NoError(t, err, "line #%d: %s", lineIdx, message)

		switch {
		case len(dataByte) == 7:
			assert.Equal(t, "short squitter", squitter.Name(), "line #%d: %s", lineIdx, message)
		case len(dataByte) == 14:
			assert.Equal(t, "extended squitter", squitter.Name(), "line #%d: %s", lineIdx, message)
		default:
			assert.True(t, false, "wrong size of message, line #%d: %s", lineIdx, message)
		}

		lineIdx++
	}
}
