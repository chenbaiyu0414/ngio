package byteutil

import "math"

// little-endian helper func

//GetInt16LE returns little-endian coded int16 value from specified buf and original index.
func GetInt16LE(buf []byte, index int) (v int16) {
	_ = buf[index+1] // bounds check
	v |= int16(buf[index+0])
	v |= int16(buf[index+1]) << 8
	return
}

//GetUint16LE returns little-endian coded uint16 value from specified buf and original index.
func GetUint16LE(buf []byte, index int) (v uint16) {
	_ = buf[index+1] // bounds check
	v |= uint16(buf[index+0])
	v |= uint16(buf[index+1]) << 8
	return
}

//GetInt32LE returns little-endian coded int32 value from specified buf and original index.
func GetInt32LE(buf []byte, index int) (v int32) {
	_ = buf[index+3] // bounds check
	v |= int32(buf[index+0])
	v |= int32(buf[index+1]) << 8
	v |= int32(buf[index+2]) << 16
	v |= int32(buf[index+3]) << 24
	return
}

//GetUint32LE returns little-endian coded uint32 value from specified buf and original index.
func GetUint32LE(buf []byte, index int) (v uint32) {
	_ = buf[index+3] // bounds check
	v |= uint32(buf[index+0])
	v |= uint32(buf[index+1]) << 8
	v |= uint32(buf[index+2]) << 16
	v |= uint32(buf[index+3]) << 24
	return
}

//GetInt64LE returns little-endian coded int64 value from specified buf and original index.
func GetInt64LE(buf []byte, index int) (v int64) {
	_ = buf[index+7] // bounds check
	v |= int64(buf[index+0])
	v |= int64(buf[index+1]) << 8
	v |= int64(buf[index+2]) << 16
	v |= int64(buf[index+3]) << 24
	v |= int64(buf[index+4]) << 32
	v |= int64(buf[index+5]) << 40
	v |= int64(buf[index+6]) << 48
	v |= int64(buf[index+7]) << 56
	return
}

//GetUint64LE returns little-endian coded uint64 value from specified buf and original index.
func GetUint64LE(buf []byte, index int) (v uint64) {
	_ = buf[index+7] // bounds check
	v |= uint64(buf[index+0])
	v |= uint64(buf[index+1]) << 8
	v |= uint64(buf[index+2]) << 16
	v |= uint64(buf[index+3]) << 24
	v |= uint64(buf[index+4]) << 32
	v |= uint64(buf[index+5]) << 40
	v |= uint64(buf[index+6]) << 48
	v |= uint64(buf[index+7]) << 56
	return
}

//GetFloat32LE returns little-endian coded float32 value from specified buf and original index.
func GetFloat32LE(buf []byte, index int) float32 {
	return math.Float32frombits(GetUint32LE(buf, index))
}

//GetFloat64LE returns little-endian coded float64 value from specified buf and original index.
func GetFloat64LE(buf []byte, index int) float64 {
	return math.Float64frombits(GetUint64LE(buf, index))
}

//SetInt16LE sets little-endian coded int16 value into specified buf with original index.
func SetInt16LE(buf []byte, index int, v int16) {
	_ = buf[index+1] // bounds check
	buf[index+0] = byte(v)
	buf[index+1] = byte(v >> 8)
}

//SetUint16LE sets little-endian coded uint16 value into specified buf with original index.
func SetUint16LE(buf []byte, index int, v uint16) {
	_ = buf[index+1] // bounds check
	buf[index+0] = byte(v)
	buf[index+1] = byte(v >> 8)
}

//SetInt32LE sets little-endian coded int32 value into specified buf with original index.
func SetInt32LE(buf []byte, index int, v int32) {
	_ = buf[index+3] // bounds check
	buf[index+0] = byte(v)
	buf[index+1] = byte(v >> 8)
	buf[index+2] = byte(v >> 16)
	buf[index+3] = byte(v >> 24)
}

//SetUint32LE sets little-endian coded uint32 value into specified buf with original index.
func SetUint32LE(buf []byte, index int, v uint32) {
	_ = buf[index+3] // bounds check
	buf[index+0] = byte(v)
	buf[index+1] = byte(v >> 8)
	buf[index+2] = byte(v >> 16)
	buf[index+3] = byte(v >> 24)
}

//SetInt64LE sets little-endian coded int64 value into specified buf with original index.
func SetInt64LE(buf []byte, index int, v int64) {
	_ = buf[index+7] // bounds check
	buf[index+0] = byte(v)
	buf[index+1] = byte(v >> 8)
	buf[index+2] = byte(v >> 16)
	buf[index+3] = byte(v >> 24)
	buf[index+4] = byte(v >> 32)
	buf[index+5] = byte(v >> 40)
	buf[index+6] = byte(v >> 48)
	buf[index+7] = byte(v >> 56)
}

//SetUint64LE sets little-endian coded uint64 value into specified buf with original index.
func SetUint64LE(buf []byte, index int, v uint64) {
	_ = buf[index+7] // bounds check
	buf[index+0] = byte(v)
	buf[index+1] = byte(v >> 8)
	buf[index+2] = byte(v >> 16)
	buf[index+3] = byte(v >> 24)
	buf[index+4] = byte(v >> 32)
	buf[index+5] = byte(v >> 40)
	buf[index+6] = byte(v >> 48)
	buf[index+7] = byte(v >> 56)
}

//SetFloat32LE sets little-endian coded float32 value into specified buf with original index.
func SetFloat32LE(buf []byte, index int, v float32) {
	SetUint32LE(buf, index, math.Float32bits(v))
}

