package progressbar_test

import (
	"testing"

	"github.com/r-egorov/otus_golang/hw07_file_copying/progressbar"
	"github.com/stretchr/testify/require"
)

func TestProgressBar(t *testing.T) {
	t.Run("counts percent right", func(t *testing.T) {
		total := int64(100)
		bar := progressbar.NewBar(total)

		step := int64(25)
		expected := int64(0)

		for i := 0; i < 4; i++ {
			bar.Progress(step)
			expected += step
			require.Equal(t, expected, bar.GetPercent())
		}
	})
}
