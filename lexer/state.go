package lexer

import (
	"unicode/utf8"
)

const (
	leftParenthesis  = '('
	rightParenthesis = ')'
)

const operatorLen = 1

type stateFn func(*lexer) stateFn

func lexUnknown(l *lexer) stateFn {
	// omit whitespace
	l.acceptRun(" \r\n")
	l.ignore()

	switch n := l.peek(); {
	case n == leftParenthesis:
		return lexLeftParenthesis
	case n == rightParenthesis:
		return lexRightParenthesis
	case n == '/' || n == '*' || n == '%':
		return lexOperator
	case '0' <= n && n <= '9':
		return lexInteger
	case n == '+' || n == '-':
		nLen := utf8.RuneLen(n)
		if r, _ := utf8.DecodeRuneInString(l.input[l.pos+nLen:]); r == ' ' {
			return lexOperator
		}
		return lexInteger
	case n == eof:
		l.emit(itemEOF)
		return nil
	default:
		return l.errorf("unexpected rune: %q", n)
	}
}

func lexOperator(l *lexer) stateFn {
	l.pos += operatorLen
	l.emit(itemOperator)
	return lexUnknown
}

func lexLeftParenthesis(l *lexer) stateFn {
	l.pos += utf8.RuneLen(leftParenthesis)
	l.emit(itemLeftParenthesis)
	return lexUnknown
}

func lexRightParenthesis(l *lexer) stateFn {
	l.pos += utf8.RuneLen(rightParenthesis)
	l.emit(itemRightParenthesis)
	return lexUnknown
}

func lexInteger(l *lexer) stateFn {
	l.accept("+-")
	l.acceptRun("0123456789")
	if l.start == l.pos {
		return l.errorf("expected number but got empty string")
	}
	l.emit(itemNumber)
	return lexUnknown
}
