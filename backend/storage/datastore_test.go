package storage

// TODO: add a constructor and a more standardised way to use the store

import (
	"testing"
)

type TestData struct {
	AString string
	AnInt   int
}

func TestStore(t *testing.T) {
	testStore := DataStore{
		Items: make(map[string]interface{}),
	}

	// testing for a faulty get on the empty data store
	item, ok := testStore.Get("randomKey")
	if ok || item != nil {
		t.Error("Loaded an item even though the store was empty", item)
	}

	// testing setting a value for a key that is not in use
	testStore.Set("randomKey", TestData{"foo", 42})
	item, ok = testStore.Get("randomKey")
	if !ok || item == nil {
		t.Error("Could not load item after saving it", item)
	}
	if item.(TestData).AString != "foo" || item.(TestData).AnInt != 42 {
		t.Error("Seems like the values to messed up when retrieving", item)
	}

	// testing updating a value for a key that is already in use
	testStore.Set("randomKey", TestData{"bar", 99})
	item, ok = testStore.Get("randomKey")
	if !ok || item == nil {
		t.Error("Could not retrieve the item after saving over it", item)
	}
	if item.(TestData).AString != "bar" || item.(TestData).AnInt != 99 {
		t.Error("Values were not updated when saving over an existing key", item)
	}
}
