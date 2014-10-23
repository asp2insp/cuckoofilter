package cuckoofilter

import (
	"testing"

	"github.com/ca-geo/go-misc/testutils"
)

func TestPut(t *testing.T) {
	m := NewCuckooTable(5, 100, 1)
	testutils.CheckInt(0, m.Size(), t)

	m.Insert([]byte("John Doe"))
	testutils.CheckInt(1, m.Size(), t)

	m.Insert([]byte("Fred Wilson"))
	testutils.CheckInt(2, m.Size(), t)

	testutils.ExpectTrue(m.Lookup("John Doe"), "Table should contain John Doe", t)
	testutils.ExpectTrue(m.Lookup("Fred Wilson"), "Table should contain Fred Wilson", t)
	testutils.ExpectFalse(m.Lookup("George"), "Table shouldn't contain George", t)
}

func TestMapRemove(t *testing.T) {
	m := NewCuckooTable(5, 100, 2)

	m.Insert([]byte("Red"))
	m.Insert([]byte("Blood"))
	m.Insert([]byte("Green"))
	m.Insert([]byte("Money"))
	m.Insert([]byte("Black"))
	m.Insert([]byte("Night"))
	m.Insert([]byte("Yellow"))
	m.Insert([]byte("Sun"))

	testutils.CheckInt(8, m.Size(), t)
	testutils.ExpectTrue(m.Lookup("Red"), "Table should contain Red", t)

	m.Remove([]byte("Red"))
	testutils.ExpectFalse(m.Lookup([]byte("Red")), "Red was removed from map", t)
	testutils.CheckInt(7, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Blood")), "Table should contain Blood", t)
	m.Remove([]byte("Blood"))
	testutils.ExpectFalse(m.Lookup([]byte("Blood")), "Blood was removed from map", t)
	testutils.CheckInt(6, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Green")), "Table should contain Green", t)
	m.Remove([]byte("Green"))
	testutils.ExpectFalse(m.Lookup([]byte("Green")), "Green was removed from map", t)
	testutils.CheckInt(5, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Money")), "Table should contain Money", t)
	m.Remove([]byte("Money"))
	testutils.ExpectFalse(m.Lookup([]byte("Money")), "Money was removed from map", t)
	testutils.CheckInt(4, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Black")), "Table should contain Black", t)
	m.Remove([]byte("Black"))
	testutils.ExpectFalse(m.Lookup([]byte("Black")), "Black was removed from map", t)
	testutils.CheckInt(3, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Night")), "Table should contain Night", t)
	m.Remove([]byte("Night"))
	testutils.ExpectFalse(m.Lookup([]byte("Night")), "Night was removed from map", t)
	testutils.CheckInt(2, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Yellow")), "Table should contain Yellow", t)
	m.Remove([]byte("Yellow"))
	testutils.ExpectFalse(m.Lookup([]byte("Yellow")), "Yellow was removed from map", t)
	testutils.CheckInt(1, m.Size(), t)

	testutils.ExpectTrue(m.Lookup([]byte("Sun")), "Table should contain Sun", t)
	m.Remove([]byte("Sun"))
	testutils.ExpectFalse(m.Lookup([]byte("Sun")), "Sun was removed from map", t)
	testutils.CheckInt(0, m.Size(), t)
}
