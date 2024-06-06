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
)

func BenchmarkHash(b *testing.B) {

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
	fmt.Println("Poids du modèle:", poids.LinearWeight[0][1])
	sfa := NewSecureFixPointArithmetic(999999999999999999, 8)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {

		preImage := make([]byte, 31)
		preImage1 := make([]byte, 31)

		for i := 0; i <= 30; i++ {
			preImage[i] = byte(sfa.Representation(poids.LinearWeight[0][i]))
			preImage1[i] = byte(sfa.Representation(poids.LinearWeight[7][582+i]))
		}
		t1 := time.Now()
		hash := mimcHash(preImage)
		elapsed1 := time.Since(t1)
		fmt.Println("Temps un appel à hash : ", elapsed1.Nanoseconds())
		hash2 := mimcHash(preImage1)
		elapsed2 := time.Since(t1)
		fmt.Println("Temps 2 appels à hash : ", elapsed2.Nanoseconds())

		// Write the hash to a file
		hashFile, err := os.Create("hash.txt")
		if err != nil {
			fmt.Println("Erreur lors de la création du fichier:", err)
			return
		}
		defer hashFile.Close()

		_, err = hashFile.WriteString(hash)
		if err != nil {
			fmt.Println("Erreur lors de l'écriture du fichier:", err)
			return
		}

		fmt.Println("Hash written to hash.txt")

		t3 := time.Now()
		var circuit Circuit
		ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
		pk, vk, err := groth16.Setup(ccs)
		elapsed3 := time.Since(t3)
		fmt.Println("Temps création circuit : ", elapsed3)

		assignment := &Circuit{
			PreImage:  preImage,
			PreImage1: preImage1,
			Hash:      hash,
			Hash1:     hash2,
		}

		witness, _ := frontend.NewWitness(assignment, ecc.BN254.ScalarField())

		t4 := time.Now()
		proof, err := groth16.Prove(ccs, pk, witness)
		elapsed4 := time.Since(t4)
		fmt.Println("Temps création preuve : ", elapsed4)

		if err != nil {
			fmt.Printf("Prove failed: %v\n", err)
			return
		}

		publicWitness, err := witness.Public()
		error := groth16.Verify(proof, vk, publicWitness)

		if error != nil {
			fmt.Printf("verification failed: %v\n", err)
			panic(error)
		}
		fmt.Printf("verification succeded\n")
	}
}
