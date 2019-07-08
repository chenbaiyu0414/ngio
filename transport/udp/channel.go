package udp

import (
	"bytes"
	"net"
	"ngio/buffer"
	"ngio/internal"
	"ngio/logger"
	"strconv"
	"sync/atomic"
)

var udpChannelId uint32

type channel struct {
	id         uint32
	isActive   bool
	conn       *net.UDPConn
	pipeline   *internal.ChannelPipeline
	attributes *internal.ChannelAttributes
	quitC      chan error
	log        logger.Logger
}

func newChannel(conn *net.UDPConn) *channel {
	ch := &channel{
		id:         atomic.AddUint32(&udpChannelId, 1),
		isActive:   false,
		conn:       conn,
		attributes: internal.NewChannelAttributes(),
		quitC:      make(chan error, 1),
		log:        logger.DefaultLogger(),
	}

	ch.pipeline = internal.NewChannelPipeline(ch)
	return ch
}

func (ch *channel) Id() uint32 {
	return ch.id
}

func (ch *channel) IsActive() bool {
	return ch.isActive
}

func (ch *channel) Pipeline() internal.IChannelPipeline {
	return ch.pipeline
}

func (ch *channel) LocalAddress() net.Addr {
	if ch.conn == nil {
		return nil
	}

	return ch.conn.LocalAddr()
}

func (ch *channel) RemoteAddress() net.Addr {
	panic("not support")
}

func (ch *channel) Attributes() internal.IChannelAttributes {
	return ch.attributes
}

func (ch *channel) Serve() (err error) {
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

			go ch.pipeline.FireChannelReadHandler(packet)
		}
	}
}

func (ch *channel) Write(msg interface{}) {
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

func (ch *channel) Close() {
	if !ch.isActive {
		return
	}

	ch.isActive = false

	ch.log.Infof("[network: %v, local: %v, remote: %v] stop listening", ch.RemoteAddress().Network(), ch.LocalAddress().Network(), ch.RemoteAddress())

	ch.quitC <- ch.conn.Close()

	ch.log.Infof("[network: %v, local: %v, remote %v] listen stopped", ch.RemoteAddress().Network(), ch.LocalAddress().Network(), ch.RemoteAddress())

}

func (ch *channel) String() string {
	buf := bytes.Buffer{}

	buf.WriteString("channel id: ")
	buf.WriteString(strconv.FormatInt(int64(ch.id), 10))
	buf.WriteString("network: ")
	buf.WriteString(ch.LocalAddress().Network())
	buf.WriteString(", remote: ")
	buf.WriteString(ch.LocalAddress().String())
	buf.WriteString(", active: ")
	buf.WriteString(strconv.FormatBool(ch.isActive))

	return buf.String()
}
