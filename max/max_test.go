package main

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func BenchmarkMaxCircuit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var circuit MaxCircuit
		ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

		// groth16 zkSNARK: Setup
		pk, vk, err := groth16.Setup(ccs)

		// witness definition
		assignment := MaxCircuit{X: 3, Y: 4, Z: 4}
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
