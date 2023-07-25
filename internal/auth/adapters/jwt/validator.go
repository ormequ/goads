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
)

type Validator struct {
	publicKey *rsa.PublicKey
}

// Validate receives JWT token and returns user's ID extracted from there
func (v Validator) Validate(ctx context.Context, token string) (int64, error) {
	if ctx.Err() != nil {
		return -1, ctx.Err()
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return "", fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return v.publicKey, nil
	})
	if err != nil {
		return -1, errors.Join(app.ErrInvalidToken, err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return -1, app.ErrInvalidToken
	}
	id, ok := claims["dat"].(float64)
	if !ok {
		return -1, app.ErrInvalidToken
	}
	return int64(id), nil
}

func NewValidator(publicKey []byte) (Validator, error) {
	block, _ := pem.Decode(publicKey)
	k, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return Validator{}, err
	}
	return Validator{k.(*rsa.PublicKey)}, nil
}
