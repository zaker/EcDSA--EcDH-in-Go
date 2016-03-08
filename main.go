package main

import (
	"./ecc"
)

func main() {

	print("hoho\n")
	C := ecc.NewCurve()
	print("hoho\n")
	// 	C.GetCurve("secp256r1")

	// 	C.Print()
	C.Test()
}
