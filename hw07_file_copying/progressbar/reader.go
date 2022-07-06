package progressbar

import "io"

// Reader is a wrapper around the std reader with a Bar embedded.
type Reader struct {
	io.Reader
	bar *Bar
}

// Read spies on the Reader and adds to the Bar.
func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err != nil {
		return 0, err
	}
	r.bar.Progress(int64(n))
	return n, nil
}
