package selectors

import "fmt"

type PathDescendantsType int

const (
	PDTAll PathDescendantsType = iota
	PDTContext
)

type PathTokenType int

const (
	PTTNull PathTokenType = iota
	PTTNumber
	PTTInteger
	PTTString
	PTTLeftParen
	PTTRightParen
	PTTComma
	PTTEquality
	PTTPlus
	PTTMinus
	PTTAsterisk
	PTTBackSlash
	PTTForwardSlash
	PTTCaret
	PTTTilde
	PTTBang
	PTTQuestionMark
	PTTColon
	PTTLeftSquare
	PTTRightSequre
	PTTAttribute
	PTTName
	PTTUnsignedInteger
	PTTDot
	PTTAmpersand
	PTTPipe
	PTTForwardArrow
	PTTBackArrow
)

func (p PathTokenType) Rune() rune {
	switch p {
	case PTTLeftParen:
		return '('
	case PTTRightParen:
		return ')'
	case PTTComma:
		return ','
	case PTTEquality:
		return '='
	case PTTPlus:
		return '+'
	case PTTMinus:
		return '-'
	case PTTAsterisk:
		return '*'
	case PTTBackSlash:
		return '\\'
	case PTTForwardSlash:
		return '/'
	case PTTCaret:
		return '^'
	case PTTTilde:
		return '~'
	case PTTBang:
		return '!'
	case PTTQuestionMark:
		return '?'
	case PTTColon:
		return ':'
	case PTTLeftSquare:
		return '['
	case PTTRightSequre:
		return ']'
	case PTTAttribute:
		return '@'
	case PTTDot:
		return '.'
	case PTTAmpersand:
		return '&'
	case PTTPipe:
		return '|'
	case PTTForwardArrow:
		return '>'
	case PTTBackArrow:
		return '<'
	}
	panic("Invalid rune type")
}

func (p PathTokenType) String() string {
	switch p {
	case PTTNumber:
		return "Number"
	case PTTInteger:
		return "Int"
	case PTTString:
		return "String"
	case PTTLeftParen:
		return "("
	case PTTRightParen:
		return ")"
	case PTTComma:
		return ","
	case PTTEquality:
		return "="
	case PTTPlus:
		return "+"
	case PTTMinus:
		return "-"
	case PTTAsterisk:
		return "*"
	case PTTBackSlash:
		return "\\"
	case PTTForwardSlash:
		return "/"
	case PTTCaret:
		return "^"
	case PTTTilde:
		return "~"
	case PTTBang:
		return "!"
	case PTTQuestionMark:
		return "?"
	case PTTColon:
		return ":"
	case PTTLeftSquare:
		return "["
	case PTTRightSequre:
		return "]"
	case PTTAttribute:
		return "@"
	case PTTName:
		return "Name"
	case PTTUnsignedInteger:
		return "Uint"
	case PTTDot:
		return "."
	case PTTAmpersand:
		return "&"
	case PTTPipe:
		return "|"
	case PTTForwardArrow:
		return ">"
	case PTTBackArrow:
		return "<"
	}
	return ""
}

func PathTokenTypes() []PathTokenType {
	return []PathTokenType{
		PTTLeftParen,
		PTTRightParen,
		PTTComma,
		PTTEquality,
		PTTPlus,
		PTTMinus,
		PTTAsterisk,
		PTTBackSlash,
		PTTForwardSlash,
		PTTCaret,
		PTTTilde,
		PTTBang,
		PTTQuestionMark,
		PTTColon,
		PTTLeftSquare,
		PTTRightSequre,
		PTTAttribute,
		PTTDot,
		PTTAmpersand,
		PTTPipe,
		PTTForwardArrow,
		PTTBackArrow,
	}
}

type PathToken struct {
	tokenType PathTokenType
	val       string
}

func MakePathToken(tokenType PathTokenType, val string) PathToken {
	return PathToken{
		tokenType: tokenType,
		val:       val,
	}
}

func (p PathToken) Type() PathTokenType {
	return p.tokenType
}

func (p PathToken) Val() string {
	return p.val
}

func (p PathToken) String() string {
	return fmt.Sprintf("(TokenType:%s, Value:%q)", p.tokenType.String(), p.val)
}
