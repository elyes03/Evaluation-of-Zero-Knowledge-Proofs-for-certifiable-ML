package main

import (
	"fmt"
	"testing"
)

func BenchmarkConvCircuit(b *testing.B) {

	image := []int{
		1, 2, 3, 4, 5, 6, 7, 8,
		9, 10, 11, 12, 13, 14, 15, 16,
		17, 18, 19, 20, 21, 22, 23, 24,
		25, 26, 27, 28, 29, 30, 31, 32,
		33, 34, 35, 36, 37, 38, 39, 40,
		41, 42, 43, 44, 45, 46, 47, 48,
		49, 50, 51, 52, 53, 54, 55, 56,
		57, 58, 59, 60, 61, 62, 63, 64,
	}

	kernel := []int{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}

	// Initialize the gnark constraint system
	b.ResetTimer()
	// Define the circuit inputs
	for i := 0; i < b.N; i++ {

		var imageSize = 8
		var kernelSize = 4
		resultSize := imageSize - kernelSize + 1
		result := make([]int, resultSize*resultSize)

		for i := 0; i < resultSize; i++ {
			for j := 0; j < resultSize; j++ {
				sum := 0
				for k := 0; k < kernelSize; k++ {
					for l := 0; l < kernelSize; l++ {
						sum += image[(i+k)*imageSize+j+l] * kernel[k*kernelSize+l]
					}
				}
				result[i*resultSize+j] = sum
			}
		}

		fmt.Println(result)
	}
}
