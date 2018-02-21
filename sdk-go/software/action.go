package software

import common "github.com/richardlt/matrix/sdk-go/common"

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
