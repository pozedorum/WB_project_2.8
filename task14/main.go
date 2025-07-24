package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Синтаксический анализатор жалуется на то, что контекст отменяется не во всех ветках,
// но это нужно, чтобы закрывать все горутины при заверщении одной из них

func or(channels ...<-chan interface{}) <-chan interface{} {
	var once sync.Once
	res := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())

	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			defer cancel()
			select {
			case <-ctx.Done():
				return
			case _, ok := <-ch:
				if !ok {
					once.Do(func() {
						close(res)
					})
				}
			}
		}(ch)
	}

	return res
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
