package handlers

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	rand "github.com/maxim-kuderko/fast-random"
	"github.com/valyala/fasthttp"
	"io"
	"sync"
)

func queriesParam(ctx *fasthttp.RequestCtx) int {
	n := ctx.Request.URI().QueryArgs().GetUintOrZero("queries")
	if n < 1 {
		n = 1
	} else if n > maxWorlds {
		n = maxWorlds
	}

	return n
}

func randomWorldNum() int {
	return rand.Intn(worldsCount) + 1
}

var encoderPool = sync.Pool{New: func() interface{} { return newJsonReusableEncode() }}

type jsonReusableEncoder struct {
	buffer  *bytes.Buffer
	encoder *jsoniter.Encoder
}

func newJsonReusableEncode() *jsonReusableEncoder {
	b := bytes.NewBuffer(nil)
	return &jsonReusableEncoder{buffer: b, encoder: jsoniter.NewEncoder(b)}
}

func (j *jsonReusableEncoder) Encode(v interface{}) error {
	return j.encoder.Encode(v)
}

func (j *jsonReusableEncoder) WriteTo(w io.Writer) (int64, error) {
	return j.buffer.WriteTo(w)
}
func (j *jsonReusableEncoder) Reset() {
	j.buffer.Reset()
}

func acquireJsonEncoder() *jsonReusableEncoder {
	return encoderPool.Get().(*jsonReusableEncoder)
}

func releaseJsonEncoder(x *jsonReusableEncoder) {
	x.Reset()
	encoderPool.Put(x)
}
