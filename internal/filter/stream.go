package filter

import (
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// OpenFile opens a log file for reading, transparently decompressing .gz files.
func OpenFile(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		return &gzipReadCloser{gz: gr, f: f}, nil
	}

	return f, nil
}

// gzipReadCloser wraps a gzip.Reader and the underlying file so both are closed.
type gzipReadCloser struct {
	gz *gzip.Reader
	f  *os.File
}

func (g *gzipReadCloser) Read(p []byte) (int, error) {
	return g.gz.Read(p)
}

func (g *gzipReadCloser) Close() error {
	gzErr := g.gz.Close()
	fErr := g.f.Close()
	if gzErr != nil {
		return gzErr
	}
	return fErr
}

// FilterFile is a convenience wrapper that opens path and runs Filter.
func FilterFile(path string, w io.Writer, opts Options) (Result, error) {
	r, err := OpenFile(path)
	if err != nil {
		return Result{}, err
	}
	defer r.Close()

	return Filter(r, w, opts)
}
