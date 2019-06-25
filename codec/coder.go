package codec

import (
	"ngio"
	"ngio/buffer"
)

type MessageToByteEncoder interface {
	Encode(ctx ngio.Context, in interface{}) buffer.ByteBuffer
}

type ByteToMessageDecoder interface {
	Decode(ctx ngio.Context, in buffer.ByteBuffer) interface{}
}

type MessageToMessageEncoder interface {
	Encode(ctx ngio.Context, in interface{}) []interface{}
}

type MessageToMessageDecoder interface {
	Decode(ctx ngio.Context, in interface{}) []interface{}
}
