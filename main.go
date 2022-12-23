package main

import (
	"fmt"
	"strconv"
	"sync"

	sds "github.com/RishabhKatiyar/SynchronizedDataStructures/synchronizeddatastructures"
)

// user define read
func simpleRead(derivedMapContainer *sds.DerivedMapContainer, readOb *sds.ReadObject) {
	response, exists := derivedMapContainer.ReadMap(readOb.Key)
	payload := sds.ReadResponse{Result: response, Found: exists}
	readOb.Resp <- payload
}

// user defined update
func simpleUpdate(derivedMapContainer *sds.DerivedMapContainer, updateOb *sds.UpdateObject) {
	derivedMapContainer.UpdateMap(updateOb.Key, updateOb.Val)
	updateOb.Resp <- true
}

// user defined delete
func simpleDelete(derivedMapContainer *sds.DerivedMapContainer, deleteOb *sds.DeleteObject) {
	derivedMapContainer.DeleteMap(deleteOb.Key.([]string)[0])
	deleteOb.Resp <- true
}

func main() {
	//region example 1 - map[string]int

	simpleMapContainer := sds.NewSynchronizedMap(simpleRead, simpleUpdate, simpleDelete)

	var wg sync.WaitGroup
	// Create

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(param int) {
			defer wg.Done()
			// define object to be created
			// pass the create object to create operation
			// get response
			updateOb := sds.UpdateObject{Key: strconv.Itoa(param), Val: param * 10, Resp: make(chan bool)}
			simpleMapContainer.UpdateOperation <- updateOb
			<-updateOb.Resp
		}(i)
	}
	wg.Wait()
	fmt.Println("Created New Map :: ", simpleMapContainer.DerivedMap)

	// Read
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(param int) {
			defer wg.Done()
			// define object to be read
			// pass the read object to read operation
			// get response
			readOb := sds.ReadObject{Key: strconv.Itoa(param), Resp: make(chan sds.ReadResponse)}
			simpleMapContainer.ReadOperation <- readOb
			response := <-readOb.Resp
			if response.Found {
				val := response.Result.(int)
				fmt.Printf("Value for key %s is %d \n", strconv.Itoa(param), val)
			}
		}(i)
	}
	wg.Wait()

	// Update
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(param int) {
			defer wg.Done()
			// define object to be updated
			// pass the update object to update operation
			// get response
			updateOb := sds.UpdateObject{Key: strconv.Itoa(param), Val: param*10 - 1, Resp: make(chan bool)}
			simpleMapContainer.UpdateOperation <- updateOb
			<-updateOb.Resp
		}(i)
	}
	wg.Wait()
	fmt.Println("Updated Map :: ", simpleMapContainer.DerivedMap)

	// Delete
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(param int) {
			defer wg.Done()
			if param%2 != 0 {
				// define object to be deleted
				// pass the delete object to delete operation
				// get response
				deleteOb := sds.DeleteObject{Key: []string{strconv.Itoa(param)}, Resp: make(chan bool)}
				simpleMapContainer.DeleteOperation <- deleteOb
				<-deleteOb.Resp
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("Deleted Odd Keys, Map :: ", simpleMapContainer.DerivedMap)

	// read key that is not present
	key := "11"
	readOb := sds.ReadObject{Key: key, Resp: make(chan sds.ReadResponse)}
	simpleMapContainer.ReadOperation <- readOb
	response := <-readOb.Resp
	if response.Found {
		val := response.Result.(int)
		fmt.Printf("Value for key %s is %d \n", key, val)
	} else {
		fmt.Printf("Value for key %s not found \n", key)
	}

	// read key that is present
	key = "4"
	readOb = sds.ReadObject{Key: key, Resp: make(chan sds.ReadResponse)}
	simpleMapContainer.ReadOperation <- readOb
	response = <-readOb.Resp

	if response.Found {
		val := response.Result.(int)
		fmt.Printf("Value for key %s is %d \n", key, val)
	} else {
		fmt.Printf("Value for key %s not found \n", key)
	}

	//endregion
}
