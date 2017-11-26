package main_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/oklog/run"
	"github.com/socialpoint-labs/bsk/contextx"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

const delay = time.Millisecond * 500

func sumOf(till int) int {
	var sum int
	for i := 1; i <= till; i++ {
		printlnc(till, "sum", i, "of", till)
		sum += i
		time.Sleep(delay)
	}

	fmt.Printf("=> sum of first %d is %d\n", till, sum)
	return sum
}

func gauss(n int) int {
	return n * (n + 1) / 2
}

// #####################################################################################################

func TestDoOne(t *testing.T) {
	assert := assert.New(t)

	times := 10
	assert.Equal(gauss(times), sumOf(times))
}

// #####################################################################################################

var numbers = []int{1, 5, 10}

func TestDoSeveralOK(t *testing.T) {
	assert := assert.New(t)

	for _, number := range numbers {
		assert.Equal(gauss(number), sumOf(number))
	}
}

// #####################################################################################################

func timer() {
	defer println("finished timer")
	for {
		<-time.After(time.Second)
		println(time.Now().Format("2006-01-02 15:04:05"))
	}
}

func TestDoSeveralConcurrentWrong(t *testing.T) {
	assert := assert.New(t)

	go timer()
	for _, number := range numbers {
		number := number
		go func() {
			assert.Equal(gauss(number), sumOf(number))
		}()
	}
}

func TestDoSeveralConcurrentClientSide(t *testing.T) {
	assert := assert.New(t)

	go timer()

	var wg sync.WaitGroup
	wg.Add(len(numbers))
	for _, number := range numbers {
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		number := number
		go func() {
			assert.Equal(gauss(number), sumOf(number))
			wg.Done()
		}()
	}
	wg.Wait()
}

// #####################################################################################################

func TestDoSeveralConcurrentServerSideStuck(t *testing.T) {
	assert := assert.New(t)

	csum := concurrentSummer{
		make(chan int, len(numbers)),
		make(chan int),
	}

	go timer()
	go csum.do()

	for _, number := range numbers {
		println("put", number)
		csum.in <- number
	}
	for range csum.out {
	}
	assert.True(true)
}

type concurrentSummer struct {
	in  chan int
	out chan int
}

func (cs concurrentSummer) do() {
	println("finished concurrentSummer")
	for {
		select {
		case in := <-cs.in:
			go func() {
				cs.out <- sumOf(in)
			}()
		}
	}
}

// #####################################################################################################

func TestDoSeveralConcurrentSignalClosesChan(t *testing.T) {
	assert := assert.New(t)

	csum := concurrentSummer{
		make(chan int, len(numbers)),
		make(chan int),
	}

	done := signalNotifier()
	go timerDone(done)
	go csum.run(done)

	for _, number := range numbers {
		println("put", number)
		csum.in <- number
	}
	assert.Equal(71, sumFromChan(csum.out))
}

func sumFromChan(ch chan int) int {
	var i, totalSum int
	for sum := range ch {
		totalSum += sum
		i++
	}
	return totalSum
}

func timerDone(done chan struct{}) {
	defer println("finished timerDone")
	for {
		select {
		case <-time.After(time.Second):
			println(time.Now().Format("2006-01-02 15:04:05"))
		case <-done:
			println("done timerDone")
			return
		}
	}
}

func (cs concurrentSummer) run(done chan struct{}) {
	defer func() {
		println("finished concurrentSummer")
		close(cs.out)
	}()
	for {
		select {
		case in := <-cs.in:
			go func() {
				cs.out <- sumOf(in)
			}()
		case <-done:
			println("done concurrentSummer")
			return
		}
	}
}

// simplified version of https://gobyexample.com/signals
func signalNotifier() chan struct{} {
	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(done)
	}()

	return done
}

// #####################################################################################################

func TestDoSeveralConcurrentSignalCancelsContext(t *testing.T) {
	assert := assert.New(t)

	csum := concurrentSummer{
		make(chan int, len(numbers)),
		make(chan int),
	}

	ctx := context.Background()
	ctx = cancelWithSignal(ctx)
	go timerCtx(ctx)
	go csum.runCtx(ctx)

	for _, number := range numbers {
		println("put", number)
		csum.in <- number
	}
	assert.Equal(71, sumFromChan(csum.out))
}

