package main

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The witness has a
//   - secret part --> known to the prover only
//   - public part --> known to the prover and the verifier
type Circuit struct {
	Secret frontend.Variable // pre-image of the hash secret known to the prover only
	Hash   frontend.Variable `gnark:",public"` // hash of the secret known to all
}

// Define declares the circuit logic. The compiler then produces a list of constraints
// which must be satisfied (valid witness) in order to create a valid zk-SNARK
// This circuit proves knowledge of a pre-image such that hash(secret) == hash
func (circuit *Circuit) Define(api frontend.API) error {
	// hash function
	mimc, _ := mimc.NewMiMC(api)

	// hash the secret
	mimc.Write(circuit.Secret)

	// ensure hashes match
	api.AssertIsEqual(circuit.Hash, mimc.Sum())

	return nil
}


/*-- witness.json --{
    "x": 3,
    "y": 5,
    "Z": 8
}

http://play.gnark.io/?id=oilfm3eosg 
