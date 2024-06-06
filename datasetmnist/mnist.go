package main

import (
	"fmt"

	"github.com/petar/GoMNIST"
)

func main() {
	// Téléchargement du dataset MNIST
	trainImages, testImages, err := GoMNIST.Load("./data./MNIST./raw")

	if err != nil {
		fmt.Println("Erreur lors du chargement du dataset MNIST:", err)
		return
	}

	// Utilisation des données (par exemple, affichage de la taille des ensembles de données)
	fmt.Println("Nombre d'images dans l'ensemble d'entraînement:", trainImages.Count())
	fmt.Println("Nombre d'images dans l'ensemble de test:", testImages.Count())

	firstImage := trainImages.Images[0]

	// Print the shape of the image
	fmt.Println("Shape of the image:", len(firstImage), "x", len(firstImage))

	// Print the pixel values of the first few rows
	for i := 0; i < 200; i++ {
		fmt.Println(firstImage[i])
	}
}
