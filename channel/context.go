package channel

import (
	"bytes"
	"fmt"
	"ngio/logger"
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
	Name() string
	fmt.Stringer
}

type DefaultContext struct {
	name       string
	next, prev Context
	pipeline   *Pipeline
	handler    Handler
	flag       Flag
	log        logger.Logger
}

func NewDefaultContext(pipeline *Pipeline, name string, handler interface{}) *DefaultContext {
	if handler == nil {
		panic("handler is nil")
	}

	return &DefaultContext{
		name:     name,
		pipeline: pipeline,
		handler:  NewHandlerAdapter(name, handler),
		log:      logger.DefaultLogger(),
	}
}

func (ctx *DefaultContext) FireActiveHandler() {
	if c := findInboundContext(ctx, Active); c != nil {
		try(c, func() {
			ctx.log.Debugf("[%v] => [%v] fire active", ctx, c)
			c.Handler().ChannelActive(c)
		})
	}
}

func (ctx *DefaultContext) FireInActiveHandler() {
	if c := findInboundContext(ctx, InActive); c != nil {
		try(c, func() {
			ctx.log.Debugf("[%v] => [%v] fire inactive", ctx, c)
			c.Handler().ChannelInActive(c)
		})
	}
}

func (ctx *DefaultContext) FireReadHandler(msg interface{}) {
	if c := findInboundContext(ctx, Read); c != nil {
		try(c, func() {
			ctx.log.Debugf("[%v] => [%v] fire read", ctx, c)
			c.Handler().ChannelRead(c, msg)
		})
	}
}

func (ctx *DefaultContext) Write(msg interface{}) {
	if c := findOutboundContext(ctx, Write); c != nil {
		try(c, func() {
			ctx.log.Debugf("[%v] => [%v] fire write", ctx, c)
			c.Handler().Write(c, msg)
		})
	}
}

func (ctx *DefaultContext) FireRecoverHandler(v interface{}) {
	if c := findInboundContext(ctx, Recover); c != nil {
		try(c, func() {
			ctx.log.Debugf("[%v] => [%v] fire recover", ctx, c)
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

func (ctx *DefaultContext) Name() string {
	return ctx.name
}

func (ctx *DefaultContext) String() string {
	buf := bytes.Buffer{}

	buf.WriteString("context: ")
	buf.WriteString(ctx.name)
	buf.WriteString(", channel: ")
	buf.WriteString(fmt.Sprintf("%p", ctx.pipeline.ch))

	return buf.String()
}

type HeadContext struct {
	name       string
	next, prev Context
	pipeline   *Pipeline
	handler    Handler
	log        logger.Logger
}

func NewHeadContext(pipeline *Pipeline, name string) *HeadContext {
	return &HeadContext{
		name:     name,
		pipeline: pipeline,
		handler:  NewHandlerAdapter(name, &headHandler{}),
		log:      logger.DefaultLogger(),
	}
}

func (ctx *HeadContext) FireActiveHandler() {
	if c := findInboundContext(ctx, Active); c != nil {
		ctx.log.Debugf("[%v] => [%v] fire active", ctx, c)
		c.Handler().ChannelActive(c)
	}
}

func (ctx *HeadContext) FireInActiveHandler() {
	if c := findInboundContext(ctx, InActive); c != nil {
		ctx.log.Debugf("[%v] => [%v] fire inactive", ctx, c)
		c.Handler().ChannelInActive(c)
	}
}

func (ctx *HeadContext) FireReadHandler(msg interface{}) {
	if c := findInboundContext(ctx, Read); c != nil {
		ctx.log.Debugf("[%v] => [%v] fire read", ctx, c)
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

func (ctx *HeadContext) Name() string {
	return ctx.name
}

func (ctx *HeadContext) String() string {
	buf := bytes.Buffer{}

	buf.WriteString("context: ")
	buf.WriteString(ctx.name)
	buf.WriteString(", channel: ")
	buf.WriteString(fmt.Sprintf("%p", ctx.pipeline.ch))

	return buf.String()
}

type headHandler struct{}

func (*headHandler) Write(ctx Context, msg interface{}) {
	ctx.Pipeline().Channel().Write(msg)
}

type TailContext struct {
	name       string
	next, prev Context
	pipeline   *Pipeline
	handler    Handler
	log        logger.Logger
}

func NewTailContext(pipeline *Pipeline, name string) *TailContext {
	return &TailContext{
		name:     name,
		pipeline: pipeline,
		handler:  NewHandlerAdapter(name, nil),
		log:      logger.DefaultLogger(),
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
		ctx.log.Debugf("[%v] => [%v] fire write", ctx, c)
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

func (ctx *TailContext) Name() string {
	return ctx.name
}

func (ctx *TailContext) String() string {
	buf := bytes.Buffer{}

	buf.WriteString("context: ")
	buf.WriteString(ctx.name)
	buf.WriteString(", channel: ")
	buf.WriteString(fmt.Sprintf("%p", ctx.pipeline.ch))

	return buf.String()
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
