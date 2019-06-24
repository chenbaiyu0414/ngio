package buffer

import (
	"fmt"
	"io"
	"math"
	"ngio/buffer/internal/byteutil"
)

const (
	//SizeBool is the size of bool.
	SizeBool = 1

	//SizeByte is the size of byte.
	SizeByte = 1

	//SizeInt8 is the size of int8.
	SizeInt8 = 1

	//SizeInt16 is the size of int16.
	SizeInt16 = 2

	//SizeInt32 is the size of int32.
	SizeInt32 = 4

	//SizeInt64 is the size of int64.
	SizeInt64 = 8

	//SizeUint8 is the size of uint8.
	SizeUint8 = 1

	//SizeUint16 is the size of uint16.
	SizeUint16 = 2

	//SizeUint32 is the size of uint32.
	SizeUint32 = 4

	//SizeUint64 is the size of uint64.
	SizeUint64 = 8

	//SizeFloat32 is the size of float32.
	SizeFloat32 = 4

	//SizeFloat64 is the size of float64.
	SizeFloat64 = 8
)

type ByteBuffer interface {
	Capacity() int

	ReadableBytes() int
	ReaderIndex() int
	SetReaderIndex(index int)
	MarkReaderIndex()
	ResetReaderIndex()

	WritableBytes() int
	WriterIndex() int
	SetWriterIndex(index int)
	MarkWriterIndex()
	ResetWriterIndex()

	WriteByte(v byte)
	WriteInt8(v int8)
	WriteUint8(v uint8)
	WriteBool(v bool)
	WriteInt16(v int16)
	WriteInt16LE(v int16)
	WriteUint16(v uint16)
	WriteUint16LE(v uint16)
	WriteInt32(v int32)
	WriteInt32LE(v int32)
	WriteUint32(v uint32)
	WriteUint32LE(v uint32)
	WriteInt64(v int64)
	WriteInt64LE(v int64)
	WriteUint64(v uint64)
	WriteUint64LE(v uint64)
	WriteFloat32(v float32)
	WriteFloat32LE(v float32)
	WriteFloat64(v float64)
	WriteFloat64LE(v float64)
	WriteBytes(buffer []byte)
	WriteSlice(buffer ByteBuffer)
	io.WriterTo

	ReadByte() byte
	ReadInt8() int8
	ReadUint8() uint8
	ReadBool() bool
	ReadInt16() int16
	ReadInt16LE() int16
	ReadUint16() uint16
	ReadUint16LE() uint16
	ReadInt32() int32
	ReadInt32LE() int32
	ReadUint32() uint32
	ReadUint32LE() uint32
	ReadInt64() int64
	ReadInt64LE() int64
	ReadUint64() uint64
	ReadUint64LE() uint64
	ReadFloat32() float32
	ReadFloat32LE() float32
	ReadFloat64() float64
	ReadFloat64LE() float64
	ReadBytes(length int) []byte
	ReadSlice(length int) ByteBuffer
	io.ReaderFrom

	GetByte(index int) byte
	GetInt8(index int) int8
	GetUint8(index int) uint8
	GetBool(index int) bool
	GetInt16(index int) int16
	GetInt16LE(index int) int16
	GetUint16(index int) uint16
	GetUint16LE(index int) uint16
	GetInt32(index int) int32
	GetInt32LE(index int) int32
	GetUint32(index int) uint32
	GetUint32LE(index int) uint32
	GetInt64(index int) int64
	GetInt64LE(index int) int64
	GetUint64(index int) uint64
	GetUint64LE(index int) uint64
	GetFloat32(index int) float32
	GetFloat32LE(index int) float32
	GetFloat64(index int) float64
	GetFloat64LE(index int) float64
	GetBytes(index, length int) []byte

	SetByte(index int, v byte)
	SetInt8(index int, v int8)
	SetUint8(index int, v uint8)
	SetBool(index int, v bool)
	SetInt16(index int, v int16)
	SetInt16LE(index int, v int16)
	SetUint16(index int, v uint16)
	SetUint16LE(index int, v uint16)
	SetInt32(index int, v int32)
	SetInt32LE(index int, v int32)
	SetUint32(index int, v uint32)
	SetUint32LE(index int, v uint32)
	SetInt64(index int, v int64)
	SetInt64LE(index int, v int64)
	SetUint64(index int, v uint64)
	SetUint64LE(index int, v uint64)
	SetFloat32(index int, v float32)
	SetFloat32LE(index int, v float32)
	SetFloat64(index int, v float64)
	SetFloat64LE(index int, v float64)
	SetBytes(index int, v []byte)

	Skip(length int)
	DiscardReadBytes()
	Buffer() []byte
	Clear()
	String() string
}

