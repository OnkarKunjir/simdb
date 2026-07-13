package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

func generateDataset(randGenerator *rand.Rand, n int) []*Vector {
	vectors := make([]*Vector, 0, n)
	for id := range n {
		values := [VectorDimension]float64{}
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

func knn(of *Vector, vectors []*Vector, k int) []*ScoredVector {
	queue := make(VectorQueue, 0, k)
	heap.Init(&queue)
	for _, vector := range vectors {
		scoredVector := &ScoredVector{vector, similarity(of, vector)}
		if len(queue) < k {
			heap.Push(&queue, scoredVector)
			continue
		}

		if queue[k-1].score < scoredVector.score {
			heap.Pop(&queue)
			heap.Push(&queue, scoredVector)
		}
	}

	sort.Sort(sort.Reverse(&queue))
	return queue
}

func main() {
	randGenerator := rand.New(rand.NewSource(42069))
	vectors := generateDataset(randGenerator, 10)
	fmt.Println(knn(vectors[0], vectors, 3))
}
