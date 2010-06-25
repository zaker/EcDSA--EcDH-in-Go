package ecc

import (
	"big"
	"fmt"
)


type Point struct {
	X *big.Int

	Y   *big.Int
	inf bool
}

func (R *Point) Copy(P *Point) *Point{
	
	R.inf = P.inf
	R.X = P.X
	R.Y = P.Y
	
	
	return R
}

func (P *Point) Print() {

	fmt.Printf("%s", P.String())

}

func (P *Point) String() string {
	s := ""

	if P.inf {
		s = "Point at infinity"
	} else {
		s = "\nPoint: (\nX:\t" + P.X.String() + ",\nY:\t" + P.Y.String() + ")\n"
	}
	return s
}

func (P *Point) SetString(X,Y string,base int) *Point{
// 	P = new(Point)
	

	P.X = new(big.Int)
	P.Y = new(big.Int)
	
	
	P.inf = false
	P.X.SetString(X,base)
	P.Y.SetString(Y,base)
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
	return P.X.Cmp(Q.X) == 0 && P.X.Cmp(Q.X) == 0

}


/* R = 2P */
func (R *Point) double(P *Point, C *Curve) *Point {
	// 	void point_doubling(point R, point P, domain_parameters curve)
	// {
	// 	printf("double:\n");
	// 	point_print(P);point_print(R);
	//If at infinity
	if P.inf {
		R.inf = true
	} else {
		//Calculate slope
		//s = (3*Px² + a) / (2*Py) mod p
		// 		number_theory_exp_modp_ui(t1, P->x, 2, curve->p);
		t1 := new(big.Int).Exp(P.X, big.NewInt(2), C.p)

		s := new(big.Int).Mul(t1, big.NewInt(3))
		s.Mod(s, C.p)
		s.Add(s, C.a)
		s.Mod(s, C.p)

		sd := new(big.Int).Mul(big.NewInt(2), P.Y)

		sd.ModInverse(sd, C.p)

		s.Mul(s, sd)
		s.Mod(s, C.p)
		// 		printf("helo\n");//t1 = Px² mod p

		// 		  mpz_mul_ui(t2, t1, 3);	;//t2 = 3 * t1
		// 		mpz_mod(t3, t2, curve->p);		//t3 = t2 mod p
		// 		mpz_add(t4, t3, curve->a);		;//t4 = t3 + a
		// 		mpz_mod(t5, t4, curve->p);			//t5 = t4 mod p
		//
		// 		mpz_mul_ui(t1, P->y, 2);			//t1 = 2*Py
		// 		number_theory_inverse(t2, t1, curve->p);		//t2 = t1^-1 mod p
		// 		mpz_mul(t1, t5, t2);				//t1 = t5 * t2
		// 		mpz_mod(s, t1, curve->p);			//s = t1 mod p
		//

		t := new(big.Int).Exp(s, big.NewInt(2), C.p)

		t2 := new(big.Int).Mul(P.X, big.NewInt(2))

		t2.Mod(t2, C.p)

		t.Sub(t, t2)

		R.X.Mod(t, C.p)
		//Calculate Rx
		//Rx = s² - 2*Px mod p
		// 		number_theory_exp_modp_ui(t1, s, 2, curve->p);//t1 = s² mod p
		// 		mpz_mul_ui(t2, P->x, 2);		//t2 = Px*2
		// 		mpz_mod(t3, t2, curve->p);		//t3 = t2 mod p
		// 		mpz_sub(t4, t1, t3);			//t4 = t1 - t3
		// 		mpz_mod(R->x, t4, curve->p);	//Rx = t4 mod p
		// 		printf("helo\n");
		//Calculate Ry using algorithm shown to the right of the commands
		//Ry = s(Px-Rx) - Py mod p

		// 		R.Y.Sub(s,R.X.Sub(P.X,R.X)).Mod(P.Y,C.p)

		t.Sub(P.X, R.X)
		// 		mpz_sub(t1, P->x, R->x);			//t1 = Px - Rx
		t.Mul(s, t)
		// 		mpz_mul(t2, s, t1);					//t2 = s*t1
		t.Sub(t, P.Y)
		// 		mpz_sub(t3, t2, P->y);				//t3 = t2 - Py

		R.Y.Mod(t, C.p)
		// 		mpz_mod(R->y, t3, curve->p);	//Ry = t3 mod p

	}
	return R
}


