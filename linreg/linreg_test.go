package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/petar/GoMNIST"
)

func BenchmarkLogisticRegressionPrediction(b *testing.B) {

	trainImages, testImages, err := GoMNIST.Load("./data./MNIST./raw")
	if err != nil {
		fmt.Println("Erreur lors du chargement du dataset MNIST:", err)
		return
	}

	fmt.Println("Nombre d'images dans l'ensemble d'entraînement:", trainImages.Count())
	fmt.Println("Nombre d'images dans l'ensemble de test:", testImages.Count())

	firstImage := trainImages.Images[0]

	fichier, err := os.Open("poids1.json")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return
	}
	defer fichier.Close()

	// Décodez les données JSON dans une structure de poids Go
	var poids ModelWeights
	decodeur := json.NewDecoder(fichier)
	if err := decodeur.Decode(&poids); err != nil {
		fmt.Println("Erreur lors du décodage du fichier JSON:", err)
		return
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var x, max float64 = 0, 0
		var index = 0
		for i := 0; i < 10; i++ {
			x = 0
			for j := range 784 {
				x += float64(firstImage[j]) * poids.LinearWeight[i][j]
			}
			if x > max {
				max = x
				index = i
			}
		}

		fmt.Println(max)
		fmt.Println("index", index)
	}
}
