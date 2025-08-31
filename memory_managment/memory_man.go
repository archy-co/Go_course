package main

import (
	"fmt"
	"os"
	"sort"
)

type FreeBlock struct {
	start, size int
}

type MemoryManager struct {
	memory       []int
	blockCounter int
	freeList     []FreeBlock
}

func NewMemoryManager(numCells int) *MemoryManager {
	return &MemoryManager{
		memory:       make([]int, numCells),
		blockCounter: 1,
		freeList:     []FreeBlock{{0, numCells}},
	}
}

func (mm *MemoryManager) allocate(numCells int) int {
	for i := 0; i < len(mm.freeList); i++ {
		block := &mm.freeList[i]
		if block.size >= numCells {
			start := block.start
			for j := start; j < start+numCells; j++ {
				mm.memory[j] = mm.blockCounter
			}
			mm.blockCounter++

			block.start += numCells
			block.size -= numCells
			if block.size == 0 {
				mm.freeList = append(mm.freeList[:i], mm.freeList[i+1:]...)
			}
			return start
		}
	}
	return -1
}

func (mm *MemoryManager) free(blockID int) {
	start, size := -1, 0
	for i := 0; i < len(mm.memory); i++ {
		if mm.memory[i] == blockID {
			if start == -1 {
				start = i
			}
			mm.memory[i] = 0
			size++
		} else if start != -1 {
			break
		}
	}

	mm.freeList = append(mm.freeList, FreeBlock{start, size})
	mm.coalesceFreeBlocks()
}

func (mm *MemoryManager) coalesceFreeBlocks() {
	if len(mm.freeList) < 2 {
		return
	}
	sort.Slice(mm.freeList, func(i, j int) bool {
		return mm.freeList[i].start < mm.freeList[j].start
	})

	merged := []FreeBlock{mm.freeList[0]}
	for i := 1; i < len(mm.freeList); i++ {
		last := &merged[len(merged)-1]
		current := mm.freeList[i]
		if last.start+last.size == current.start {
			last.size += current.size
		} else {
			merged = append(merged, current)
		}
	}
	mm.freeList = merged
}

func (mm *MemoryManager) print(maxPerLine int) {
	for i := 0; i < len(mm.memory); i++ {
		fmt.Printf("%d ", mm.memory[i])
		if (i+1)%maxPerLine == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func main() {
	fmt.Print("******** MEMORY MANAGER ********\n")

	var numCells, maxPerLine int
	fmt.Print("Number of memory cells\n> ")
	fmt.Scan(&numCells)
	fmt.Print("Max output width\n> ")
	fmt.Scan(&maxPerLine)

	mm := NewMemoryManager(numCells)

	for {
		fmt.Print("> ")
		var command string
		fmt.Scan(&command)

		switch command {
		case "allocate":
			var numCells int
			fmt.Scan(&numCells)
			blockID := mm.allocate(numCells)
			if blockID == -1 {
				fmt.Println("Not enough memory to allocate")
			} else {
				fmt.Printf("Block allocated at index %d\n", blockID)
			}
		case "free":
			var blockID int
			fmt.Scan(&blockID)
			mm.free(blockID)
			fmt.Println("Block freed")
		case "print":
			mm.print(maxPerLine)
		case "exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		case "help":
			fmt.Println("Commands:")
			fmt.Println(" allocate <num_cells> - allocate a memory block of <num_cells> cells")
			fmt.Println(" free <block_id> - free a block with id <block_id>")
			fmt.Println(" print - print the current memory layout")
			fmt.Println(" exit - exit the program")
		default:
			fmt.Println("Unknown command. Type 'help' for the list of commands.")
		}
	}
}
