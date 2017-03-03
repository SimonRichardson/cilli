package selectors

type Descendants interface {
	Descendants() PathExpression
}

type Branch interface {
	Left() PathExpression
	Right() PathExpression
}
