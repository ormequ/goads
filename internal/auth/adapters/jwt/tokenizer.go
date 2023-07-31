package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"goads/internal/auth/app"
	"goads/internal/pkg/errwrap"
	"time"
)

type Tokenizer struct {
	privateKey *rsa.PrivateKey
	expires    time.Duration
}

func (t Tokenizer) Generate(ctx context.Context, id int64) (string, error) {
	const op = "jwt.Generate"

	if ctx.Err() != nil {
		return "", errwrap.New(ctx.Err(), app.ServiceName, op)
	}

	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["dat"] = id                        // user data - id
	claims["exp"] = now.Add(t.expires).Unix() // expires
	claims["iat"] = now.Unix()                // issued at
	claims["nbf"] = now.Unix()                // not before

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(t.privateKey)
	if err != nil {
		err = errwrap.New(errors.New("cannot sign token"), app.ServiceName, op).WithDetails(err.Error())
	}
	return token, err
}

func NewTokenizer(expires time.Duration, privateKey []byte) (Tokenizer, error) {
	const op = "jwt.NewTokenizer"

	block, _ := pem.Decode(privateKey)
	fmt.Println(block)
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return Tokenizer{}, errwrap.New(err, app.ServiceName, op)
	}
	return Tokenizer{
		expires:    expires,
		privateKey: k.(*rsa.PrivateKey),
	}, nil
}
