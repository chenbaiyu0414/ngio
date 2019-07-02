package codec

import (
	"bytes"
	"fmt"
	"ngio/buffer"
	"ngio/channel"
)

type LineBasedFrameDecoder struct {
	maxLength      int
	stripDelimiter bool
}

func NewLineBasedFrameDecoder(maxLength int, stripDelimiter bool) *LineBasedFrameDecoder {
	return &LineBasedFrameDecoder{
		maxLength:      maxLength,
		stripDelimiter: stripDelimiter,
	}
}

func (decoder *LineBasedFrameDecoder) Decode(ctx *channel.Context, in buffer.ByteBuffer) (out interface{}) {
	delimiterLength := 1

	i := bytes.IndexByte(in.GetBytes(in.ReaderIndex(), in.ReadableBytes()), '\n')
	if i == -1 {
		return
	}

	if i > 0 && in.GetByte(in.ReaderIndex()+i-1) == '\r' {
		delimiterLength = 2
		i--
	}

	if i > decoder.maxLength {
		in.Skip(i + delimiterLength)
		ctx.FireChannelErrorHandler(fmt.Errorf("[LineBasedFrameDecoder] max length exceeds"))
		return
	}

	if decoder.stripDelimiter {
		out = in.ReadSlice(i)
		in.Skip(delimiterLength)
	} else {
		out = in.ReadSlice(i + delimiterLength)
	}

	return
}
