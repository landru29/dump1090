package rtl28xxx_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/landru29/dump1090/internal/mocks"
	"github.com/landru29/dump1090/internal/source/rtl28xxx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReader(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProcessor := mocks.NewMockProcesser(ctrl)

	mockProcessor.EXPECT().Process(gomock.Any()).AnyTimes()

	file, err := os.Open("../../../testdata/modes1.bin")
	require.NoError(t, err)

	defer func(closer io.Closer) {
		require.NoError(t, closer.Close())
	}(file)

	reader := rtl28xxx.NewReader(file, mockProcessor)
	require.NotNil(t, reader)

	require.NoError(t, reader.Start(context.Background()))
}
