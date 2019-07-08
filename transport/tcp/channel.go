package tcp

import (
	"bytes"
	"io"
	"net"
	"ngio/buffer"
	"ngio/internal"
	"ngio/logger"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var tcpChannelId uint32

type channel struct {
	id                  uint32
	isActive            bool
	conn                net.Conn
	closeC              chan struct{}
	quitC               chan error
	writeC              chan buffer.ByteBuffer
	wg                  sync.WaitGroup
	pipeline            *internal.ChannelPipeline
	recvAllocator       *buffer.RecvByteBufAllocator
	attributes          *internal.ChannelAttributes
	writeDeadlinePeriod time.Duration
	readDeadlinePeriod  time.Duration
	log                 logger.Logger
}

func newChannel(conn net.Conn, writeDeadlinePeriod, readDeadlinePeriod time.Duration) *channel {
	ch := &channel{
		id:                  atomic.AddUint32(&tcpChannelId, 1),
		isActive:            false,
		conn:                conn,
		closeC:              make(chan struct{}),
		quitC:               make(chan error, 1),
		writeC:              make(chan buffer.ByteBuffer, 16),
		wg:                  sync.WaitGroup{},
		recvAllocator:       buffer.NewRecvByteBufAllocator(buffer.DefaultMinimum, buffer.DefaultMaximum, buffer.DefaultInitial),
		attributes:          internal.NewChannelAttributes(),
		writeDeadlinePeriod: writeDeadlinePeriod,
		readDeadlinePeriod:  readDeadlinePeriod,
		log:                 logger.DefaultLogger(),
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
	if ch.conn == nil {
		return nil
	}

	return ch.conn.RemoteAddr()
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

	go ch.read()
	go ch.write()

	ch.isActive = true

	ch.log.Debugf("[%v] serve", ch)

	ch.pipeline.FireChannelActiveHandler()

	return <-ch.quitC
}

func (ch *channel) read() {
	ch.wg.Add(1)
	defer ch.wg.Done()

	for {
		select {
		case <-ch.closeC:
			return
		default:
			buf := ch.recvAllocator.Allocate()

			// set read timeout
			if ch.readDeadlinePeriod > 0 {
				if err := ch.conn.SetReadDeadline(time.Now().Add(ch.readDeadlinePeriod)); err != nil {
					ch.Close()
					return
				}
			}

			n, err := ch.conn.Read(buf.Buffer())

			if err == nil {
				ch.recvAllocator.Record(n)
				buf.SetWriterIndex(n)
				ch.pipeline.FireChannelReadHandler(buf)
				continue
			}

			if err == io.EOF || err == io.ErrUnexpectedEOF {
				ch.Close()
			}

			return
		}
	}
}

func (ch *channel) write() {
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

					// set write timeout
					if ch.writeDeadlinePeriod > 0 {
						if err := ch.conn.SetWriteDeadline(time.Now().Add(ch.writeDeadlinePeriod)); err != nil {
							ch.Close()
							return
						}
					}

					n, err := buf.WriteTo(ch.conn)

					if err == nil {
						unwritten -= n
						continue
					}

					if err == io.EOF || err == io.ErrUnexpectedEOF {
						ch.Close()
					}

					return
				}
			}
		}
	}
}

func (ch *channel) Write(msg interface{}) {
	if !ch.isActive {
		// todo: logger
		return
	}

	if buf, ok := msg.(buffer.ByteBuffer); ok {
		ch.writeC <- buf
	}
	// todo: handle msg isn't ByteBuffer
}

func (ch *channel) Close() {
	if !ch.isActive {
		return
	}

	ch.isActive = false
	ch.pipeline.FireChannelInActiveHandler()

	go func() {
		// broadcast close signal
		close(ch.closeC)
		close(ch.writeC)

		err := ch.conn.Close()

		ch.wg.Wait()

		ch.quitC <- err
	}()
}

func (ch *channel) String() string {
	buf := bytes.Buffer{}

	buf.WriteString("channel id: ")
	buf.WriteString(strconv.FormatInt(int64(ch.id), 10))
	buf.WriteString(", network: ")
	buf.WriteString(ch.RemoteAddress().Network())
	buf.WriteString(", remote: ")
	buf.WriteString(ch.RemoteAddress().String())
	buf.WriteString(", active: ")
	buf.WriteString(strconv.FormatBool(ch.isActive))

	return buf.String()
}
