package event

import (
	"sync"
	"time"
)

type EventType int

const (
	EventMouseMove EventType = iota
	EventMouseClick
	EventMouseRelease
	EventMouseWheel
	EventKeyPress
	EventKeyRelease
	EventChar
	EventPaint
	EventResize
	EventClose
)

type MouseEvent struct {
	X, Y   int32
	Button int // 1: Left, 2: Right, 3: Middle
	Delta  int // For mouse wheel
}

type KeyEvent struct {
	VirtualKeyCode uint32
	Rune           rune
	Modifiers      uint32 // Bitmask: 1=Shift, 2=Ctrl, 4=Alt
}

const (
	ModShift = 1
	ModCtrl  = 2
	ModAlt   = 4
)

type Event struct {
	Type      EventType
	Timestamp int64
	Data      interface{}
}

type EventBus interface {
	Subscribe(eventType EventType) <-chan Event
	Publish(event Event)
	Close()
}

type Bus struct {
	subscribers map[EventType][]chan Event
	mu          sync.RWMutex
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[EventType][]chan Event),
	}
}

func (b *Bus) Subscribe(eventType EventType) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, 10)
	b.subscribers[eventType] = append(b.subscribers[eventType], ch)
	return ch
}

func (b *Bus) Publish(event Event) {
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().UnixNano()
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	if chans, ok := b.subscribers[event.Type]; ok {
		for _, ch := range chans {
			select {
			case ch <- event:
			default:
				// Buffer full, drop event
			}
		}
	}
}

func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, chs := range b.subscribers {
		for _, ch := range chs {
			close(ch)
		}
	}
	b.subscribers = nil
}
