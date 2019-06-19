package tcp

import (
	"net"
	"ngio/buffer"
	"ngio/channel"
	"sync"
)

// tcpChannel is a connection between server and client
type tcpChannel struct {
	isActive             bool
	rwc                  *net.TCPConn
	closeC               chan struct{}
	quitC                chan struct{}
	writeC               chan buffer.ByteBuffer
	wg                   sync.WaitGroup
	pipeline             *channel.Pipeline
	recvByteBufAllocator *buffer.RecvByteBufAllocator
	attributes           channel.Attributes
}

func newChannel(conn *net.TCPConn) *tcpChannel {
	c := &tcpChannel{
		isActive:             true,
		rwc:                  conn,
		closeC:               make(chan struct{}, 1),
		quitC:                make(chan struct{}, 1),
		writeC:               make(chan buffer.ByteBuffer, 16),
		recvByteBufAllocator: buffer.NewRecvByteBufAllocator(buffer.DefaultMinimum, buffer.DefaultMaximum, buffer.DefaultInitial),
		attributes:           channel.NewDefaultAttributes(),
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

func (ch *tcpChannel) Serve() <-chan struct{} {
	go ch.read()
	go ch.write()
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
			buf := ch.recvByteBufAllocator.Allocate()
			n, err := ch.rwc.Read(buf.Buffer())

			if err == nil {
				buf.SetWriterIndex(n)
				ch.recvByteBufAllocator.Record(n)

				ch.pipeline.FireReadHandler(buf)
			} else {
				go ch.Close()
				return
			}
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
				total := buf.ReadableBytes()
				written := 0

				for {
					n, err := ch.rwc.Write(buf.Buffer()[buf.ReaderIndex():buf.WriterIndex()])
					if err == nil {
						written += n
						buf.SetReaderIndex(buf.ReaderIndex() + n)

						if written >= total {
							break
						}
					} else {
						go ch.Close()
						return
					}
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
	defer func() {
		ch.quitC <- struct{}{}
	}()

	ch.isActive = false

	ch.pipeline.FireInActiveHandler()

	ch.closeC <- struct{}{}

	ch.wg.Wait()

	close(ch.closeC)
	close(ch.writeC)

	return ch.rwc.Close()
}
