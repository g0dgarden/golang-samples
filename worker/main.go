// see: https://gist.github.com/lestrrat/c9b78369cf9b9c5d9b0c909ed1e2452e

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		<-sigCh
		cancel()
	}()

	d := NewDispatcher(3)
	for i := 0; i < 100; i++ {
		u := fmt.Sprintf("Goroutine:%d, http://www.lancers.jp/work/detail/%d", runtime.NumGoroutine(), i)
		d.Work(ctx, func(ctx context.Context) {
			log.Printf("start processing %s", u)

			t := time.NewTimer(time.Duration(rand.Intn(10)) * time.Second)
			defer t.Stop()

			select {
			case <-ctx.Done():
				log.Printf("cancel work func %s", u)
				return
			case <-t.C:
				log.Printf("done processing  %s", u)
				return
			}
		})
	}

	d.Wait()
}

type Dispatcher struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

type WorkFunc func(context.Context)

func NewDispatcher(max int) *Dispatcher {
	return &Dispatcher{
		sem: make(chan struct{}, max),
	}
}

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

func (d *Dispatcher) Work(ctx context.Context, proc WorkFunc) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.work(ctx, proc)
	}()
}

func (d *Dispatcher) work(ctx context.Context, proc WorkFunc) {
	select {
	case <-ctx.Done():
		log.Printf("cancel work")
		return
	case d.sem <- struct{}{}:
		// got semaphore
		defer func() { <-d.sem }()
	}

	proc(ctx)
}
