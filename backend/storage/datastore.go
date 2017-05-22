// Package storage handles storing a sessionId to our Data struct
package storage

type DataStore struct {
	// using interface{} instead of `Data` because a data store is generally abstract
	// ie: in this case, why would the data store need to care about what Data is?
	// in the future: if we want more custom behaviour for Data we could add it
	// as a function on the struct itself
	Items map[string]interface{}
}

func (d *DataStore) Set(key string, newData interface{}) {
	d.Items[key] = newData
}

// had to be extra cheeky here to make it work
func (d *DataStore) Get(key string) (interface{}, bool) {
	if item, ok := d.Items[key]; ok {
		return item, true
	} else {
		return nil, false
	}
}
