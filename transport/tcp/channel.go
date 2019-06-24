package tcp

import (
	"net"
	"ngio/buffer"
	"ngio/channel"
	"sync"
)

// tcpChannel is a connection between server and client
type tcpChannel struct {
	isActive      bool
	rwc           *net.TCPConn
	closeC        chan struct{}
	quitC         chan error
	writeC        chan buffer.ByteBuffer
	wg            sync.WaitGroup
	pipeline      *channel.Pipeline
	recvAllocator *buffer.RecvByteBufAllocator
	attributes    channel.Attributes
}

func newChannel(rwc *net.TCPConn) *tcpChannel {
	c := &tcpChannel{
		rwc:           rwc,
		closeC:        make(chan struct{}, 1),
		quitC:         make(chan error, 1),
		writeC:        make(chan buffer.ByteBuffer, 16),
		recvAllocator: buffer.NewRecvByteBufAllocator(buffer.DefaultMinimum, buffer.DefaultMaximum, buffer.DefaultInitial),
		attributes:    channel.NewDefaultAttributes(),
	}

	c.pipeline = channel.NewPipeline(c)
	return c
}

func (ch *tcpChannel) IsActive() bool {
	return ch.isActive
}

func (ch *tcpChannel) Pipeline() *channel.Pipeline {
	return ch.pipeline
}

func (ch *tcpChannel) LocalAddress() net.Addr {
	if ch.rwc == nil {
		return nil
	}

	return ch.rwc.LocalAddr()
}

func (ch *tcpChannel) RemoteAddress() net.Addr {
	if ch.rwc == nil {
		return nil
	}

	return ch.rwc.RemoteAddr()
}

func (ch *tcpChannel) Attributes() channel.Attributes {
	return ch.attributes
}

func (ch *tcpChannel) Serve() <-chan error {
	go ch.read()
	go ch.write()

	ch.isActive = true
	ch.pipeline.FireActiveHandler()

	return ch.quitC
}

func (ch *tcpChannel) read() {
	ch.wg.Add(1)
	defer ch.wg.Done()

	for {
		select {
		case <-ch.closeC:
			return
		default:
			buf := ch.recvAllocator.Allocate()
			n, err := buf.ReadFrom(ch.rwc)

			if err != nil {
				_ = ch.Close()
				return
			}

			ch.recvAllocator.Record(n)
			ch.pipeline.FireReadHandler(buf)
		}
	}
}

func (ch *tcpChannel) write() {
	ch.wg.Add(1)
	defer ch.wg.Done()

	for {
		select {
		case <-ch.closeC:
			return
		default:
			for buf := range ch.writeC {
				unwritten := int64(buf.ReadableBytes())

				for unwritten > 0 {
					n, err := buf.WriteTo(ch.rwc)

					if err != nil {
						_ = ch.Close()
						return
					}

					unwritten -= n
				}
			}
		}
	}
}

func (ch *tcpChannel) Write(msg interface{}) {
	if !ch.isActive {
		return
	}

	if buf, ok := msg.(buffer.ByteBuffer); ok {
		ch.writeC <- buf
	}
}

func (ch *tcpChannel) Close() error {
	ch.isActive = false
	ch.pipeline.FireInActiveHandler()

	ch.closeC <- struct{}{}
	ch.wg.Wait()

	close(ch.closeC)
	close(ch.writeC)

	err := ch.rwc.Close()
	ch.quitC <- err

	return err
}
