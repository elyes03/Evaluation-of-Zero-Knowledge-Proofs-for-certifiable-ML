package main

import (
	"fmt"
	"math"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type SecureFixPointArithmetic struct {
	q *big.Int // prime finite field
	m int      // precision bits
	M int      // upper bound on maximum number encountered during training
	K int      // statistical security parameter
}

// NewSecureFixPointArithmetic initializes a new instance of SecureFixPointArithmetic
func NewSecureFixPointArithmetic(q *big.Int, m, M, K int) *SecureFixPointArithmetic {
	return &SecureFixPointArithmetic{
		q: q,
		m: m,
		M: M,
		K: K,
	}
}

// Representation converts a real number to its field element representation
func (sfa *SecureFixPointArithmetic) Representation(a float64) *big.Int {
	aScaled := a * math.Pow(10, float64(sfa.m)) // scale the real number
	aMod := new(big.Int).SetInt64(int64(aScaled))
	return aMod.Mod(aMod, sfa.q)
}

func (sfa *SecureFixPointArithmetic) SecureAddition(a, b *big.Int) *big.Int {
	sum := new(big.Int).Add(a, b)
	sum.Mod(sum, sfa.q)
	return sum
}

type AddCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:"y"`       // x  --> secret visibility (default)
	Z frontend.Variable `gnark:",public"` // Y  --> public visibility
}

// Define declares the circuit constraints
func (circuit *AddCircuit) Define(api frontend.API) error {
	var result frontend.Variable

	result = api.Add(circuit.X, circuit.Y)

	api.AssertIsEqual(circuit.Z, result)
	return nil
}

func main() {

	// compiles our circuit into a R1CS
	var circuit AddCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)

	// witness definition
	var assignment AddCircuit

	assignment.X = 1

	assignment.Y = 2

	assignment.Z = 3

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, err := witness.Public()

	// groth16: Prove & Verify
	proof, err := groth16.Prove(ccs, pk, witness)
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}

	fmt.Println("Zero knowledge proof generated and verified successfully!")

}
