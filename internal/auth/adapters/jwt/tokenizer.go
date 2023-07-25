package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Tokenizer struct {
	privateKey *rsa.PrivateKey
	expires    time.Duration
}

func (t Tokenizer) Generate(ctx context.Context, id int64) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["dat"] = id                        // user data - id
	claims["exp"] = now.Add(t.expires).Unix() // expires
	claims["iat"] = now.Unix()                // issued at
	claims["nbf"] = now.Unix()                // not before

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(t.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func NewTokenizer(expires time.Duration, privateKey []byte) (Tokenizer, error) {
	block, _ := pem.Decode(privateKey)
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return Tokenizer{}, err
	}
	return Tokenizer{
		expires:    expires,
		privateKey: k.(*rsa.PrivateKey),
	}, nil
}
