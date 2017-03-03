package parselets

import (
	"github.com/SimonRichardson/cilli/expressions"
	s "github.com/SimonRichardson/cilli/selectors"
)

type pathInstance struct{}

func MakePathInstance() s.PathInfixParselet {
	return pathInstance{}
}

func (p pathInstance) Parse(parser s.PathParser, expr s.PathExpression, token s.PathToken) (s.PathExpression, error) {
	right, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}

	return expressions.MakePathInstance(expr, right), nil
}

func (p pathInstance) Precedence() s.PathPrecedence {
	return s.PPPostfix
}
