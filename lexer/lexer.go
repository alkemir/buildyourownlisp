package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const eof = rune(26)

type lexer struct { // parse a string and send tokens through a channel
	name  string       // name of the lexer, for debugging purposes
	input string       // input to be lexed
	start int          // the position inside the input at which the current item starts
	pos   int          // the position inside the input at which the lexer is
	width int          // the width in bytes of the last decoded rune
	items chan LexItem // channel to send tokens
}

// Return a new lexer that will parse the input string
func Lex(name, input string) (*lexer, chan LexItem) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan LexItem),
	}

	go l.run()

	return l, l.items
}

// Start parsing the lexer
func (l *lexer) run() {
	for state := lexUnknown; state != nil; {
		state = state(l)
	}

	close(l.items) // no more tokens will be emitted
}

func (l *lexer) emit(t itemType) {
	l.items <- LexItem{typ: t, val: l.input[l.start:l.pos]}
	l.start = l.pos
}

// Return the next rune and consume it
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// Skip over the last item
func (l *lexer) ignore() {
	l.start = l.pos
}

// Unread the last rune
func (l *lexer) backup() {
	l.pos -= l.width
}

// Return the next rune but don't consume it
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// Consume the next rune if it is part of the valid set
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}

	l.backup()
	return false
}

// Consume a run of runes part of the valid set
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}

	l.backup()
}

// Notify of an error during lexing
func (l *lexer) errorf(format string, args ...any) stateFn {
	l.items <- LexItem{typ: itemError, val: fmt.Sprintf(format, args...)}
	return nil
}
