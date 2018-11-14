package leeloo


import (
	"fmt"
)
//CowHashTable.go
//Author:zhibinwang
//Data:2018

//interface handler
type HashTableHandler interface {
	Hash_key(key interface{}) int
	Eq_key(key interface{}, item interface{}) bool
	Get_key(item interface{}) interface{}
}

type CowHashTable struct {
	_table   []*Bucket
	_nbucket int
	_factor  int
	_nitem   int
	_handler HashTableHandler
}

type Bucket struct {
	_ref  int
	_item interface{}
	_next *Bucket
}

func (b *Bucket) isMutable() bool {
	return b._ref == 1
}
func (b *Bucket) decRef() {
	b._ref--
}
func (b *Bucket) incRef() {
	b._ref++
}


func (ht *CowHashTable) chain_dec(begin *Bucket, end *Bucket) {

	for ; begin != end; begin = begin._next {
		begin.decRef()
	}

}
func (ht *CowHashTable) chain_inc(begin *Bucket, end *Bucket) {

	for ; begin != end; begin = begin._next {
		begin.incRef()
	}

}

func NewCowHashTable() *CowHashTable {
	return &CowHashTable{_nbucket: 0, _factor: 0, _nitem: 0}
}

func (ht *CowHashTable) Init(nbucket int, factor int, handler HashTableHandler) {
	ht._nbucket = nbucket
	ht._factor = factor
	ht._table = make([]*Bucket, ht._nbucket + 1)
	ht._handler = handler
}
func (ht *CowHashTable) isInit() bool {

	if nil != ht._table {
		return true
	} else {
		return false
	}

}
func CowHashTableCopy(rht *CowHashTable) *CowHashTable {
	ht := NewCowHashTable()
	ht.CowCopy(rht)
	return ht
}
//copy cowhashtable another one, share the same memory
func (ht *CowHashTable) CowCopy(rht *CowHashTable) *CowHashTable {
	if !rht.isInit() {
		return nil
	}
	ht._nbucket = rht._nbucket
	ht._factor = rht._factor
	ht._nitem = rht._nitem
	ht._handler = rht._handler
	ht._table = make([]*Bucket, rht._nbucket + 1)

	for i := 0; i < ht._nbucket; i++ {
		ht._table[i] = rht._table[i]
		ht.chain_inc(ht._table[i], nil)
	}
	return ht

}

//insert one item to CowHashTable
//key: any type
//item: any type
//Return  insert item
func (ht *CowHashTable) Insert(key interface{}, item interface{}) interface{} {

	bkt := ht._handler.Hash_key(key) % ht._nbucket

	bucket := ht._table[bkt]
	//find not equal from head to end bucket,
	for nil != bucket && !ht._handler.Eq_key(key, bucket._item) {
		bucket = bucket._next
	}

	//not find equal insert table head
	if nil == bucket {

		if ht.isResize() { // re-calculate bucket
			bkt = ht._handler.Hash_key(key) % ht._nbucket
		}

		head_bucket := &Bucket{_ref: 1, _item: item}

		head_bucket._next = ht._table[bkt]

		ht._table[bkt] = head_bucket

		ht._nitem++

		return head_bucket._item
	}
	//if not CowCopy replace new item
	if bucket.isMutable() {
		bucket._item = item
		return bucket._item
	}

	new_bucket := &Bucket{_ref: 1, _item: item}

	new_bucket._next = ht._table[bkt]

	prior_next := &new_bucket._next

	share_bucket := ht._table[bkt]

	for share_bucket != bucket && share_bucket.isMutable() {
		prior_next = &(share_bucket._next)
		share_bucket = share_bucket._next
	}

	//copy ref > 2 bucket until was found item
	for share_bucket != bucket {
		cow_bucket := &Bucket{_ref: 1, _item: share_bucket._item}
		share_bucket.decRef()
		*prior_next = cow_bucket
		prior_next = &(cow_bucket._next)
		share_bucket = share_bucket._next
	}

	*prior_next = bucket._next
	bucket.decRef()
	ht._table[bkt] = new_bucket
	return new_bucket._item
}

