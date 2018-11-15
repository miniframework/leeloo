package leeloo

import (
    "testing"
    "fmt"
)

type mapHandler struct {
}
type tuple struct {
	key   int
	cluster int
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
func (m *mapHandler) Eq_cluster(cluster interface{}, item interface{}) bool {

	if cluster.(int) == item.(*tuple).cluster {
		return true
	}
	return false
}
func (m *mapHandler) Get_key(item interface{}) interface{} {

	return item.(*tuple).key
}

func Benchmark_HashTable_Insert(b *testing.B) {
	
	handler := &mapHandler{}
	hashTable := NewCowHashTable()
	hashTable.Init(10, 10, handler)
    for i := 0; i < 1000; i++ { //use b.N for looping 
        item := &tuple{i, i+1, fmt.Sprintf("abc%d", i) }
		hashTable.Insert(i,i+1, item)
    }
}
func Benchmark_Map_Insert(b *testing.B) {
	
	var m map[int]string
	m = make(map[int]string)
    for i := 0; i < 1000; i++ { //use b.N for looping 
		m[i]=fmt.Sprintf("abc%d", i)
    }
}

func Test_HashTable(t *testing.T) {
	item1 := &tuple{1,3, "aaa"}
	item2 := &tuple{1,4, "bbb"}
	item7 := &tuple{1,7, "resize"}
	handler := &mapHandler{}
	hashTable := NewCowHashTable()
	hashTable.Init(10, 10, handler)
	hashTable.Insert(1,3, item1)
	hashTable.Insert(1,4, item2)
	hashTable.Insert(1,7, item7)
	hashTable.Find(1, 3)
	t.Log(hashTable.PrintInfo())
	
	if(hashTable.Find(1, 3) == nil) {
		t.Error("not find (1,3)")
	}
	hashTable.Erase(1, 3) 
	t.Log(hashTable.PrintInfo())
	if(hashTable.Find(1, 3) != nil) {
		t.Error("not Erase (1,3)")
	}
	list := hashTable.Seek(1)
	if(list == nil) {
		t.Error("not find unique key (1)")
	}
	t.Log("find items ", list)
	hashTable.Remove(1)
	t.Log(hashTable.PrintInfo())
	list1 := hashTable.Seek(1)
	if(len(list1) != 0) {
		t.Error("not remove all by unique key(1)")
	}
	
	
	
}