type ByteBuf struct {
	// underlying bytes buffer
	buf []byte

	// capacity of buf
	cap int

	// buf read and write index, expect 0 <= r <= w <= len(buf)
	r, w int

	// marked read and writ index
	mr, mw int
}

func NewByteBufSize(size int) *ByteBuf {
	return NewByteBuf(make([]byte, size))
}

func NewByteBuf(buffer []byte) *ByteBuf {
	return newByteBuf(buffer, 0, 0)
}

func newByteBuf(buffer []byte, r, w int) *ByteBuf {
	capacity := cap(buffer)

	if capacity == 0 {
		panic("new byte buf error")
	}

	return &ByteBuf{
		buf: buffer,
		r:   r,
		w:   w,
		cap: capacity,
	}
}

func (bf *ByteBuf) Capacity() int {
	return bf.cap
}

func (bf *ByteBuf) ReadableBytes() int {
	return bf.w - bf.r
}

func (bf *ByteBuf) ReaderIndex() int {
	return bf.r
}

func (bf *ByteBuf) SetReaderIndex(index int) {
	if index < 0 || index > bf.w {
		panic(fmt.Errorf("reader index out of range. set:%d (except: 0 <= readerIndex <= writerIndex(%d))", index, bf.w))
	}

	bf.r = index
}

func (bf *ByteBuf) MarkReaderIndex() {
	bf.mr = bf.r
}

func (bf *ByteBuf) ResetReaderIndex() {
	bf.SetReaderIndex(bf.mr)
}

func (bf *ByteBuf) WritableBytes() int {
	return bf.cap - bf.w
}

func (bf *ByteBuf) WriterIndex() int {
	return bf.w
}

func (bf *ByteBuf) SetWriterIndex(index int) {
	if index < bf.r || index > bf.cap {
		panic(fmt.Errorf("writer index out of range. set:%d (except: 0 <= readerIndex(%d) <= writerIndex <= capacity(%d))", index, bf.r, bf.cap))
	}

	bf.w = index
}

func (bf *ByteBuf) MarkWriterIndex() {
	bf.mw = bf.w
}

func (bf *ByteBuf) ResetWriterIndex() {
	bf.SetWriterIndex(bf.mw)
}

func (bf *ByteBuf) WriteByte(v byte) {
	bf.ensureWritableBytes(SizeByte)
	bf.buf[bf.w] = v
	bf.w += SizeByte
}

func (bf *ByteBuf) WriteInt8(v int8) {
	bf.ensureWritableBytes(SizeInt8)
	bf.buf[bf.w] = byte(v)
	bf.w += SizeInt8
}

func (bf *ByteBuf) WriteUint8(v uint8) {
	bf.ensureWritableBytes(SizeUint8)
	bf.buf[bf.w] = v
	bf.w += SizeUint8
}

func (bf *ByteBuf) WriteBool(v bool) {
	bf.ensureWritableBytes(SizeBool)
	if v {
		bf.buf[bf.w] = 1
	} else {
		bf.buf[bf.w] = 0
	}
	bf.w += SizeBool
}

func (bf *ByteBuf) WriteInt16(v int16) {
	bf.ensureWritableBytes(SizeInt16)
	byteutil.SetInt16(bf.buf, bf.w, v)
	bf.w += SizeInt16
}

func (bf *ByteBuf) WriteInt16LE(v int16) {
	bf.ensureWritableBytes(SizeInt16)
	byteutil.SetInt16LE(bf.buf, bf.w, v)
	bf.w += SizeInt16
}

