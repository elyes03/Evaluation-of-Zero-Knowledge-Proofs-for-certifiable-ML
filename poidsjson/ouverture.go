package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ModelWeights struct {
	LinearWeight [][]float64 `json:"linear.weight"`
}

func main() {
	// Ouvrir le fichier JSON contenant les poids
	fichier, err := os.Open("poids.json")
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

	// Utiliser les poids comme nécessaire
	fmt.Println("Poids du modèle:", poids.LinearWeight)
}
