package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func BenchmarkLinRegCircuit(b *testing.B) {

	file, err := os.Open("data/cifar-10-batches-bin/test_batch.bin")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the entire file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	fmt.Println("data", len(data))
	// Extract the first image
	label := data[3*recordSize]
	imageData := data[3*recordSize+1 : 4*recordSize]

	fmt.Printf("Saved image with label %d\n", label)

	// Ouvrir le fichier JSON contenant les poids
	fichier, err := os.Open("poidscifar.json")
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
	fmt.Println("weights", len(poids.LinearWeight[0]))
	fmt.Println("Poids du modèle:", poids.LinearWeight[0][1])
	sfa := NewSecureFixPointArithmetic(999999999999999999, 8)
	b.ResetTimer()

	//b.ResetTimer()
	// Define the circuit inputs
	for n := 0; n < 5; n++ {

		t := time.Now()
		var assignment LogisticRegressionCircuit
		ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &assignment) // compile circuit
		pk, vk, err := groth16.Setup(ccs)                                                   // Initialize the gnark constraint system
		elapsed := time.Since(t)
		fmt.Println("Temps création circuit : ", elapsed)

		for i := range assignment.X {
			assignment.X[i] = imageData[i]
		}
		for i := 0; i < 10; i++ {
			for j := range 3072 {

				assignment.Weights[3072*i+j] = sfa.Representation(poids.LinearWeight[i][j])

			}
		}

		tt := time.Now()
		var x, max int64 = 0, 0
		var index int = 0
		for i := 0; i < 10; i++ {
			x = 0
			for j := range 3072 {
				x += int64(imageData[j]) * sfa.Representation(poids.LinearWeight[i][j])
			}
			fmt.Println(i, ":", x)
			if x > max {
				max = x
				index = i
			}
		}

		elapsed10 := time.Since(tt)
		fmt.Println("Temps native : ", elapsed10.Nanoseconds())
		assignment.Bias = 0

		fmt.Println("this cooresponds to a ", index)
		assignment.Y = index

		//fmt.Println(assignment.Y)

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

		//fmt.Println(vk, proof)
		if err != nil {
			panic(err)
		}

		fmt.Println("Zero knowledge proof generated and verified successfully!")

	}
}
