package parselets

import (
	"errors"

	"github.com/SimonRichardson/cilli/expressions"
	s "github.com/SimonRichardson/cilli/selectors"
)

var (
	ErrInvalidIndexAccess = errors.New("Invalid Index Access")
)

type pathDescendants struct{}

func MakePathDescendants() s.PathPrefixParselet {
	return pathDescendants{}
}

func (p pathDescendants) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	context := s.PDTContext

	if parser.Match(s.PTTForwardSlash) {
		context = s.PDTAll
	}

	right, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}

	return expressions.MakePathDescendants(context, right), nil
}

type pathNameDescendants struct{}

func MakePathNameDescendants() s.PathInfixParselet {
	return pathNameDescendants{}
}

func (p pathNameDescendants) Parse(parser s.PathParser, expr s.PathExpression, token s.PathToken) (s.PathExpression, error) {
	right, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}

	return expressions.MakePathNameDescendants(expr, right), nil
}

func (p pathNameDescendants) Precedence() s.PathPrecedence {
	return s.PPPostfix
}

type pathIndexAccess struct{}

func MakePathIndexAccess() s.PathInfixParselet {
	return pathIndexAccess{}
}

func (p pathIndexAccess) Parse(parser s.PathParser, expr s.PathExpression, token s.PathToken) (s.PathExpression, error) {
	if !parser.Match(s.PTTRightSquare) {
		param, err := parser.ParseExpression()
		if err != nil {
			return nil, err
		}

		if _, err := parser.ConsumeToken(s.PTTRightSquare); err != nil {
			return nil, err
		}

		return expressions.MakePathIndexAccess(expr, param), nil
	}
	return nil, ErrInvalidIndexAccess
}

func (p pathIndexAccess) Precedence() s.PathPrecedence {
	return s.PPCall
}
