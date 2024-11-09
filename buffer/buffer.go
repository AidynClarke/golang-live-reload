package buffer

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

type Buffer struct {
	timer *time.Timer
	lastEvent fsnotify.Event

	trigger func(event fsnotify.Event)
}

func NewBuffer(trigger func(event fsnotify.Event)) *Buffer {
	timer := time.NewTimer(time.Millisecond)
	<-timer.C
	
	buffer := &Buffer{timer: timer, trigger: trigger}

	go buffer.startEventLoop()

	return buffer
}

func (b *Buffer) NewEvent(event fsnotify.Event) {
	b.lastEvent = event
	b.timer.Reset(time.Millisecond * 10)
}

func (b *Buffer) startEventLoop() {
	for range b.timer.C {
		b.trigger(b.lastEvent)
	}
}

