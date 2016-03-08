package ecc

import (
	// 	"./ecc"
	sha "crypto/sha256"
	"math/big"
	"testing"
	// 	"fmt"
	// 	"hash"
)

func TestPointOperations(t *testing.T) {

	C := NewCurve()
	curve_name := "secp256k1"
	C.GetCurve(curve_name)

	if C.name != "secp256k1" {
		t.Errorf("Failed loading curve %s", curve_name)

	}

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
	I.Inv(C.G, C)

	if !Cmp(I, Qi) {

		t.Errorf("Additive inverse failed // I = -K")
	}

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

func TestSigning(t *testing.T) {

	secret := "secret0"

	h := sha.New()

	h.Write([]byte(secret))
	secnum := new(big.Int).SetBytes(h.Sum())

	C := NewCurve()
	curve_name := "secp256k1"
	C.GetCurve(curve_name)

	E := NewPair()

	E.EccdsaKeyGen(secnum, C)

	if E.d.String() != secnum.String() {
		t.Errorf("PrivateKey fail")

	}

	if !E.EccdsaKeyValidate() {
		t.Errorf("KeyPair Fail")

	}
	Qc := new(Point)

	Qc = Qc.SetString("109837834369274504389664124593965957096098981290804729253410145109750445901603",
		"109557846963727535798669146173709773735187162620328251093267341541062606906157", 10)

	if !Cmp(E.Q, Qc) {

		t.Errorf("Wrong PrivateKey")

	}

	if C.name != E.C.name {
		t.Errorf("Wrong Curve")
	}

	h.Reset()

	m := "ABCD"

	h.Write([]byte(m))
	md := new(big.Int).SetBytes(h.Sum())

	S := NewSign()

	S = E.EccdsaSign(md)

	if !E.EccdsaVerify(md, S) {

		t.Errorf("Signature failed")
	}

}

func TestDiffieHellman(t *testing.T) {

	secret := "secret0"

	h := sha.New()

	h.Write([]byte(secret))
	secnum := new(big.Int).SetBytes(h.Sum())

	C := NewCurve()
	curve_name := "secp256k1"
	C.GetCurve(curve_name)

	A := NewPair()

	A.EccdsaKeyGen(secnum, C)

	secret2 := "secret032"

	h.Reset()

	h.Write([]byte(secret2))
	secnum2 := new(big.Int).SetBytes(h.Sum())

	B := NewPair()

	B.EccdsaKeyGen(secnum2, C)

	if A.EccDH(B.Q).Cmp(B.EccDH(A.Q)) != 0 {
		t.Errorf("Diffie Hellman failed")

	}

}