//erase the value match key from CowHashTable
//Return
//		 true:one item was earsed ,
//		 false:nothing match
func (ht *CowHashTable) Erase(key interface{}) bool {

	bkt := ht._handler.Hash_key(key) % ht._nbucket
	bucket := ht._table[bkt]
	prior_next := &(ht._table[bkt])

	//find next equal item
	for nil != bucket && !ht._handler.Eq_key(key, bucket._item) {
		prior_next = &bucket._next
		bucket = bucket._next
	}
	//no equal item
	if nil == bucket {
		return false
	}
	if bucket.isMutable() {
		*prior_next = bucket._next
		bucket.decRef()
		ht._nitem--
		return true
	}

	head_bucket := ht._table[bkt]
	prior_next = &head_bucket
	share_bucket := ht._table[bkt]

	for share_bucket != bucket && share_bucket.isMutable() {
		prior_next = &(share_bucket._next)
		share_bucket = share_bucket._next
	}

	for share_bucket != bucket {
		cow_bucket := &Bucket{_ref: 1, _item: share_bucket._item}
		share_bucket.decRef()
		*prior_next = cow_bucket
		prior_next = &(cow_bucket._next)
		share_bucket = share_bucket._next
	}
	*prior_next = bucket._next
	bucket.decRef()
	ht._nitem--
	ht._table[bkt] = head_bucket

	return true
}

// ratios between neighbor primes in find_prime
// are generally near but less than 2, using 190 rather than 200
// here guarantees primes in the sequence are not skipped.
func (ht *CowHashTable) isResize() bool {
	
	if (ht._nbucket * ht._factor) < (ht._nitem * 100) {
		return ht.resize(ht._nitem * 190 / ht._factor)
	}
	return false
}

// Change number of buckets
// Params:
//   nbucket2  intended number of buckets
// Returns: resized or not
func (ht *CowHashTable) resize(nbucket2 int) bool {
	nbucket2 = find_prime(nbucket2)
	//no space malloc
	if ht._nbucket == nbucket2 {
		return false
	}
	newTable := make([]*Bucket, nbucket2+1)

	if ht._table != nil {

		for i := 0; i < ht._nbucket; i++ {

			bucket := ht._table[i]

			for bucket != nil {

				next_bucket := bucket._next
				key := ht._handler.Get_key(bucket._item)
				bkt := ht._handler.Hash_key(key) % nbucket2

				if bucket.isMutable() {
					bucket._next = newTable[bkt]
					newTable[bkt] = bucket
				} else {
					if nil == newTable[bkt] && nil == bucket._next {
						newTable[bkt] = bucket
					} else {
						cow_bucket := &Bucket{_ref: 1, _item: bucket._item}
						cow_bucket._next = newTable[bkt]
						newTable[bkt] = cow_bucket
						bucket.decRef()
					}
				}

				bucket = next_bucket
			}
		}
		ht._table = nil
	}
	ht._table = newTable
	ht._nbucket = nbucket2
	return true
}

// Erase all items
func (ht *CowHashTable) Clear() {
	if nil != ht._table {
		for i := 0; i < ht._nbucket; i++ {
			ht.chain_inc(ht._table[i], nil)
		}
		ht._table = nil
	}
}

//search match key from CowHashTable
//Return item address
func (ht *CowHashTable) Seek(key interface{}) interface{} {

	bkt := ht._handler.Hash_key(key) % ht._nbucket
	bucket := ht._table[bkt]
	for nil != bucket {
		if ht._handler.Eq_key(key, bucket._item) {
			return bucket._item
		}
		bucket = bucket._next
	}

	return nil
}
func (ht *CowHashTable) PrintInfo() string{
	return fmt.Sprintf("nbucket:%d ,factor:%d ,nitem:%d", ht._nbucket, ht._factor, ht._nitem)
}
//find next prime number
func find_prime(nbucket int) int {

	//copy from c++ stl
	var prime = [...]int{
		53, 97, 193, 389, 769,
		1543, 3079, 6151, 12289, 24593,
		49157, 98317, 196613, 393241, 786433,
		1572869, 3145739, 6291469, 12582917, 25165843,
		50331653, 100663319, 201326611, 402653189, 805306457,
		1610612741, 3221225473, 4294967291}

	len_prime := len(prime)
	for i := 0; i < len_prime; i++ {
		if prime[i] >= nbucket {
			return prime[i]
		}
	}
	return nbucket
}
