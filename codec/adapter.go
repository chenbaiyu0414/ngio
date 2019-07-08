package codec

import (
	"errors"
	"ngio"
	"ngio/buffer"
)

var (
	ErrEncoderIsNil             = errors.New("encoder adapter: encoder is nil")
	ErrDecoderIsNil             = errors.New("decoder adapter: decoder is nil")
	ErrAtLeastProduceOneMessage = errors.New("encoder adapter: message to message encoder must at least produce one message")
)

//MessageToByteEncoderAdapter
type MessageToByteEncoderAdapter struct {
	encoder MessageToByteEncoder
}

func NewMessageToByteEncoderAdapter(encoder MessageToByteEncoder) *MessageToByteEncoderAdapter {
	if encoder == nil {
		panic(ErrEncoderIsNil)
	}

	return &MessageToByteEncoderAdapter{
		encoder: encoder,
	}
}

func (adapter *MessageToByteEncoderAdapter) Write(ctx ngio.ChannelContext, msg interface{}) {
	out := adapter.encoder.Encode(ctx, msg)
	if out != nil {
		ctx.Write(out)
	}
}

//ByteToMessageDecoderAdapter
type ByteToMessageDecoderAdapter struct {
	decoder  ByteToMessageDecoder
	remained buffer.ByteBuffer // store last remaining ByteBuffer buf
}

func NewByteToMessageDecoderAdapter(decoder ByteToMessageDecoder) *ByteToMessageDecoderAdapter {
	if decoder == nil {
		panic(ErrDecoderIsNil)
	}

	return &ByteToMessageDecoderAdapter{
		decoder: decoder,
	}
}

func (adapter *ByteToMessageDecoderAdapter) ChannelRead(ctx ngio.ChannelContext, in interface{}) {
	r, ok := in.(buffer.ByteBuffer)

	if !ok {
		ctx.FireChannelReadHandler(in)
		return
	}

	if adapter.remained == nil {
		adapter.remained = r
	} else {
		adapter.remained.WriteSlice(r)
	}

	for adapter.remained.ReadableBytes() > 0 {
		out := adapter.decoder.Decode(ctx, adapter.remained)

		if out != nil {
			ctx.FireChannelReadHandler(out)
		} else {
			break
		}
	}

	if adapter.remained.ReadableBytes() == 0 {
		adapter.remained = nil
	}
}

//MessageToMessageEncoderAdapter
type MessageToMessageEncoderAdapter struct {
	encoder MessageToMessageEncoder
}

func NewMessageToMessageEncoderAdapter(encoder MessageToMessageEncoder) *MessageToMessageEncoderAdapter {
	if encoder == nil {
		panic(ErrEncoderIsNil)
	}

	return &MessageToMessageEncoderAdapter{
		encoder: encoder,
	}
}

func (adapter *MessageToMessageEncoderAdapter) Write(ctx ngio.ChannelContext, msg interface{}) {
	outs := adapter.encoder.Encode(ctx, msg)
	if len(outs) > 0 {
		for _, out := range outs {
			ctx.Write(out)
		}
	} else {
		panic(ErrAtLeastProduceOneMessage)
	}
}

//MessageToMessageDecoderAdapter
type MessageToMessageDecoderAdapter struct {
	decoder MessageToMessageDecoder
}

func NewMessageToMessageDecoderAdapter(decoder MessageToMessageDecoder) *MessageToMessageDecoderAdapter {
	if decoder == nil {
		panic(ErrDecoderIsNil)
	}

	return &MessageToMessageDecoderAdapter{
		decoder: decoder,
	}
}

func (adapter *MessageToMessageDecoderAdapter) ChannelRead(ctx ngio.ChannelContext, msg interface{}) {
	outs := adapter.decoder.Decode(ctx, msg)
	for _, out := range outs {
		ctx.FireChannelReadHandler(out)
	}
}
