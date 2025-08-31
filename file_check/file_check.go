package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
)

type HashTable struct {
	buckets [][]string
	size    int
}

func NewHashTable(size int) *HashTable {
	return &HashTable{
		buckets: make([][]string, size),
		size:    size,
	}
}

func (ht *HashTable) hash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % ht.size
}

func (ht *HashTable) Add(key string) {
	index := ht.hash(key)
	for _, v := range ht.buckets[index] {
		if v == key {
			return
		}
	}
	ht.buckets[index] = append(ht.buckets[index], key)
}

func (ht *HashTable) Contains(key string) bool {
	index := ht.hash(key)
	for _, v := range ht.buckets[index] {
		if v == key {
			return true
		}
	}
	return false
}

func readLinesToHashTable(filename string, ht *HashTable) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ht.Add(scanner.Text())
	}
	return scanner.Err()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <file1> <file2>\n", os.Args[0])
		return
	}

	file1 := os.Args[1]
	file2 := os.Args[2]

	ht1 := NewHashTable(100)
	ht2 := NewHashTable(100)

	// O(N1)
	if err := readLinesToHashTable(file1, ht1); err != nil {
		fmt.Printf("Error while reading %s: %v\n", file1, err)
		return
	}

	// O(N2)
	if err := readLinesToHashTable(file2, ht2); err != nil {
		fmt.Printf("Error while reading %s: %v\n", file2, err)
		return
	}

	// Амортизована складність O(1)
	for i := 0; i < 100; i++ {
		for _, v := range ht1.buckets[i] {
			if !ht2.Contains(v) {
				fmt.Println("Sets are not the same")
				return
			}
		}
	}

	// Амортизована складність O(1)
	for i := 0; i < 100; i++ {
		for _, v := range ht2.buckets[i] {
			if !ht1.Contains(v) {
				fmt.Println("Sets are not the same")
				return
			}
		}
	}

	fmt.Println("Sets are the same")
}