func (bf *ByteBuf) WriteUint16(v uint16) {
	bf.ensureWritableBytes(SizeUint16)
	byteutil.SetUint16(bf.buf, bf.w, v)
	bf.w += SizeUint16
}

func (bf *ByteBuf) WriteUint16LE(v uint16) {
	bf.ensureWritableBytes(SizeUint16)
	byteutil.SetUint16LE(bf.buf, bf.w, v)
	bf.w += SizeUint16
}

func (bf *ByteBuf) WriteInt32(v int32) {
	bf.ensureWritableBytes(SizeInt32)
	byteutil.SetInt32(bf.buf, bf.w, v)
	bf.w += SizeInt32
}

func (bf *ByteBuf) WriteInt32LE(v int32) {
	bf.ensureWritableBytes(SizeInt32)
	byteutil.SetInt32LE(bf.buf, bf.w, v)
	bf.w += SizeInt32
}

func (bf *ByteBuf) WriteUint32(v uint32) {
	bf.ensureWritableBytes(SizeUint32)
	byteutil.SetUint32(bf.buf, bf.w, v)
	bf.w += SizeUint32
}

func (bf *ByteBuf) WriteUint32LE(v uint32) {
	bf.ensureWritableBytes(SizeUint32)
	byteutil.SetUint32LE(bf.buf, bf.w, v)
	bf.w += SizeUint32
}

func (bf *ByteBuf) WriteInt64(v int64) {
	bf.ensureWritableBytes(SizeInt64)
	byteutil.SetInt64(bf.buf, bf.w, v)
	bf.w += SizeInt64
}

func (bf *ByteBuf) WriteInt64LE(v int64) {
	bf.ensureWritableBytes(SizeInt64)
	byteutil.SetInt64LE(bf.buf, bf.w, v)
	bf.w += SizeInt64
}

func (bf *ByteBuf) WriteUint64(v uint64) {
	bf.ensureWritableBytes(SizeUint64)
	byteutil.SetUint64(bf.buf, bf.w, v)
	bf.w += SizeUint64
}

func (bf *ByteBuf) WriteUint64LE(v uint64) {
	bf.ensureWritableBytes(SizeUint64)
	byteutil.SetUint64LE(bf.buf, bf.w, v)
	bf.w += SizeUint64
}

func (bf *ByteBuf) WriteFloat32(v float32) {
	bf.ensureWritableBytes(SizeFloat32)
	byteutil.SetFloat32(bf.buf, bf.w, v)
	bf.w += SizeFloat32
}

func (bf *ByteBuf) WriteFloat32LE(v float32) {
	bf.ensureWritableBytes(SizeFloat32)
	byteutil.SetFloat32LE(bf.buf, bf.w, v)
	bf.w += SizeFloat32
}

func (bf *ByteBuf) WriteFloat64(v float64) {
	bf.ensureWritableBytes(SizeFloat64)
	byteutil.SetFloat64(bf.buf, bf.w, v)
	bf.w += SizeFloat64
}

func (bf *ByteBuf) WriteFloat64LE(v float64) {
	bf.ensureWritableBytes(SizeFloat64)
	byteutil.SetFloat64LE(bf.buf, bf.w, v)
	bf.w += SizeFloat64
}

func (bf *ByteBuf) WriteBytes(buffer []byte) {
	bf.ensureWritableBytes(len(buffer))
	copy(bf.buf[bf.w:], buffer)
	bf.w += len(buffer)
}

func (bf *ByteBuf) WriteSlice(buffer ByteBuffer) {
	bf.WriteBytes(buffer.GetBytes(buffer.ReaderIndex(), buffer.ReadableBytes()))
}

// WriteTo writes all readable bytes into io.Writer. It returns the number of bytes written.
// And the reader index will increased n.
func (bf *ByteBuf) WriteTo(w io.Writer) (n int64, err error) {
	if readable := bf.ReadableBytes(); readable > 0 {
		n, err := w.Write(bf.buf[bf.r:bf.w])
		if n > readable {
			panic("buffer.ByteBuf.WriteTo: invalid write count")
		}

		bf.r += n

		if err != nil {
			return int64(n), err
		}

		if n != readable {
			return int64(n), io.ErrShortWrite
		}
	}

	bf.DiscardReadBytes()
	return n, nil
}

