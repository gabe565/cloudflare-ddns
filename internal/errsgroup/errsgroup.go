package errsgroup

import (
	"errors"
	"sync"
)

type Group struct {
	wg   sync.WaitGroup
	errs []error
	mu   sync.Mutex
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.mu.Lock()
			g.errs = append(g.errs, err)
			g.mu.Unlock()
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return errors.Join(g.errs...)
}
