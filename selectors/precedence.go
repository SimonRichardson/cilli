package selectors

type PathPrecedence int

const (
	PPEquality PathPrecedence = iota
	PPConditional
	PPSum
	PPProduct
	PPExponent
	PPPrefix
	PPPostfix
	PPCall
)

func (p PathPrecedence) Val() int {
	return int(p)
}
