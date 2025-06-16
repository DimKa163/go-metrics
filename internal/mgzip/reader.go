package mgzip

import (
	"compress/gzip"
	"io"
)

type GzipReader struct {
	reader io.ReadCloser
	gz     *gzip.Reader
}

func NewGZIPReader(r io.ReadCloser) (*GzipReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		reader: r,
		gz:     gz,
	}, nil
}

func (g *GzipReader) Read(p []byte) (n int, err error) {
	return g.gz.Read(p)
}

func (g *GzipReader) Close() error {
	if err := g.reader.Close(); err != nil {
		return err
	}
	return g.gz.Close()
}
