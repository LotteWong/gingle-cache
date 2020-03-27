package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Calls struct {
	mux  sync.Mutex
	dict map[string]*call
}

func (cs *Calls) Do(key string, function func() (interface{}, error)) (interface{}, error) {
	cs.mux.Lock()
	if cs.dict == nil {
		cs.dict = make(map[string]*call) // TODO: 延迟初始化
	}

	// TODO: 高并发时键值重复，等待避免重入
	if c, ok := cs.dict[key]; ok {
		cs.mux.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	cs.dict[key] = c
	cs.mux.Unlock()

	c.val, c.err = function()
	c.wg.Done()

	cs.mux.Lock()
	delete(cs.dict, key)
	cs.mux.Unlock()

	return c.val, c.err
}
