package encryption

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type Error string

func (e Error) Error() string {

	return string(e)
}

const (
	ErrInvalidRSABits = Error("Invalid number of bits for RSA private key generation")
	ErrInvalidRSAPath = Error("Invalid path for private key")
)

// This function generates a new private RSA key provided its path to save on disk
// the path should be specified, otherwise it returns an error
// the bits are optional, defaulting to 4096 bits
func GenerateRSAKey(path string, bits int) (*rsa.PrivateKey, error) {
	if path == "" {
		return nil, ErrInvalidRSAPath
	}

	if bits < 0 {
		return nil, ErrInvalidRSABits
	}

	var _bits int = 0

	if bits == 0 {
		_bits = 4096
	}

	privKey, err := rsa.GenerateKey(rand.Reader, _bits)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(path)
	if err != nil {
		return privKey, err
	}

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	)

	_, err = file.Write(pemdata)
	if err != nil {
		return privKey, err
	}

	return privKey, file.Close()
}

// This function tries to load a private key storage at the path variable
func LoadRSAKey(path string) (*rsa.PrivateKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(buffer.Bytes())

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, file.Close()
}
