package channel

import (
	"fmt"
	"ngio/internal/logger"
)

type Context interface {
	FireActiveHandler()
	FireInActiveHandler()
	FireReadHandler(msg interface{})
	Write(msg interface{})
	FireRecoverHandler(v interface{})
	Next() Context
	SetNext(next Context)
	Prev() Context
	SetPrev(prev Context)
	Pipeline() *Pipeline
	Handler() Handler
}

type DefaultContext struct {
	next, prev Context
	pipeline   *Pipeline
	handler    Handler
	flag       Flag
}

func NewDefaultContext(pipeline *Pipeline, handler interface{}) *DefaultContext {
	if handler == nil {
		panic("handler is nil")
	}

	return &DefaultContext{
		pipeline: pipeline,
		handler:  NewHandlerWrapper(handler),
	}
}

func (ctx *DefaultContext) FireActiveHandler() {
	if c := findInboundContext(ctx, Active); c != nil {
		try(c, func() {
			logger.Debugf("channel[%v] is active", c.Pipeline().Channel().RemoteAddress())
			c.Handler().ChannelActive(c)
		})
	}
}

func (ctx *DefaultContext) FireInActiveHandler() {
	if c := findInboundContext(ctx, InActive); c != nil {
		try(c, func() {
			logger.Debugf("channel[%v] is active", c.Pipeline().Channel().RemoteAddress())
			c.Handler().ChannelInActive(c)
		})
	}
}

func (ctx *DefaultContext) FireReadHandler(msg interface{}) {
	if c := findInboundContext(ctx, Read); c != nil {
		try(c, func() {
			c.Handler().ChannelRead(c, msg)
		})
	}
}

func (ctx *DefaultContext) Write(msg interface{}) {
	if c := findOutboundContext(ctx, Write); c != nil {
		try(c, func() {
			c.Handler().Write(c, msg)
		})
	}
}

func (ctx *DefaultContext) FireRecoverHandler(v interface{}) {
	if c := findInboundContext(ctx, Recover); c != nil {
		try(c, func() {
			c.Handler().ChannelRecovered(c, v)
		})
	}
}

func (ctx *DefaultContext) Next() Context {
	return ctx.next
}

func (ctx *DefaultContext) SetNext(next Context) {
	ctx.next = next
}

func (ctx *DefaultContext) Prev() Context {
	return ctx.prev
}

func (ctx *DefaultContext) SetPrev(prev Context) {
	ctx.prev = prev
}

func (ctx *DefaultContext) Pipeline() *Pipeline {
	return ctx.pipeline
}

func (ctx *DefaultContext) Handler() Handler {
	return ctx.handler
}

func (ctx *DefaultContext) Flag() Flag {
	return ctx.flag
}

type HeadContext struct {
	next, prev Context
	pipeline   *Pipeline
	handler    Handler
}

func NewHeadContext(pipeline *Pipeline) *HeadContext {
	return &HeadContext{
		pipeline: pipeline,
		handler:  NewHandlerWrapper(&headHandler{}),
	}
}

func (ctx *HeadContext) FireActiveHandler() {
	if c := findInboundContext(ctx, Active); c != nil {
		logger.Debugf("channel[%v] is active", c.Pipeline().Channel().RemoteAddress())
		c.Handler().ChannelActive(c)
	}
}

func (ctx *HeadContext) FireInActiveHandler() {
	if c := findInboundContext(ctx, InActive); c != nil {
		logger.Debugf("channel[%v] is inactive", c.Pipeline().Channel().RemoteAddress())
		c.Handler().ChannelInActive(c)
	}
}

func (ctx *HeadContext) FireReadHandler(msg interface{}) {
	if c := findInboundContext(ctx, Read); c != nil {
		c.Handler().ChannelRead(c, msg)
	}
}

func (ctx *HeadContext) Write(msg interface{}) {

}

func (ctx *HeadContext) FireRecoverHandler(v interface{}) {

}

func (ctx *HeadContext) Next() Context {
	return ctx.next
}

func (ctx *HeadContext) SetNext(next Context) {
	ctx.next = next
}

func (ctx *HeadContext) Prev() Context {
	return ctx.prev
}

func (ctx *HeadContext) SetPrev(prev Context) {
	ctx.prev = prev
}

func (ctx *HeadContext) Pipeline() *Pipeline {
	return ctx.pipeline
}

func (ctx *HeadContext) Handler() Handler {
	return ctx.handler
}

type headHandler struct{}

func (*headHandler) Write(ctx Context, msg interface{}) {
	ctx.Pipeline().Channel().Write(msg)
}

type TailContext struct {
	next, prev Context
	pipeline   *Pipeline
	handler    Handler
}

func NewTailContext(pipeline *Pipeline) *TailContext {
	return &TailContext{
		pipeline: pipeline,
		handler:  NewHandlerWrapper(nil),
	}
}

func (*TailContext) FireActiveHandler() {

}

func (*TailContext) FireInActiveHandler() {

}

func (*TailContext) FireReadHandler(msg interface{}) {

}

func (ctx *TailContext) Write(msg interface{}) {
	if c := findOutboundContext(ctx, Write); c != nil {
		c.Handler().Write(c, msg)
	}
}

func (*TailContext) FireRecoverHandler(v interface{}) {

}

func (ctx *TailContext) Next() Context {
	return ctx.next
}

func (ctx *TailContext) SetNext(next Context) {
	ctx.next = next
}

func (ctx *TailContext) Prev() Context {
	return ctx.prev
}

func (ctx *TailContext) SetPrev(prev Context) {
	ctx.prev = prev
}

func (ctx *TailContext) Pipeline() *Pipeline {
	return ctx.pipeline
}

func (ctx *TailContext) Handler() Handler {
	return ctx.handler
}

func findInboundContext(current Context, flag Flag) Context {
	for c := current.Next(); c != nil; c = c.Next() {
		if c.Handler().Flag()&flag == flag {
			return c
		}
	}

	return nil
}

func findOutboundContext(current Context, flag Flag) Context {
	for c := current.Next(); c != nil; c = c.Prev() {
		if c.Handler().Flag()&flag == flag {
			return c
		}
	}

	return nil
}

func try(ctx Context, doHandler func()) {
	defer func() {
		v := recover()
		if v == nil {
			return
		}

		if !(ctx.Handler().Flag()&Recover == Recover) {
			logger.Errorf("unhandled error: %v", v)
			return
		}

		switch v.(type) {
		case error:
			ctx.Handler().ChannelRecovered(ctx, v.(error))
		default:
			ctx.Handler().ChannelRecovered(ctx, fmt.Errorf("%v", v))
		}
	}()

	doHandler()
}
