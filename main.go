package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

const M = 2
const efConstruction = 10
const efSearch = 10

func generateDataset(randGenerator *rand.Rand, size, n int) []*Vector {
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

func main() {
	nsw := &NavigableSmallWorld{}
	// randGenerator := rand.New(rand.NewSource(42069))
	// for _, vector := range generateDataset(randGenerator, 2, 10) {
	// 	nsw.Insert(vector, M, efConstruction)
	// }

	vectorA := &Vector{id: "0", values: []float64{0, 0}}
	vectorB := &Vector{id: "1", values: []float64{1, 0}}
	nsw.Insert(vectorA, M, efConstruction)
	nsw.Insert(vectorB, M, efConstruction)
	nsw.Insert(&Vector{"2", []float64{2, 0}}, M, efConstruction)

	fmt.Println(nsw)
	fmt.Println(nsw.Search(&Vector{id: "2", values: []float64{2, 0}}, 2, efSearch))
}
