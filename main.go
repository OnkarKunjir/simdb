package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"strconv"
)

const M = 16
const efconstruction = M * 4
const efsearch = 40
const vectorSize = 500

var randGenerator = rand.New(rand.NewSource(42069))

func generateDataset(size, n int) []*Vector {
	vectors := make([]*Vector, 0, n)
	for id := range n {
		values := make([]float64, size)
		for i := range values {
			values[i] = randGenerator.Float64()

		}

		vectors = append(vectors, &Vector{
			id:     strconv.Itoa(id),
			values: values,
		})
	}
	return vectors
}

func knn(vector *Vector, hnsw *HierarchicalNavigableSmallWorld, k int) []*Document {
	queue := make(PriorityQueue, 0, k)
	for _, node := range hnsw.nodes {
		queuedNode := &QueuedNode{node, -vector.Distance(node.vector)}
		if queue.Len() < k {
			heap.Push(&queue, queuedNode)
		} else if queue[0].score < queuedNode.score {
			heap.Pop(&queue)
			heap.Push(&queue, queuedNode)
		}
	}

	documents := make([]*Document, len(queue))
	for index := len(queue) - 1; index >= 0; index-- {
		queuedNode := heap.Pop(&queue).(*QueuedNode)
		documents[index] = &Document{queuedNode.node.vector, -queuedNode.score}
	}

	return documents
}

func benchmark(hnsw *HierarchicalNavigableSmallWorld, n, k int) {
	var averageRecall float64

	for _, vector := range generateDataset(vectorSize, n) {
		expectedIds := make(map[string]struct{})
		for _, document := range knn(vector, hnsw, k) {
			expectedIds[document.vector.id] = struct{}{}
		}

		count := 0
		for _, document := range hnsw.Search(vector, k, efsearch) {
			if _, ok := expectedIds[document.vector.id]; ok {
				count++
			}
		}
		averageRecall += float64(count) / 10
	}

	averageRecall = averageRecall / float64(n)
	fmt.Printf("Average recall: %0.2f\n", averageRecall)
}

func main() {
	hnsw := &HierarchicalNavigableSmallWorld{}

	fmt.Println("Inserting data")
	for _, vector := range generateDataset(vectorSize, 500) {
		hnsw.Insert(vector, M, efconstruction)
	}
	fmt.Println("Data insert finished")

	benchmark(hnsw, 1000, 10)

	// // fmt.Println(hnsw)
	// toSearch := &Vector{id: "2", values: []float64{2, 0}}
	// fmt.Println(hnsw.Search(toSearch, 10, efsearch))
	// fmt.Println(knn(toSearch, hnsw, 10))
}
