package blowcbc

import (
	"bytes"
	"testing"
)

func TestPadd(t *testing.T) {

	if len(padd(make([]byte, 3), byte(8))) != 8 {
		t.Errorf("Padding of 3 Failed")

	}
	if len(padd(make([]byte, 8), byte(8))) != 16 {
		t.Errorf("Padding of 8 Failed")
	}

	if len(depadd(padd(make([]byte, 3), byte(8)))) != 3 {
		t.Errorf("dePadding of 3 Failed")

	}
	if len(depadd(padd(make([]byte, 8), byte(8)))) != 8 {
		t.Errorf("dePadding of 8 Failed")
	}

	if len(depadd(padd(make([]byte, 4321), byte(8)))) != 4321 {
		t.Errorf("dePadding of 4321 Failed")
	}

	if len(depadd(padd(make([]byte, 4328), byte(8)))) != 4328 {
		t.Errorf("dePadding of 4328 Failed")
	}

}

func TestCBC(t *testing.T) {

	iv := []byte{0, 0, 0, 0, 0, 0, 0, 1}

	s := "TOP SECRET: Meeting at the docks"

	bs := bytes.NewBufferString(s)

	pt := bs.Bytes()

	key := []byte{4, 3, 2, 1, 8, 7, 6, 5}

	ct := CBCEncrypt(pt, key, iv)

	pt2 := CBCDecrypt(ct, key, iv)

	bs2 := bytes.NewBuffer(pt2)

	s2 := bs2.String()

	if s != s2 {
		t.Errorf("Encrypt/Decrypt Failed")
		println(s2)
	}
}