func (bf *ByteBuf) ReadByte() (v byte) {
	bf.checkReadableBytes(SizeByte)
	v = bf.buf[bf.r]
	bf.r += SizeByte
	return
}

func (bf *ByteBuf) ReadInt8() (v int8) {
	bf.checkReadableBytes(SizeInt8)
	v = int8(bf.buf[bf.r])
	bf.r += SizeInt8
	return
}

func (bf *ByteBuf) ReadUint8() (v uint8) {
	bf.checkReadableBytes(SizeUint8)
	v = bf.buf[bf.r]
	bf.r += SizeUint8
	return
}

func (bf *ByteBuf) ReadBool() (v bool) {
	bf.checkReadableBytes(SizeBool)
	v = bf.buf[bf.r] != 0
	bf.r += SizeBool
	return
}

func (bf *ByteBuf) ReadInt16() (v int16) {
	bf.checkReadableBytes(SizeInt16)
	v = byteutil.GetInt16(bf.buf, bf.r)
	bf.r += SizeInt16
	return
}

func (bf *ByteBuf) ReadInt16LE() (v int16) {
	bf.checkReadableBytes(SizeInt16)
	v = byteutil.GetInt16LE(bf.buf, bf.r)
	bf.r += SizeInt16
	return
}

func (bf *ByteBuf) ReadUint16() (v uint16) {
	bf.checkReadableBytes(SizeUint16)
	v = byteutil.GetUint16(bf.buf, bf.r)
	bf.r += SizeUint16
	return
}

func (bf *ByteBuf) ReadUint16LE() (v uint16) {
	bf.checkReadableBytes(SizeUint16)
	v = byteutil.GetUint16LE(bf.buf, bf.r)
	bf.r += SizeUint16
	return
}

func (bf *ByteBuf) ReadInt32() (v int32) {
	bf.checkReadableBytes(SizeInt32)
	v = byteutil.GetInt32(bf.buf, bf.r)
	bf.r += SizeInt32
	return
}

func (bf *ByteBuf) ReadInt32LE() (v int32) {
	bf.checkReadableBytes(SizeInt32)
	v = byteutil.GetInt32LE(bf.buf, bf.r)
	bf.r += SizeInt32
	return
}

func (bf *ByteBuf) ReadUint32() (v uint32) {
	bf.checkReadableBytes(SizeUint32)
	v = byteutil.GetUint32(bf.buf, bf.r)
	bf.r += SizeUint32
	return
}

func (bf *ByteBuf) ReadUint32LE() (v uint32) {
	bf.checkReadableBytes(SizeUint32)
	v = byteutil.GetUint32LE(bf.buf, bf.r)
	bf.r += SizeUint32
	return
}

func (bf *ByteBuf) ReadInt64() (v int64) {
	bf.checkReadableBytes(SizeInt64)
	v = byteutil.GetInt64(bf.buf, bf.r)
	bf.r += SizeInt64
	return
}

func (bf *ByteBuf) ReadInt64LE() (v int64) {
	bf.checkReadableBytes(SizeInt64)
	v = byteutil.GetInt64LE(bf.buf, bf.r)
	bf.r += SizeInt64
	return
}

func (bf *ByteBuf) ReadUint64() (v uint64) {
	bf.checkReadableBytes(SizeUint64)
	v = byteutil.GetUint64(bf.buf, bf.r)
	bf.r += SizeUint64
	return
}

func (bf *ByteBuf) ReadUint64LE() (v uint64) {
	bf.checkReadableBytes(SizeUint64)
	v = byteutil.GetUint64LE(bf.buf, bf.r)
	bf.r += SizeUint64
	return
}

func (bf *ByteBuf) ReadFloat32() (v float32) {
	bf.checkReadableBytes(SizeFloat32)
	v = byteutil.GetFloat32(bf.buf, bf.r)
	bf.r += SizeFloat32
	return
}

