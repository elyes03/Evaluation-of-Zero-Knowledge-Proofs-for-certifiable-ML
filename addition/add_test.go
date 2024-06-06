package main

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func BenchmarkAddCircuit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639747", 10)
		sfa := NewSecureFixPointArithmetic(q, 2, 10, 128)

		// Example numbers
		a := 10.54
		b := 3.71

		// Representation
		repA := sfa.Representation(a)
		repB := sfa.Representation(b)

		// Secure operations
		sum := sfa.SecureAddition(repA, repB)

		// compiles our circuit into a R1CS
		var circuit AddCircuit
		ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

		// groth16 zkSNARK: Setup
		pk, vk, err := groth16.Setup(ccs)

		// witness definition
		assignment := AddCircuit{X: repA, Y: repB, Z: sum}
		witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
		publicWitness, err := witness.Public()

		// groth16: Prove & Verify
		proof, err := groth16.Prove(ccs, pk, witness)
		err = groth16.Verify(proof, vk, publicWitness)
		if err != nil {
			panic(err)
		}

	}
}

/*func TinyOperation(a, b uint64) uint64 {
	return a * b
}

var result uint64

func BenchmarkTinyOperation(b *testing.B) {
	var s uint64
	for i := 0; i < b.N; i++ {
		TinyOperation(1, 2)
	}
	result = s

}*/
