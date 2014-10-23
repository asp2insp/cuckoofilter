package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/asp2insp/cuckoofilter/cuckoofilter"
)

func main() {
	dict := cuckoofilter.NewCuckooTable(260000, 500, 1)
	fd, err := os.Open("/usr/share/dict/web2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)
	collisions := 0
	total := 0
	missed := 0
	for scanner.Scan() {
		s := []byte(scanner.Text())
		if dict.Lookup(s) {
			collisions += 1
		}
		dict.Insert(s)
		if !dict.Lookup(s) {
			missed += 1
		}
		total += 1
	}
	collisionPercent := float64(collisions) / float64(total)
	fmt.Printf("Added %v items to the filter with %v false positives, %v missed, and final utilization of %v\n",
		total, collisionPercent, missed, dict.Utilization())
}
