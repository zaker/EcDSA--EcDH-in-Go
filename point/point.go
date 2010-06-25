package point

import(
	"big"
	"fmt"
)


type Point struct{
	X *big.Int

	Y *big.Int
	inf bool
}

func (P *Point)Print() {

	if P.inf{
		fmt.Printf("Point at infinity")
	}else{
		fmt.Printf("\nPoint: (\n\t");
		print(P.X);
		fmt.Printf("\n,\n\t");
		print(P.Y);
		fmt.Printf("\n)\n");
	}
	
	
}


func Add(J,K Point, c Curve) Point {
	var R Point


	return R
}

func Inv(J Point, C Curve) Point{
	R:=J
	
	R.Y = -J.Y


	return R
}

func Multi(k big.Int, P Point) Point{
	var L Point


	return L
}

