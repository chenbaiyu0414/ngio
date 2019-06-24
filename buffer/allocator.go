package buffer

import "ngio/buffer/internal/mathutil"

const (
	DefaultMinimum = 64
	DefaultInitial = 1024
	DefaultMaximum = 65536
	IndexIncrement = 4
	IndexDecrement = 1
)

var sizeTable = [27]int64{
	0x00010, 0x00020, 0x00030, 0x00040, 0x00050, 0x00060, 0x00070, 0x00080, 0x00090,
	0x00100, 0x00110, 0x00120, 0x00130, 0x00140, 0x00150, 0x00160, 0x00170, 0x00180,
	0x00190, 0x00200, 0x00400, 0x00800, 0x01000, 0x02000, 0x04000, 0x08000, 0x10000,
}

func getSizeTableIndex(size int64) int {
	for low, high := 0, len(sizeTable)-1; ; {
		if high < low {
			return low
		}
		if high == low {
			return high
		}

		mid := int(uint(low+high) >> 1)

		a := sizeTable[mid]
		b := sizeTable[mid+1]

		if size > b {
			low = mid + 1
		} else if size < a {
			high = mid - 1
		} else if size == a {
			return mid
		} else {
			return mid + 1
		}
	}
}

type RecvByteBufAllocator struct {
	nextRecvBufferSize        int64
	index, minIndex, maxIndex int
	decreaseNow               bool
}

func NewRecvByteBufAllocator(minimum, maximum, initial int64) *RecvByteBufAllocator {
	var minIndex, maxIndex int

	min := getSizeTableIndex(minimum)
	if sizeTable[min] < minimum {
		minIndex = min + 1
	} else {
		minIndex = min
	}

	max := getSizeTableIndex(maximum)
	if sizeTable[max] > maximum {
		maxIndex = max - 1
	} else {
		maxIndex = max
	}

	initialIndex := getSizeTableIndex(initial)

	return &RecvByteBufAllocator{
		minIndex:           minIndex,
		maxIndex:           maxIndex,
		index:              initialIndex,
		nextRecvBufferSize: sizeTable[initialIndex],
	}
}

func (r *RecvByteBufAllocator) Allocate() ByteBuffer {
	buf := make([]byte, r.nextRecvBufferSize)
	return NewByteBuf(buf)
}

func (r *RecvByteBufAllocator) Record(readBytes int64) {
	if readBytes <= sizeTable[mathutil.Max(0, r.index-IndexDecrement-1)] {
		if r.decreaseNow {
			r.index = mathutil.Max(r.index-IndexDecrement, r.minIndex)
			r.nextRecvBufferSize = sizeTable[r.index]
			r.decreaseNow = false
		} else {
			r.decreaseNow = true
		}
	} else if readBytes >= r.nextRecvBufferSize {
		r.index = mathutil.Min(r.index+IndexIncrement, r.maxIndex)
		r.nextRecvBufferSize = sizeTable[r.index]
		r.decreaseNow = false
	}
}
