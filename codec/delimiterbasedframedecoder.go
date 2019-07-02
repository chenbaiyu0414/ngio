package codec

import (
	"bytes"
	"fmt"
	"ngio/buffer"
	"ngio/channel"
)

type DelimiterBasedFrameDecoder struct {
	maxLength      int
	stripDelimiter bool
	delimiters     [][]byte
}

func NewDelimiterBasedFrameDecoder(maxLength int, stripDelimiter bool, delimiters ...[]byte) *DelimiterBasedFrameDecoder {
	return &DelimiterBasedFrameDecoder{
		maxLength:      maxLength,
		stripDelimiter: stripDelimiter,
		delimiters:     delimiters,
	}
}

func (decoder *DelimiterBasedFrameDecoder) Decode(ctx *channel.Context, in buffer.ByteBuffer) (out interface{}) {
	minLength := decoder.maxLength
	shortestDelimLength := 0

	for _, delim := range decoder.delimiters {
		i := bytes.Index(in.GetBytes(in.ReaderIndex(), in.ReadableBytes()), delim)

		if i > 0 && i < minLength {
			minLength = i
			shortestDelimLength = len(delim)
		}
	}

	if shortestDelimLength == 0 {
		return
	}

	if minLength > decoder.maxLength {
		in.Skip(minLength + shortestDelimLength)
		ctx.FireChannelErrorHandler(fmt.Errorf("[DelimiterBasedFrameDecoder] max length exceeds"))
		return
	}

	if decoder.stripDelimiter {
		out = in.ReadSlice(minLength)
		in.Skip(shortestDelimLength)
	} else {
		out = in.ReadSlice(minLength + shortestDelimLength)
	}

	return
}
