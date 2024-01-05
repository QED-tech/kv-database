package compute

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidChar = errors.New("[database] invalid character for parse")

type parseState int

const (
	stateWordsString parseState = iota
	stateWordByWord
)

type event int

const (
	eventSpace event = iota
	eventLetter
	eventEndOfString
)

type parseStateMachine struct {
	state        parseState
	buffer       []string
	currentToken strings.Builder
}

func newParseStateMachine() *parseStateMachine {
	return &parseStateMachine{
		state: stateWordByWord,
	}
}

func (m *parseStateMachine) parse(query string) ([]string, error) {
	for _, char := range query {
		switch {
		case char == '\n':
			m.handleEvent(eventEndOfString, ' ')
		case unicode.IsSpace(char):
			m.handleEvent(eventSpace, char)
		case unicode.IsLetter(char) || char == '_':
			m.handleEvent(eventLetter, char)
		case char == '"' || char == '\'':
			m.state = stateWordsString
		default:
			return nil, ErrInvalidChar
		}
	}

	m.handleEvent(eventEndOfString, ' ')

	return m.buffer, nil
}

func (m *parseStateMachine) handleEvent(event event, char rune) {
	switch m.state {
	case stateWordByWord:
		switch event {
		case eventLetter:
			m.currentToken.WriteRune(char)
		case eventSpace:
			m.addCurrentToken()
		case eventEndOfString:
			m.addCurrentToken()
		}
	case stateWordsString:
		switch event {
		case eventLetter:
			m.currentToken.WriteRune(char)
		case eventSpace:
			m.currentToken.WriteRune(char)
		case eventEndOfString:
			m.addCurrentToken()
		}
	}
}

func (m *parseStateMachine) addCurrentToken() {
	if m.currentToken.Len() == 0 {
		return
	}

	m.buffer = append(m.buffer, m.currentToken.String())
	m.currentToken.Reset()
}
