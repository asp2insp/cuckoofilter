package cuckoofilter

import (
	"testing"

	"github.com/ca-geo/go-misc/testutils"
)

func TestOneToOne(t *testing.T) {
	m := NewCuckooTable(1, 0, 1, 1)
	testutils.CheckUint(0, m.Size(), t)

	m.Insert([]byte("John Doe"))
	testutils.CheckUint(1, m.Size(), t)
	testutils.ExpectTrue(m.Lookup([]byte("John Doe")), "Table should contain John Doe", t)
	m.Delete([]byte("John Doe"))
	testutils.CheckUint(0, m.Size(), t)
	testutils.ExpectFalse(m.Lookup([]byte("John Doe")), "John Doe should have been deleted", t)
}

func TestOneBucket(t *testing.T) {
	m := NewCuckooTable(1, 0, 3, 4)
	testutils.CheckUint(0, m.Size(), t)

	m.Insert([]byte("John Doe"))
	testutils.CheckUint(1, m.Size(), t)
	testutils.ExpectTrue(m.Lookup([]byte("John Doe")), "Table should contain John Doe", t)

	m.Insert([]byte("Mary had a little lamb"))
	testutils.CheckUint(2, m.Size(), t)
	testutils.ExpectTrue(m.Lookup([]byte("Mary had a little lamb")), "Table should contain Mary", t)

	m.Insert([]byte("Fred Wilson"))
	testutils.CheckUint(3, m.Size(), t)
	testutils.ExpectTrue(m.Lookup([]byte("Fred Wilson")), "Table should contain Fred Wilson", t)

	err := m.Insert([]byte("Full!"))
	if err == nil {
		t.Error("Should have thrown an error inserting into full table")
	}

	m.Delete([]byte("Mary had a little lamb"))
	testutils.CheckUint(2, m.Size(), t)
	testutils.ExpectFalse(m.Lookup([]byte("Mary had a little lamb")), "Table shouldn't contain Mary", t)

	m.Delete([]byte("Fred Wilson"))
	testutils.CheckUint(1, m.Size(), t)
	testutils.ExpectFalse(m.Lookup([]byte("Fred Wilson")), "Table shouldn't contain Fred Wilson", t)
}

func TestInsert(t *testing.T) {
	m := NewCuckooTable(5, 100, 1, 4)
	testutils.CheckUint(0, m.Size(), t)

	m.Insert([]byte("John Doe"))
	testutils.CheckUint(1, m.Size(), t)

	m.Insert([]byte("Fred Wilson"))
	testutils.CheckUint(2, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("John Doe")), "Table should contain John Doe", t)
	testutils.ExpectTrue(m.Lookup([]byte("Fred Wilson")), "Table should contain Fred Wilson", t)
	testutils.ExpectFalse(m.Lookup([]byte("George")), "Table shouldn't contain George", t)
}

func TestDelete(t *testing.T) {
	m := NewCuckooTable(13, 100, 1, 3)

	var err error
	err = m.Insert([]byte("Red"))
	testutils.ExpectNil(err, "Insertion shouldn't cause error", t)
	err = m.Insert([]byte("Blood"))
	testutils.ExpectNil(err, "Insertion shouldn't cause error", t)
	err = m.Insert([]byte("Green"))
	testutils.ExpectNil(err, "Insertion shouldn't cause error", t)
	err = m.Insert([]byte("Money"))
	testutils.ExpectNil(err, "Insertion shouldn't cause error", t)
	err = m.Insert([]byte("Black"))
	testutils.ExpectNil(err, "Insertion shouldn't cause error", t)
	err = m.Insert([]byte("Night"))
	testutils.ExpectNil(err, "Insertion shouldn't cause error", t)
	err = m.Insert([]byte("Yellow"))
	testutils.ExpectNil(err, "Insertion shouldn't cause error", t)
	err = m.Insert([]byte("Sun"))
	testutils.ExpectNil(err, "Insertion shouldn't cause error", t)

	testutils.CheckUint(8, m.Size(), t)
	testutils.ExpectTrue(m.Lookup([]byte("Red")), "Table should contain Red", t)

	m.Delete([]byte("Red"))
	testutils.ExpectFalse(m.Lookup([]byte("Red")), "Red was removed from map", t)
	testutils.CheckUint(7, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Blood")), "Table should contain Blood", t)
	m.Delete([]byte("Blood"))
	testutils.ExpectFalse(m.Lookup([]byte("Blood")), "Blood was removed from map", t)
	testutils.CheckUint(6, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Green")), "Table should contain Green", t)
	m.Delete([]byte("Green"))
	testutils.ExpectFalse(m.Lookup([]byte("Green")), "Green was removed from map", t)
	testutils.CheckUint(5, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Money")), "Table should contain Money", t)
	m.Delete([]byte("Money"))
	testutils.ExpectFalse(m.Lookup([]byte("Money")), "Money was removed from map", t)
	testutils.CheckUint(4, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Black")), "Table should contain Black", t)
	m.Delete([]byte("Black"))
	testutils.ExpectFalse(m.Lookup([]byte("Black")), "Black was removed from map", t)
	testutils.CheckUint(3, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Night")), "Table should contain Night", t)
	m.Delete([]byte("Night"))
	testutils.ExpectFalse(m.Lookup([]byte("Night")), "Night was removed from map", t)
	testutils.CheckUint(2, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Yellow")), "Table should contain Yellow", t)
	m.Delete([]byte("Yellow"))
	testutils.ExpectFalse(m.Lookup([]byte("Yellow")), "Yellow was removed from map", t)
	testutils.CheckUint(1, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Sun")), "Table should contain Sun", t)
	m.Delete([]byte("Sun"))
	testutils.ExpectFalse(m.Lookup([]byte("Sun")), "Sun was removed from map", t)
	testutils.CheckUint(0, m.Size(), t)
}

func TestBucketFunctions(t *testing.T) {
	b := make([]Fingerprint, 1)
	f := myFingerprintFunc([]byte("John Doe"), 4)

	testutils.ExpectTrue(insert(f, b), "insert f failed", t)
	testutils.CheckByteSlice(f, b[0], t)
	testutils.ExpectTrue(contains(f, b), "Bucket should contain f", t)
	testutils.ExpectTrue(insert(f, b), "Insert of f failed", t)

	f2 := myFingerprintFunc([]byte("Fred"), 4)
	testutils.ExpectFalse(insert(f2, b), "Insert into full succeeded", t)

	deleteF(f, b)
	testutils.ExpectTrue(b[0] == nil, "Bucket should be empty", t)
	testutils.ExpectTrue(insert(f2, b), "Insert of f2 failed", t)
	testutils.ExpectTrue(contains(f2, b), "Bucket should contain f2", t)

	b = append(b, nil)
	testutils.ExpectTrue(insert(f2, b), "Second insert of f2 failed", t)
	testutils.ExpectTrue(contains(f2, b), "Bucket should contain f2", t)
	testutils.ExpectFalse(contains(f, b), "Bucket shouldn't contain f", t)
	testutils.ExpectTrue(insert(f, b), "insert f failed", t)
	testutils.ExpectTrue(contains(f, b), "Bucket should contain f", t)
}
