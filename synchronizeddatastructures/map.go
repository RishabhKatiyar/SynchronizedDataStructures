//Reference - https://gobyexample.com/stateful-goroutines

package synchronizeddatastructures

type ReadResponse struct {
	Result interface{}
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

type DeleteObject struct {
	Key  interface{}
	Resp chan bool
}

type DerivedMapContainer struct {
	ReadOperation   chan ReadObject
	UpdateOperation chan UpdateObject
	DeleteOperation chan DeleteObject
	DerivedMap      map[interface{}]interface{}
}

// use this method only when defining customRead method for NewSynchronizedMap
func (c *DerivedMapContainer) ReadMap(key interface{}) (interface{}, bool) {
	val, exists := c.DerivedMap[key]
	return val, exists
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
