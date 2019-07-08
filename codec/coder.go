package codec

import (
	"ngio"
	"ngio/buffer"
)

type MessageToByteEncoder interface {
	Encode(ctx ngio.ChannelContext, in interface{}) buffer.ByteBuffer
}

type ByteToMessageDecoder interface {
	Decode(ctx ngio.ChannelContext, in buffer.ByteBuffer) interface{}
}

type MessageToMessageEncoder interface {
	Encode(ctx ngio.ChannelContext, in interface{}) []interface{}
}

type MessageToMessageDecoder interface {
	Decode(ctx ngio.ChannelContext, in interface{}) []interface{}
}
