package leeloo


type DefaultHashMapHandler struct {
}

func (d *DefaultHashMapHandler) Eq_cluster(cluster interface{}, item interface{}) bool {
	return true
}

type HashMap struct {
	hashTable *CowHashTable
}

func NewHashMap() *HashMap {
	return &HashMap{hashTable: &CowHashTable{_nbucket: 0, _factor: 0, _nitem: 0}}
}

func (hm *HashMap) Init(nbucket int, factor int, handler HashTableHandler) {
	hm.hashTable.Init(nbucket, factor, handler)
}

func (hm *HashMap) Insert(key interface{}, item interface{}) interface{} {
	return hm.hashTable.Insert(key, nil, item)
}
func (hm *HashMap) Erase(key interface{}) bool {
	return hm.hashTable.Erase(key, nil)
}

func (hm *HashMap) Find(key interface{}) interface{} {
	return hm.hashTable.Find(key, nil)
}
func (hm *HashMap) Clear() {
	hm.hashTable.Clear()
}
func (hm *HashMap) PrintInfo() {
	hm.hashTable.PrintInfo()
}
