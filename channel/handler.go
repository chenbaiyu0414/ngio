package channel

import "ngio/logger"

type Flag int

const (
	None Flag = 1 << iota
	Active
	InActive
	Read
	Write
	Recover
)

type Handler interface {
	ActiveHandler
	InActiveHandler
	ReadHandler
	WriteHandler
	RecoverHandler
	Flag() Flag
}

type ActiveHandler interface {
	ChannelActive(ctx Context)
}

type InActiveHandler interface {
	ChannelInActive(ctx Context)
}

type ReadHandler interface {
	ChannelRead(ctx Context, msg interface{})
}

type WriteHandler interface {
	Write(ctx Context, msg interface{})
}

type RecoverHandler interface {
	ChannelRecovered(ctx Context, v interface{})
}

type HandlerAdapter struct {
	name            string
	activeHandler   ActiveHandler
	inActiveHandler InActiveHandler
	readHandler     ReadHandler
	writeHandler    WriteHandler
	recoverHandler  RecoverHandler
	flag            Flag
	log             logger.Logger
}

func NewHandlerAdapter(name string, handler interface{}) *HandlerAdapter {
	wrapper := &HandlerAdapter{
		name: name,
		flag: None,
		log:  logger.DefaultLogger(),
	}

	if h, ok := handler.(ActiveHandler); ok {
		wrapper.activeHandler = h
		wrapper.flag |= Active
	}

	if h, ok := handler.(InActiveHandler); ok {
		wrapper.inActiveHandler = h
		wrapper.flag |= InActive
	}

	if h, ok := handler.(ReadHandler); ok {
		wrapper.readHandler = h
		wrapper.flag |= Read
	}

	if h, ok := handler.(WriteHandler); ok {
		wrapper.writeHandler = h
		wrapper.flag |= Write
	}

	if h, ok := handler.(RecoverHandler); ok {
		wrapper.recoverHandler = h
		wrapper.flag |= Recover
	}

	return wrapper
}

func (adapter *HandlerAdapter) ChannelActive(ctx Context) {
	if adapter.activeHandler != nil {
		adapter.log.Debugf("[handler: %s] call channel active", adapter.name)
		adapter.activeHandler.ChannelActive(ctx)
	}
}

func (adapter *HandlerAdapter) ChannelInActive(ctx Context) {
	if adapter.inActiveHandler != nil {
		adapter.log.Debugf("[handler: %s] call channel inactive", adapter.name)
		adapter.inActiveHandler.ChannelInActive(ctx)
	}
}

func (adapter *HandlerAdapter) ChannelRead(ctx Context, msg interface{}) {
	if adapter.readHandler != nil {
		adapter.log.Debugf("[handler: %s] call channel read", adapter.name)
		adapter.readHandler.ChannelRead(ctx, msg)
	}
}

func (adapter *HandlerAdapter) Write(ctx Context, msg interface{}) {
	if adapter.writeHandler != nil {
		adapter.log.Debugf("[handler: %s] call channel write", adapter.name)
		adapter.writeHandler.Write(ctx, msg)
	}
}

func (adapter *HandlerAdapter) ChannelRecovered(ctx Context, v interface{}) {
	if adapter.recoverHandler != nil {
		adapter.log.Debugf("[handler: %s] call channel recover", adapter.name)
		adapter.recoverHandler.ChannelRecovered(ctx, v)
	}
}

func (adapter *HandlerAdapter) Flag() Flag {
	return adapter.flag
}
