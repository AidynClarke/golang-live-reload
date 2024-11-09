package buffer

import (
	"fmt"
	"sync"
	"time"
)

type Colour string

const (
	RESET Colour = "\033[0m"
	RED   Colour = "\033[31m"
	GREEN Colour = "\033[32m"
	YELLOW Colour = "\033[33m"
	BLUE Colour = "\033[34m"
	PURPLE Colour = "\033[35m"
	CYAN Colour = "\033[36m"
)

func Log(colour Colour, message string) {
	// fmt.Printf("%s%s%s", string(colour), message, RESET)
}

func Logln(colour Colour, message string) {
	// Log(colour, message + "\n")
}

type CmdBuffer struct {
	duration time.Duration
	fn func()

	lock *sync.Mutex
	isActive bool
	rerun bool
}

func NewCmdBuffer(duration time.Duration, fn func()) *CmdBuffer {
	return &CmdBuffer{
		duration: duration,
		fn: fn,
		lock: &sync.Mutex{},
		isActive: false,
		rerun: false,
	}
}

func (b *CmdBuffer) Call() {
	b.lock.Lock()
	Logln(BLUE, "Call()")
	Logln(YELLOW, fmt.Sprintf("b.isActive: %t", b.isActive))
	Logln(YELLOW, fmt.Sprintf("b.rerun: %t", b.rerun))

	if b.isActive {
		Logln(RED, "Already active")
		b.rerun = true
		b.lock.Unlock()
		return
	}

	Logln(GREEN, "Calling...")

	b.isActive = true
	b.lock.Unlock()
	b.call()
}

func (b *CmdBuffer) call() {
	b.fn()

	time.AfterFunc(b.duration, func() {
		Logln(GREEN, "Finished")
		b.lock.Lock()

		Logln(PURPLE, fmt.Sprintf("b.rerun: %t", b.rerun))

		if b.rerun {
			b.rerun = false
			b.lock.Unlock()
			b.call()
			return
		} 

		b.isActive = false

		Logln(GREEN, "Inactive")
		b.lock.Unlock()
	})
}
