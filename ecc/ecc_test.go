package ecc

import (
	// 	"./ecc"
	sha "crypto/sha256"
	"big"
	"testing"
	// 	"fmt"
	// 	"hash"
)


func TestPointOperations(t *testing.T) {

	secret := "secret0"

	h := sha.New()

	h.Write([]byte(secret))
	// 	print("pass:",secret,"\nhash:", h.Sum(),"\n")
	secnum := new(big.Int).SetBytes(h.Sum())
	println(h.Size(), " ", secnum.String())

	C := NewCurve()
	curve_name := "secp256k1"
	C.GetCurve(curve_name)

	if C.name != "secp256k1" {
		t.Errorf("Failed loading curve %s", curve_name)

	}
	println("Setting Checkpoint:")

	k := new(big.Int)
	Qi := new(Point)
	Qa := new(Point)
	Qd := new(Point)
	Qk := new(Point)
	G1 := new(Point)

	k.SetString("1234567890", 10)

	Qi.SetString("55066263022277343669578718895168534326250603453777594175500187360389116729240",
		"83121579216557378445487899878180864668798711284981320763518679672151497189239", 10)

	Qa.SetString("19554613901239451535242065419393388143161175917962243373192119368079661378784",
		"86824435966238386053712309710725852390399392108259125416532031684543684637357", 10)

	Qd.SetString("89565891926547004231252920425935692360644145829622209833684329913297188986597",
		"12158399299693830322967808612713398636155367887041628176798871954788371653930", 10)

	Qk.SetString("19635924277356798752105674083697999930996555344818160161847497917044432760610",
		"21218882238660449272792211265489841951893738252848232230063147580786068364204", 10)

	G1.SetString("55066263022277343669578718895168534326250603453777594175500187360389116724240",
		"32670510020758816978083085130507043184471273380659243275938904335757337488424", 10)

	I := NewPoint()
	println("Inverting Generator:")
	I.Inv(C.G, C)
	C.G.Print()
	if !Cmp(I, Qi) {

		t.Errorf("Additive inverse failed // I = -K")
	}

	println("Adding G + (-G)")

	T := NewPoint()
	T.Add(I, C.G, C)

	if !T.inf {
		T.Print()
		t.Errorf("Addition didn't set infinite")
	}

	T = T.Add(G1, C.G, C)

	if !Cmp(T, Qa) {
		t.Errorf("Addition Failed")
	}

	T.Add(C.G, C.G, C)
	if !Cmp(T, Qd) {
		T.Print()
		t.Errorf("Doubling Failed")
	}

	T.Mul(k, C.G, C)

	if !Cmp(T, Qk) {
		T.Print()
		t.Errorf("Multiplication Failed")
	}

}
