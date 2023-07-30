package shutdown

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Listen(ctx context.Context) error
}

func Gracefully(eg *errgroup.Group, ctx context.Context, servers ...Server) {
	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	for _, s := range servers {
		s := s
		eg.Go(func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("server panic with error: %v", r)
				}
			}()
			return s.Listen(ctx)
		})
	}
}