/* R = P + Q mod C.p*/
func (R *Point) Add(P, Q *Point, C *Curve) *Point {

	//If Q is at infinity, set result to P
	// 	printf("Px:%s\nPy:%s\nQx:%s\nQy:%s\n",mpz_get_str(NULL,16,P->x),mpz_get_str(NULL,16,P->y),mpz_get_str(NULL,16,Q->x),mpz_get_str(NULL,16,Q->y));
	// 	printf("Pi:%d\nQi:%d\n",P->infinity,Q->infinity);
	// 	printf("add:\n");
	// 	point_print(P);point_print(Q);

	print("start tests\n",Cmp(P,Q),"\n")
	switch {
	case Q.inf:
		print("start test1\n")
		return P

	case P.inf:
		print("start test2\n")
		return Q

	case Cmp(P, Q):
		print("start test3\n")
		return R.double(Q, C)

	}
	print("done tests\n")
	//Calculate the inverse point

	//If it is the inverse
	if Cmp(Q.Inv(Q, C), P) {
		//result must be point at infinity
		R.inf = true
	} else {
		/*
			Modulo algebra rules:
			(b1 + b2) mod  n = (b2 mod n) + (b1 mod n) mod n
			(b1 * b2) mod  n = (b2 mod n) * (b1 mod n) mod n
		*/

		//Calculate slope
		//s = (Py - Qy)/(Px-Qx) mod p
		s := new(big.Int).Sub(P.Y, Q.Y)
		sd := new(big.Int).Sub(P.X, Q.X)

		sd.ModInverse(sd, C.p)
		// 			mpz_sub(t1, P->y, Q->y);
		// 			mpz_sub(t2, P->x, Q->x);
		//Using Modulo to stay within the group!

		// 			number_theory_inverse(t3, t2, curve->p); //Handle errors
		// 			mpz_invert(t3, t2, curve->p);

		// 			mpz_mul(t4, t1, t3);
		// 			mpz_mod(s, t4, curve->p);


		//Calculate Rx using algorithm shown to the right of the commands
		//Rx = s² - Px - Qx = (s² mod p) - (Px mod p) - (Qx mod p) mod p
		x := new(big.Int).Exp(s, big.NewInt(2), C.p)
		// 		number_theory_exp_modp_ui(t1, s, 2, curve->p);	//t1 = s² mod p
		// 			mpz_powm_ui(t1, s, 2, curve->p);
		x.Sub(x, P.X)
		x.Sub(x, Q.X)

		R.X.Mod(x, C.p)

		// 		mpz_mod(t2, P->x, curve->p);		//t2 = Px mod p
		// 		mpz_mod(t3, Q->x, curve->p);		//t3 = Qx mod p
		// 		mpz_sub(t4, t1, t2);				//t4 = t1 - t2
		// 		mpz_sub(t5, t4, t3);				//t5 = t4 - t3
		// 		mpz_mod(result->x, t5, curve->p);	//R->x = t5 mod p

		//Calculate Ry using algorithm shown to the right of the commands
		//Ry = s(Px-Rx) - Py mod p
		y := new(big.Int).Sub(P.X, R.X)

		y.Mul(s, y)

		y.Sub(y, P.Y)

		R.Y.Mod(y, C.p)
		// 		mpz_sub(t1, P->x, result->x);		//t1 = Px - Rx
		// 		mpz_mul(t2, s, t1);					//t2 = s*t1
		// 		mpz_sub(t3, t2, P->y);				//t3 = t2 - Py
		// 		mpz_mod(result->y, t3, curve->p);	//Ry = t3 mod p

	}

	return R
}

/* R = -J mod C.p, R=(Jx,-Jy mod C.p)*/
func (R *Point) Inv(J *Point, C *Curve) *Point {
	R = J

	if !J.inf {
		R.Y.Sub(C.p, J.Y)
	}

	return R
}

/* R = kP  */
func (R *Point) Mul(k *big.Int, P *Point, C *Curve) *Point {
	//If at infinity R is also at infinity

	// 	printf("mult(%s)\n",mpz_get_str(NULL,16,multiplier));
	// 	point_print(P);point_print(R);
	print("multiplying\n")
	P.Print()
	if P.inf {

		R.inf = true
	} else {
		//Set R = point at infinity
		// 		point_at_infinity(R);
		
		/*
			Loops through the integer bit per bit, if a bit is 1 then x is added to the result. Looping through the multiplier in this manner allows us to use as many point doubling operations as possible. No reason to say 5P=P+P+P+P+P, when you might as well just use 5P=2(2P)+P.
			This is not the most effecient method of point multiplication, but it's faster than P+P+P+... which is not computational feasiable.
		*/
		T := new(Point)
		T.Copy(P)
		for i := int64(0); i < int64(k.BitLen()); i++ {
			// 		while(bit <= bits)
			// 		{
			// 			printf("%ld\n",bit);
			// 			bi := big.NewInt(i)
			if new(big.Int).And(k, big.NewInt(i)) != big.NewInt(0) {
				
				R.Copy(T)
				R.Add(T, R, C)
			}
			T.double(T, C)
		}
	}
	return R
}
