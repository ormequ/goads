package generator

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type repo struct {
	sz int64
}

func (r repo) SizeApprox(_ context.Context) (int64, error) {
	return r.sz, nil
}

func FuzzGenerator_Generate(f *testing.F) {
	f.Fuzz(func(t *testing.T, sz int64) {
		if sz < 0 {
			return
		}
		gen := Generator{repo{sz}}
		alias, err := gen.Generate(context.Background())
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(alias), 4)
		require.LessOrEqual(t, len(alias), 11) // log[62, 2^64-1] ~ 10.75
	})
}
