package cuckoofilter

import (
	"testing"

	"github.com/ca-geo/go-misc/testutils"
)

func TestInsert(t *testing.T) {
	m := NewCuckooTable(5, 100, 1)
	testutils.CheckInt(0, m.Size(), t)

	m.Insert([]byte("John Doe"))
	testutils.CheckInt(1, m.Size(), t)

	m.Insert([]byte("Fred Wilson"))
	testutils.CheckInt(2, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("John Doe")), "Table should contain John Doe", t)
	testutils.ExpectTrue(m.Lookup([]byte("Fred Wilson")), "Table should contain Fred Wilson", t)
	testutils.ExpectFalse(m.Lookup([]byte("George")), "Table shouldn't contain George", t)
}

func TestDelete(t *testing.T) {
	m := NewCuckooTable(10, 100, 1)

	m.Insert([]byte("Red"))
	m.Insert([]byte("Blood"))
	m.Insert([]byte("Green"))
	m.Insert([]byte("Money"))
	m.Insert([]byte("Black"))
	m.Insert([]byte("Night"))
	m.Insert([]byte("Yellow"))
	m.Insert([]byte("Sun"))

	testutils.CheckInt(8, m.Size(), t)
	testutils.ExpectTrue(m.Lookup([]byte("Red")), "Table should contain Red", t)

	m.Delete([]byte("Red"))
	testutils.ExpectFalse(m.Lookup([]byte("Red")), "Red was removed from map", t)
	testutils.CheckInt(7, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Blood")), "Table should contain Blood", t)
	m.Delete([]byte("Blood"))
	testutils.ExpectFalse(m.Lookup([]byte("Blood")), "Blood was removed from map", t)
	testutils.CheckInt(6, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Green")), "Table should contain Green", t)
	m.Delete([]byte("Green"))
	testutils.ExpectFalse(m.Lookup([]byte("Green")), "Green was removed from map", t)
	testutils.CheckInt(5, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Money")), "Table should contain Money", t)
	m.Delete([]byte("Money"))
	testutils.ExpectFalse(m.Lookup([]byte("Money")), "Money was removed from map", t)
	testutils.CheckInt(4, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Black")), "Table should contain Black", t)
	m.Delete([]byte("Black"))
	testutils.ExpectFalse(m.Lookup([]byte("Black")), "Black was removed from map", t)
	testutils.CheckInt(3, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Night")), "Table should contain Night", t)
	m.Delete([]byte("Night"))
	testutils.ExpectFalse(m.Lookup([]byte("Night")), "Night was removed from map", t)
	testutils.CheckInt(2, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Yellow")), "Table should contain Yellow", t)
	m.Delete([]byte("Yellow"))
	testutils.ExpectFalse(m.Lookup([]byte("Yellow")), "Yellow was removed from map", t)
	testutils.CheckInt(1, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Sun")), "Table should contain Sun", t)
	m.Delete([]byte("Sun"))
	testutils.ExpectFalse(m.Lookup([]byte("Sun")), "Sun was removed from map", t)
	testutils.CheckInt(0, m.Size(), t)
}
