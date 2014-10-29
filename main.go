package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/asp2insp/cuckoofilter/cuckoofilter"
)

func main() {
	runWordsBench()
}

func runWordsBench() {
	dict := cuckoofilter.NewCuckooTable(600000, 500, 4, 2)
	fd, err := os.Open("/usr/share/dict/web2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)
	total := 0
	missed := 0
	errors := 0
	for scanner.Scan() {
		s := []byte(scanner.Text())
		err = dict.Insert(s)
		if err != nil {
			errors += 1
		} else if !dict.Lookup(s) {
			missed += 1
			fmt.Println("FALSE NEGATIVE:", string(s))
		}
		total += 1
	}
	fd, err = os.Open("/usr/share/dict/web2a")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	collisions := 0
	otherSetSize := 0
	scanner = bufio.NewScanner(fd)
	for scanner.Scan() {
		s := []byte(scanner.Text())
		if dict.Lookup(s) {
			collisions += 1
		}
		otherSetSize += 1
	}
	collisionPercent := 100.0 * float64(collisions) / float64(otherSetSize)
	fmt.Printf("Added %v/%v items to the filter with %v/%v (%.2f%%) false positives, %v false negatives, %v errors.\n",
		dict.Size(), total, collisions, otherSetSize, collisionPercent, missed, errors)
	ut, rr, cr := dict.Stats()
	fmt.Printf("Final utilization %.2v, rebuckets %.2f, compression ratio %.3v:1\n", ut, rr, cr)
}
