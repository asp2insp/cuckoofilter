package cuckoofilter

import (
	"fmt"
	"math"
	"math/rand"
)

type fingerprint []byte

type CuckooTable interface {
	Insert([]byte) error                                           // O(1)
	Lookup([]byte) bool                                            // O(1)
	Delete([]byte)                                                 // O(1)
	Size() int                                                     // O(1)
	Stats() (utilization, rebucketRatio, compressionRatio float64) // O(1)
}

type ConfigurableCuckooTable struct {
	data          [][]fingerprint
	fingerprinter func([]byte, uint) fingerprint
	capacity,
	rebuckets,
	bucketSize,
	maxRetries,
	size,
	fingerprintSize,
	bytesRepresented uint
}

func NewCuckooTable(capacity, maxRetries, bucketSize, fingerprintSize uint) *ConfigurableCuckooTable {
	ret := new(ConfigurableCuckooTable)
	ret.size = 0

	// Allocate all memory up front
	ret.data = make([][]fingerprint, capacity)
	var i uint
	for i = 0; i < capacity; i++ {
		ret.data[i] = make([]fingerprint, bucketSize)
	}
	ret.fingerprinter = myFingerprintFunc
	ret.maxRetries = maxRetries
	ret.bucketSize = bucketSize
	ret.fingerprintSize = fingerprintSize
	ret.capacity = capacity
	ret.rebuckets = 0
	ret.bytesRepresented = 0
	return ret
}

func (t *ConfigurableCuckooTable) Insert(item []byte) error {
	f, i1, i2 := t.getIndices(item)
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
		var n uint
		tempf := make([]byte, t.fingerprintSize)
		for n = 0; n < t.maxRetries; n++ {
			// Swap into the bucket for a random element
			randElem := rand.Uint32() % uint32(t.bucketSize)
			copy(tempf, t.data[i][randElem])
			copy(t.data[i][randElem], f)
			copy(f, tempf)

			// look in the alternate location for that random element
			i = (i ^ hash(f)) % t.capacity
			t.rebuckets += 1
			if insert(f, t.data[i]) {
				success = true
				break
			}
		}
	}
	if success {
		// fmt.Printf("Inserting %v(%v) %v/%v\n", string(item), f, i1, i2)
		t.size += 1
		t.bytesRepresented += uint(len(item))
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
			f2 := make([]byte, len(f))
			copy(f2, f)
			b[i] = f2
			return true
		}
	}
	return false
}

func (t *ConfigurableCuckooTable) Lookup(item []byte) bool {
	f, i1, i2 := t.getIndices(item)
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
	f, i1, i2 := t.getIndices(item)
	// fmt.Printf("Deleting %v %v/%v\n", string(item), t.data[i1], t.data[i2])

	if contains(f, t.data[i1]) {
		deleteF(f, t.data[i1])

		t.size -= 1
		t.bytesRepresented -= uint(len(item))
	} else if contains(f, t.data[i2]) {
		deleteF(f, t.data[i2])
		t.size -= 1
		t.bytesRepresented -= uint(len(item))
	}
}

func (t *ConfigurableCuckooTable) getIndices(item []byte) (f []byte, i1, i2 uint) {
	f = t.fingerprinter(item, t.fingerprintSize)
	i1 = hash(item)
	i2 = (i1 ^ hash(f)) % t.capacity
	i1 = i1 % t.capacity
	return
}

func deleteF(f fingerprint, b []fingerprint) {
	for i, v := range b {
		if equal(v, f) {
			b[i] = nil
		}
	}
}

func (t *ConfigurableCuckooTable) Size() uint {
	return t.size
}

func (t *ConfigurableCuckooTable) Stats() (utilization, rebucketRatio, compressionRatio float64) {
	utilization = float64(t.size) / float64(t.capacity*t.bucketSize)
	rebucketRatio = float64(t.rebuckets) / float64(t.size)
	compressionRatio = float64(t.bytesRepresented) / float64(t.capacity*t.bucketSize*t.fingerprintSize)
	return
}

// Adapted from https://moinakg.wordpress.com/tag/rabin-fingerprint/
func myFingerprintFunc(value []byte, fingerprintSize uint) fingerprint {
	windowsize := 8
	mod := 255
	prime := 71
	pow := int(math.Pow(float64(prime), float64(windowsize))) % mod
	ret := make([]byte, fingerprintSize)
	for j, _ := range ret {
		// Seed with prime number
		rollhash := 149 // 1001 0101
		outbyte := 0
		for i, v := range value {
			// rollhash = (rollhash * PRIME + inbyte - outbyte * POW) % MODULUS
			// Include j so that each iteration is different
			rollhash = (rollhash*prime + int(v) + j*prime - outbyte*pow) % mod
			if i > windowsize {
				outbyte = int(v)
			}
		}
		ret[j] = byte(rollhash)
	}
	return ret
}

func hash(f fingerprint) uint {
	var h uint = 5381

	for i := 0; i < len(f); i++ {
		h = ((h << 5) + h) + uint(f[i])
	}
	return h
}
