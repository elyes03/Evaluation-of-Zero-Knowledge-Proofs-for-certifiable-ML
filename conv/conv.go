package main

import (
	"fmt"

	"github.com/petar/GoMNIST"
)

func main() {

	trainImages, testImages, err := GoMNIST.Load("./data./MNIST./raw")

	if err != nil {
		fmt.Println("Erreur lors du chargement du dataset MNIST:", err)
		return
	}

	// Utilisation des données (par exemple, affichage de la taille des ensembles de données)
	fmt.Println("Nombre d'images dans l'ensemble d'entraînement:", trainImages.Count())
	fmt.Println("Nombre d'images dans l'ensemble de test:", testImages.Count())

	firstImage := trainImages.Images[0]

	var image [784]int
	for i := range 784 {
		image[i] = int(firstImage[i])
	}
	kernel := []int{
		1, 0, 0, 1}

	// Appliquer la convolution
	result := convolution(image, kernel, 28, 2)

	// Afficher le résultat
	fmt.Println("Résultat de la convolution:")
	for _, val := range result {
		fmt.Print(val, " ")
	}
	fmt.Println()
}

// convolution applique la convolution à chaque sous-matrice 4x4 de l'image avec le noyau donné
func convolution(image [784]int, kernel []int, imageSize, kernelSize int) []int {
	resultSize := imageSize / kernelSize
	result := make([]int, resultSize*resultSize)

	for i := 0; i < resultSize; i++ {
		for j := 0; j < resultSize; j++ {
			sum := 0
			for k := 0; k < kernelSize; k++ {
				for l := 0; l < kernelSize; l++ {
					sum += image[(2*i+k)*imageSize+2*j+l] * kernel[k*kernelSize+l]
				}
			}
			result[i*resultSize+j] = sum
		}
	}

	var sum1 int = 0
	for i := range result {
		sum1 += result[i]
	}
	fmt.Println(sum1)
	return result
}
