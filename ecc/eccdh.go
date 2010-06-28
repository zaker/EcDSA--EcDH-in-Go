package ecc

import (
	"big"
)


func (E *EccKeyPair) EccDH(q *Point) *big.Int {

	key := new(Point)

	key.Mul(E.d, q, E.C)

	return key.x
}
