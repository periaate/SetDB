package encconn

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net"
	"testing"
)

type localError string

func (e localError) Error() string {

	return string(e)
}

const ErrDifferentHashes = localError("hashes were different/inconsistent while testing")

type fakeConn struct {
	net.Conn
	ch chan []byte
}

func (f *fakeConn) Read(b []byte) (int, error) {
	data := <-f.ch
	return io.ReadFull(bytes.NewBuffer(data), b)
}

func (f *fakeConn) Write(b []byte) (int, error) {
	f.ch <- b

	return len(b), nil
}

func TestAES(t *testing.T) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("foda-se")

	// The channel must have enough bytes so it won't deadlock
	// It's only for the test tho
	f := fakeConn{
		ch: make(chan []byte, 5000),
	}

	aesconn, err := NewAESConn(context.Background(), &f, key)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	_, err = aesconn.Write(plaintext)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	buffer := make([]byte, len(plaintext))
	_, err = aesconn.Read(buffer)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	h1 := md5.New()
	h2 := md5.New()

	h1.Write(buffer)
	h2.Write(plaintext)

	if string(h1.Sum(nil)) != string(h2.Sum(nil)) {
		t.Error()
	}
}
