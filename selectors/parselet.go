package selectors

type PathPrefixParselet interface {
	Parse(PathParser, PathToken) (PathExpression, error)
}

type PathInfixParselet interface {
	Parse(PathParser, PathExpression, PathToken) (PathExpression, error)
	Precedence() PathPrecedence
}
