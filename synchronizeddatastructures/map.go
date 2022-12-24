//Reference - https://gobyexample.com/stateful-goroutines

package synchronizeddatastructures

type ReadObject struct {
	Key  interface{}
	Resp chan interface{}
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
