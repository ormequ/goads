package bcrypt

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"goads/internal/auth/app"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestBCrypt_Generate(t *testing.T) {
	type fields struct {
		Cost int
	}
	type args struct {
		ctx      context.Context
		password string
	}

	bgCtx := context.Background()
	canceledCtx, cancel := context.WithCancel(bgCtx)
	cancel()

	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "correct generation",
			fields: fields{Cost: 10},
			args: args{
				ctx:      bgCtx,
				password: "qwe123!!~....утф😊",
			},
			want: hash("qwe123!!~....утф😊", 10),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
		{
			name:   "canceled context",
			fields: fields{Cost: 10},
			args: args{
				ctx:      bgCtx,
				password: "qwe123!!~....утф😊",
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, canceledCtx.Err())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BCrypt{
				cost: tt.fields.Cost,
			}
			got, err := b.Generate(tt.args.ctx, tt.args.password)
			if !tt.wantErr(t, err, fmt.Sprintf("Generate(%v, %v)", tt.args.ctx, tt.args.password)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Generate(%v, %v)", tt.args.ctx, tt.args.password)
		})
	}
}

func hash(s string, cost int) string {
	res, _ := bcrypt.GenerateFromPassword([]byte(s), cost)
	return string(res)
}

func TestBCrypt_Compare(t *testing.T) {
	type fields struct {
		Cost int
	}
	type args struct {
		ctx      context.Context
		hash     string
		password string
	}

	bgCtx := context.Background()
	canceledCtx, cancel := context.WithCancel(bgCtx)
	cancel()

	tests := [...]struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "correct comparison",
			fields: fields{Cost: 10},
			args: args{
				ctx:      bgCtx,
				hash:     hash("qwe123!!~....утф😊", 10),
				password: "qwe123!!~....утф😊",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
		},
		{
			name:   "incorrect password",
			fields: fields{Cost: 10},
			args: args{
				ctx:      bgCtx,
				hash:     hash("qwe123!!~....утф😊", 10),
				password: "qwe123!!~....ут😊ф",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, app.ErrIncorrectCredentials)
			},
		},
		{
			name:   "context canceled",
			fields: fields{Cost: 10},
			args: args{
				ctx:      canceledCtx,
				hash:     hash("qwe123!!~....утф😊", 10),
				password: "qwe123!!~....утф😊",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, canceledCtx.Err())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BCrypt{
				cost: tt.fields.Cost,
			}
			tt.wantErr(t, b.Compare(tt.args.ctx, tt.args.hash, tt.args.password), fmt.Sprintf("Compare(%v, %v, %v)", tt.args.ctx, tt.args.hash, tt.args.password))
		})
	}
}
