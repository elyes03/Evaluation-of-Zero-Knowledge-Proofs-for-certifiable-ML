package main

import (
	"fmt"

	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/math/cmp"
)

type MaxCircuit struct {
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:"y"`       // x  --> secret visibility (default)
	Z frontend.Variable `gnark:",public"` // Y  --> public visibility
}

func (circuit *MaxCircuit) Define(api frontend.API) error {

	cmp16bit := cmp.NewBoundedComparator(api, big.NewInt(1<<16-1), false)
	x := cmp16bit.Min(api.Neg(circuit.X), api.Neg(circuit.Y))
	api.AssertIsEqual(api.Neg(x), circuit.Z)

	return nil
}

func main() {

	// compiles our circuit into a R1CS
	var circuit MaxCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)

	// witness definition
	assignment := MaxCircuit{X: 2, Y: 23, Z: 23}
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
