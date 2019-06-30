package codec

import (
	"ngio/buffer"
	"ngio/channel"
)

type MessageToByteEncoder interface {
	Encode(ctx channel.Context, in interface{}) buffer.ByteBuffer
}

type ByteToMessageDecoder interface {
	Decode(ctx channel.Context, in buffer.ByteBuffer) interface{}
}

type MessageToMessageEncoder interface {
	Encode(ctx channel.Context, in interface{}) []interface{}
}

type MessageToMessageDecoder interface {
	Decode(ctx channel.Context, in interface{}) []interface{}
}
