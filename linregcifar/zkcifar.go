package main

import (
	"encoding/json"
	"fmt"

	/*"image"
	"image/color"
	"image/png"*/
	"io/ioutil"
	"log"

	"math"
	"math/big"
	"os"

	//"io/ioutil"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/math/cmp"
)

type SecureFixPointArithmetic struct {
	q int64 // prime finite field
	m int   // precision bits
}

type ModelWeights struct {
	LinearWeight [][]float64 `json:"linear.weight"`
}

func NewSecureFixPointArithmetic(q int64, m int) *SecureFixPointArithmetic {
	return &SecureFixPointArithmetic{
		q: q,
		m: m,
	}
}

func (sfa *SecureFixPointArithmetic) Representation(a float64) int64 {
	scaled := int64(a*math.Pow(10, float64(sfa.m))) % sfa.q
	return scaled
}

func (sfa *SecureFixPointArithmetic) InverseRepresentation(representation int64) float64 {
	a := float64(representation) / math.Pow(10, float64(sfa.m))
	return a
}

func (sfa *SecureFixPointArithmetic) SecureAddition(a, b int64) int64 {
	sum := (a + b) % sfa.q
	return sum
}

// 3.1*2.1 ==> 31*21=651
func (sfa *SecureFixPointArithmetic) SecureMultiplication(a, b int64) int64 {
	product := (a * b) % sfa.q
	return product
}

// 3.1*2.1=6.51 ==> 651
func (sfa *SecureFixPointArithmetic) SecureMultiplicationTrue(a, b float64) int64 {
	product := a * b
	x := sfa.Representation(product) % sfa.q
	if x < 0 {
		x += sfa.q
	}
	return x
}

type LogisticRegressionCircuit struct {
	// Public inputs
	X       [3072]frontend.Variable  `gnark:"x"`
	Y       frontend.Variable        `gnark:",public"`
	Weights [30720]frontend.Variable `gnark:"w"`
	Bias    frontend.Variable        `gnark:"b"`
}

// Define the constraints for the logistic regression circuit
func (circuit *LogisticRegressionCircuit) Define(api frontend.API) error {

	var prod0, prod1, prod2, prod3, prod4, prod5, prod6, prod7, prod8, prod9 frontend.Variable = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

	//var image frontend.Variable

	cmp16bit := cmp.NewBoundedComparator(api, big.NewInt(1000000000000000000), false)
	//for i := range 2 {

	for j := range 3072 {
		prod0 = api.Add(prod0, api.Mul(circuit.Weights[j], circuit.X[j]))
		prod1 = api.Add(prod1, api.Mul(circuit.Weights[3072*1+j], circuit.X[j]))
		prod2 = api.Add(prod2, api.Mul(circuit.Weights[3072*2+j], circuit.X[j]))
		prod3 = api.Add(prod3, api.Mul(circuit.Weights[3072*3+j], circuit.X[j]))
		prod4 = api.Add(prod4, api.Mul(circuit.Weights[3072*4+j], circuit.X[j]))
		prod5 = api.Add(prod5, api.Mul(circuit.Weights[3072*5+j], circuit.X[j]))
		prod6 = api.Add(prod6, api.Mul(circuit.Weights[3072*6+j], circuit.X[j]))
		prod7 = api.Add(prod7, api.Mul(circuit.Weights[3072*7+j], circuit.X[j]))
		prod8 = api.Add(prod8, api.Mul(circuit.Weights[3072*8+j], circuit.X[j]))
		prod9 = api.Add(prod9, api.Mul(circuit.Weights[3072*9+j], circuit.X[j]))

	}

	a := api.Select(cmp16bit.IsLess(prod0, prod1), frontend.Variable(1), frontend.Variable(0))
	aa := api.Neg(cmp16bit.Min(api.Neg(prod0), api.Neg(prod1)))
	b := api.Select(cmp16bit.IsLess(aa, prod2), frontend.Variable(2), a)
	bb := api.Neg(cmp16bit.Min(api.Neg(aa), api.Neg(prod2)))
	c := api.Select(cmp16bit.IsLess(bb, prod3), frontend.Variable(3), b)
	cc := api.Neg(cmp16bit.Min(api.Neg(bb), api.Neg(prod3)))
	d := api.Select(cmp16bit.IsLess(cc, prod4), frontend.Variable(4), c)
	dd := api.Neg(cmp16bit.Min(api.Neg(cc), api.Neg(prod4)))
	e := api.Select(cmp16bit.IsLess(dd, prod5), frontend.Variable(5), d)
	ee := api.Neg(cmp16bit.Min(api.Neg(dd), api.Neg(prod5)))
	f := api.Select(cmp16bit.IsLess(ee, prod6), frontend.Variable(6), e)
	ff := api.Neg(cmp16bit.Min(api.Neg(ee), api.Neg(prod6)))
	g := api.Select(cmp16bit.IsLess(ff, prod7), frontend.Variable(7), f)
	gg := api.Neg(cmp16bit.Min(api.Neg(ff), api.Neg(prod7)))
	h := api.Select(cmp16bit.IsLess(gg, prod8), frontend.Variable(8), g)
	hh := api.Neg(cmp16bit.Min(api.Neg(gg), api.Neg(prod8)))
	index := api.Select(cmp16bit.IsLess(hh, prod9), frontend.Variable(9), h)
	//ll := api.Neg(cmp16bit.Min(api.Neg(h), api.Neg(tab[9])))

	api.AssertIsEqual(circuit.Y, index)

	return nil
}

