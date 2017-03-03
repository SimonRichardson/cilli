package parselets

import (
	"github.com/SimonRichardson/cilli/expressions"
	s "github.com/SimonRichardson/cilli/selectors"
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
