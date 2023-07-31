package shutdown

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

var errDone = errors.New("done")

type serverErrCreator struct{}

func (s serverErrCreator) Listen(ctx context.Context) error {
	t := make(chan struct{})
	go func(c chan struct{}) {
		time.Sleep(time.Second)
		c <- struct{}{}
	}(t)
	select {
	case <-ctx.Done():
		return nil
	case <-t:
		return errDone
	}
}

type serverPanicCreator struct{}

const panicTxt = "server didn't stop"

func (s serverPanicCreator) Listen(ctx context.Context) error {
	t := make(chan struct{})
	go func(c chan struct{}) {
		time.Sleep(time.Second * 3)
		c <- struct{}{}
	}(t)
	select {
	case <-ctx.Done():
		return nil
	case <-t:
		panic(panicTxt)
	}
}

func TestGracefully(t *testing.T) {
	ctx := context.Background()
	assert.NotPanics(t, func() {
		eg, ctx := errgroup.WithContext(ctx)
		Gracefully(eg, ctx, serverErrCreator{}, serverPanicCreator{})
		err := eg.Wait()
		assert.ErrorIs(t, err, errDone)
	})
	assert.NotPanics(t, func() {
		eg, ctx := errgroup.WithContext(ctx)
		Gracefully(eg, ctx, serverPanicCreator{})
		err := eg.Wait()
		assert.Equal(t, err.Error(), fmt.Sprintf("server panic with error: %s", panicTxt))
	})
	assert.NotPanics(t, func() {
		ctx, cancel := context.WithCancel(ctx)
		eg, ctx := errgroup.WithContext(ctx)
		Gracefully(eg, ctx, serverPanicCreator{})
		time.Sleep(time.Second)
		cancel()
		before := time.Now()
		err := eg.Wait()
		after := time.Now()
		assert.Equal(t, before.Truncate(time.Second), after.Truncate(time.Second))
		assert.NoError(t, err)
	})
}
