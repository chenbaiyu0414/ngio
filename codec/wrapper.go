package codec

import (
	"ngio/buffer"
	"ngio/channel"
)

type MessageToByteEncoder interface {
	Encode(ctx channel.Context, in interface{}) buffer.ByteBuffer
}

type MessageToMessageEncoder interface {
	Encode(ctx channel.Context, in interface{}) []interface{}
}

type ByteToMessageDecoder interface {
	Decode(ctx channel.Context, in buffer.ByteBuffer) interface{}
}

type ByteToMessageDecoderWrapper struct {
	decoder  ByteToMessageDecoder
	remained buffer.ByteBuffer // store last remaining ByteBuffer buf
}

func NewByteToMessageDecoderWrapper(decoder ByteToMessageDecoder) *ByteToMessageDecoderWrapper {
	if decoder == nil {
		panic("decoder is nil")
	}

	return &ByteToMessageDecoderWrapper{
		decoder: decoder,
	}
}

func (wrapper *ByteToMessageDecoderWrapper) ChannelRead(ctx channel.Context, in interface{}) {
	r, ok := in.(buffer.ByteBuffer)

	if !ok {
		ctx.FireReadHandler(in)
		return
	}

	if wrapper.remained == nil {
		wrapper.remained = r
	} else {
		wrapper.remained.WriteSlice(r)
	}

	for wrapper.remained.ReadableBytes() > 0 {
		out := wrapper.decoder.Decode(ctx, wrapper.remained)

		if out != nil {
			ctx.FireReadHandler(out)
		} else {
			break
		}
	}

	if wrapper.remained.ReadableBytes() == 0 {
		wrapper.remained = nil
	}
}

type MessageToByteEncoderWrapper struct {
	encoder MessageToByteEncoder
}

func NewMessageToByteEncoderWrapper(encoder MessageToByteEncoder) *MessageToByteEncoderWrapper {
	if encoder == nil {
		panic("encoder is nil")
	}

	return &MessageToByteEncoderWrapper{
		encoder: encoder,
	}
}

func (wrapper *MessageToByteEncoderWrapper) Write(ctx channel.Context, msg interface{}) {
	out := wrapper.encoder.Encode(ctx, msg)
	if out != nil {
		ctx.Write(out)
	}
}

type MessageToMessageEncoderWrapper struct {
	encoder MessageToMessageEncoder
}

func NewMessageToMessageEncoderWrapper(encoder MessageToMessageEncoder) *MessageToMessageEncoderWrapper {
	if encoder == nil {
		panic("encoder is nil")
	}

	return &MessageToMessageEncoderWrapper{
		encoder: encoder,
	}
}

func (wrapper *MessageToMessageEncoderWrapper) Write(ctx channel.Context, msg interface{}) {
	outs := wrapper.encoder.Encode(ctx, msg)
	if len(outs) > 0 {
		for _, out := range outs {
			ctx.Write(out)
		}
	} else {
		panic("at least produce one message")
	}
}
