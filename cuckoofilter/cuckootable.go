package cuckoofilter

import (
	"fmt"
	"math/rand"
)

type fingerprint []byte

type CuckooTable interface {
	Insert([]byte) error  // O(1)
	Lookup([]byte) bool   // O(1)
	Delete([]byte)        // O(1)
	Size() int            // O(1)
	Utilization() float64 // O(1)
}

type ConfigurableCuckooTable struct {
	size, maxRetries, bucketSize int
	data                         [][]fingerprint
	fingerprinter                func([]byte) fingerprint
	capacity                     uint
}

func NewCuckooTable(capacity, maxRetries, bucketSize int) *ConfigurableCuckooTable {
	ret := new(ConfigurableCuckooTable)
	ret.size = 0
	ret.data = make([][]fingerprint, capacity)
	for i := 0; i < capacity; i++ {
		ret.data[i] = make([]fingerprint, bucketSize)
	}
	ret.fingerprinter = fingerprinter
	ret.maxRetries = maxRetries
	ret.bucketSize = bucketSize
	ret.capacity = uint(capacity)
	return ret
}

func (t *ConfigurableCuckooTable) Insert(item []byte) error {
	f := fingerprint(item)
	i1 := hash(item) % t.capacity
	i2 := (i1 ^ hash(f)) % t.capacity
	success := false
	if insert(f, t.data[i1]) {
		success = true
	} else if insert(f, t.data[i2]) {
		success = true
	} else {
		// Start moving things around
		var i uint
		if rand.Int()%2 == 0 {
			i = i1
		} else {
			i = i2
		}
		var f2 = f
		for n := 0; n < t.maxRetries; n++ {
			tempf := make([]byte, 4)
			randElem := rand.Uint32() % uint32(t.bucketSize)
			copy(tempf, t.data[i][randElem])
			copy(t.data[i][randElem], f2)
			copy(f2, tempf)

			i = (i ^ hash(f2)) % t.capacity
			if insert(f2, t.data[i]) {
				success = true
			}
		}
	}
	if success {
		t.size += 1
		return nil
	} else {
		return fmt.Errorf("Couldn't find a home for %v", f)
	}
}

func insert(f fingerprint, b []fingerprint) bool {
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
	i1 := hash(item) % t.capacity
	i2 := (i1 ^ hash(f)) % t.capacity
	if contains(f, t.data[i1]) || contains(f, t.data[i2]) {
		return true
	}
	return false
}

func contains(f fingerprint, b []fingerprint) bool {
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
	i1 := hash(item) % t.capacity
	i2 := (i1 ^ hash(f)) % t.capacity
	if contains(f, t.data[i1]) {
		delete(f, t.data[i1])
		t.size -= 1
	}
	if contains(f, t.data[i2]) {
		delete(f, t.data[i2])
		t.size -= 1
	}
}

func delete(f fingerprint, b []fingerprint) {
	for i, v := range b {
		if equal(v, f) {
			b[i] = nil
		}
	}
}

func (t *ConfigurableCuckooTable) Size() int {
	return t.size
}

func (t *ConfigurableCuckooTable) Utilization() float64 {
	return float64(t.size) / float64(len(t.data)*t.bucketSize)
}

func fingerprinter(value []byte) fingerprint {
	var fingerprintSize = 4

	var ret = make([]byte, fingerprintSize)
	for i, v := range value {
		if ret[i] == 0 {
			ret[i] = 31
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
