package controlFlow

import (
	"errors"
	"fmt"
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
	var errs []error
	for i := 0; i < len(funcs); i++ {
		if err := <-ch; err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