const (
	imageWidth  = 32
	imageHeight = 32
	numChannels = 3
	recordSize  = imageWidth*imageHeight*numChannels + 1
)

func main() {

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
	//label := data[3*recordSize]
	imageData := data[3*recordSize+1 : 4*recordSize]

	// Create an empty image
	/*img := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	// Fill the image with data and print pixel values
	pixels := make([]float64, imageWidth*imageHeight*numChannels)
	idx := 0
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			r := imageData[y*imageWidth+x]
			g := imageData[imageWidth*imageHeight+y*imageWidth+x]
			b := imageData[2*imageWidth*imageHeight+y*imageWidth+x]
			img.SetNRGBA(x, y, color.NRGBA{r, g, b, 255})

			// Normalize pixel values to range [0, 1]
			pixels[idx] = float64(r)
			pixels[idx+1] = float64(g)
			pixels[idx+2] = float64(b)
			idx += 3
		}
	}

	// Save the image to a PNG file (optional, to visualize the image)
	outputFile, err := os.Create("output_image.png")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, img)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	fmt.Printf("Saved first image with label %d\n", label)

	// Access the pixel values for linear regression or other processing
	// For example, print the first 10 pixel values
	fmt.Println("First 10 pixel values (normalized):", pixels[:10])*/

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
	sfa := NewSecureFixPointArithmetic(999999999999999999, 12)

	var assignment LogisticRegressionCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &assignment) // compile circuit
	pk, vk, err := groth16.Setup(ccs)                                                   // Initialize the gnark constraint system

	// Define the circuit inputs

	for i := range assignment.X {
		assignment.X[i] = imageData[i]
	}
	for i := 0; i < 10; i++ {
		for j := range 3072 {

			assignment.Weights[3072*i+j] = sfa.Representation(poids.LinearWeight[i][j])

		}
	}

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

	assignment.Bias = 0

	fmt.Println("this cooresponds to a ", index)
	assignment.Y = index

	fmt.Println(assignment.Y)

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField()) // witness

	publicWitness, err := witness.Public()

	proof, err := groth16.Prove(ccs, pk, witness)  // generate the proof
	err = groth16.Verify(proof, vk, publicWitness) // verify the proof

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

/*
0: Airplane
1: Automobile
2: Bird
3: Cat
4: Deer
5: Dog
6: Frog
7: Horse
8: Ship
9: Truck
*/
