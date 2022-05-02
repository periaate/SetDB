package encconn

import (
	"context"
	"crypto/rsa"
	"net"
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
