package mutex

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

// TestDebug enables verbose output if -v flag is given to test
func TestDebug(*testing.T) {
	testDebug = testing.Verbose()
}

func drain(which string, c chan int) {
	if testDebug {
		fmt.Printf("enter drain %s len %d cap %d\n", which, len(c), cap(c))
	}
	for i := len(c); i > 0; i-- {
		<-c
	}
	if testDebug {
		fmt.Printf("exit  drain %s len %d cap %d\n", which, len(c), cap(c))
	}
}

var N int = 5
var started = make(chan int, N)
var done = make(chan int, N)
var mutex *Mutex = NewMutex()

func recurse(n int) {
	space := spaces(n)
	if testDebug {
		defer scopedTrace(space + fmt.Sprintf(">>%d<<", n))()
	}
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
	if testDebug {
		defer scopedTrace(fmt.Sprintf("ChanWait max %d len %d cap %d", max, len(c), cap(c)))()
	}
	for i := 0; i < max; i++ {
		if testDebug {
			fmt.Println("ChanWait i", i, "len", len(c), "cap", cap(c))
		}
		select {
		case <-c:
			if testDebug {
				fmt.Println("ChanWait i", i, "len", len(c), "cap", cap(c))
			}
		}
	}
}

func WaitForNRoutinesToStart(which string) {
	if testDebug {
		defer scopedTrace(which)()
	}
	ChanWait(started, cap(started))
}

func TestMonitor0(t *testing.T) {
	if testDebug {
		fmt.Println("Entering Function TestMonitor0")
	}
	rand.Seed(42)
	runtime.GOMAXPROCS(runtime.NumCPU())
	{
		for i := 0; i < N; i++ {
			go func(n int) {
				defer monitor(n)
			}(i)
		}
		started <- 0
		WaitForNRoutinesToStart("TestMonitor0")
		TimedWait()
		{
			defer mutex.Monitor()()
			recurse(3)
			done <- 0
		}
	}
	if testDebug {
		fmt.Println("Exiting Function TestMonitor0")
	}
}

func TestFirst(t *testing.T) {
	drain("started", started)
	drain("done", done)
}

func start(i int, m func(...interface{}) func()) {
	defer m(i)()
	if testDebug {
		fmt.Printf("%s Entering start %d\n", spaces(i), i)
		defer scopedTrace(spaces(i) + fmt.Sprintf(" start i %d cap %d len %d", i, len(started), cap(started)))()
	}
	started <- i
	if testDebug {
		fmt.Printf("%s Exiting start %d\n", spaces(i), i)
	}
}

func run1(i int, m func(...interface{}) func()) {
	if testDebug {
		defer scopedTrace("run1")()
	}
	defer m(i)()
	recurse(3)
	done <- 0
	if testDebug {
		fmt.Println("exiting function run1")
	}
}

func TestMonitor1(t *testing.T) {
	if testDebug {
		fmt.Println("Entering Function TestMonitor1")
	}
	m := NewMonitor()
	rand.Seed(42)
	runtime.GOMAXPROCS(runtime.NumCPU())
	{
		for i := 0; i < N; i++ {
			start(i, m)
		}
		WaitForNRoutinesToStart("TestMonitor1")
		TimedWait()
		{
			run1(-1, m)
		}
	}
	if testDebug {
		fmt.Println("Exiting Function TestMonitor1")
	}
}

func TestLast(t *testing.T) {
	if testDebug {
		defer scopedTrace("TestLast")()
	}
	drain("started", started)
	drain("done", done)
}
