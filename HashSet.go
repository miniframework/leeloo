package leeloo

import (
	"unsafe"
)

type DefaultHashSetHandler struct {
}

func (d *DefaultHashSetHandler) Eq_cluster(cluster interface{}, item interface{}) bool {
	return true
}

func (m *DefaultHashSetHandler) Get_key(item interface{}) interface{} {

	return item
}
func (m *DefaultHashSetHandler) Eq_key(key interface{}, item interface{}) bool {

	return key == item
}
func (m *DefaultHashSetHandler) Hash_key(item interface{}) int {

	 return int(*(*byte)(unsafe.Pointer(&item)))
	  
}

type HashSet struct {
	hashTable *CowHashTable
}

func NewHashSet() *HashSet {
	return &HashSet{hashTable: &CowHashTable{_nbucket: 0, _factor: 0, _nitem: 0}}
}

func (hm *HashSet) Init(nbucket int, factor int, handler HashTableHandler) {
	hm.hashTable.Init(nbucket, factor, handler)
}

func (hm *HashSet) Insert(item interface{}) interface{} {
	return hm.hashTable.Insert(item, nil, item)
}
func (hm *HashSet) Erase(item interface{}) bool {
	return hm.hashTable.Erase(item, nil)
}

func (hm *HashSet) Find(item interface{}) interface{} {
	return hm.hashTable.Find(item, nil)
}
func (hm *HashSet) Seek() interface{} {
	return hm.hashTable.GetAll()
}

func (hm *HashSet) Clear() {
	hm.hashTable.Clear()
}
func (hm *HashSet) PrintInfo() {
	hm.hashTable.PrintInfo()
}
