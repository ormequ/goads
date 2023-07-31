package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type correctCfg struct {
	ID   int64  `env:"ID" env-required:"true"`
	Foo  string `env:"FOO"`
	Some string `env:"SOME"`
}

type incorrectCfg struct {
	NotExists string `env:"NOT_EXISTS" env-required:"true"`
}

func TestWithCorrectCfg(t *testing.T) {
	type args struct {
		path string
	}
	type testCase[T any] struct {
		name      string
		args      args
		want      *T
		wantPanic bool
	}
	tests := []testCase[correctCfg]{
		{
			name: "correct",
			args: args{path: "test.env"},
			want: &correctCfg{
				ID:   100500,
				Foo:  "bar",
				Some: "string",
			},
		},
		{
			name:      "file does not exist",
			args:      args{path: "not_exists.env"},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					MustLoadENV[correctCfg](tt.args.path)
				})
			} else {
				assert.Equal(t, tt.want, MustLoadENV[correctCfg](tt.args.path))
			}
		})
	}
}

func TestWithIncorrectCfg(t *testing.T) {
	type args struct {
		path string
	}
	type testCase[T any] struct {
		name      string
		args      args
		want      *T
		wantPanic bool
	}
	tests := []testCase[incorrectCfg]{
		{
			name:      "panics on incorrect cfg",
			args:      args{path: "test.env"},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					MustLoadENV[incorrectCfg](tt.args.path)
				})
			} else {
				assert.Equal(t, tt.want, MustLoadENV[incorrectCfg](tt.args.path))
			}
		})
	}
}
