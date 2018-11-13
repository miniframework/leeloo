package main

import (
	"fmt"
	"leeloo"
)

type mapHandler struct {
}
type tuple struct {
	key   int
	value string
}

func (m *mapHandler) Hash_key(key interface{}) int {
	nkey := key.(int) % 2
	return nkey
}
func (m *mapHandler) Eq_key(key interface{}, item interface{}) bool {

	if key.(int) == item.(*tuple).key {
		return true
	}
	return false
}
func (m *mapHandler) Get_key(item interface{}) interface{} {

	return item.(*tuple).key
}

func main() {

	mh1 := &tuple{3, "aaa"}
	mh3 := &tuple{5, "ccc"}
	//mh4 := &tuple{6, "dddddd"}

	handler := &mapHandler{}
	hashTable := leeloo.NewCowHashTable()
	hashTable.Init(10, 10, handler)
	hashTable.Insert(3, mh1)
	hashTable.Insert(5, mh3)

	hashTable1 := leeloo.NewCowHashTable()

	hashTable1.CowCopy(hashTable)

	hashTable1.Erase(3)
	item := hashTable.Seek(3)
	fmt.Println(item)
	fmt.Println("end.")

}
