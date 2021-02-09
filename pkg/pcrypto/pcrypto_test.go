package pcrypto

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"testing"
)

func TestPcrypto_Encrypto(t *testing.T){
	pp := new(Pcrypto)
	pp.Init([]byte("Test"))
	res, err := pp.Encrypto([]byte("test content"))
	if err != nil{
		t.Error(err)
	}

	fmt.Printf("[%x]\n", res)
}

func TestPcrypto_Decrypto(t *testing.T) {
	content := []byte("test content")
	pp := new(Pcrypto)
	pp.Init([]byte("Test"))
	res, err := pp.Encrypto(content)
	if err != nil{
		t.Error(err)
	}

	res, err = pp.Decrypto(res)

	if err != nil{
		t.Error(err)
	}

	if !bytes.Equal(content, res) {
		t.Error("decrypt error")
	}

	fmt.Printf("[%s]\n", string(res))
}

func TestPKCS7Padding(t *testing.T) {
	ltt := []byte("test content")
	ltt = PKCS7Padding(ltt, aes.BlockSize)
	fmt.Printf("[%x]\n", ltt)
}

func TestPKCS7UnPadding(t *testing.T) {
	oldLtt := []byte("test content")
	ltt := PKCS7Padding(oldLtt, aes.BlockSize)
	ltt = PKCS7UnPadding(ltt)
	if !bytes.Equal(oldLtt, ltt){
		t.Error("unpadding error")
	}
	fmt.Printf("[%x]\n", ltt)
}