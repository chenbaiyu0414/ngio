package internal

import "sync"

type ChannelAttributes struct {
	m sync.Map
}

func NewChannelAttributes() *ChannelAttributes {
	return &ChannelAttributes{}
}

func (attr *ChannelAttributes) Set(key, value interface{}) {
	attr.m.Store(key, value)
}

func (attr *ChannelAttributes) Get(key interface{}) interface{} {
	v, ok := attr.m.Load(key)
	if ok {
		return v
	}

	return nil
}

func (attr *ChannelAttributes) Has(key interface{}) bool {
	_, ok := attr.m.Load(key)
	return ok
}

func (attr *ChannelAttributes) Del(key interface{}) {
	attr.m.Delete(key)
}
