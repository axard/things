package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	NumberOfAbortSignals = 2
)

type Action interface {
	Do()
}

type ActionFunc func()

func (this ActionFunc) Do() {
	this()
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ticker := time.NewTicker(1 * time.Second)

	hook := Hook{}
	hook.Append(ActionFunc(cancelFunc))
	hook.Append(ActionFunc(ticker.Stop))

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case t := <-ticker.C:
				str := t.Format("15:04:05")
				fmt.Printf("\r%s\r%s", strings.Repeat(" ", len(str)), str)
			}
		}
	}()

	go func() {
		defer wg.Done()

		sigchan := make(chan os.Signal, NumberOfAbortSignals)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

		<-sigchan

		// TODO: добавить вызов хука
		hook.Do()
	}()

	wg.Wait()
}
