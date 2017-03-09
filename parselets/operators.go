package parselets

import (
	"errors"

	"github.com/SimonRichardson/cilli/expressions"
	s "github.com/SimonRichardson/cilli/selectors"
)

var (
	ErrInvalidEqualityProperty = errors.New("Invalid Equality Property")
)

type pathEquality struct{}

func MakePathEquality() s.PathInfixParselet {
	return pathEquality{}
}

func (p pathEquality) Parse(parser s.PathParser, expr s.PathExpression, token s.PathToken) (s.PathExpression, error) {
	if _, err := parser.ConsumeToken(s.PTTEquality); err != nil {
		return nil, err
	}

	if expr.Type() != s.PETName {
		return nil, ErrInvalidEqualityProperty
	}

	right, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}

	return expressions.MakePathEquality(expr, right), nil
}

func (p pathEquality) Precedence() s.PathPrecedence {
	return s.PPPostfix
}
