package buffer

import "net"

type DatagramPacket struct {
	bf    ByteBuffer
	raddr *net.UDPAddr
}

func NewDatagramPacket(raddr *net.UDPAddr, bf ByteBuffer) *DatagramPacket {
	return &DatagramPacket{
		bf:    bf,
		raddr: raddr,
	}
}

func (packet *DatagramPacket) ByteBuf() ByteBuffer {
	return packet.bf
}

func (packet *DatagramPacket) RemoteAddress() *net.UDPAddr {
	return packet.raddr
}
