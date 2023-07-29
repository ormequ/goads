package generator

import (
	"context"
	"fmt"
	"math"
	"math/rand"
)

type Repository interface {
	SizeApprox(ctx context.Context) (int64, error)
}

type Generator struct {
	Repo Repository
}

const symbols = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Generate generates alias for url with dynamic length = max(4, 1.5*log[62, n]), where
// n = estimating of repository size. With this function probability of duplication lower
// than 1/62^2 = 1/3844. With n > 15e6 P(duplication) = 1/sqrt(n)
func (g Generator) Generate(ctx context.Context) (string, error) {
	sz, err := g.Repo.SizeApprox(ctx)
	if err != nil {
		return "", err
	}
	if sz < 1 {
		sz = 1 // prevent log(x), x <= 0
	}
	deg := math.Log(float64(sz)) / math.Log(float64(len(symbols)))
	l := int(math.Max(4, deg))
	fmt.Println(deg, l, sz)
	res := make([]byte, l)
	for i := range res {
		res[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(res), nil
}

func New(repo Repository) Generator {
	return Generator{repo}
}
