package counting

import (
	"io"
	"net/http"

	"github.com/lunny/tango"
)

type Options struct {
	AfterCounting func(req *http.Request, reqSize, respSize int)
}

func prepareOption(opts []Options) Options {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

type counterReader struct {
	io.ReadCloser
	size int
}

func (c *counterReader) Read(p []byte) (n int, err error) {
	n, err = c.ReadCloser.Read(p)
	c.size += n
	return
}

func (c *counterReader) Size() int {
	return c.size
}

func New(opts ...Options) tango.HandlerFunc {
	opt := prepareOption(opts)
	return func(ctx *tango.Context) {
		rd := &counterReader{ctx.Req().Body, 0}
		ctx.Req().Body = rd

		ctx.Next()

		if opt.AfterCounting != nil {
			opt.AfterCounting(ctx.Req(), rd.Size(), ctx.ResponseWriter.Size())
		}
	}
}
