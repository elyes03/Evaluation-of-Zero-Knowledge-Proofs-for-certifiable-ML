package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"image"
	"image/color"

	"math"
	"math/big"
	"os"

	//"io/ioutil"

	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/cmp"
	"github.com/petar/GoMNIST"
)

type ModelWeights struct {
	LinearWeight [][]float64 `json:"linear.weight"`
}

type SecureFixPointArithmetic struct {
	q int64 // prime finite field
	m int   // precision bits
}

func NewSecureFixPointArithmetic(q int64, m int) *SecureFixPointArithmetic {
	return &SecureFixPointArithmetic{
		q: q,
		m: m,
	}
}

func (sfa *SecureFixPointArithmetic) Representation(a float64) int64 {
	scaled := int64(a*math.Pow(10, float64(sfa.m))) % sfa.q
	/*for scaled < 0 {
		scaled += sfa.q
	}*/
	return scaled
}

func (sfa *SecureFixPointArithmetic) InverseRepresentation(representation int64) float64 {
	/*for representation < 0 {
		representation += sfa.q
	}*/
	a := float64(representation) / math.Pow(10, float64(sfa.m))
	return a
}

func (sfa *SecureFixPointArithmetic) SecureAddition(a, b int64) int64 {
	sum := (a + b) % sfa.q
	/*+if sum < 0 {
		sum += sfa.q
	}*/
	return sum
}

// 3.1*2.1 ==> 31*21=651
func (sfa *SecureFixPointArithmetic) SecureMultiplication(a, b int64) int64 {
	product := (a * b) % sfa.q
	/*if product < 0 {
		product += sfa.q
	}*/
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
	PreImage  frontend.Variable
	Hash      frontend.Variable `gnark:",public"`
	PreImage1 frontend.Variable
	Hash1     frontend.Variable `gnark:",public"`

	X       [784]frontend.Variable  `gnark:"x"`
	Y       frontend.Variable       `gnark:",public"`
	Weights [7840]frontend.Variable `gnark:"w"`
	Bias    frontend.Variable       `gnark:"b"`
}

// Define the constraints for the logistic regression circuit
func (circuit *LogisticRegressionCircuit) Define(api frontend.API) error {

	var prod0, prod1, prod2, prod3, prod4, prod5, prod6, prod7, prod8, prod9 frontend.Variable = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

	//var image frontend.Variable

	cmp16bit := cmp.NewBoundedComparator(api, big.NewInt(1000000000000000000), false)
	//for i := range 2 {

	for j := range 784 {
		prod0 = api.Add(prod0, api.Mul(circuit.Weights[j], circuit.X[j]))
		prod1 = api.Add(prod1, api.Mul(circuit.Weights[784*1+j], circuit.X[j]))
		prod2 = api.Add(prod2, api.Mul(circuit.Weights[784*2+j], circuit.X[j]))
		prod3 = api.Add(prod3, api.Mul(circuit.Weights[784*3+j], circuit.X[j]))
		prod4 = api.Add(prod4, api.Mul(circuit.Weights[784*4+j], circuit.X[j]))
		prod5 = api.Add(prod5, api.Mul(circuit.Weights[784*5+j], circuit.X[j]))
		prod6 = api.Add(prod6, api.Mul(circuit.Weights[784*6+j], circuit.X[j]))
		prod7 = api.Add(prod7, api.Mul(circuit.Weights[784*7+j], circuit.X[j]))
		prod8 = api.Add(prod8, api.Mul(circuit.Weights[784*8+j], circuit.X[j]))
		prod9 = api.Add(prod9, api.Mul(circuit.Weights[784*9+j], circuit.X[j]))

	}
	//prodmax = api.Neg(cmp16bit.Min(api.Neg(prod), api.Neg(prodmax)))
	//a := cmp16bit.IsLess(prod, prodmax)
	//v = api.Select(cmp16bit.IsLess(prod, prodmax), v, frontend.Variable(i))
	//index = v
	//api.IsZero(cmp16bit.IsLess(prod, prodmax))
	//fmt.Println("aa", api.IsZero(cmp16bit.IsLess(prod, prodmax)))

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

	mimc, _ := mimc.NewMiMC(api)
	mimc.Write(circuit.PreImage)
	api.AssertIsEqual(circuit.Hash, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage1)
	api.AssertIsEqual(circuit.Hash1, mimc.Sum())

	return nil
}

func mimcHash(data []byte) string {
	f := bn254.NewMiMC()
	f.Write(data)
	hash := f.Sum(nil)
	hashInt := big.NewInt(0).SetBytes(hash)
	return hashInt.String()
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

	img := image.NewGray(image.Rect(0, 0, 28, 28))
	for y := 0; y < 28; y++ {
		for x := 0; x < 28; x++ {
			img.SetGray(x, y, color.Gray{Y: uint8(firstImage[y*28+x])})
		}
	}

	// Save the image to a file
	/*outputFile, err := os.Create("second_image.png")
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return
	}
	defer outputFile.Close()

	// Encode the image as PNG and write it to the file
	err = png.Encode(outputFile, img)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Image saved as second_image.png")*/

	// Print the shape of the image
	fmt.Println("Shape of the image:", len(firstImage), "x", len(firstImage))

	// Print the pixel values of the first few rows

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
	fmt.Println("weights", len(poids.LinearWeight))
	fmt.Println("Poids du modèle:", poids.LinearWeight[0][1])
	sfa := NewSecureFixPointArithmetic(999999999999999999, 8)

	var assignment LogisticRegressionCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &assignment) // compile circuit
	pk, vk, err := groth16.Setup(ccs)                                                   // Initialize the gnark constraint system

	// Define the circuit inputs

	preImage := make([]byte, 32)
	preImage1 := make([]byte, 32)

	for i := 0; i <= 31; i++ {
		preImage[i] = byte(sfa.Representation(poids.LinearWeight[0][i]))
		preImage1[i] = byte(sfa.Representation(poids.LinearWeight[7][585+i]))
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

	fmt.Println("This cooresponds to a ", index)
	assignment.Y = index

	fmt.Println(assignment.Y)

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField()) // witness

	publicWitness, err := witness.Public()

	/*schema, _ := frontend.NewSchema(&assignment)
	ret, _ := publicWitness.ToJSON(schema)
	var b bytes.Buffer
	json.Indent(&b, ret, "", "\t")
	fmt.Println(b.String())*/

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
