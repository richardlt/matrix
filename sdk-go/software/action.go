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
	slots   map[uint64](map[common.Button]bool)
}

func (m *mutliPress) SendAction(slot uint64, cmd common.Command) {
	if _, ok := m.slots[slot]; !ok {
		m.slots[slot] = map[common.Button]bool{}
	}

	button, pressed := common.CommandToButtonState(cmd)

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

	if allPressed && m.actionCallback != nil {
		m.actionCallback(slot)
	}
}

// NewLongPress returns a long press generator.
func NewLongPress(button common.Button, triggerDelay,
	fireDelay time.Duration) ActionGenerator {
	return &longPress{
		button:       button,
		triggerDelay: triggerDelay,
		fireDelay:    fireDelay,
		timers:       map[uint64]*time.Timer{},
		tickers:      map[uint64]*time.Ticker{},
	}
}

type longPress struct {
	passThrough
	button                  common.Button
	triggerDelay, fireDelay time.Duration
	timersLock, tickersLock sync.RWMutex
	timers                  map[uint64]*time.Timer
	tickers                 map[uint64]*time.Ticker
}

func (l *longPress) SendAction(slot uint64, cmd common.Command) {
	button, pressed := common.CommandToButtonState(cmd)
	if button != l.button {
		return
	}

	l.timersLock.RLock()
	timer, okTimer := l.timers[slot]
	l.timersLock.RUnlock()

	l.tickersLock.RLock()
	ticker, okTicker := l.tickers[slot]
	l.tickersLock.RUnlock()

	// if the button was just pressed
	if (!okTimer || timer == nil) && (!okTicker || ticker == nil) && pressed {
		newTimer := time.NewTimer(l.triggerDelay)

		l.timersLock.Lock()
		l.timers[slot] = newTimer
		l.timersLock.Unlock()

		go func() {
			<-newTimer.C
			l.action(slot)

			newTicker := time.NewTicker(l.fireDelay)

			l.tickersLock.Lock()
			l.tickers[slot] = newTicker
			l.tickersLock.Unlock()

			for _ = range newTicker.C {
				l.action(slot)
			}
		}()
	} else if !pressed {
		if timer != nil {
			timer.Stop()
			l.timersLock.Lock()
			l.timers[slot] = nil
			l.timersLock.Unlock()
		}
		if ticker != nil {
			ticker.Stop()
			l.tickersLock.Lock()
			l.tickers[slot] = nil
			l.tickersLock.Unlock()
		}
	}
}

func (l *longPress) action(slot uint64) {
	if l.actionCallback != nil {
		go l.actionCallback(slot)
	}
}
