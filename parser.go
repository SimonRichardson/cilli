package cilli

import (
	"errors"

	"github.com/SimonRichardson/cilli/parselets"
	s "github.com/SimonRichardson/cilli/selectors"
)

var (
	ErrParsePrefixError = errors.New("Parse Prefix Error")
	ErrParseInfixError  = errors.New("Parse Infix Error")
	ErrBufferUnderflow  = errors.New("Buffer Underflow")
	ErrBufferOverflow   = errors.New("Buffer Overflow")
	ErrUnexpectedToken  = errors.New("Unexpected Token")
)

type pathParser struct {
	tokens s.PathLexerIterator
	prefix map[s.PathTokenType]s.PathPrefixParselet
	infix  map[s.PathTokenType]s.PathInfixParselet
	stream []s.PathToken
}

func NewPathParser(iter s.PathLexerIterator) s.PathParser {
	return &pathParser{
		tokens: iter,
		prefix: map[s.PathTokenType]s.PathPrefixParselet{
			s.PTTName:         parselets.MakePathName(),
			s.PTTNumber:       parselets.MakePathNumber(),
			s.PTTString:       parselets.MakePathString(),
			s.PTTAsterisk:     parselets.MakePathWildcard(),
			s.PTTForwardSlash: parselets.MakePathDescendants(),
			s.PTTLeftParen:    parselets.MakePathGroup(),
		},
		infix: map[s.PathTokenType]s.PathInfixParselet{
			s.PTTDot:          parselets.MakePathInstance(),
			s.PTTForwardSlash: parselets.MakePathNameDescendants(),
			s.PTTLeftSquare:   parselets.MakePathIndexAccess(),
		},
		stream: []s.PathToken{},
	}
}

func (p *pathParser) ParseExpression() (s.PathExpression, error) {
	return p.ParseExpressionBy(0)
}

func (p *pathParser) ParseExpressionBy(precedence s.PathPrecedence) (s.PathExpression, error) {
	token, err := p.Consume()
	if err != nil {
		return nil, err
	}

	// fmt.Println("Prefix", token)

	prefix, ok := p.prefix[token.Type()]
	if !ok {
		return nil, ErrParsePrefixError
	}

	expression, err := prefix.Parse(p, token)
	if err != nil {
		return nil, err
	}

	for {
		if next, err := p.nextTokenPrecedence(); err != nil {
			if err == ErrBufferOverflow {
				return expression, nil
			}
			return nil, err
		} else if precedence < next {
			token, err = p.Consume()
			if err != nil {
				return nil, err
			}

			// fmt.Println("Infix", token)

			infix, ok := p.infix[token.Type()]
			if !ok {
				return nil, ErrParseInfixError
			}

			expression, err = infix.Parse(p, expression, token)
			if err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	return expression, nil
}

func (p *pathParser) Match(expected s.PathTokenType) bool {
	token, err := p.advance(0)
	if err != nil {
		return false
	}

	if token.Type() != expected {
		return false
	}

	if _, err := p.Consume(); err != nil {
		return false
	}
	return true
}

func (p *pathParser) Consume() (s.PathToken, error) {
	if _, err := p.advance(0); err != nil {
		return s.PathToken{}, err
	}

	if len(p.stream) < 1 {
		return s.PathToken{}, ErrBufferUnderflow
	}

	res := p.stream[0]
	p.stream = p.stream[1:]
	return res, nil
}

func (p *pathParser) ConsumeToken(expected s.PathTokenType) (s.PathToken, error) {
	token, err := p.advance(0)
	if err != nil {
		return s.PathToken{}, err
	}
	if token.Type() != expected {
		return s.PathToken{}, ErrUnexpectedToken
	}

	return p.Consume()
}

func (p *pathParser) advance(distance int) (s.PathToken, error) {
	for {
		if distance >= len(p.stream) {
			if !p.tokens.HasNext() {
				return s.PathToken{}, ErrBufferOverflow
			}

			next, err := p.tokens.Next()
			if err != nil {
				return s.PathToken{}, err
			}
			p.stream = append(p.stream, next)

			continue
		}

		break
	}

	return p.stream[distance], nil
}

func (p *pathParser) nextTokenPrecedence() (s.PathPrecedence, error) {
	token, err := p.advance(0)
	if err != nil {
		return -1, err
	}

	parselet, ok := p.infix[token.Type()]
	if ok {
		return parselet.Precedence(), nil
	}
	return 0, nil
}
