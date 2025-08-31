package main

import (
	"fmt"
)

type Graph struct {
	adjList map[int][]int // map with integer keys and list of integer values
}

func newGraph() *Graph {
	return &Graph{adjList: make(map[int][]int)}
}

func (g *Graph) addEdge(v, u int) {
	g.adjList[v] = append(g.adjList[v], u)
}

// set name and list of integers (vertices)
type NamedSets map[string][]int

func dependencyIndex(g *Graph, setM, setN []int) int {
	count := 0
	for _, v := range setM {
		for _, u := range setN {
			for _, neighbor := range g.adjList[v] {
				if neighbor == u {
					count++
				}
			}
		}
	}
	return count
}

func generateWeightedGraph(g *Graph, sets NamedSets) map[string]map[string]int {
	// nested map. outer map with string keys (set M) and values which are maps with string keys (set N) and int values (dependency index)
	weightedGraph := make(map[string]map[string]int)
	for mName, setM := range sets {
		for nName, setN := range sets {
			if mName != nName {
				index := dependencyIndex(g, setM, setN)
				if index > 0 {
					// check if there's already entry for mName in the weightedGraph map. If not, initialize it
					if _, exists := weightedGraph[mName]; !exists {
						weightedGraph[mName] = make(map[string]int)
					}
					weightedGraph[mName][nName] = index
				}
			}
		}
	}
	return weightedGraph
}

func main() {
	var n, m int

	fmt.Print("Кількість вершин: ")
	fmt.Scan(&n)

	fmt.Print("Кількість ребер: ")
	fmt.Scan(&m)

	graph := newGraph()

	fmt.Println("Ребра графа (v u):")
	for i := 0; i < m; i++ {
		var v, u int
		fmt.Scan(&v, &u)
		graph.addEdge(v, u)
	}

	sets := make(NamedSets)
	var setName string
	for {
		fmt.Print("Назва множини (або 'finish' щоб завершити): ")
		fmt.Scan(&setName)
		if setName == "finish" {
			break
		}
		var k int
		fmt.Printf("Кількість елементів у множині %s: ", setName)
		fmt.Scan(&k)
		set := make([]int, k)
		fmt.Printf("Улементи множини %s: ", setName)
		for i := 0; i < k; i++ {
			fmt.Scan(&set[i])
		}
		sets[setName] = set
	}

	weightedGraph := generateWeightedGraph(graph, sets)

	fmt.Println("Зважений орієнтований граф:")
	for mName, neighbors := range weightedGraph {
		for nName, weight := range neighbors {
			fmt.Printf("%s -> %s (вага: %d)\n", mName, nName, weight)
		}
	}
}
