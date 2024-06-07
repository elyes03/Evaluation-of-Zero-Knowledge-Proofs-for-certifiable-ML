package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/petar/GoMNIST"
)

func BenchmarkLinRegCircuit(b *testing.B) {

	sfa := NewSecureFixPointArithmetic(999999999999999999, 8)

	trainImages, testImages, err := GoMNIST.Load("./data./MNIST./raw")

	fmt.Println("Nombre d'images dans l'ensemble d'entraînement:", trainImages.Count())
	fmt.Println("Nombre d'images dans l'ensemble de test:", testImages.Count())

	firstImage := trainImages.Images[0]

	if err != nil {
		fmt.Println("Erreur lors du chargement du dataset MNIST:", err)
		return
	}
	fichier, err := os.Open("poids.json")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return
	}
	defer fichier.Close()

	var poids ModelWeights
	decodeur := json.NewDecoder(fichier)
	if err := decodeur.Decode(&poids); err != nil {
		fmt.Println("Erreur lors du décodage du fichier JSON:", err)
		return
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {

		t := time.Now()
		var assignment LogisticRegressionCircuit
		ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &assignment) // compile circuit
		pk, vk, err := groth16.Setup(ccs)                                                   // Initialize the gnark constraint system
		elapsed := time.Since(t)
		fmt.Println("Temps création circuit : ", elapsed)

		preImage := make([]byte, 31)
		preImage1 := make([]byte, 31)

		for i := 0; i <= 30; i++ {
			preImage[i] = byte(sfa.Representation(poids.LinearWeight[0][i]))
			preImage1[i] = byte(sfa.Representation(poids.LinearWeight[0][31+i]))
		}

		hash := mimcHash(preImage)
		hash1 := mimcHash(preImage1)

		assignment.PreImage = preImage
		assignment.Hash = hash
		assignment.PreImage1 = preImage1
		assignment.Hash1 = hash1

		for i := range assignment.X {
			assignment.X[i] = firstImage[i]
		}
		for i := 0; i < 10; i++ {
			for j := range 784 {
				assignment.Weights[784*i+j] = sfa.Representation(poids.LinearWeight[i][j])
			}
		}

		var x, max int64 = 0, 0
		var index int = 0
		for i := 0; i < 10; i++ {
			x = 0
			for j := range 784 {
				x += int64(firstImage[j]) * sfa.Representation(poids.LinearWeight[i][j])

			}
			if x > max {
				max = x
				index = i
			}
		}

		assignment.Bias = 0
		assignment.Y = index
		fmt.Println(assignment.Y)

		t1 := time.Now()
		witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField()) // witness
		publicWitness, err := witness.Public()
		elapsed1 := time.Since(t1)
		fmt.Println("Temps création witness : ", elapsed1)

		t2 := time.Now()
		proof, err := groth16.Prove(ccs, pk, witness) // generate the proof
		elapsed2 := time.Since(t2)
		fmt.Println("Temps création preuve : ", elapsed2)

		t3 := time.Now()
		err = groth16.Verify(proof, vk, publicWitness) // verify the proof
		elapsed3 := time.Since(t3)
		fmt.Println("Temps vérification preuve : ", elapsed3.Nanoseconds())

		//fmt.Println(proof, vk)
		if err != nil {
			panic(err)
		}

		fmt.Println("Zero knowledge proof generated successfully!")

	}
}
