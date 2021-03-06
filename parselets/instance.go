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

type pathGroup struct{}

func MakePathGroup() s.PathPrefixParselet {
	return pathGroup{}
}

func (p pathGroup) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	exprs := make([]s.PathExpression, 0)

	for !parser.Match(s.PTTRightParen) {
		expr, err := parser.ParseExpression()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}

	return expressions.MakePathGroup(exprs), nil
}

type pathAttribute struct{}

func MakePathAttribute() s.PathPrefixParselet {
	return pathAttribute{}
}

func (p pathAttribute) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	return expressions.MakePathAttribute(), nil
}

type pathInfixAttribute struct{}

func MakePathInfixAttribute() s.PathInfixParselet {
	return pathInfixAttribute{}
}

func (p pathInfixAttribute) Parse(parser s.PathParser, expr s.PathExpression, token s.PathToken) (s.PathExpression, error) {
	right, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}

	return expressions.MakePathInfixAttribute(expr, right), nil
}

func (p pathInfixAttribute) Precedence() s.PathPrecedence {
	return s.PPPostfix
}
