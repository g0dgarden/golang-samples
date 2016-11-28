package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
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
		u := fmt.Sprintf("http://minna-no-go/%d", i)
		d.Work(ctx, func(ctx context.Context) {
			log.Printf("start processing %s", u)
			// 本当だったらここでURL取りに行くとかするけど、ダミーだから
			// 適当な時間待ちます
			t := time.NewTimer(time.Duration(rand.Intn(5)) * time.Second)
			defer t.Stop()

			select {
			case <-ctx.Done():
				log.Printf("cancel work func %s", u)
				return
			case <-t.C:
				log.Printf("done processing %s", u)
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