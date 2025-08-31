package main

import (
	"fmt"
	"os"
)

type MemoryManager struct {
	memory       []int
	blockCounter int
}

func NewMemoryManager(numCells int) *MemoryManager {
	return &MemoryManager{
		memory:       make([]int, numCells),
		blockCounter: 1,
	}
}

func (mm *MemoryManager) allocate(numCells int) int {
	start := -1
	count := 0

	for i := 0; i < len(mm.memory); i++ {
		if mm.memory[i] == 0 {
			if count == 0 {
				start = i
			}
			count++
		} else {
			start = -1
			count = 0
		}

		if count == numCells {
			for j := start; j < start+numCells; j++ {
				mm.memory[j] = mm.blockCounter
			}
			mm.blockCounter++
			return start
		}
	}

	return -1
}

func (mm *MemoryManager) free(blockID int) {
	for i := 0; i < len(mm.memory); i++ {
		if mm.memory[i] == blockID {
			mm.memory[i] = 0
		}
	}
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
