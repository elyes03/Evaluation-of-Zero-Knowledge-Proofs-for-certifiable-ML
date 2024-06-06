package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"math"
	"math/big"
	"os"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"github.com/consensys/gnark/std/hash/mimc"
)

type Circuit struct {
	PreImage   frontend.Variable
	Hash       frontend.Variable `gnark:",public"`
	PreImage1  frontend.Variable
	Hash1      frontend.Variable `gnark:",public"`
	PreImage2  frontend.Variable
	Hash2      frontend.Variable `gnark:",public"`
	PreImage3  frontend.Variable
	Hash3      frontend.Variable `gnark:",public"`
	PreImage4  frontend.Variable
	Hash4      frontend.Variable `gnark:",public"`
	PreImage5  frontend.Variable
	Hash5      frontend.Variable `gnark:",public"`
	PreImage6  frontend.Variable
	Hash6      frontend.Variable `gnark:",public"`
	PreImage7  frontend.Variable
	Hash7      frontend.Variable `gnark:",public"`
	PreImage8  frontend.Variable
	Hash8      frontend.Variable `gnark:",public"`
	PreImage9  frontend.Variable
	Hash9      frontend.Variable `gnark:",public"`
	PreImage10 frontend.Variable
	Hash10     frontend.Variable `gnark:",public"`
	PreImage11 frontend.Variable
	Hash11     frontend.Variable `gnark:",public"`
	PreImage12 frontend.Variable
	Hash12     frontend.Variable `gnark:",public"`
	PreImage13 frontend.Variable
	Hash13     frontend.Variable `gnark:",public"`
	PreImage14 frontend.Variable
	Hash14     frontend.Variable `gnark:",public"`
	PreImage15 frontend.Variable
	Hash15     frontend.Variable `gnark:",public"`
}

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
	return scaled
}

