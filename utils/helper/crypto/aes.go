package crypto

import (
	"crypto/cipher"
	"crypto/aes"
	"bytes"
)

var (
	gDefalutAesKey = []byte{
		0xF3, 0x62, 0x12, 0x05, 0x13, 0xE3, 0x89, 0xFF,
		0x23, 0x11, 0xD7, 0x36, 0x01, 0x23, 0x10, 0x07,
		0x05, 0xA2, 0x10, 0x00, 0x7A, 0xCC, 0x02, 0x3C,
		0x39, 0x01, 0xDA, 0x2E, 0xCB, 0x12, 0x44, 0x8B,
	}
	gAesIV = []byte{
		0x15, 0xFF, 0x01, 0x00, 0x34, 0xAB, 0x4C, 0xD3,
		0x55, 0xFE, 0xA1, 0x22, 0x08, 0x4F, 0x13, 0x07,
	}
)

// Crypto ...
type Crypto interface {
	Encrypt(src []byte) []byte
	Decrypt(src []byte) ([]byte,error)
}

type aesCrypto struct {
	cb cipher.Block
}

// NewAesCrypto ...
func NewAesCrypto() Crypto {
	cblock, err := aes.NewCipher(gDefalutAesKey)
	if err != nil {
		panic(err)
	}

	return &aesCrypto{
		cb: cblock,
	}
}

func (cpt *aesCrypto)BlockSize() int {
	return cpt.cb.BlockSize()
}

// Encrypt ...
func (cpt *aesCrypto) Encrypt(src []byte) []byte{
	stream := cipher.NewCFBEncrypter(cpt.cb, gAesIV)
	src = PKCS7Padding(src,cpt.cb.BlockSize())
	dst := make([]byte,len(src))
	stream.XORKeyStream(dst, src)
	return dst
}

// Decrypt ...
func (cpt *aesCrypto) Decrypt(src []byte) ([]byte,error) {
	stream := cipher.NewCFBDecrypter(cpt.cb, gAesIV)
	dst := make([]byte,len(src))
	stream.XORKeyStream(dst, src)
	dst = PKCS7UnPadding(dst,cpt.cb.BlockSize())
	return dst,nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}