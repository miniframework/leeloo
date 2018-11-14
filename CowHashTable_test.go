package leeloo

import (
    "testing"
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



func Test_HashTable(t *testing.T) {
	item1 := &tuple{3, "aaa"}
	item2 := &tuple{4, "bbb"}
	item7 := &tuple{7, "resize"}
	handler := &mapHandler{}
	hashTable := NewCowHashTable()
	hashTable.Init(10, 10, handler)
	hashTable.Insert(3, item1)
	hashTable.Insert(4, item2)
	hashTable.Insert(7, item7)
	value := hashTable.Seek(3)
	t.Log(hashTable.PrintInfo())
	r := hashTable.Erase(4)
	t.Log(hashTable.PrintInfo())
	if( r == false) {
		t.Error("Erase false")
	} else {
		t.Log("Erase item2 sucess!")
	}
	if(value == nil || value.(*tuple).value != "aaa" ) {
		t.Error("value == nil || value != aaa")
	} else {
		t.Log("Find hashTable item1 sucess!")
	}
	hashTable1 := CowHashTableCopy(hashTable)
	t.Log(hashTable1.PrintInfo())
	value2 := hashTable1.Seek(3)
	if(value2 == nil || value2.(*tuple).value != "aaa" ) {
		t.Error("value2 == nil || value != aaa")
	} else {
		t.Log("Find hashTable1 item1 sucess!")
	}
	item3 := &tuple{5, "dddd"}
	hashTable1.Insert(5, item3)
	t.Log(hashTable1.PrintInfo())
	t.Log(hashTable.PrintInfo())
	hashTable1.Erase(3)
	t.Log(hashTable1.PrintInfo())
	t.Log(hashTable.PrintInfo())
	value3 := hashTable1.Seek(3)
	if(value3 != nil ) {
		t.Error("value3 != nil")
	} else {
		t.Log("seek  value3  sucess!")
	}
	value4 := hashTable.Seek(3)
	if(value4 == nil || value4.(*tuple).value != "aaa" ) {
		t.Error("value4 == nil || value != aaa")
	} else {
		t.Log("Find hashTable1 item1 sucess!")
	}
}