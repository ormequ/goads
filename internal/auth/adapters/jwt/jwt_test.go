package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func setup(t require.TestingT, expires time.Duration) (Tokenizer, Validator) {
	private, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	require.NoError(t, err)
	return Tokenizer{
		privateKey: private,
		expires:    expires,
	}, Validator{publicKey: private.Public().(*rsa.PublicKey)}
}

func FuzzCorrectJWT(f *testing.F) {
	tok, val := setup(f, time.Hour)
	f.Add(int64(0))
	f.Fuzz(func(t *testing.T, id int64) {
		token, err := tok.Generate(context.Background(), id)
		require.NoError(t, err)
		got, err := val.Validate(context.Background(), token)
		require.NoError(t, err)
		require.Equal(t, id, got)
	})
}

func TestJWT(t *testing.T) {
}
