package controlFlow

import (
	"errors"
	"fmt"
	"sync"
)

func NewParallel(funcs ...func() error) *Parallel {
	return &Parallel{
		funcs: funcs,
	}
}

type Parallel struct {
	funcs []func() error
}

func (pt *Parallel) Add(funcs ...func() error) {
	pt.funcs = append(pt.funcs, funcs...)
}

func (pt *Parallel) Run() error {
	return parallel(pt.funcs...)
}

func (pt *Parallel) RunBreakOnError() error {
	return parallelBreakOnError(pt.funcs...)
}

func parallel(funcs ...func() error) error {
	var ch = make(chan error, len(funcs))
	defer close(ch)
	for _, fn := range funcs {
		go func(f func() error) {
			defer func() {
				if e := recover(); e != nil {
					ch <- errors.New(fmt.Sprintf("%v", e))
				}
			}()
			ch <- f()
		}(fn)
	}
	var err error
	for i := 0; i < len(funcs); i++ {
		if e := <-ch; e != nil {
			err = e
		}
	}
	return err
}

func parallelBreakOnError(funcs ...func() error) error {
	if len(funcs) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	ch := make(chan error, len(funcs))

	for _, fn := range funcs {
		wg.Add(1)
		go func(f func() error) {
			defer wg.Done()
			defer func() {
				if e := recover(); e != nil {
					ch <- errors.New(fmt.Sprintf("%v", e))
				}
			}()
			if err := f(); err != nil {
				ch <- err
			}
		}(fn)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for err := range ch {
		return err
	}
	return nil
}
