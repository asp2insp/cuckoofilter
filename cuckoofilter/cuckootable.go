package cuckoofilter

const sBUCKET_SIZE

type fingerprint []byte
type bucket [sBUCKET_SIZE]fingerprint

type CuckooTable interface {
	Insert([]byte) error  // O(1)
	Lookup([]byte) bool   // O(1)
	Delete([]byte)        // O(1)
	Size() int            // O(1)
	Utilization() float64 // O(1)
}

type ConfigurableCuckooTable struct {
	size          int
	data          []bucket
	fingerprinter func([]byte) fingerprint
}

func (t *ConfigurableCuckooTable) Insert(item []byte) error {
	return nil
}

func (t *ConfigurableCuckooTable) Lookup(item []byte) bool {

	return false
}

func (t *ConfigurableCuckooTable) Delete(item []byte) {

}

func (t *ConfigurableCuckooTable) Size() int {
	return size
}

func (t *ConfigurableCuckooTable) Utilization() float64 {
	return float64(size) / float64(len(data)*sBUCKET_SIZE)
}


func fingerprint(value []byte) []byte {
	var fingerprintSize = 4
	var hash uint = 5381

	var toReturn = make([fingerprintSize]byte)
	for var i, v := range value {
		
	}
	return toReturn
}