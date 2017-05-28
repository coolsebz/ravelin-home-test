// Package storage handles storing a sessionId to our Data struct
//   using a mutex to solve potential early concurrency issues
package storage

import (
	"sync"
)

type DataStore struct {
	sync.RWMutex
	// using interface{} instead of `Data` because a data store is generally abstract
	// ie: in this case, why would the data store need to care about what Data is?
	// in the future: if we want more custom behaviour for Data we could add it
	// as a function on the struct itself
	Items map[string]interface{}
}

func (d *DataStore) Set(key string, newData interface{}) {
	d.Lock()
	d.Items[key] = newData
	d.Unlock()
}

// had to be extra cheeky here to make it work
func (d *DataStore) Get(key string) (interface{}, bool) {
	d.RLock()
	defer d.RUnlock()
	if item, ok := d.Items[key]; ok {
		return item, true
	} else {
		return nil, false
	}
}

func New() DataStore {
	return DataStore{
		Items: make(map[string]interface{}),
	}
}
