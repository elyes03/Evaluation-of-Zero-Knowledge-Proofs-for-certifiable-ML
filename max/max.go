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

// CubicCircuit defines a simple circuit
// x**3 + x + 5 == y
type MaxCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:"y"`       // x  --> secret visibility (default)
	Z frontend.Variable `gnark:",public"` // Y  --> public visibility
}

// Define declares the circuit constraints
func (circuit *MaxCircuit) Define(api frontend.API) error {
	//fmt.Println(api.Cmp(circuit.X, circuit.Y))
	cmp16bit := cmp.NewBoundedComparator(api, big.NewInt(1<<16-1), false)
	x := cmp16bit.Min(api.Neg(circuit.X), api.Neg(circuit.Y))
	//fmt.Println(x)
	api.AssertIsEqual(api.Neg(x), circuit.Z)
	v := api.Select(cmp16bit.IsLess(6, 5), 1, 0)
	fmt.Println(v)
	api.AssertIsEqual(v, 0)

	/*if api.Cmp(api.Neg(circuit.X), api.Neg(circuit.Y) == frontend.Variable(1) {
		api.Println("hhh", one)
	} else {
		api.Println(api.Add(api.Cmp(circuit.X, circuit.Y), 1))
	}*/

	return nil
}

func main() {

	// compiles our circuit into a R1CS
	var circuit MaxCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)

	// witness definition
	assignment := MaxCircuit{X: 1000, Y: 22, Z: 1000}
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
