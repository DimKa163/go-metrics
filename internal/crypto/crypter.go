package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type Encrypter struct {
	key *rsa.PublicKey
}

func NewEncrypter(path string) (*Encrypter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("x509: failed to parse RSA public key")
	}
	cert, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	pub, ok := cert.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not RSA public key")
	}
	return &Encrypter{key: pub}, nil
}

func (e *Encrypter) Encrypt(plaintext []byte) ([]byte, error) {
	cipherData, err := rsa.EncryptPKCS1v15(rand.Reader, e.key, plaintext)
	if err != nil {
		return nil, err
	}
	return cipherData, nil
}

type Decrypter struct {
	key *rsa.PrivateKey
}

func NewDecrypter(path string) (*Decrypter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid PEM block type")
	}
	cert, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &Decrypter{key: cert}, nil
}

func (d *Decrypter) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, d.key, ciphertext)
}
