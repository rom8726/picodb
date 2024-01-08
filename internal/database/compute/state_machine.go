package compute

import (
	"strings"
)

const (
	foundLetterEvent = iota
	foundWhiteSpaceEvent
	// must be last
	eventsNumber
)

const (
	initialState = iota
	wordState
	whiteSpaceState
	invalidState
	// must be last
	statesNumber
)

type transition struct {
	jump   func(byte) int
	action func()
}

type stateMachine struct {
	transitions [statesNumber][eventsNumber]transition
	state       int

	tokens []string
	sb     strings.Builder
}

func newStateMachine() *stateMachine {
	machine := &stateMachine{
		state: initialState,
	}

	machine.transitions = [statesNumber][eventsNumber]transition{
		initialState: {
			foundLetterEvent:     transition{jump: machine.appendLetterJump},
			foundWhiteSpaceEvent: transition{jump: machine.skipWhiteSpaceJump},
		},
		wordState: {
			foundLetterEvent:     transition{jump: machine.appendLetterJump},
			foundWhiteSpaceEvent: transition{jump: machine.skipWhiteSpaceJump, action: machine.addTokenAction},
		},
		whiteSpaceState: {
			foundLetterEvent:     transition{jump: machine.appendLetterJump},
			foundWhiteSpaceEvent: transition{jump: machine.skipWhiteSpaceJump},
		},
		invalidState: {},
	}

	return machine
}

func (sm *stateMachine) parse(query string) ([]string, error) {
	for i := 0; i < len(query); i++ {
		symbol := query[i]
		switch {
		case isWhiteSpace(symbol):
			sm.processEvent(foundWhiteSpaceEvent, symbol)
		case isLetter(symbol):
			sm.processEvent(foundLetterEvent, symbol)
		default:
			return nil, ErrInvalidSymbol
		}
	}

	sm.processEvent(foundWhiteSpaceEvent, ' ')

	return sm.tokens, nil
}

func (sm *stateMachine) processEvent(event int, symbol byte) {
	ts := sm.transitions[sm.state][event]
	sm.state = ts.jump(symbol)
	if ts.action != nil {
		ts.action()
	}
}

func (sm *stateMachine) appendLetterJump(letter byte) int {
	sm.sb.WriteByte(letter)

	return wordState
}

func (sm *stateMachine) skipWhiteSpaceJump(byte) int {
	return whiteSpaceState
}

func (sm *stateMachine) addTokenAction() {
	sm.tokens = append(sm.tokens, sm.sb.String())
	sm.sb.Reset()
}
