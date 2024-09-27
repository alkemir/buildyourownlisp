package lexer

import "fmt"

type itemType int

const (
	itemError itemType = iota

	itemLeftParenthesis
	itemRightParenthesis
	itemOperator
	itemNumber
	itemEOF
)

type LexItem struct {
	typ itemType
	val string
}

func (i LexItem) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%v %.10q", i.typ, i.val)
	}

	return fmt.Sprintf("%v %q", i.typ, i.val)
}
