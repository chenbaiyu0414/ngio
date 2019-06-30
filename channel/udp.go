package channel

import (
	"bytes"
	"fmt"
	"net"
	"ngio/buffer"
	"ngio/logger"
	"strconv"
)

type UDPChannel struct {
	isActive   bool
	conn       *net.UDPConn
	pipeline   *Pipeline
	attributes Attributes
	quitC      chan error
	log        logger.Logger
}

func NewUDPChannel(conn *net.UDPConn) *UDPChannel {
	ch := &UDPChannel{
		isActive:   false,
		conn:       conn,
		attributes: NewDefaultAttributes(),
		quitC:      make(chan error, 1),
		log:        logger.DefaultLogger(),
	}

	ch.pipeline = NewPipeline(ch)
	return ch
}

func (ch *UDPChannel) IsActive() bool {
	return ch.isActive
}

func (ch *UDPChannel) Pipeline() *Pipeline {
	return ch.pipeline
}

func (ch *UDPChannel) LocalAddress() net.Addr {
	if ch.conn == nil {
		return nil
	}

	return ch.conn.LocalAddr()
}

func (ch *UDPChannel) RemoteAddress() net.Addr {
	panic("not support")
}

func (ch *UDPChannel) Attributes() Attributes {
	return ch.attributes
}

func (ch *UDPChannel) Serve() (err error) {
	defer func() {
		if err != nil {
			ch.log.Debugf("[%v] close\r\n %v", ch, err)
		} else {
			ch.log.Debugf("[%v] close", ch)
		}
	}()

	ch.isActive = true
	ch.log.Debugf("[%v] serve", ch)

	for {
		select {
		case exitErr := <-ch.quitC:
			return exitErr
		default:
			// todo: size by config
			buf := make([]byte, 1024)

			r, raddr, err := ch.conn.ReadFromUDP(buf)
			if err != nil {
				ch.Close()
				return err
			}

			bf := buffer.NewByteBuf(buf, 0, r)
			packet := buffer.NewDatagramPacket(raddr, bf)

			go ch.pipeline.FireReadHandler(packet)
		}
	}
}

func (ch *UDPChannel) Write(msg interface{}) {
	if !ch.isActive {
		// todo: logger
		return
	}

	if packet, ok := msg.(*buffer.DatagramPacket); ok {
		shouldWrite := packet.ByteBuf().ReadableBytes()

		w, err := ch.conn.WriteToUDP(packet.ByteBuf().ReadBytes(shouldWrite), packet.RemoteAddress())
		if err != nil {
			// todo: logger
		}

		if w != shouldWrite {
			// todo: logger
		}
	}
}

func (ch *UDPChannel) Close() {
	if !ch.isActive {
		return
	}

	ch.isActive = false

	ch.log.Infof("[network: %v, local: %v, remote: %v] stop listening", ch.RemoteAddress().Network(), ch.LocalAddress().Network(), ch.RemoteAddress())

	ch.quitC <- ch.conn.Close()

	ch.log.Infof("[network: %v, local: %v, remote %v] listen stopped", ch.RemoteAddress().Network(), ch.LocalAddress().Network(), ch.RemoteAddress())

}

func (ch *UDPChannel) String() string {
	buf := bytes.Buffer{}

	buf.WriteString("channel: ")
	buf.WriteString(fmt.Sprintf("%p", ch))
	buf.WriteString("network: ")
	buf.WriteString(ch.LocalAddress().Network())
	buf.WriteString(", remote: ")
	buf.WriteString(ch.LocalAddress().String())
	buf.WriteString(", active: ")
	buf.WriteString(strconv.FormatBool(ch.isActive))

	return buf.String()
}
