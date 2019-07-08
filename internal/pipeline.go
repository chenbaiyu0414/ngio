package internal

import (
	"fmt"
	"reflect"
	"sync"
)

type ChannelPipeline struct {
	ch         IChannel
	head, tail *channelContext
	contexts   map[string]*channelContext
	mu         sync.Mutex
}

func NewChannelPipeline(ch IChannel) *ChannelPipeline {
	pipeline := &ChannelPipeline{
		ch:       ch,
		contexts: make(map[string]*channelContext, 8),
	}

	pipeline.head = newContext("HEAD", &headHandler{}, pipeline)
	pipeline.tail = newContext("TAIL", nil, pipeline)

	pipeline.head.next = pipeline.tail
	pipeline.tail.prev = pipeline.head

	return pipeline
}

func (pipeline *ChannelPipeline) AddFirst(name string, handler interface{}) {
	pipeline.ensureHandlerType(handler)
	added := newContext(name, handler, pipeline)

	pipeline.mu.Lock()
	pipeline.insertBetween(added, pipeline.head, pipeline.head.next)
	pipeline.mu.Unlock()
}

func (pipeline *ChannelPipeline) AddLast(name string, handler interface{}) {
	pipeline.ensureHandlerType(handler)
	added := newContext(name, handler, pipeline)

	pipeline.mu.Lock()
	pipeline.insertBetween(added, pipeline.tail.prev, pipeline.tail)
	pipeline.mu.Unlock()
}

func (pipeline *ChannelPipeline) AddAfter(baseName, name string, handler interface{}) {
	pipeline.ensureHandlerType(handler)
	added := newContext(name, handler, pipeline)

	pipeline.mu.Lock()
	defer pipeline.mu.Unlock()

	if baseCtx, ok := pipeline.contexts[baseName]; ok {
		pipeline.insertBetween(added, baseCtx, baseCtx.next)
	} else {
		panic(fmt.Errorf(`non-existent channelContext with name "%s"`, baseName))
	}
}

func (pipeline *ChannelPipeline) AddBefore(baseName, name string, handler interface{}) {
	pipeline.ensureHandlerType(handler)
	added := newContext(name, handler, pipeline)

	pipeline.mu.Lock()
	defer pipeline.mu.Unlock()

	if base, ok := pipeline.contexts[baseName]; ok {
		pipeline.insertBetween(added, base.prev, base)
	} else {
		panic(fmt.Errorf(`non-existent channelContext with name "%s"`, baseName))
	}
}

func (pipeline *ChannelPipeline) Remove(name string) {
	pipeline.mu.Lock()
	defer pipeline.mu.Unlock()

	if deleted, ok := pipeline.contexts[name]; ok {
		prev, next := deleted.prev, deleted.prev
		deleted.prev.next = next
		deleted.next.prev = prev
		delete(pipeline.contexts, name)
	} else {
		panic(fmt.Errorf(`non-existent channelContext with name "%s"`, name))
	}
}

func (pipeline *ChannelPipeline) Replace(oldName, newName string, newHandler interface{}) {
	pipeline.ensureHandlerType(newHandler)
	replaced := newContext(newName, newHandler, pipeline)

	pipeline.mu.Lock()
	defer pipeline.mu.Unlock()

	if old, ok := pipeline.contexts[oldName]; ok {
		if _, ok := pipeline.contexts[newName]; ok {
			panic(fmt.Errorf(`repeated new name "%s"`, newName))
		}

		replaced.prev, replaced.next = old.prev, old.next
		old.prev.next = replaced
		old.next.prev = replaced

		delete(pipeline.contexts, oldName)
		pipeline.contexts[newName] = replaced
	} else {
		panic(fmt.Errorf(`non-existent channelContext with name "%s"`, oldName))
	}
}

func (pipeline *ChannelPipeline) Channel() IChannel {
	return pipeline.ch
}

func (pipeline *ChannelPipeline) FireChannelActiveHandler() {
	pipeline.head.FireChannelActiveHandler()
}

func (pipeline *ChannelPipeline) FireChannelInActiveHandler() {
	pipeline.head.FireChannelInActiveHandler()
}

func (pipeline *ChannelPipeline) FireChannelReadHandler(msg interface{}) {
	pipeline.head.FireChannelReadHandler(msg)
}

func (pipeline *ChannelPipeline) Write(msg interface{}) {
	pipeline.tail.Write(msg)
}

func (pipeline *ChannelPipeline) FireChannelErrorHandler(err error) {
	pipeline.head.FireChannelErrorHandler(err)
}

func (pipeline *ChannelPipeline) insertBetween(added, prev, next *channelContext) {
	if _, ok := pipeline.contexts[added.name]; !ok {
		pipeline.contexts[added.name] = added
	} else {
		panic(fmt.Errorf(`[name: %s] repeat name. "head" and "tail" are retained by default`, added.name))
	}

	added.prev, added.next = prev, next
	prev.next = added
	next.prev = added
}

func (pipeline *ChannelPipeline) ensureHandlerType(handler interface{}) {
	switch handler.(type) {
	case ActiveHandler, InActiveHandler, ReadHandler, WriteHandler, ErrorHandler:
	default:
		panic(fmt.Errorf(`invalid handler type: %s`, reflect.TypeOf(handler).Elem().Name()))
	}
}

func (pipeline *ChannelPipeline) String() string {
	return "p"
}

type headHandler struct {
}

func (*headHandler) Write(ctx IChannelContext, msg interface{}) {
	ctx.Pipeline().Channel().Write(msg)
}
