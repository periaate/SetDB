package encryption

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"io"
	"net"
	encconn "setdb/encryption/conn"
)

type RSAListener struct {
	net.Listener
	ctx     context.Context
	privKey rsa.PrivateKey
}

func NewRSAListener(ctx context.Context, listener net.Listener, privKey rsa.PrivateKey) (*RSAListener, error) {
	return &RSAListener{
		Listener: listener,
		ctx:      ctx,
		privKey:  privKey,
	}, nil
}

func _send_pbkey(w io.Writer, pb *rsa.PublicKey) (int, error) {
	messageSizeBuffer := make([]byte, 4)

	data := x509.MarshalPKCS1PublicKey(pb)
	datalen := len(data)
	binary.BigEndian.PutUint32(messageSizeBuffer, uint32(datalen))

	for t := uint32(0); t < 4; {
		n, err := w.Write(data[t:])
		if err != nil {
			return int(t), err
		}
		t += uint32(n)
	}

	for t := uint32(0); t < uint32(datalen); {
		n, err := w.Write(data[t:])
		if err != nil {
			return int(t), err
		}
		t += uint32(n)
	}

	return datalen + 4, nil
}

func _recv_pbkey(w io.Writer, pb *rsa.PublicKey) (int, error) {
	messageSizeBuffer := make([]byte, 4)

	for t := uint32(0); t < 4; {
		n, err := w.Write(messageSizeBuffer[t:])
		if err != nil {
			return int(t), err
		}
		t += uint32(n)
	}

	datalen := binary.BigEndian.Uint32(messageSizeBuffer)
	data := make([]byte, datalen)

	for t := uint32(0); t < uint32(datalen); {
		n, err := w.Write(data[t:])
		if err != nil {
			return int(t), err
		}
		t += uint32(n)
	}

	p, err := x509.ParsePKCS1PublicKey(data)
	if err != nil {
		return -1, err
	}

	*pb = *p

	return int(datalen) + 4, nil
}

func (l *RSAListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	key, err := l.AgreeAESKey(conn)
	if err != nil {
		return nil, err
	}

	return encconn.NewAESConn(l.ctx, conn, key)
}

func (l *RSAListener) AgreeAESKey(rw io.ReadWriter) ([]byte, error) {
	ch := make(chan []byte)
	errch := make(chan error)

	pubKey := &l.privKey.PublicKey

	phrase := make([]byte, 128)
	if _, err := io.ReadFull(rand.Reader, phrase); err != nil {
		return nil, err
	}

	go func() {
		_, err := _send_pbkey(rw, pubKey)
		if err != nil {
			errch <- err
		}

		messageSizeBuffer := make([]byte, 4)

		for t := 0; t < 4; {
			n, err := rw.Read(messageSizeBuffer[t:])
			if err != nil {
				errch <- err
			}
			t += n
		}

		messageSize := binary.BigEndian.Uint32(messageSizeBuffer)
		data := make([]byte, messageSize)

		for t := uint32(0); t < messageSize; {
			n, err := rw.Read(data[t:])
			if err != nil {
				errch <- err
			}
			t += uint32(n)
		}

		data, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, &l.privKey, data, nil)
		if err != nil {
			errch <- err
			return
		}

		result := make([]byte, 128)
		for i := 0; i < 128; i++ {
			result[i] = data[i] | phrase[i]
		}

		ch <- result

	}()

	go func() {
		pb := &rsa.PublicKey{}
		if _, err := _recv_pbkey(rw, pb); err != nil {
			errch <- err
			return
		}

		data, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pb, phrase, nil)
		if err != nil {
			errch <- err
			return
		}

		messageSize := len(data)
		messageSizeBuffer := make([]byte, 4)
		binary.BigEndian.PutUint32(messageSizeBuffer, uint32(messageSize))

		for t := 0; t < 4; {
			n, err := rw.Write(messageSizeBuffer[t:])
			if err != nil {
				errch <- err
			}
			t += n
		}
		for t := 0; t < messageSize; {
			n, err := rw.Write(data[t:])
			if err != nil {
				errch <- err
			}
			t += n
		}
	}()

	select {
	case data := <-ch:
		return data, nil
	case err := <-errch:
		return nil, err
	case <-l.ctx.Done():
		return nil, l.ctx.Err()
	}
}
