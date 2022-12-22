//Reference - https://gobyexample.com/stateful-goroutines

package synchronizeddatastructures

import (
	"encoding/json"
)

type ReadResponse struct {
	Result []byte
	Found  bool
}

type ReadObject struct {
	Key  interface{}
	Resp chan ReadResponse
}

type UpdateObject struct {
	Key  interface{}
	Val  interface{}
	Resp chan bool
}

// derived map can be map of maps so an operation must contain all the keys (nested)
// e.g. [key1][key2]...Val
type DeleteObject struct {
	Key  []interface{}
	Resp chan bool
}

type DerivedMapContainer struct {
	ReadOperation   chan ReadObject
	UpdateOperation chan UpdateObject
	DeleteOperation chan DeleteObject
	DerivedMap      map[interface{}]interface{}
}

// use this method only when defining customRead method for NewSynchronizedMap
func (c *DerivedMapContainer) ReadMap(key interface{}) ([]byte, bool) {
	val, exists := c.DerivedMap[key]
	response, _ := json.Marshal(val)
	return response, exists
}

// use this method only when defining customUpdate method for NewSynchronizedMap
func (c *DerivedMapContainer) UpdateMap(key interface{}, val interface{}) {
	c.DerivedMap[key] = val
}

// use this method only when defining customerDelete method for NewSynchronizedMap
func (c *DerivedMapContainer) DeleteMap(key interface{}) {
	delete(c.DerivedMap, key)
}

type customRead func(derivedMapContainer *DerivedMapContainer, readOb *ReadObject)
type customUpdate func(derivedMapContainer *DerivedMapContainer, readOb *UpdateObject)
type customerDelete func(derivedMapContainer *DerivedMapContainer, deleteOb *DeleteObject)

func NewSynchronizedMap(customRead customRead, customUpdate customUpdate, customDelete customerDelete) *DerivedMapContainer {
	derivedMapContainer := DerivedMapContainer{
		ReadOperation:   make(chan ReadObject),
		UpdateOperation: make(chan UpdateObject),
		DeleteOperation: make(chan DeleteObject),
		DerivedMap:      make(map[interface{}]interface{}),
	}

	go func() {
		for {
			select {
			case readOb := <-derivedMapContainer.ReadOperation:
				customRead(&derivedMapContainer, &readOb)
			case updateOb := <-derivedMapContainer.UpdateOperation:
				customUpdate(&derivedMapContainer, &updateOb)
			case deleteOb := <-derivedMapContainer.DeleteOperation:
				customDelete(&derivedMapContainer, &deleteOb)
			}
		}
	}()

	return &derivedMapContainer
}
