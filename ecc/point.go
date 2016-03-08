package ecc

import (
	"fmt"
	"math/big"
)

type Point struct {
	x *big.Int

	y   *big.Int
	inf bool
}

func NewPoint() *Point {

	P := new(Point)

	P.x = big.NewInt(0)
	P.y = big.NewInt(0)

	return P
}

func (R *Point) Set(P *Point) *Point {

	R.inf = P.inf

	R.x = new(big.Int)
	R.y = new(big.Int)

	R.x.Set(P.x)
	R.y.Set(P.y)

	return R
}

func (R *Point) SetInf() *Point {

	R.inf = true

	R.x = big.NewInt(0)
	R.y = big.NewInt(0)

	return R
}

func (P *Point) Print() {

	fmt.Printf("%s", P.String())

}

func (P *Point) String() string {
	s := ""

	if P.inf {
		s = "Point at infinity\n"
	} else {
		s = "\nPoint: (\nx:\t" + P.x.String() + ",\ny:\t" + P.y.String() + ")\n"
	}
	return s
}

func (P *Point) SetString(x, y string, base int) *Point {
	// 	P = new(Point)

	P.x = new(big.Int)
	P.y = new(big.Int)

	P.inf = false
	P.x.SetString(x, base)
	P.y.SetString(y, base)
	return P
}

/* Compares two points returns true if equal*/
func Cmp(P, Q *Point) bool {
	//If at infinity
	if P.inf && Q.inf {
		return true
	}
	if P.inf || Q.inf {
		return false
	}
	return P.y.Cmp(Q.y) == 0 && P.x.Cmp(Q.x) == 0

}

/* R = 2P */
func (R *Point) double(P *Point, C *Curve) *Point {
	//If at infinity
	if P.inf {
		R.SetInf()
	} else {
		//Calculate slope
		//s = (3*Px² + a) / (2*Py) mod p
		t1 := new(big.Int).Exp(P.x, big.NewInt(2), C.p) // t1 = Px² mod p

		s := new(big.Int).Mul(t1, big.NewInt(3)) // s = 3 * t1

		s.Mod(s, C.p) // s = s mod p

		s.Add(s, C.a) // s = s + a mod p
		s.Mod(s, C.p)

		sd := new(big.Int).Mul(big.NewInt(2), P.y)

		sd.ModInverse(sd, C.p) // sd = 1 / 2*Py mod p

		s.Mul(s, sd) // s = (s / sd) mod p
		s.Mod(s, C.p)

		//Calculate Rx
		//Rx = s² - 2*Px mod p
		s2 := new(big.Int).Exp(s, big.NewInt(2), C.p) //s2 = s²
		s2.Mod(s2, C.p)

		t := new(big.Int).Mul(P.x, big.NewInt(2)) // t = 2*Px mod p
		t.Mod(t, C.p)

		t.Sub(s2, t) // Rx = s2 - t mod p
		R.x.Mod(t, C.p)

		//Calculate Ry using algorithm shown to the right of the commands
		//Ry = s(Px-Rx) - Py mod p

		t.Sub(P.x, R.x) //t = Px - Rx mod p
		t.Mod(t, C.p)

		t.Mul(s, t) //t = s * t mod p
		t.Mod(t, C.p)

		t.Sub(t, P.y) //t = t - Py mod p
		R.y.Mod(t, C.p)

		// 		mpz_mod(R->y, t3, curve->p);	//Ry = t3 mod p

	}
	return R
}

