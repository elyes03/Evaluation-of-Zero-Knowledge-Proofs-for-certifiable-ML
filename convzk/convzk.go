package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/petar/GoMNIST"
)

// CubicCircuit defines a simple circuit
// x**3 + x + 5 == y
type ConvCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X [784]frontend.Variable `gnark:"x"`
	Y [4]frontend.Variable   `gnark:"y"`       // x  --> secret visibility (default)
	Z frontend.Variable      `gnark:",public"` // Y  --> public visibility
}

// Define declares the circuit constraints
func (circuit *ConvCircuit) Define(api frontend.API) error {

	result := make([]frontend.Variable, 196)
	var sum10 frontend.Variable = 0
	for j := 0; j < 14; j++ {
		for i := 0; i < 14; i++ {
			var sum frontend.Variable = 0
			sum = api.Add(api.Mul(circuit.X[28*j*2+2*i], circuit.Y[0]), api.Mul(circuit.X[28*j*2+1+2*i], circuit.Y[1]),
				api.Mul(circuit.X[28*(2*j+1)+2*i], circuit.Y[2]), api.Mul(circuit.X[28*(2*j+1)+1+2*i], circuit.Y[3]))
			result[j*i] = sum
			sum10 = api.Add(sum10, sum)
		}
	}
	/*for i := 0; i < 2; i++ {
		var sum frontend.Variable = 0

		sum = api.Add(api.Mul(circuit.X[8+2*i], circuit.Y[0]), api.Mul(circuit.X[9+2*i], circuit.Y[1]),
			api.Mul(circuit.X[12+2*i], circuit.Y[2]), api.Mul(circuit.X[13+2*i], circuit.Y[3]))

		result[2+i] = sum
		sum10 = api.Add(sum10, sum)

	}*/
	api.AssertIsEqual(circuit.Z, sum10)

	return nil
}

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

	t := time.Now()
	var assignment ConvCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &assignment) // compile circuit
	pk, vk, err := groth16.Setup(ccs)                                                   // Initialize the gnark constraint system
	elapsed := time.Since(t)
	fmt.Println("Temps création circuit : ", elapsed)

	var image [784]int
	for i := range 784 {
		assignment.X[i] = firstImage[i]
		image[i] = int(firstImage[i])
	}

	kernel := []int{
		1, 0, 0, 1}

	kernel1 := [4]frontend.Variable{
		1, 0, 0, 1}

	assignment.Y = kernel1

	resultSize := 28 / 2
	result := make([]int, resultSize*resultSize)

	for i := 0; i < resultSize; i++ {
		for j := 0; j < resultSize; j++ {
			// Calculate the starting indices of the current 4x4 block
			startIndexI := i * 2
			startIndexJ := j * 2

			// Calculate the convolution sum for the current block
			sum := 0
			for k := 0; k < 2; k++ {
				for l := 0; l < 2; l++ {
					sum += image[(startIndexI+k)*28+startIndexJ+l] * kernel[k*2+l]

				}
			}

			// Store the convolution result in the result slice
			result[i*resultSize+j] = sum
		}
	}

	var sum1 int = 0
	for i := range result {
		sum1 += result[i]
	}
	fmt.Println(sum1)
	// Afficher le résultat
	fmt.Println("Résultat de la convolution:")
	for _, val := range result {
		fmt.Print(val, " ")
	}

	assignment.Z = sum1

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
	if err != nil {
		panic(err)
	}

	fmt.Println("Zero knowledge proof generated and verified successfully!")

	//Convert the proof to JSON format
	proofJSON, err := json.Marshal(proof)
	if err != nil {
		fmt.Println("Error converting proof to JSON:", err)
		return
	}

	// Write the JSON data to a file
	err = ioutil.WriteFile("proof.json", proofJSON, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Proof exported to proof.json successfully.")
	// Get the file size
	fileInfo, err := os.Stat("proof.json")
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	fileSize := fileInfo.Size()
	fmt.Printf("Size of the proof file: %d bytes\n", fileSize)

}