func (circuit *Circuit) Define(api frontend.API) error {

	t3 := time.Now()
	mimc, _ := mimc.NewMiMC(api)
	mimc.Write(circuit.PreImage)
	api.AssertIsEqual(circuit.Hash, mimc.Sum())
	fmt.Println("temps 1 hash zk", time.Since(t3).Nanoseconds())
	mimc.Reset()
	mimc.Write(circuit.PreImage1)
	api.AssertIsEqual(circuit.Hash1, mimc.Sum())
	fmt.Println("temps 2 hash zk", time.Since(t3).Nanoseconds())
	mimc.Reset()
	mimc.Write(circuit.PreImage2)
	api.AssertIsEqual(circuit.Hash2, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage3)
	api.AssertIsEqual(circuit.Hash3, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage4)
	api.AssertIsEqual(circuit.Hash4, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage5)
	api.AssertIsEqual(circuit.Hash5, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage6)
	api.AssertIsEqual(circuit.Hash6, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage7)
	api.AssertIsEqual(circuit.Hash7, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage8)
	api.AssertIsEqual(circuit.Hash8, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage9)
	api.AssertIsEqual(circuit.Hash9, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage10)
	api.AssertIsEqual(circuit.Hash10, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage11)
	api.AssertIsEqual(circuit.Hash11, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage12)
	api.AssertIsEqual(circuit.Hash12, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage13)
	api.AssertIsEqual(circuit.Hash13, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage14)
	api.AssertIsEqual(circuit.Hash14, mimc.Sum())
	mimc.Reset()
	mimc.Write(circuit.PreImage15)
	api.AssertIsEqual(circuit.Hash15, mimc.Sum())

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
	preImage := make([]byte, 31)
	preImage1 := make([]byte, 31)
	preImage2 := make([]byte, 31)
	preImage3 := make([]byte, 31)
	preImage4 := make([]byte, 31)
	preImage5 := make([]byte, 31)
	preImage6 := make([]byte, 31)
	preImage7 := make([]byte, 31)
	preImage8 := make([]byte, 31)
	preImage9 := make([]byte, 31)
	preImage10 := make([]byte, 31)
	preImage11 := make([]byte, 31)
	preImage12 := make([]byte, 31)
	preImage13 := make([]byte, 31)
	preImage14 := make([]byte, 31)
	preImage15 := make([]byte, 31)

	for i := 0; i <= 30; i++ {
		preImage[i] = byte(sfa.Representation(poids.LinearWeight[0][i]))
		preImage1[i] = byte(sfa.Representation(poids.LinearWeight[7][582+i]))
		preImage2[i] = byte(sfa.Representation(poids.LinearWeight[4][592+i]))
		preImage3[i] = byte(sfa.Representation(poids.LinearWeight[2][362+i]))
		preImage4[i] = byte(sfa.Representation(poids.LinearWeight[3][i]))
		preImage5[i] = byte(sfa.Representation(poids.LinearWeight[8][45+i]))
		preImage6[i] = byte(sfa.Representation(poids.LinearWeight[1][107+i]))
		preImage7[i] = byte(sfa.Representation(poids.LinearWeight[2][62+i]))
		preImage8[i] = byte(sfa.Representation(poids.LinearWeight[1][i]))
		preImage9[i] = byte(sfa.Representation(poids.LinearWeight[2][i]))
		preImage10[i] = byte(sfa.Representation(poids.LinearWeight[4][i]))
		preImage11[i] = byte(sfa.Representation(poids.LinearWeight[7][362+i]))
		preImage12[i] = byte(sfa.Representation(poids.LinearWeight[8][i]))
		preImage13[i] = byte(sfa.Representation(poids.LinearWeight[9][45+i]))
		preImage14[i] = byte(sfa.Representation(poids.LinearWeight[3][107+i]))
		preImage15[i] = byte(sfa.Representation(poids.LinearWeight[8][62+i]))
	}

	hash := mimcHash(preImage)
	hash1 := mimcHash(preImage1)
	hash2 := mimcHash(preImage2)
	hash3 := mimcHash(preImage3)
	hash4 := mimcHash(preImage4)
	hash5 := mimcHash(preImage5)
	hash6 := mimcHash(preImage6)
	hash7 := mimcHash(preImage7)
	hash8 := mimcHash(preImage8)
	hash9 := mimcHash(preImage9)
	hash10 := mimcHash(preImage10)
	hash11 := mimcHash(preImage11)
	hash12 := mimcHash(preImage12)
	hash13 := mimcHash(preImage13)
	hash14 := mimcHash(preImage14)
	hash15 := mimcHash(preImage15)

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
	//_, err = hashFile.WriteString(hash2)

	fmt.Println("Hash written to hash.txt")

	/*hashFromFile, err := os.ReadFile("hash.txt")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier hash.txt:", err)
		return
	}

	// Convert the hash from file to a big.Int
	hashInt := new(big.Int)
	hashInt.SetString(string(hashFromFile), 10)

	// Convert hash to frontend.Variable
	hashVariable := frontend.Variable(hashInt)*/

	t3 := time.Now()
	var circuit Circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	pk, vk, err := groth16.Setup(ccs)
	elapsed3 := time.Since(t3)
	fmt.Println("Temps création circuit : ", elapsed3)

	assignment := Circuit{
		PreImage:   preImage,
		PreImage1:  preImage1,
		PreImage2:  preImage2,
		PreImage3:  preImage3,
		PreImage4:  preImage4,
		PreImage5:  preImage5,
		PreImage6:  preImage6,
		PreImage7:  preImage7,
		PreImage8:  preImage8,
		PreImage9:  preImage9,
		PreImage10: preImage10,
		PreImage11: preImage11,
		PreImage12: preImage12,
		PreImage13: preImage13,
		PreImage14: preImage14,
		PreImage15: preImage15,
		Hash:       hash,
		Hash1:      hash1,
		Hash2:      hash2,
		Hash3:      hash3,
		Hash4:      hash4,
		Hash5:      hash5,
		Hash6:      hash6,
		Hash7:      hash7,
		Hash8:      hash8,
		Hash9:      hash9,
		Hash10:     hash10,
		Hash11:     hash11,
		Hash12:     hash12,
		Hash13:     hash13,
		Hash14:     hash14,
		Hash15:     hash15,
	}

	t1 := time.Now()
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField()) // witness
	publicWitness, err := witness.Public()
	elapsed1 := time.Since(t1)
	fmt.Println("Temps création witness : ", elapsed1.Nanoseconds())

	t4 := time.Now()
	proof, err := groth16.Prove(ccs, pk, witness)
	elapsed4 := time.Since(t4)
	fmt.Println("Temps création preuve : ", elapsed4)

	t2 := time.Now()
	err = groth16.Verify(proof, vk, publicWitness) // verify the proof
	elapsed2 := time.Since(t2)
	fmt.Println("Temps vérification preuve : ", elapsed2.Nanoseconds())
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