/* R = P + Q mod C.p*/
func (R *Point) Add(P, Q *Point, C *Curve) *Point {

	//If Q is at infinity, set result to P
	// 	R.Set(P)
	// 	R.SetInf()

	IP := new(Point)

	IP.Inv(P, C)
	// 	R.Print()
	// 	P.Print()
	// 	Q.Print()
	// 	println("P==Q : ",Cmp(P, Q), "-P==Q : ", Cmp(IP, Q))

	switch {
	case Q.inf:
		R.Set(P)
		// 		return R

	case P.inf:
		R.Set(Q)
		// 		return R

	case Cmp(P, Q):
		R.double(Q, C)
		// 		return R
		//If it is the inverse
	case Cmp(IP, Q):
		//result must be point at infinity
		R.SetInf()
		// 		R.Print()
		// 		return R
		// 		print("should nevever\n")

	default:

		// 		print("sdfsdfds\n")
		R.inf = false
		/*
			Modulo algebra rules:
			(b1 + b2) mod  n = (b2 mod n) + (b1 mod n) mod n
			(b1 * b2) mod  n = (b2 mod n) * (b1 mod n) mod n
		*/

		//Calculate slope
		//s = (Py - Qy)/(Px-Qx) mod p
		dy := new(big.Int).Sub(P.y, Q.y) // s = Py - Qy mod p
		dx := new(big.Int).Sub(P.x, Q.x) // sd = Px - Qx mod p
		// 		println("S: " + s.String())
		// 		println("Sd: " + sd.String() +"\n mod :" + C.p.String())
		dx.Mod(dx, C.p)
		dy.Mod(dy, C.p)
		// 		println("Sd: " + sd.String())
		// 		s.Mod(s,C.p)
		dxi := new(big.Int).ModInverse(dx, C.p) //Using Modulo to stay within the Field!
		// 		println("Sd: " + sd.String())
		// 		println("S: " + s.String())
		s := new(big.Int)
		s.Mul(dy, dxi)
		// 		println("S: " + s.String())
		s.Mod(s, C.p)
		// 		println("S: " + s.String())
		//Calculate Rx using algorithm shown to the right of the commands
		//Rx = s² - Px - Qx = (s² mod p) - (Px mod p) - (Qx mod p) mod p
		x := new(big.Int).Exp(s, big.NewInt(2), C.p) //x  = s² mod p
		x.Mod(x, C.p)
		x.Sub(x, P.x) // x = x - Px
		x.Mod(x, C.p)
		x.Sub(x, Q.x) // x = x - Qx

		R.x.Mod(x, C.p) // Rx = x mod p
		// 		println("Rx: " + R.x.String())
		//Calculate Ry using algorithm shown to the right of the commands
		//Ry = s(Px-Rx) - Py mod p
		dx.Sub(P.x, R.x) //y = Px - Rx
		dx.Mod(dx, C.p)

		dx.Mul(s, dx) //dx = s*dx mod p
		dx.Mod(dx, C.p)
		R.y.Sub(dx, P.y) //y = y - Py

		R.y.Mod(R.y, C.p) //Ry = y mod p
	}

	return R
}

/* R = -J mod C.p, R=(Jx,-Jy mod C.p)*/
func (R *Point) Inv(J *Point, C *Curve) *Point {
	R.Set(J)

	if !J.inf {
		R.y.Sub(C.p, J.y)
	}

	return R
}

/* BitCheck returns true if there is a bit at position n*/

func bitCheck(k *big.Int, n uint) bool {
	b := false

	// 	check :=

	lsh := new(big.Int).Lsh(big.NewInt(1), uint(n))

	bit := new(big.Int).And(k, lsh).Cmp(big.NewInt(0))
	if bit > 0 {
		b = true
	}

	return b
}

/* R = kP  */
func (R *Point) Mul(k *big.Int, P *Point, C *Curve) *Point {
	//If at infinity R is also at infinity

	if P.inf {

		R.inf = true
	} else {

		x := new(Point)
		t := new(Point)

		x.Set(P)
		t.Set(x)

		R.inf = true

		/*
			Loops through the integer bit per bit, if a bit is 1 then x is added to the result. Looping through the multiplier in this manner allows us to use as many point doubling operations as possible. No reason to say 5P=P+P+P+P+P, when you might as well just use 5P=2(2P)+P.
			This is not the most effecient method of point multiplication, but it's faster than P+P+P+... which is not computational feasiable.
		*/
		// 		T := new(Point)
		// 		R.SetString("0","0",10)
		for i := int64(0); i < int64(k.BitLen()); i++ {

			// 			 big.NewInt(1).Lsh(big.NewInt(i))
			if bitCheck(k, uint(i)) {

				t.Add(x, R, C)
				R.Set(t)
				// 				R.Print()
			}
			// 			point_doubling(t, x, curve)
			t.double(x, C)

			// 			point_copy(x, t);
			x.Set(t)

		}
	}
	return R
}