func timerCtx(ctx context.Context) {
	defer println("finished timerCtx")
	for {
		select {
		case <-time.After(time.Second):
			println(time.Now().Format("2006-01-02 15:04:05"))
		case <-ctx.Done():
			println("ctx done timerCtx")
			return
		}
	}
}

func (cs concurrentSummer) runCtx(ctx context.Context) {
	defer func() {
		println("finished runCtx")
		close(cs.out)
	}()
	for {
		select {
		case in := <-cs.in:
			go func() {
				cs.out <- sumOf(in)
			}()
		case <-ctx.Done():
			println("ctx done runCtx concurrentSummer")
			return
		}
	}
}

func cancelWithSignal(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	return ctx
}

// #####################################################################################################

func TestDoSeveralConcurrentWithBSK(t *testing.T) {
	assert := assert.New(t)

	csum := concurrentSummer{
		make(chan int, len(numbers)),
		make(chan int),
	}

	ctx := context.Background()
	go timerCtx(ctx)
	go contextx.PosixSignalsAdapter().Adapt(csum).Run(ctx)

	for _, number := range numbers {
		println("put", number)
		csum.in <- number
	}
	assert.Equal(71, sumFromChan(csum.out))
}

// Run implements BSK's contextx.Runner interface
func (cs concurrentSummer) Run(ctx context.Context) {
	defer func() {
		println("finished Run concurrentSummer")
		close(cs.out)
	}()
	for {
		select {
		case in := <-cs.in:
			go func() {
				defer func() {
					if r := recover(); r != nil {
						//println("send to out recovered")
					}
				}()
				cs.out <- sumOf(in)
			}()
		case <-ctx.Done():
			println("ctx done Run concurrentSummer")
			return
		}
	}
}

// #####################################################################################################

func timerCtxErr(ctx context.Context) error {
	defer println("finished timerCtx")
	var i int
	for {
		select {
		case <-time.After(time.Second):
			i++
			if i == 3 {
				err := errors.New("TIMER GENERATED ERROR")
				println(err.Error())
				return err
			}
			println(time.Now().Format("2006-01-02 15:04:05"))
		case <-ctx.Done():
			println("ctx done timerCtx")
			return nil
		}
	}
	return nil
}

func TestDoSeveralConcurrentWithErrgroup(t *testing.T) {
	assert := assert.New(t)

	csum := concurrentSummer{
		make(chan int, len(numbers)),
		make(chan int),
	}

	ctx := context.Background()
	ctx = cancelWithSignal(ctx)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := timerCtxErr(ctx)
		return err
	})
	g.Go(func() error {
		csum.Run(ctx)
		return nil
	})

	for _, number := range numbers {
		println("put", number)
		csum.in <- number
	}

	err := g.Wait()
	assert.NoError(err)
	assert.Equal(71, sumFromChan(csum.out))
}

// #####################################################################################################

func TestDoSeveralConcurrentWithOKLogRun(t *testing.T) {
	assert := assert.New(t)

	csum := concurrentSummer{
		make(chan int, len(numbers)),
		make(chan int),
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx = cancelWithSignal(ctx)

	var g run.Group
	g.Add(func() error {
		return timerCtxErr(ctx)
	}, func(error) {
		cancel()
	})
	g.Add(func() error {
		csum.Run(ctx)
		return nil
	}, func(error) {
		cancel()
	})

	go func() {
		for _, number := range numbers {
			println("put", number)
			csum.in <- number
		}
	}()

	err := g.Run()
	assert.NoError(err)
	assert.Equal(71, sumFromChan(csum.out))
}

// #####################################################################################################

var mu sync.Mutex

func printlnc(i int, params ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	var c color.Attribute
	switch i {
	case 1:
		c = color.FgCyan
	case 5:
		c = color.FgGreen
	case 10:
		c = color.FgRed
	}
	color.Set(c)
	fmt.Println(params...)
	color.Unset()
}
