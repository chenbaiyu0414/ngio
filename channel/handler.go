package channel

import "ngio/logger"

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
	ChannelActive(ctx *Context)
}

type InActiveHandler interface {
	ChannelInActive(ctx *Context)
}

type ReadHandler interface {
	ChannelRead(ctx *Context, msg interface{})
}

type WriteHandler interface {
	Write(ctx *Context, msg interface{})
}

type ErrorHandler interface {
	HandleError(ctx *Context, err error)
}

type HandlerAdapter struct {
	name            string
	activeHandler   ActiveHandler
	inActiveHandler InActiveHandler
	readHandler     ReadHandler
	writeHandler    WriteHandler
	errorHandler    ErrorHandler
	flag            Flag
	log             logger.Logger
}

func NewHandlerAdapter(name string, handler interface{}) *HandlerAdapter {
	adapter := &HandlerAdapter{
		name: name,
		flag: None,
		log:  logger.DefaultLogger(),
	}

	if h, ok := handler.(ActiveHandler); ok {
		adapter.activeHandler = h
		adapter.flag |= Active
	}

	if h, ok := handler.(InActiveHandler); ok {
		adapter.inActiveHandler = h
		adapter.flag |= InActive
	}

	if h, ok := handler.(ReadHandler); ok {
		adapter.readHandler = h
		adapter.flag |= Read
	}

	if h, ok := handler.(WriteHandler); ok {
		adapter.writeHandler = h
		adapter.flag |= Write
	}

	if h, ok := handler.(ErrorHandler); ok {
		adapter.errorHandler = h
		adapter.flag |= HandleError
	}

	return adapter
}

func (adapter *HandlerAdapter) ChannelActive(ctx *Context) {
	if adapter.activeHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel active", adapter.name)
		adapter.activeHandler.ChannelActive(ctx)
	}
}

func (adapter *HandlerAdapter) ChannelInActive(ctx *Context) {
	if adapter.inActiveHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel inactive", adapter.name)
		adapter.inActiveHandler.ChannelInActive(ctx)
	}
}

func (adapter *HandlerAdapter) ChannelRead(ctx *Context, msg interface{}) {
	if adapter.readHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel read", adapter.name)
		adapter.readHandler.ChannelRead(ctx, msg)
	}
}

func (adapter *HandlerAdapter) Write(ctx *Context, msg interface{}) {
	if adapter.writeHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel write", adapter.name)
		adapter.writeHandler.Write(ctx, msg)
	}
}

func (adapter *HandlerAdapter) HandleError(ctx *Context, err error) {
	if adapter.errorHandler != nil {
		adapter.log.Debugf("[handler: %s] invoke channel handle error", adapter.name)
		adapter.errorHandler.HandleError(ctx, err)
	}
}
