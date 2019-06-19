package channel

import "sync"

type Attributes interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Has(key interface{}) bool
	Del(key interface{})
}

type DefaultAttributeMap struct {
	m sync.Map
}

func NewDefaultAttributes() *DefaultAttributeMap {
	return &DefaultAttributeMap{}
}

func (d *DefaultAttributeMap) Set(key, value interface{}) {
	d.m.Store(key, value)
}

func (d *DefaultAttributeMap) Get(key interface{}) interface{} {
	v, ok := d.m.Load(key)
	if ok {
		return v
	}

	return nil
}

func (d *DefaultAttributeMap) Has(key interface{}) bool {
	_, ok := d.m.Load(key)
	return ok
}

func (d *DefaultAttributeMap) Del(key interface{}) {
	d.m.Delete(key)
}
