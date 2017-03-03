package selectors

type PathParser interface {
	ParseExpression() (PathExpression, error)
	ParseExpressionBy(PathPrecedence) (PathExpression, error)

	Match(PathTokenType) bool
	Consume() (PathToken, error)
	ConsumeToken(PathTokenType) (PathToken, error)
}
