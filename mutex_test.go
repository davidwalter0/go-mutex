package mutex

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

var N int = 5
var started = make(chan int, N)
var done = make(chan int, N)
var mutex *Mutex = NewMutex()

func ScopedTrace(text string) func() {
	fmt.Println(text)
	return func() {
		fmt.Println(text)
	}
}

func Spaces(n int) string {
	return fmt.Sprintf("%*s", 2*n, " ")
}

func recurse(n int) {
	spaces := Spaces(n)

	defer ScopedTrace(spaces + fmt.Sprintf(">>%d<<", n))()
	if n > 0 {
		if n%2 == 1 {
			recurse(n - 1)
		} else {
			recurse(n - 2)
		}
	}
}

func monitor(n int) {
	started <- 0
	TimedWait()
	{
		defer mutex.Monitor()()
		recurse(n + 2)
	}
	done <- 0
}

func TimedWait() {
	t := rand.Intn(250)
	select {
	case <-time.After(time.Duration(t) * time.Millisecond):
	}
}

func ChanWait(c chan int, max int) {
	for i := 0; i < max; i++ {
		select {
		case <-c:
		}
	}
}

func WaitStart() {
	ChanWait(started, cap(started))
}

func TestMonitor0(t *testing.T) {
	rand.Seed(42)
	runtime.GOMAXPROCS(runtime.NumCPU())
	{
		for i := 0; i < N; i++ {
			go func(n int) {
				defer monitor(n)
			}(i)
		}
		started <- 0
		WaitStart()
		TimedWait()
		{
			defer mutex.Monitor()()
			recurse(3)
			fmt.Println("Exiting Function TestMonitor0")
			done <- 0
		}
	}
}

func TestLast(t *testing.T) {
	ChanWait(done, cap(done))
}
