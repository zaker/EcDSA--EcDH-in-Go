package blowcbc

import (
	"crypto/blowfish"
)

func printBytes(pb []byte) {

	for i := range pb {
		print(pb[i], ",")
	}
	println("")
}

func padd(blocks []byte, bz byte) []byte {
	// 	blocks length
	bl := len(blocks)

	p := bz - (byte(bl) % bz)
	if p == 0 {
		p = bz
	}

	padded := make([]byte, bl+int(p))
	copy(padded, blocks)

	for i := 0; i < int(p); i++ {
		padded[bl+i] = byte(p)
	}

	return padded
}

func depadd(blocks []byte) []byte {

	depadded := blocks[0 : len(blocks)-int(blocks[len(blocks)-1])]
	return depadded
}

func CBCEncrypt(pt, key, iv []byte) []byte {

	blow, _ := blowfish.NewCipher(key)
	defer blow.Reset()
	bz := blow.BlockSize()

	padded := padd(pt, byte(bz))
	encrypted := make([]byte, len(padded))

	for i := 0; i < len(padded); i += bz {

		for j := 0; j < bz; j++ {
			padded[i+j] ^= iv[j]

		}

		blow.Encrypt(padded[i:i+bz], encrypted[i:i+bz])

		iv = encrypted[i : i+bz]
	}

	return encrypted
}

func CBCDecrypt(ct, key, iv []byte) []byte {

	blow, _ := blowfish.NewCipher(key)
	defer blow.Reset()
	bz := blow.BlockSize()

	dc := make([]byte, len(ct))
	for i := 0; i < len(ct); i += bz {

		blow.Decrypt(ct[i:i+bz], dc[i:i+bz])

		for j := 0; j < bz; j++ {
			dc[i+j] ^= iv[j]

		}
		iv = ct[i : i+bz]

	}

	decrypted := depadd(dc)
	return decrypted
}