func (bf *ByteBuf) ReadFloat32LE() (v float32) {
	bf.checkReadableBytes(SizeFloat32)
	v = byteutil.GetFloat32LE(bf.buf, bf.r)
	bf.r += SizeFloat32
	return
}

func (bf *ByteBuf) ReadFloat64() (v float64) {
	bf.checkReadableBytes(SizeFloat64)
	v = byteutil.GetFloat64(bf.buf, bf.r)
	bf.r += SizeFloat64
	return
}

func (bf *ByteBuf) ReadFloat64LE() (v float64) {
	bf.checkReadableBytes(SizeFloat64)
	v = byteutil.GetFloat64LE(bf.buf, bf.r)
	bf.r += SizeFloat64
	return
}

func (bf *ByteBuf) ReadBytes(length int) (v []byte) {
	bf.checkReadableBytes(length)
	v = bf.buf[bf.r : bf.r+length]
	bf.r += length
	return
}

func (bf *ByteBuf) ReadSlice(length int) (v ByteBuffer) {
	v = newByteBuf(bf.ReadBytes(length), 0, 0)
	return
}

// ReadFrom reads bytes from io.Reader into underlying buf. It returns the number of bytes read.
// And the writer index will increased n.
func (bf *ByteBuf) ReadFrom(r io.Reader) (int64, error) {
	writable := bf.WritableBytes()
	if writable == 0 {
		// grow buffer
		bf.grow(calculateNewCapacity(bf.cap, bf.cap))
	}

	n, err := r.Read(bf.buf[bf.w:bf.cap])
	bf.w += n
	return int64(n), err
}

func (bf *ByteBuf) GetByte(index int) byte {
	bf.checkIndex(index, SizeByte)
	return bf.buf[index]
}

func (bf *ByteBuf) GetInt8(index int) int8 {
	bf.checkIndex(index, SizeInt8)
	return int8(bf.buf[index])
}

func (bf *ByteBuf) GetUint8(index int) uint8 {
	bf.checkIndex(index, SizeUint8)
	return uint8(bf.buf[index])
}

func (bf *ByteBuf) GetBool(index int) bool {
	return bf.GetByte(index) != 0
}

func (bf *ByteBuf) GetInt16(index int) int16 {
	bf.checkIndex(index, SizeInt16)
	return byteutil.GetInt16(bf.buf, index)
}

func (bf *ByteBuf) GetInt16LE(index int) int16 {
	bf.checkIndex(index, SizeInt16)
	return byteutil.GetInt16LE(bf.buf, index)
}

func (bf *ByteBuf) GetUint16(index int) uint16 {
	bf.checkIndex(index, SizeUint16)
	return byteutil.GetUint16(bf.buf, index)
}

func (bf *ByteBuf) GetUint16LE(index int) uint16 {
	bf.checkIndex(index, SizeUint16)
	return byteutil.GetUint16LE(bf.buf, index)
}

func (bf *ByteBuf) GetInt32(index int) int32 {
	bf.checkIndex(index, SizeInt32)
	return byteutil.GetInt32(bf.buf, index)
}

func (bf *ByteBuf) GetInt32LE(index int) int32 {
	bf.checkIndex(index, SizeInt32)
	return byteutil.GetInt32LE(bf.buf, index)
}

func (bf *ByteBuf) GetUint32(index int) uint32 {
	bf.checkIndex(index, SizeUint32)
	return byteutil.GetUint32(bf.buf, index)
}

func (bf *ByteBuf) GetUint32LE(index int) uint32 {
	bf.checkIndex(index, SizeUint32)
	return byteutil.GetUint32LE(bf.buf, index)
}

func (bf *ByteBuf) GetInt64(index int) int64 {
	bf.checkIndex(index, SizeInt64)
	return byteutil.GetInt64(bf.buf, index)
}

func (bf *ByteBuf) GetInt64LE(index int) int64 {
	bf.checkIndex(index, SizeInt64)
	return byteutil.GetInt64LE(bf.buf, index)
}

