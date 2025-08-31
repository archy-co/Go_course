package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func readMatrix(file *os.File) ([][]int, error) {
	scanner := bufio.NewScanner(file)
	var matrix [][]int
	for scanner.Scan() {
		row := strings.Fields(scanner.Text())
		intRow := make([]int, len(row))
		for i, v := range row {
			num, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			intRow[i] = num
		}
		matrix = append(matrix, intRow)
	}
	return matrix, scanner.Err()
}

func multiplyMatrices(matA, matB [][]int) [][]int {
	n, m, p := len(matA), len(matA[0]), len(matB[0])
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, p)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		for j := 0; j < p; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				sum := 0
				for k := 0; k < m; k++ {
					sum += matA[i][k] * matB[k][j]
				}
				result[i][j] = sum
			}(i, j)
		}
	}
	wg.Wait()
	return result
}

func writeMatrix(file *os.File, matrix [][]int) error {
	writer := bufio.NewWriter(file)
	for _, row := range matrix {
		for _, val := range row {
			_, err := writer.WriteString(fmt.Sprintf("%d ", val))
			if err != nil {
				return err
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <matrixA.txt> <matrixB.txt> <result.txt>")
		return
	}
	matrixAFile := os.Args[1]
	matrixBFile := os.Args[2]
	resultFile := os.Args[3]

	fileA, err := os.Open(matrixAFile)
	if err != nil {
		fmt.Println("Error opening matrixA.txt:", err)
		return
	}
	defer fileA.Close()

	fileB, err := os.Open(matrixBFile)
	if err != nil {
		fmt.Println("Error opening matrixB.txt:", err)
		return
	}
	defer fileB.Close()

	matA, err := readMatrix(fileA)
	if err != nil {
		fmt.Println("Error reading matrix A:", err)
		return
	}

	matB, err := readMatrix(fileB)
	if err != nil {
		fmt.Println("Error reading matrix B:", err)
		return
	}

	if len(matA[0]) != len(matB) {
		fmt.Println("Matrices cannot be multiplied")
		return
	}

	result := multiplyMatrices(matA, matB)

	fmt.Println("Resulting Matrix:")
	for _, row := range result {
		fmt.Println(row)
	}

	outFile, err := os.Create(resultFile)
	if err != nil {
		fmt.Println("Error creating result.txt:", err)
		return
	}
	defer outFile.Close()

	err = writeMatrix(outFile, result)
	if err != nil {
		fmt.Println("Error writing to result.txt:", err)
	}
}
