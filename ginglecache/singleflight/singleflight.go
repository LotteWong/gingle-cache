package singleflight

import "sync"

// call maintains information of each connect
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Calls maintains information of many connects
type Calls struct {
	mux  sync.Mutex
	dict map[string]*call
}

// Do uses sync APIs to implement singleflight strategy
func (cs *Calls) Do(key string, function func() (interface{}, error)) (interface{}, error) {
	cs.mux.Lock()

	// Delay initialization
	if cs.dict == nil {
		cs.dict = make(map[string]*call)
	}

	if c, ok := cs.dict[key]; ok {
		cs.mux.Unlock()
		c.wg.Wait() // If existed, wait instead of excute
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1) // If not existed, excute only once
	cs.dict[key] = c

	cs.mux.Unlock()
	c.val, c.err = function()
	c.wg.Done() // If not existed, excute only once

	cs.mux.Lock()
	delete(cs.dict, key) // Thread-safely delete map kvpair
	cs.mux.Unlock()

	return c.val, c.err
}
