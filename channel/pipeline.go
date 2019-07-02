package channel

import (
	"fmt"
	"sync"
)

type Pipeline struct {
	ch         Channel
	head, tail *Context
	contexts   map[string]*Context
	mu         sync.Mutex
}

func NewPipeline(ch Channel) *Pipeline {
	pipeline := &Pipeline{
		ch:       ch,
		contexts: make(map[string]*Context, 8),
	}

	pipeline.head = NewContext("HEAD", &headHandler{}, pipeline)
	pipeline.tail = NewContext("TAIL", &tailHandler{}, pipeline)

	pipeline.head.next = pipeline.tail
	pipeline.tail.prev = pipeline.head

	return pipeline
}

func (pipeline *Pipeline) insertBetween(added, prev, next *Context) {
	if _, ok := pipeline.contexts[added.name]; !ok {
		pipeline.contexts[added.name] = added
	} else {
		panic(fmt.Errorf(`[name: %s] repeat name. "head" and "tail" are retained by default`, added.name))
	}

	added.prev, added.next = prev, next
	prev.next = added
	next.prev = added
}

func (pipeline *Pipeline) AddFirst(name string, handler interface{}) {
	added := NewContext(name, handler, pipeline)

	pipeline.mu.Lock()
	pipeline.insertBetween(added, pipeline.head, pipeline.head.next)
	pipeline.mu.Unlock()
}

func (pipeline *Pipeline) AddLast(name string, handler interface{}) {
	added := NewContext(name, handler, pipeline)

	pipeline.mu.Lock()
	pipeline.insertBetween(added, pipeline.tail.prev, pipeline.tail)
	pipeline.mu.Unlock()
}

func (pipeline *Pipeline) AddAfter(basename, name string, handler interface{}) {
	added := NewContext(name, handler, pipeline)

	pipeline.mu.Lock()
	defer pipeline.mu.Unlock()

	if baseCtx, ok := pipeline.contexts[basename]; ok {
		pipeline.insertBetween(added, baseCtx, baseCtx.next)
	} else {
		panic(fmt.Errorf(`non-existent context with name "%s"`, basename))
	}
}

func (pipeline *Pipeline) AddBefore(basename, name string, handler interface{}) {
	added := NewContext(name, handler, pipeline)

	pipeline.mu.Lock()
	defer pipeline.mu.Unlock()

	if base, ok := pipeline.contexts[basename]; ok {
		pipeline.insertBetween(added, base.prev, base)
	} else {
		panic(fmt.Errorf(`non-existent context with name "%s"`, basename))
	}
}

func (pipeline *Pipeline) Remove(name string) {
	pipeline.mu.Lock()
	defer pipeline.mu.Unlock()

	if deleted, ok := pipeline.contexts[name]; ok {
		prev, next := deleted.prev, deleted.prev
		deleted.prev.next = next
		deleted.next.prev = prev
		delete(pipeline.contexts, name)
	} else {
		panic(fmt.Errorf(`non-existent context with name "%s"`, name))
	}
}

func (pipeline *Pipeline) Replace(oldName, newName string, newHandler interface{}) {
	replaced := NewContext(newName, newHandler, pipeline)

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
		panic(fmt.Errorf(`non-existent context with name "%s"`, oldName))
	}
}

func (pipeline *Pipeline) Channel() Channel {
	return pipeline.ch
}

func (pipeline *Pipeline) FireActiveHandler() {
	pipeline.head.FireActiveHandler()
}

func (pipeline *Pipeline) FireInActiveHandler() {
	pipeline.head.FireInActiveHandler()
}

func (pipeline *Pipeline) FireReadHandler(msg interface{}) {
	pipeline.head.FireReadHandler(msg)
}

func (pipeline *Pipeline) FireWriteHandler(msg interface{}) {
	pipeline.tail.Write(msg)
}

func (pipeline *Pipeline) FireErrorHandler(err error) {
	pipeline.head.FireChannelErrorHandler(err)
}

type headHandler struct {
}

func (*headHandler) Write(ctx *Context, msg interface{}) {
	ctx.pipeline.ch.Write(msg)
}

type tailHandler struct{}
