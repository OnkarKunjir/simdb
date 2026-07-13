package main

import (
	"fmt"
	"math"
)

const VectorDimension = 10

type Vector struct {
	id     string
	values [VectorDimension]float64
}

func similarity(a, b *Vector) float64 {
	var dotProduct, lenA, lenB float64
	for i := range VectorDimension {
		dotProduct += a.values[i] * b.values[i]
		lenA += a.values[i] * a.values[i]
		lenB += b.values[i] * b.values[i]
	}

	lenA = math.Sqrt(lenA)
	lenB = math.Sqrt(lenB)

	return dotProduct / (lenA * lenB)
}

type ScoredVector struct {
	vector *Vector
	score  float64
}

func (scoredVector *ScoredVector) String() string {
	return fmt.Sprintf("ScoredVector{id: %s, score: %0.4f}", scoredVector.vector.id, scoredVector.score)
}

type VectorQueue []*ScoredVector

func (queue VectorQueue) Len() int { return len(queue) }

func (queue VectorQueue) Less(i, j int) bool {
	return queue[i].score < queue[j].score
}

func (queue VectorQueue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

func (queue *VectorQueue) Push(x any) {
	*queue = append(*queue, x.(*ScoredVector))
}

func (queue *VectorQueue) Pop() any {
	old := *queue
	n := len(old)
	item := old[n-1]
	old[n-1] = nil

	*queue = old[0 : n-1]
	return item
}
