package encconn

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"io"
	"net"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const ErrInvalidNonceSize = Error("invalid nonce size attempting to encrypt/decrypt")
const ErrCipherDataTooShort = Error("invalid cipher data length (too short)")

type AESConn struct {
	net.Conn
	ctx   context.Context
	block cipher.Block
}

func NewAESConn(ctx context.Context, conn net.Conn, key []byte) (*AESConn, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &AESConn{
		Conn:  conn,
		ctx:   ctx,
		block: block,
	}, nil
}

func (a *AESConn) Read(b []byte) (int, error) {
	nread := make([]byte, 4)
	n, err := a.Conn.Read(nread)
	if err != nil {
		return n, err
	}
	size := binary.BigEndian.Uint32(nread)
	ciphertext := []byte{}
	for t := uint32(0); t < size; t++ {
		tmp := make([]byte, size-t)
		n, err := a.Conn.Read(tmp)
		if err != nil && err != io.EOF {
			return n, err
		}
		ciphertext = append(ciphertext, tmp...)
		t += uint32(n)
	}

	if len(ciphertext) < aes.BlockSize {
		return -1, ErrCipherDataTooShort
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(a.block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return io.ReadFull(bytes.NewBuffer(ciphertext), b)
}

func (a *AESConn) Write(b []byte) (int, error) {
	ciphertext := make([]byte, aes.BlockSize+len(b))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return -1, err
	}

	stream := cipher.NewCFBEncrypter(a.block, iv)

	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)

	tmp := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, uint32(len(ciphertext)))
	_, err := a.Conn.Write(tmp)
	if err != nil {
		return -1, err
	}

	t := 0
	for t < len(ciphertext) {
		n, err := a.Conn.Write(ciphertext[t:])
		if err != nil {
			return n, err
		}
		t += n
	}

	return t, nil
}
