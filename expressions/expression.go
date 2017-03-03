package expressions

import (
	"bufio"

	s "github.com/SimonRichardson/cilli/selectors"
)

type PathDescendantsExpression interface {
	Descendants() s.PathExpression
}

type PathLeftRightNodeExpression interface {
	Type() s.PathExpressionType
	Left() s.PathExpression
	Right() s.PathExpression
}

type PathPrefixExpression struct {
	operator s.PathTokenType
	right    s.PathExpression
}

func MakePathPrefix(operator s.PathTokenType, right s.PathExpression) PathPrefixExpression {
	return PathPrefixExpression{
		operator: operator,
		right:    right,
	}
}

func (p PathPrefixExpression) Describe(w *bufio.Writer) error {
	if _, err := w.WriteRune('('); err != nil {
		return err
	}

	if _, err := w.WriteString(p.Type().String()); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	_, err := w.WriteRune(')')
	return err
}

func (p PathPrefixExpression) Type() s.PathExpressionType {
	return p.right.Type()
}

type PathPostfixExpression struct {
	operator s.PathTokenType
	left     s.PathExpression
}

func MakePathPostfix(operator s.PathTokenType, left s.PathExpression) PathPostfixExpression {
	return PathPostfixExpression{
		operator: operator,
		left:     left,
	}
}

func (p PathPostfixExpression) Describe(w *bufio.Writer) error {
	if _, err := w.WriteRune('('); err != nil {
		return err
	}

	if _, err := w.WriteString(p.Type().String()); err != nil {
		return err
	}

	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	_, err := w.WriteRune(')')
	return err
}

func (p PathPostfixExpression) Type() s.PathExpressionType {
	return p.left.Type()
}
