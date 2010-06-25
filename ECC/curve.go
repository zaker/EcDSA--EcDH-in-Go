package elliptic_curve

import(
	"big"
)


type Curve struct{
	name string;
	p big.Int;	//Prime
	mpz_t a;	//'a' parameter of the elliptic curve
	mpz_t b;	//'b' parameter of the elliptic curve
	G;	//Generator point of the curve, also known as base point.
	mpz_t n;
	mpz_t h;	
	
}