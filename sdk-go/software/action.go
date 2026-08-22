package software

import (
	"sync"
	"time"

	common "github.com/richardlt/matrix/sdk-go/common"
)

// ActionGenerator takes SDK actions to generate custom one.
type ActionGenerator interface {
	SendAction(slot uint64, cmd common.Command)
	OnAction(func(slot uint64))
	// Reset forgets which buttons are down. A generator only learns that a button came
	// back up from the command saying so, and a software that stops receiving commands
	// -- while the core holds a pause, for instance -- never sees it.
	Reset()
}

type passThrough struct {
	actionCallback func(uint64)
}

func (p *passThrough) SendAction(slot uint64, cmd common.Command) {
	if p.actionCallback != nil {
		p.actionCallback(slot)
	}
}

func (p *passThrough) OnAction(f func(uint64)) { p.actionCallback = f }

func (p *passThrough) Reset() {}

// NewMultiPress returns a multi press generator.
func NewMultiPress(cs ...common.Button) ActionGenerator {
	m := map[common.Button]common.Button{}
	for _, c := range cs {
		m[c] = c
	}
	return &mutliPress{buttons: m, slots: map[uint64](map[common.Button]bool){}}
}

type mutliPress struct {
	passThrough
	buttons map[common.Button]common.Button
	lock    sync.Mutex
	slots   map[uint64](map[common.Button]bool)
}

func (m *mutliPress) SendAction(slot uint64, cmd common.Command) {
	button, pressed := common.CommandToButtonState(cmd)

	m.lock.Lock()

	if _, ok := m.slots[slot]; !ok {
		m.slots[slot] = map[common.Button]bool{}
	}

	if _, ok := m.buttons[button]; ok {
		m.slots[slot][button] = pressed
	}

	allPressed := len(m.slots[slot]) == len(m.buttons)
	if allPressed {
		for _, p := range m.slots[slot] {
			if !p {
				allPressed = false
				break
			}
		}
	}

	m.lock.Unlock()

	if allPressed && m.actionCallback != nil {
		m.actionCallback(slot)
	}
}

func (m *mutliPress) Reset() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.slots = map[uint64](map[common.Button]bool){}
}

// NewLongPress returns a long press generator.
func NewLongPress(button common.Button, triggerDelay,
	fireDelay time.Duration) ActionGenerator {
	return &longPress{
		button:       button,
		triggerDelay: triggerDelay,
		fireDelay:    fireDelay,
		holds:        map[uint64]chan struct{}{},
	}
}

type longPress struct {
	passThrough
	button                  common.Button
	triggerDelay, fireDelay time.Duration
	lock                    sync.Mutex
	// holds carries one channel per slot with the button down, closed to end the repeat.
	// The timer and the ticker are waited on together with it, so releasing the button
	// leaves nothing behind.
	holds map[uint64]chan struct{}
}

func (l *longPress) SendAction(slot uint64, cmd common.Command) {
	button, pressed := common.CommandToButtonState(cmd)
	if button != l.button {
		return
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	hold, held := l.holds[slot]

	if !pressed {
		if held {
			close(hold)
			delete(l.holds, slot)
		}
		return
	}

	if held {
		return
	}

	hold = make(chan struct{})
	l.holds[slot] = hold
	go l.repeat(slot, hold)
}

// repeat fires once the button has been down for the trigger delay, then at the fire
// delay until it comes back up.
func (l *longPress) repeat(slot uint64, hold chan struct{}) {
	timer := time.NewTimer(l.triggerDelay)
	defer timer.Stop()

	select {
	case <-hold:
		return
	case <-timer.C:
	}

	ticker := time.NewTicker(l.fireDelay)
	defer ticker.Stop()

	for {
		// The end of the hold and the next tick can become ready together, and select
		// picks either one at random: look at the hold on its own so that nothing fires
		// once the button is up or the generator has been reset.
		select {
		case <-hold:
			return
		default:
		}

		l.action(slot)

		select {
		case <-hold:
			return
		case <-ticker.C:
		}
	}
}

func (l *longPress) Reset() {
	l.lock.Lock()
	defer l.lock.Unlock()

	for slot, hold := range l.holds {
		close(hold)
		delete(l.holds, slot)
	}
}

func (l *longPress) action(slot uint64) {
	if l.actionCallback != nil {
		go l.actionCallback(slot)
	}
}
