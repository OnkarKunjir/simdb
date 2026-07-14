package main

import (
	"fmt"
	"math/rand"
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

func main() {
	randGenerator := rand.New(rand.NewSource(42069))
	generateDataset(randGenerator, 10)

	nsw := &NavigableSmallWorld{}

	vectorA := &Vector{id: "0", values: [VectorDimension]float64{0, 0}}
	vectorB := &Vector{id: "1", values: [VectorDimension]float64{1, 0}}
	nsw.Insert(vectorA, 10, 1)
	nsw.Insert(vectorB, 10, 1)
	nsw.Insert(&Vector{"2", [VectorDimension]float64{2, 0}}, 10, 1)

	fmt.Println(nsw)
	fmt.Println(nsw.Search(&Vector{id: "2", values: [VectorDimension]float64{2, 0}}, 2, 1))
}
