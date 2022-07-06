package progressbar

import (
	"fmt"
	"io"
	"strings"
)

const (
	barWidth      = 50
	finishMessage = "--- Finished! ---"
)

type Bar struct {
	percent, cur, total int64
	barString, char     string
	finishChan          chan struct{}
}

// Progress takes a current value of an operation,
// adds it to the `cur` and adjusts the `barString`.
func (b *Bar) Progress(cur int64) {
	b.cur += cur
	b.percent = b.getPercent()
	b.barString = strings.Repeat(b.char, int(b.percent)/(100/barWidth))
}

func (b *Bar) getPercent() int64 {
	return int64((float32(b.cur) / float32(b.total)) * 100)
}

// Finish stops the output-goroutine, flushes the buffer and prints the finish.
func (b *Bar) Finish() {
	b.finishChan <- struct{}{}
	fmt.Printf("\n%*s\n", (barWidth+len(finishMessage))/2, finishMessage)
}

// Start starts the output-goroutine.
func (b *Bar) Start() {
	go func() {
		for {
			select {
			case <-b.finishChan:
				// We need to duplicate the `b.showProgress()` call
				// because the goroutine does not manage to output
				// when the limit is too small.
				b.showProgress()
				return
			default:
				b.showProgress()
			}
		}
	}()
}

func (b *Bar) showProgress() {
	fmt.Printf(
		"\r[%-*s]%3d%% %8d/%d",
		barWidth, b.barString, b.percent, b.cur, b.total,
	)
}

// NewProxyReader returns a wrapped reader which watches the read-calls.
func (b *Bar) NewProxyReader(r io.Reader) *Reader {
	return &Reader{r, b}
}

func NewBar(total int64) *Bar {
	finishChan := make(chan struct{})
	return &Bar{
		cur:        0,
		percent:    0,
		total:      total,
		barString:  "",
		char:       "#",
		finishChan: finishChan,
	}
}