//SetFloat64LE sets little-endian coded float64 value into specified buf with original index.
func SetFloat64LE(buf []byte, index int, v float64) {
	SetUint64LE(buf, index, math.Float64bits(v))
}

// big-endian helper func

//GetInt16 returns big-endian coded int16 value from specified buf and original index.
func GetInt16(buf []byte, index int) (v int16) {
	_ = buf[index+1] // bounds check
	v |= int16(buf[index+0]) << 8
	v |= int16(buf[index+1])
	return
}

//GetUint16 returns big-endian coded uint16 value from specified buf and original index.
func GetUint16(buf []byte, index int) (v uint16) {
	_ = buf[index+1] // bounds check
	v |= uint16(buf[index+0]) << 8
	v |= uint16(buf[index+1])
	return
}

//GetInt32 returns big-endian coded int32 value from specified buf and original index.
func GetInt32(buf []byte, index int) (v int32) {
	_ = buf[index+3] // bounds check
	v |= int32(buf[index+0]) << 24
	v |= int32(buf[index+1]) << 16
	v |= int32(buf[index+2]) << 8
	v |= int32(buf[index+3])
	return
}

//GetUint32 returns big-endian coded uint32 value from specified buf and original index.
func GetUint32(buf []byte, index int) (v uint32) {
	_ = buf[index+3] // bounds check
	v |= uint32(buf[index+0]) << 24
	v |= uint32(buf[index+1]) << 16
	v |= uint32(buf[index+2]) << 8
	v |= uint32(buf[index+3])
	return
}

//GetInt64 returns big-endian coded int64 value from specified buf and original index.
func GetInt64(buf []byte, index int) (v int64) {
	_ = buf[index+7] // bounds check
	v |= int64(buf[index+0]) << 56
	v |= int64(buf[index+1]) << 48
	v |= int64(buf[index+2]) << 40
	v |= int64(buf[index+3]) << 32
	v |= int64(buf[index+4]) << 24
	v |= int64(buf[index+5]) << 16
	v |= int64(buf[index+6]) << 8
	v |= int64(buf[index+7])
	return
}

//GetUint64 returns big-endian coded uint64 value from specified buf and original index.
func GetUint64(buf []byte, index int) (v uint64) {
	_ = buf[index+7] // bounds check
	v |= uint64(buf[index+0]) << 56
	v |= uint64(buf[index+1]) << 48
	v |= uint64(buf[index+2]) << 40
	v |= uint64(buf[index+3]) << 32
	v |= uint64(buf[index+4]) << 24
	v |= uint64(buf[index+5]) << 16
	v |= uint64(buf[index+6]) << 8
	v |= uint64(buf[index+7])
	return
}

//GetFloat32 returns big-endian coded float32 value from specified buf and original index.
func GetFloat32(buf []byte, index int) float32 {
	return math.Float32frombits(GetUint32(buf, index))
}

//GetFloat64 returns big-endian coded float64 value from specified buf and original index.
func GetFloat64(buf []byte, index int) float64 {
	return math.Float64frombits(GetUint64(buf, index))
}

//SetInt16 sets big-endian coded int16 value into specified buf with original index.
func SetInt16(buf []byte, index int, v int16) {
	_ = buf[index+1] // bounds check
	buf[index+0] = byte(v >> 8)
	buf[index+1] = byte(v)
}

//SetUint16 sets big-endian coded uint16 value into specified buf with original index.
func SetUint16(buf []byte, index int, v uint16) {
	_ = buf[index+1] // bounds check
	buf[index+0] = byte(v >> 8)
	buf[index+1] = byte(v)
}

//SetInt32 sets big-endian coded int32 value into specified buf with original index.
func SetInt32(buf []byte, index int, v int32) {
	_ = buf[index+3] // bounds check
	buf[index+0] = byte(v >> 24)
	buf[index+1] = byte(v >> 16)
	buf[index+2] = byte(v >> 8)
	buf[index+3] = byte(v)
}

//SetUint32 sets big-endian coded uint32 value into specified buf with original index.
func SetUint32(buf []byte, index int, v uint32) {
	_ = buf[index+3] // bounds check
	buf[index+0] = byte(v >> 24)
	buf[index+1] = byte(v >> 16)
	buf[index+2] = byte(v >> 8)
	buf[index+3] = byte(v)
}

//SetInt64 sets big-endian coded int64 value into specified buf with original index.
func SetInt64(buf []byte, index int, v int64) {
	_ = buf[index+7] // bounds check
	buf[index+0] = byte(v >> 56)
	buf[index+1] = byte(v >> 48)
	buf[index+2] = byte(v >> 40)
	buf[index+3] = byte(v >> 32)
	buf[index+4] = byte(v >> 24)
	buf[index+5] = byte(v >> 16)
	buf[index+6] = byte(v >> 8)
	buf[index+7] = byte(v)
}

//SetUint64 sets big-endian coded uint64 value into specified buf with original index.
func SetUint64(buf []byte, index int, v uint64) {
	_ = buf[index+7] // bounds check
	buf[index+0] = byte(v >> 56)
	buf[index+1] = byte(v >> 48)
	buf[index+2] = byte(v >> 40)
	buf[index+3] = byte(v >> 32)
	buf[index+4] = byte(v >> 24)
	buf[index+5] = byte(v >> 16)
	buf[index+6] = byte(v >> 8)
	buf[index+7] = byte(v)
}

//SetFloat32 sets big-endian coded float32 value into specified buf with original index.
func SetFloat32(buf []byte, index int, v float32) {
	SetUint32(buf, index, math.Float32bits(v))
}

//SetFloat64 sets big-endian coded float64 value into specified buf with original index.
func SetFloat64(buf []byte, index int, v float64) {
	SetUint64(buf, index, math.Float64bits(v))
}
