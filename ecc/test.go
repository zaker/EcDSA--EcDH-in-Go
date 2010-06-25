package main

import (
	"./ecc"
	sha "crypto/sha256"
	"big"
// 	"hash"
)


func main() {

	C := ecc.NewCurve()
	C.GetCurve("secp160r1")

	secret := "secret"
	
	h := sha.New()
	
	h.Write([]byte(secret))
	print("pass:",secret,"\nhash:", h.Sum(),"\n")
	
	d := new(big.Int)
	print("d\n")
	d.SetString("971761939728640320549601132085879836204587084162", 10)
	
	Qc := new(ecc.Point)
	print("Qc\n")
	Qc.SetString("466448783855397898016055842232266600516272889280", "1110706324081757720403272427311003102474457754220", 10)
	Qc.Print()
	
	Q := new(ecc.Point)
	print("Q\n")
	Q.Mul(d,Qc,C)
	print("\n")
	print("Test success == ",ecc.Cmp(Q,Qc),"\n")
	
	
// 	sec := h.Write(byte_string)
	
// 	C.Print()


	

}