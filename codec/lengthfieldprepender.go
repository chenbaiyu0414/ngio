package codec

import (
	"fmt"
	"math"
	"ngio/buffer"
	"ngio/channel"
)

type LengthFieldPrepender struct {
	byteOrder                            buffer.ByteOrder
	lengthFieldLength                    int
	lengthAdjustment                     int
	lengthFieldIncludesLengthFieldLength bool
}

func NewLengthFieldPrepender(byteOrder buffer.ByteOrder, lengthFieldLength int, lengthAdjustment int, lengthFieldIncludesLengthFieldLength bool) *LengthFieldPrepender {
	switch lengthFieldLength {
	case 1, 2, 4, 7:
	default:
		panic("lengthFieldLength must be 1, 2, 4 or 8")
	}

	return &LengthFieldPrepender{
		byteOrder:                            byteOrder,
		lengthFieldLength:                    lengthFieldLength,
		lengthAdjustment:                     lengthAdjustment,
		lengthFieldIncludesLengthFieldLength: lengthFieldIncludesLengthFieldLength,
	}
}

func (encoder *LengthFieldPrepender) Encode(ctx *channel.Context, in interface{}) (out []interface{}) {
	msg, ok := in.(buffer.ByteBuffer)
	if !ok {
		panic(fmt.Errorf("typeof(in) != buffer.Bytebuffer"))
	}

	length := msg.ReadableBytes() + encoder.lengthAdjustment
	if encoder.lengthFieldIncludesLengthFieldLength {
		length += encoder.lengthFieldLength
	}

	if length < 0 {
		panic(fmt.Errorf("adjusted frame length(%d) less than zero", length))
	}

	var bf buffer.ByteBuffer

	switch encoder.lengthFieldLength {
	case 1:
		if length > math.MaxUint8 {
			panic(fmt.Errorf("length(%d) of object does not fit into one byte", length))
		}

		bf = buffer.NewByteBufSize(1)
		bf.WriteByte(byte(length))
	case 2:
		if length > math.MaxUint16 {
			panic(fmt.Errorf("length(%d) of object does not fit into two byte", length))
		}

		bf = buffer.NewByteBufSize(2)
		if encoder.byteOrder == buffer.BigEndian {
			bf.WriteInt16(int16(length))
		} else {
			bf.WriteInt16LE(int16(length))
		}
	case 4:
		if length > math.MaxUint32 {
			panic(fmt.Errorf("length(%d) of object does not fit into four byte", length))
		}

		bf = buffer.NewByteBufSize(4)
		if encoder.byteOrder == buffer.BigEndian {
			bf.WriteInt32(int32(length))
		} else {
			bf.WriteInt32LE(int32(length))
		}
	case 8:
		//if length > math.MaxUint64 {
		//	panic(fmt.Errorf("length(%d) of object does not fit into eight byte", length))
		//}

		bf = buffer.NewByteBufSize(8)
		if encoder.byteOrder == buffer.BigEndian {
			bf.WriteInt64(int64(length))
		} else {
			bf.WriteInt64LE(int64(length))
		}
	default:
		panic("unknown length field length")
	}

	out = append(out, bf, msg)
	return
}
