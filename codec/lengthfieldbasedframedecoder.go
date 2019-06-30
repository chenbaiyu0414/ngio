package codec

import (
	"fmt"
	"ngio/buffer"
	"ngio/channel"
)

type LengthFieldBasedFrameDecoder struct {
	byteOrder            buffer.ByteOrder
	maxFrameLength       int
	lengthFieldOffset    int
	lengthFieldLength    int
	lengthAdjustment     int
	lengthFieldEndOffset int
	initialBytesToStrip  int
}

func NewLengthFieldBasedFrameDecoder(byteOrder buffer.ByteOrder, maxFrameLength int, lengthFieldOffset int, lengthFieldLength int, lengthAdjustment int, initialBytesToStrip int) *LengthFieldBasedFrameDecoder {
	if maxFrameLength <= 0 {
		panic(fmt.Errorf("out of range. maxFrameLength(%d) must be a positive integer", maxFrameLength))
	}

	if lengthFieldOffset < 0 {
		panic(fmt.Errorf("out of range. lengthFieldOffset(%d) must be a positive integer", maxFrameLength))
	}

	if initialBytesToStrip < 0 {
		panic(fmt.Errorf("out of range. initialBytesToStrip(%d) must be a positive integer", maxFrameLength))
	}

	if lengthFieldOffset+lengthFieldLength > maxFrameLength {
		panic(fmt.Errorf("out of range. expected: maxFrameLength(%d) >= lengthFieldOffset(%d) + lengthFieldLength(%d)", maxFrameLength, lengthFieldOffset, lengthFieldLength))
	}

	return &LengthFieldBasedFrameDecoder{
		byteOrder:            byteOrder,
		maxFrameLength:       maxFrameLength,
		lengthFieldOffset:    lengthFieldOffset,
		lengthFieldLength:    lengthFieldLength,
		lengthAdjustment:     lengthAdjustment,
		lengthFieldEndOffset: lengthFieldOffset + lengthFieldLength,
		initialBytesToStrip:  initialBytesToStrip,
	}
}

func (decoder *LengthFieldBasedFrameDecoder) Decode(ctx channel.Context, in buffer.ByteBuffer) interface{} {
	if in.ReadableBytes() < decoder.lengthFieldEndOffset {
		return nil
	}

	actualLengthFieldOffset := in.ReaderIndex() + decoder.lengthFieldOffset
	frameLength := decoder.GetUnadjustedFrameLength(in, actualLengthFieldOffset, decoder.lengthFieldLength, decoder.byteOrder)

	if frameLength < 0 {
		in.Skip(decoder.lengthFieldEndOffset)
		panic(fmt.Errorf("negative pre-adjustment length field:%d", frameLength))
	}

	frameLength += int64(decoder.lengthAdjustment + decoder.lengthFieldEndOffset)

	if frameLength < int64(decoder.lengthFieldEndOffset) {
		in.Skip(decoder.lengthFieldEndOffset)
		panic(fmt.Errorf("adjusted frame length (%d) is less than lengthFieldEndOffset: %d ", frameLength, decoder.lengthFieldEndOffset))
	}

	if frameLength > int64(decoder.maxFrameLength) {
		discard := frameLength - int64(in.ReadableBytes())

		if discard < 0 {
			// discard all
			in.Skip(int(frameLength))
		} else {
			in.Skip(in.ReadableBytes())
		}

		return nil
	}

	if in.ReadableBytes() < int(frameLength) {
		return nil
	}

	if int64(decoder.initialBytesToStrip) > frameLength {
		in.Skip(int(frameLength))
		panic(fmt.Errorf("adjusted frame length (%d) is less than initialBytesToStrip: %d ", frameLength, decoder.initialBytesToStrip))
	}

	in.Skip(decoder.initialBytesToStrip)

	actualFrameLength := frameLength - int64(decoder.initialBytesToStrip)
	return in.ReadSlice(int(actualFrameLength))
}

func (decoder *LengthFieldBasedFrameDecoder) GetUnadjustedFrameLength(in buffer.ByteBuffer, offset, length int, order buffer.ByteOrder) int64 {
	var frameLength int64

	switch length {
	case 1:
		frameLength = int64(in.GetByte(offset))
	case 2:
		if order == buffer.BigEndian {
			frameLength = int64(in.GetInt16(offset))
		} else {
			frameLength = int64(in.GetInt16LE(offset))
		}
	case 4:
		if order == buffer.BigEndian {
			frameLength = int64(in.GetInt32(offset))
		} else {
			frameLength = int64(in.GetInt32LE(offset))
		}
	case 8:
		if order == buffer.BigEndian {
			frameLength = int64(in.GetInt64(offset))
		} else {
			frameLength = int64(in.GetInt64LE(offset))
		}
	default:
		panic(fmt.Errorf("unsupported lengthFieldLength: %d (expected: 1, 2, 4, or 8)", decoder.lengthFieldLength))
	}

	return frameLength
}
