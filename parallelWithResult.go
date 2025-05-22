package controlFlow

import (
	"errors"
	"fmt"
)

func NewParallelWithResult[T any](funcs ...func() (*T, error)) *ParallelWithResult[T] {
	return &ParallelWithResult[T]{
		funcs: funcs,
	}
}

type result[T any] struct {
	Val *T
	Err error
}

type ParallelWithResult[T any] struct {
	funcs []func() (*T, error)
}

func (pt *ParallelWithResult[T]) Add(funcs ...func() (*T, error)) {
	pt.funcs = append(pt.funcs, funcs...)
}

func (pt *ParallelWithResult[T]) Run() ([]*T, error) {
	return parallelWithResult[T](pt.funcs...)
}

func parallelWithResult[T any](funcs ...func() (*T, error)) ([]*T, error) {
	var ch = make(chan *result[T], len(funcs))
	defer close(ch)
	for _, fn := range funcs {
		go func(f func() (*T, error)) {
			defer func() {
				if e := recover(); e != nil {
					ch <- &result[T]{
						Err: errors.New(fmt.Sprintf("%v", e)),
					}
				}
			}()
			r, err := f()
			ch <- &result[T]{
				Val: r,
				Err: err,
			}
		}(fn)
	}
	var results []*T
	var errs []error
	for i := 0; i < len(funcs); i++ {
		r := <-ch
		if r.Err != nil {
			errs = append(errs, r.Err)
		}
		results = append(results, r.Val)
	}
	return results, errors.Join(errs...)
}
