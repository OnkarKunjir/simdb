package core

import "math"

type Vector struct {
	Id     string
	Values []float64
}

func (this *Vector) Similarity(other *Vector) float64 {
	var dotProduct, lenA, lenB float64
	for i := range this.Values {
		dotProduct += this.Values[i] * other.Values[i]
		lenA += this.Values[i] * this.Values[i]
		lenB += other.Values[i] * other.Values[i]
	}

	lenA = math.Sqrt(lenA)
	lenB = math.Sqrt(lenB)
	return dotProduct / (lenA * lenB)
}

func (this *Vector) Distance(other *Vector) float64 {
	var distance float64
	for i := range this.Values {
		diff := other.Values[i] - this.Values[i]
		distance += diff * diff
	}
	return math.Sqrt(distance)
}
