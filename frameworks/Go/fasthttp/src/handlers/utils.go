package handlers

import (
	"bufio"
	"encoding/json"
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

var encoderPool = sync.Pool{New: func() interface{} { return newJsonReusableEncode(io.Discard) }}

type jsonReusableEncoder struct {
	buffer  *bufio.Writer
	encoder *json.Encoder
}

func newJsonReusableEncode(w io.Writer) *jsonReusableEncoder {
	b := bufio.NewWriter(w)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	return &jsonReusableEncoder{buffer: b, encoder: enc}
}

func (j *jsonReusableEncoder) Encode(v interface{}) error {
	return j.encoder.Encode(v)
}

func (j *jsonReusableEncoder) Reset(w io.Writer) {
	j.buffer.Reset(w)
}

func acquireJsonEncoder(w io.Writer) *jsonReusableEncoder {
	j := encoderPool.Get().(*jsonReusableEncoder)
	j.Reset(w)
	return j
}

func releaseJsonEncoder(x *jsonReusableEncoder) {
	x.buffer.Flush()
	x.Reset(io.Discard)
	encoderPool.Put(x)
}
