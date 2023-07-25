package jwt

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

func mustReadFile(file string) []byte {
	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func Test(t *testing.T) {
	tokenizer, err := NewTokenizer(720*time.Hour, mustReadFile("../../../../cert/private.key"))
	require.NoError(t, err)
	validator, err := NewValidator(mustReadFile("../../../../cert/public.key.pub"))
	require.NoError(t, err)
	fmt.Println(tokenizer, validator)
	token, err := tokenizer.Generate(context.Background(), 100500)
	require.NoError(t, err)
	fmt.Println(token)
	id, err := validator.Validate(context.Background(), token)
	require.NoError(t, err)
	fmt.Println(id)
	require.Equal(t, id, 100500)
}