func (bf *ByteBuf) GetUint64(index int) uint64 {
	bf.checkIndex(index, SizeUint64)
	return byteutil.GetUint64(bf.buf, index)
}

func (bf *ByteBuf) GetUint64LE(index int) uint64 {
	bf.checkIndex(index, SizeUint64)
	return byteutil.GetUint64LE(bf.buf, index)
}

func (bf *ByteBuf) GetFloat32(index int) float32 {
	bf.checkIndex(index, SizeFloat32)
	return byteutil.GetFloat32(bf.buf, index)
}

func (bf *ByteBuf) GetFloat32LE(index int) float32 {
	bf.checkIndex(index, SizeFloat32)
	return byteutil.GetFloat32LE(bf.buf, index)
}

func (bf *ByteBuf) GetFloat64(index int) float64 {
	bf.checkIndex(index, SizeFloat64)
	return byteutil.GetFloat64(bf.buf, index)
}

func (bf *ByteBuf) GetFloat64LE(index int) float64 {
	bf.checkIndex(index, SizeFloat64)
	return byteutil.GetFloat64LE(bf.buf, index)
}

func (bf *ByteBuf) GetBytes(index, length int) []byte {
	bf.checkIndex(index, length)
	return bf.buf[index : index+length]
}

func (bf *ByteBuf) SetByte(index int, v byte) {
	bf.checkIndex(index, SizeByte)
	bf.buf[index] = v
}

func (bf *ByteBuf) SetInt8(index int, v int8) {
	bf.checkIndex(index, SizeInt8)
	bf.buf[index] = byte(v)
}

func (bf *ByteBuf) SetUint8(index int, v uint8) {
	bf.checkIndex(index, SizeUint8)
	bf.buf[index] = byte(v)
}

func (bf *ByteBuf) SetBool(index int, v bool) {
	if v {
		bf.SetByte(index, 1)
	} else {
		bf.SetByte(index, 0)
	}
}

func (bf *ByteBuf) SetInt16(index int, v int16) {
	bf.checkIndex(index, SizeInt16)
	byteutil.SetInt16(bf.buf, index, v)
}

func (bf *ByteBuf) SetInt16LE(index int, v int16) {
	bf.checkIndex(index, SizeInt16)
	byteutil.SetInt16LE(bf.buf, index, v)
}

func (bf *ByteBuf) SetUint16(index int, v uint16) {
	bf.checkIndex(index, SizeUint16)
	byteutil.SetUint16(bf.buf, index, v)
}

func (bf *ByteBuf) SetUint16LE(index int, v uint16) {
	bf.checkIndex(index, SizeUint16)
	byteutil.SetUint16LE(bf.buf, index, v)
}

func (bf *ByteBuf) SetInt32(index int, v int32) {
	bf.checkIndex(index, SizeInt32)
	byteutil.SetInt32(bf.buf, index, v)
}

func (bf *ByteBuf) SetInt32LE(index int, v int32) {
	bf.checkIndex(index, SizeInt32)
	byteutil.SetInt32LE(bf.buf, index, v)
}

func (bf *ByteBuf) SetUint32(index int, v uint32) {
	bf.checkIndex(index, SizeUint32)
	byteutil.SetUint32(bf.buf, index, v)
}

func (bf *ByteBuf) SetUint32LE(index int, v uint32) {
	bf.checkIndex(index, SizeUint32)
	byteutil.SetUint32LE(bf.buf, index, v)
}

func (bf *ByteBuf) SetInt64(index int, v int64) {
	bf.checkIndex(index, SizeInt64)
	byteutil.SetInt64(bf.buf, index, v)
}

func (bf *ByteBuf) SetInt64LE(index int, v int64) {
	bf.checkIndex(index, SizeInt64)
	byteutil.SetInt64LE(bf.buf, index, v)
}

func (bf *ByteBuf) SetUint64(index int, v uint64) {
	bf.checkIndex(index, SizeUint64)
	byteutil.SetUint64(bf.buf, index, v)
}

func (bf *ByteBuf) SetUint64LE(index int, v uint64) {
	bf.checkIndex(index, SizeUint64)
	byteutil.SetUint64LE(bf.buf, index, v)
}

