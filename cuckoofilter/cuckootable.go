package cuckoofilter

import (
	"fmt"
	"math/rand"
)

type fingerprint []byte
type bucket []fingerprint

type CuckooTable interface {
	Insert([]byte) error  // O(1)
	Lookup([]byte) bool   // O(1)
	Delete([]byte)        // O(1)
	Size() int            // O(1)
	Utilization() float64 // O(1)
}

type ConfigurableCuckooTable struct {
	size, maxRetries, bucketSize int
	data                         []bucket
	fingerprinter                func([]byte) fingerprint
}

func NewCuckooTable(capacity, maxRetries, bucketSize int) *ConfigurableCuckooTable {
	ret := new(ConfigurableCuckooTable)
	ret.size = 0
	ret.data = make([]bucket, capacity)
	for i := 0; i < capacity; i++ {
		ret.data[i] = make([]fingerprint, bucketSize)
	}
	ret.fingerprinter = fingerprinter
	ret.maxRetries = maxRetries
	ret.bucketSize = bucketSize
	return ret
}

func (t *ConfigurableCuckooTable) Insert(item []byte) error {
	f := fingerprint(item)
	i1 := hash(item)
	i2 := Int.xor(i1, hash(f))
	if t.data[i1].insert(f) {
		return nil
	} else if t.data[i2].insert(f) {
		return nil
	} else {
		// Start moving things around
		var i int
		if rand.Int()%2 == 0 {
			i = i1
		} else {
			i = i2
		}
		var f2 = f
		for n := 0; n < t.maxRetries; n++ {
			var tempf = t.data[i][rand.Int()%t.bucketSize]
			t.data[i].delete(tempf)
			t.data[i].insert(f2)
			f2 = tempf

			i = Int.xor(i, hash(f2))
			if t.data[i].insert(f2) {
				return nil
			}
		}
	}
	return fmt.Errorf("Couldn't find a home for %v", f)
}

func (b *bucket) insert(f fingerprint) bool {
	for i, v := range b {
		if equal(v, f) {
			return true
		}
		if v == nil {
			b[i] = f
			return true
		}
	}
	return false
}

func (t *ConfigurableCuckooTable) Lookup(item []byte) bool {
	f := fingerprint(item)
	i1 := hash(item)
	i2 := Int.xor(i1, hash(f))
	if t.data[i1].contains(f) || t.data[i2].contains(f) {
		return true
	}
	return false
}

func (b *bucket) contains(f fingerprint) bool {
	for _, v := range b {
		if equal(v, f) {
			return true
		}
	}
	return false
}

func equal(f1, f2 fingerprint) bool {
	if len(f1) != len(f2) {
		return false
	}
	for i, _ := range f1 {
		if f1[i] != f2[i] {
			return false
		}
	}
	return true
}

func (t *ConfigurableCuckooTable) Delete(item []byte) {
	f := fingerprint(item)
	i1 := hash(item)
	i2 := Int.xor(i1, hash(f))
	if t.data[i1].contains(f) {
		t.data[i1].delete(f)
	}
	if t.data[i2].contains(f) {
		t.data[i2].delete(f)
	}
	return false
}

func (b *bucket) delete(f fingerprint) bool {
	for i, v := range b {
		if equal(v, f) {
			b[i] = nil
		}
	}
	return false
}

func (t *ConfigurableCuckooTable) Size() int {
	return size
}

func (t *ConfigurableCuckooTable) Utilization() float64 {
	return float64(size) / float64(len(data)*sBUCKET_SIZE)
}

func fingerprinter(value []byte) []byte {
	var fingerprintSize = 4

	var ret = make([fingerprintSize]byte, fingerprintSize)
	for i, v := range value {
		if ret[i] == 0 {
			ret[i] = 5381
		}
		ret[i] = ((ret[i] << 3) + ret[i]) + v
	}
	return ret
}

func hash(f fingerprint) uint {
	var hash uint = 5381

	for i := 0; i < len(f); i++ {
		hash = ((hash << 5) + hash) + uint(f[i])
	}
	return hash
}
