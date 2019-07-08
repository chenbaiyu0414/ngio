package internal

import (
	"fmt"
	"net"
)

type Flag int

const (
	None Flag = 1 << iota
	Active
	InActive
	Read
	Write
	HandleError
)

type ActiveHandler interface {
	ChannelActive(ctx IChannelContext)
}

type InActiveHandler interface {
	ChannelInActive(ctx IChannelContext)
}

type ReadHandler interface {
	ChannelRead(ctx IChannelContext, msg interface{})
}

type WriteHandler interface {
	Write(ctx IChannelContext, msg interface{})
}

type ErrorHandler interface {
	HandleError(ctx IChannelContext, err error)
}

type IChannelHandler interface {
	ActiveHandler
	InActiveHandler
	ReadHandler
	WriteHandler
	ErrorHandler
	Flag() Flag
}

type Initializer func(channel IChannel)

type IChannel interface {
	Id() uint32
	IsActive() bool
	Pipeline() IChannelPipeline
	LocalAddress() net.Addr
	RemoteAddress() net.Addr
	Attributes() IChannelAttributes
	Serve() error
	Write(msg interface{})
	Close()
	fmt.Stringer
}

type IChannelContext interface {
	FireChannelActiveHandler()
	FireChannelInActiveHandler()
	FireChannelReadHandler(msg interface{})
	Write(msg interface{})
	FireChannelErrorHandler(err error)

	Next() IChannelContext
	Prev() IChannelContext
	Pipeline() IChannelPipeline
	Name() string
	fmt.Stringer
}

type IChannelPipeline interface {
	AddFirst(name string, handler interface{})
	AddLast(name string, handler interface{})
	AddAfter(baseName, name string, handler interface{})
	AddBefore(baseName, name string, handler interface{})
	Remove(name string)
	Replace(oldName, newName string, newHandler interface{})
	Channel() IChannel

	FireChannelActiveHandler()
	FireChannelInActiveHandler()
	FireChannelReadHandler(msg interface{})
	Write(msg interface{})
	FireChannelErrorHandler(err error)

	fmt.Stringer
}

type IChannelAttributes interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Has(key interface{}) bool
	Del(key interface{})
}