func (bf *ByteBuf) SetFloat32(index int, v float32) {
	bf.checkIndex(index, SizeFloat32)
	byteutil.SetFloat32(bf.buf, index, v)
}

func (bf *ByteBuf) SetFloat32LE(index int, v float32) {
	bf.checkIndex(index, SizeFloat32)
	byteutil.SetFloat32LE(bf.buf, index, v)
}

func (bf *ByteBuf) SetFloat64(index int, v float64) {
	bf.checkIndex(index, SizeFloat64)
	byteutil.SetFloat64(bf.buf, index, v)
}

func (bf *ByteBuf) SetFloat64LE(index int, v float64) {
	bf.checkIndex(index, SizeFloat64)
	byteutil.SetFloat64LE(bf.buf, index, v)
}

func (bf *ByteBuf) SetBytes(index int, v []byte) {
	bf.checkIndex(index, len(v))
	copy(bf.buf[index:], v)
}

func (bf *ByteBuf) Skip(length int) {
	bf.checkReadableBytes(length)
	bf.r += length
}

func (bf *ByteBuf) DiscardReadBytes() {
	if bf.r == 0 {
		return
	}

	if bf.r != bf.w {
		copy(bf.buf[0:], bf.buf[bf.r:bf.w])
		bf.adjustMakers(bf.r)
		bf.r = 0
		bf.w -= bf.r
	} else {
		bf.adjustMakers(bf.r)
		bf.r = 0
		bf.w = 0
	}
}

func (bf *ByteBuf) Buffer() []byte {
	return bf.buf
}

func (bf *ByteBuf) Clear() {
	bf.r, bf.w = 0, 0
}

func (bf *ByteBuf) String() string {
	panic("implement me")
}

func (bf *ByteBuf) ensureWritableBytes(fieldLength int) {
	if fieldLength <= bf.WritableBytes() {
		return
	}

	if bf.r > fieldLength {
		bf.DiscardReadBytes()
	} else {
		bf.grow(calculateNewCapacity(bf.w+fieldLength, bf.cap))
	}
}

func (bf *ByteBuf) adjustMakers(decrement int) {
	if bf.mr <= decrement {
		bf.mr = 0
		if bf.mw <= decrement {
			bf.mw = 0
		} else {
			bf.mw -= decrement
		}
	} else {
		bf.mr -= decrement
		bf.mw -= decrement
	}
}

func (bf *ByteBuf) grow(newCapacity int) {
	if newCapacity == 0 {
		return
	}

	buf := make([]byte, newCapacity)
	copy(buf, bf.buf[bf.r:bf.w])

	bf.buf = buf
	bf.cap = newCapacity

	bf.adjustMakers(bf.r)

	bf.w -= bf.r
	bf.r = 0
}

func (bf *ByteBuf) checkReadableBytes(fieldLength int) {
	if bf.r+fieldLength > bf.w {
		panic(fmt.Errorf("index out of range. readerIndex(%d) + length(%d) > writerIndex(%d)", bf.r, fieldLength, bf.w))
	}
}

func (bf *ByteBuf) checkIndex(index, fieldLength int) {
	if index < 0 || fieldLength < 0 || index+fieldLength < 0 || bf.cap-(index+fieldLength) < 0 {
		panic(fmt.Errorf("index out of range. index:%d, length:%d (expected: range[0,%d])", index, fieldLength, bf.cap))
	}
}

const growCapThreshold = 4 * 1024 * 1024 // 4 Mib

func calculateNewCapacity(minNewCapacity, currentCapacity int) int {
	if minNewCapacity > currentCapacity {
		return 0
	}

	if minNewCapacity == growCapThreshold {
		return growCapThreshold
	}

	if minNewCapacity > growCapThreshold {
		return int(math.Float64bits(math.Ceil(float64(minNewCapacity/growCapThreshold)))) * growCapThreshold
	}

	// 如果 minNewCapacity 小于阈值，则直接从64开始翻倍扩容
	newCap := 64
	for newCap < minNewCapacity {
		newCap <<= 1
	}

	return newCap
}
