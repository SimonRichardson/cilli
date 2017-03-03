package parselets

import (
	"errors"
	"strconv"

	"github.com/SimonRichardson/cilli/expressions"
	s "github.com/SimonRichardson/cilli/selectors"
)

var (
	ErrInvalidBoolean = errors.New("Invalid Boolean")
	ErrInvalidNumber  = errors.New("Invalid Number")
	ErrInvalidName    = errors.New("Invalid Name")
	ErrUnexpectedNull = errors.New("Unexpected Null")
)

type pathBoolean struct{}

func MakePathBoolean() s.PathPrefixParselet {
	return pathBoolean{}
}

func (p pathBoolean) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	res, err := strconv.ParseBool(token.Val())
	if err != nil {
		return nil, ErrInvalidBoolean
	}
	return expressions.MakePathBoolean(res), nil
}

type keyword struct {
	name   string
	parser func() s.PathPrefixParselet
}

type keywords []keyword

func (k keywords) Find(val string) (keyword, bool) {
	for _, v := range k {
		if v.name == val {
			return v, true
		}
	}
	return keyword{}, false
}

type pathName struct {
	keywords keywords
}

func MakePathName() s.PathPrefixParselet {
	return pathName{
		keywords: []keyword{
			keyword{"true", MakePathBoolean},
			keyword{"false", MakePathBoolean},
			keyword{"null", MakePathNull},
		},
	}
}

func (p pathName) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	if res, ok := p.keywords.Find(token.Val()); ok {
		return res.parser().Parse(parser, token)
	}
	return expressions.MakePathName(token.Val()), nil
}

type pathNull struct{}

func MakePathNull() s.PathPrefixParselet {
	return pathNull{}
}

func (p pathNull) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	return nil, ErrUnexpectedNull
}

type pathNumber struct{}

func MakePathNumber() s.PathPrefixParselet {
	return pathNumber{}
}

func (p pathNumber) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	res, err := strconv.ParseFloat(token.Val(), 64)
	if err != nil {
		return nil, ErrInvalidNumber
	}
	return expressions.MakePathNumber(res), nil
}

type pathString struct{}

func MakePathString() s.PathPrefixParselet {
	return pathString{}
}

func (p pathString) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	return expressions.MakePathString(token.Val()), nil
}

type pathWildcard struct{}

func MakePathWildcard() s.PathPrefixParselet {
	return pathWildcard{}
}

func (p pathWildcard) Parse(parser s.PathParser, token s.PathToken) (s.PathExpression, error) {
	return expressions.MakePathWildcard(), nil
}
