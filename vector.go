package main

import "math"

type Vector struct {
	id     string
	values []float64
}

func (this *Vector) Similarity(other *Vector) float64 {
	var dotProduct, lenA, lenB float64
	for i := range this.values {
		dotProduct += this.values[i] * other.values[i]
		lenA += this.values[i] * this.values[i]
		lenB += other.values[i] * other.values[i]
	}

	lenA = math.Sqrt(lenA)
	lenB = math.Sqrt(lenB)
	return dotProduct / (lenA * lenB)
}

func (this *Vector) Distance(other *Vector) float64 {
	var distance float64
	for i := range this.values {
		diff := other.values[i] - this.values[i]
		distance += diff * diff
	}
	return math.Sqrt(distance)
}
