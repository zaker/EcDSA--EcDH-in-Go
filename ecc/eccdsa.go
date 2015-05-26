package ecc

import (
	"fmt"
	"io"
	"math/big"
	"os"
)

type EccKeyPair struct {
	d *big.Int
	Q *Point
	C *Curve
}

type Signature struct {
	r *big.Int
	s *big.Int
}

// randomNumber returns a uniform random value in [0, max).
func randomNumber(rand io.Reader, max *big.Int) (n *big.Int, err error) {
	k := (max.BitLen() + 7) / 8

	// r is the number of bits in the used in the most significant byte of
	// max.
	r := uint(max.BitLen() % 8)
	if r == 0 {
		r = 8
	}

	bytes := make([]byte, k)
	n = new(big.Int)

	for {
		_, err = io.ReadFull(rand, bytes)
		if err != nil {
			return
		}

		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bytes[0] &= uint8(int(1<<r) - 1)

		n.SetBytes(bytes)
		if n.Cmp(max) < 0 {
			return
		}
	}
}

func NewPair() *EccKeyPair {

	E := new(EccKeyPair)
	E.d = big.NewInt(0)
	E.Q = NewPoint()
	E.C = NewCurve()

	return E
}

func (S *Signature) SetRS(r, s *big.Int) *Signature {
	S.r = r
	S.s = s

	return S
}
func (S *Signature) SetString(r, s string, base int) *Signature {
	S.r.SetString(r, base)
	S.s.SetString(s, base)

	return S
}

func NewSign() *Signature {
	S := new(Signature)

	S.SetRS(big.NewInt(0), big.NewInt(0))

	return S
}

func (S *Signature) Print() {

	fmt.Printf(S.String())
}

func (S *Signature) String() string {
	s := ""

	s = "r: " + S.r.String()
	s += ",\n"
	s += "s: " + S.s.String()

	s += "\n"

	return s

}

func (E *EccKeyPair) EccdsaKeyGen(d *big.Int, C *Curve) *EccKeyPair {

	E.d = d
	E.C = C
	E.Q.Mul(d, C.G, C)

	return E
}

/* Signature Generation
For signing a message m by sender A, using A's private key dA
Calculate e = HASH (m), where HASH is a cryptographic hash function, such as SHA-1
Select a random integer k from [1,n - 1]
Calculate r = x1 (mod n), where (x1, y1) = k * G. If r = 0, go to step 2
Calculate s = k-1(e + dAr)(mod n). If s = 0, go to step 2
The signature is the pair (r, s)
*/

func (E *EccKeyPair) EccdsaSign(md *big.Int) *Signature {
	S := new(Signature)
	urand, _ := os.Open("/dev/urandom")
	r := new(big.Int)
	s := new(big.Int)
	Qs := new(Point)
	for {
		k, _ := randomNumber(urand, E.C.p)

		Qs.Mul(k, E.C.G, E.C)
		r.Mod(Qs.x, E.C.p)

		if r.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		s.Mul(r, E.d)
		s.Mod(s, E.C.n)

		s.Add(s, md)
		s.Mod(s, E.C.n)

		k.ModInverse(k, E.C.n)

		s.Mul(s, k)
		s.Mod(s, E.C.n)

		if s.Cmp(big.NewInt(0)) != 0 {
			break
		}

	}

	S.SetRS(r, s)
	return S
}

/* Signature Verification
For B to authenticate A's signature, B must have A's public key QA
Verify that r and s are integers in [1,n - 1]. If not, the signature is invalid
Calculate e = HASH (m), where HASH is the same function used in the signature generation
Calculate w = s-1 (mod n)
Calculate u1 = ew (mod n) and u2 = rw (mod n)
Calculate (x1, y1) = u1G + u2QA
The signature is valid if x1 = r(mod n), invalid otherwise
*/

func (E *EccKeyPair) EccdsaVerify(md *big.Int, S *Signature) bool {

	v := false

	if S.r.Cmp(E.C.p) >= 0 || S.s.Cmp(E.C.p) >= 0 {

		return v
	}

	w := new(big.Int).ModInverse(S.s, E.C.n)

	u1 := new(big.Int).Mul(w, md)
	u1.Mod(u1, E.C.n)

	u2 := new(big.Int).Mul(w, S.r)
	u2.Mod(u2, E.C.n)

	Gv := new(Point)
	Qv := new(Point)

	Gv.Mul(u1, E.C.G, E.C)
	Qv.Mul(u2, E.Q, E.C)

	Qv.Add(Gv, Qv, E.C)

	if S.r.Cmp(Qv.x) == 0 {
		v = true
	}
	return v
}

func (E *EccKeyPair) EccdsaKeyValidate() bool {
	v := false

	y := new(big.Int)
	x := new(big.Int)

	switch {
	case E.Q.x.Cmp(big.NewInt(1)) < 0:
		println("X less than 1")

	case E.Q.x.Cmp(E.C.p) > 0:
		println("X over p")
	case E.Q.y.Cmp(big.NewInt(1)) < 0:
		println("Y less than 1")

	case E.Q.y.Cmp(E.C.p) > 0:
		println("Y over p")
	default:
		y.Exp(E.Q.y, big.NewInt(2), E.C.p)

		x.Exp(E.Q.x, big.NewInt(3), E.C.p)

		t := new(big.Int)

		t.Mul(E.C.a, E.Q.x)
		t.Mod(t, E.C.p)

		t.Add(t, x)

		t.Add(t, E.C.b)
		t.Mod(t, E.C.p)

		return t.Cmp(y) == 0

	}

	return v
}
