package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"goads/internal/auth/app"
	"goads/internal/pkg/errwrap"
)

type Validator struct {
	publicKey *rsa.PublicKey
}

// Validate receives JWT token and returns user's ID extracted from there
func (v Validator) Validate(ctx context.Context, token string) (int64, error) {
	const op = "jwt.Validate"
	if ctx.Err() != nil {
		return -1, errwrap.New(ctx.Err(), app.ServiceName, op)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return "", fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return v.publicKey, nil
	})
	if err != nil {
		return -1, errwrap.New(app.ErrInvalidToken, app.ServiceName, op).WithDetails(err.Error())
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return -1, errwrap.New(app.ErrInvalidToken, app.ServiceName, op)
	}
	id, ok := claims["dat"].(float64)
	if !ok {
		return -1, errwrap.New(app.ErrInvalidToken, app.ServiceName, op)
	}
	return int64(id), nil
}

func NewValidator(publicKey []byte) (Validator, error) {
	const op = "jwt.NewValidator"

	block, _ := pem.Decode(publicKey)
	k, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return Validator{}, errwrap.New(err, app.ServiceName, op)
	}
	return Validator{k.(*rsa.PublicKey)}, nil
}
