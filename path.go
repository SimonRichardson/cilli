package cilli

import (
	"bufio"
	"errors"

	"github.com/SimonRichardson/cilli/expressions"
	s "github.com/SimonRichardson/cilli/selectors"
)

var (
	ErrUnexpectedExpression = errors.New("Unexpected Expression")
)

type Path struct {
	expression s.PathExpression
}

func NewPath(expression s.PathExpression) *Path {
	return &Path{
		expression: expression,
	}
}

func (p *Path) Describe(w *bufio.Writer) error {
	if x, ok := p.expression.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}
	return nil
}

func (p *Path) Execute(element s.Element) ([]s.Element, error) {
	var (
		res []s.Element

		nodes      = []s.Element{element}
		expression = p.expression
	)

	// Add context to shortcuts
	switch expression.Type() {
	case s.PETWildcard:
		expression = expressions.MakePathDescendants(s.PDTAll, expression)
	}

	// Loop through everything.
loop:
	for {
		// fmt.Println(">", expression.Type(), expression)

		switch expression.Type() {
		case s.PETWildcard:
			res = nodes
			break loop
		case s.PETAllDescendants:
			nodes = getAllChildren(nodes)
			if expr, ok := descendants(expression); ok {
				expression = expr
				continue loop
			}
			return nil, ErrUnexpectedExpression
		case s.PETDescendants:
			nodes = getContextChildren(nodes)
			if expr, ok := descendants(expression); ok {
				expression = expr
				continue loop
			}
			return nil, ErrUnexpectedExpression
		case s.PETNameDescendants:
			if x, ok := left(expression); ok {
				// fmt.Println("Left", x.Type(), x, len(nodes))
				switch x.Type() {
				case s.PETName:
					nodes = filterByName(x, nodes)
				case s.PETIndexAccess:
					if y, ok := left(x); ok {
						nodes = filterByName(y, nodes)

						if z, ok := right(x); ok {
							nodes = filterByIndex(z, nodes)
							break
						}
					}
					return nil, ErrUnexpectedExpression
				default:
					return nil, ErrUnexpectedExpression
				}

				if y, ok := right(expression); ok {
					// fmt.Println("Right", y.Type(), y, len(nodes))
					switch y.Type() {
					case s.PETName:
						expression = expressions.MakePathDescendants(s.PDTContext, y)
					case s.PETNameDescendants:
						nodes = getContextChildren(nodes)
						expression = y
					case s.PETWildcard:
						expression = y
					case s.PETIndexAccess:
						expression = expressions.MakePathDescendants(
							s.PDTContext,
							expressions.MakePathNameDescendants(y, expressions.MakePathWildcard()),
						)
					default:
						return nil, ErrUnexpectedExpression
					}
					continue loop
				}
			}
			return nil, ErrUnexpectedExpression
		case s.PETName:
			res = filterByName(expression, nodes)
			break loop
		default:
			return nil, ErrUnexpectedExpression
		}
	}
	return res, nil
}

func getAllChildren(nodes []s.Element) []s.Element {
	res := make([]s.Element, 0)
	for _, v := range nodes {
		children := v.Children()
		res = append(append(res, children...), getAllChildren(children)...)

	}
	return res
}

func getContextChildren(nodes []s.Element) []s.Element {
	res := make([]s.Element, 0)
	for _, v := range nodes {
		res = append(res, v.Children()...)
	}
	return res
}

func descendants(expression s.PathExpression) (s.PathExpression, bool) {
	if expr, ok := expression.(s.Descendants); ok {
		expression = expr.Descendants()
		return expression, true
	}
	return nil, false
}

func left(expression s.PathExpression) (s.PathExpression, bool) {
	if expr, ok := expression.(s.Branch); ok {
		expression = expr.Left()
		return expression, true
	}
	return nil, false
}

func right(expression s.PathExpression) (s.PathExpression, bool) {
	if expr, ok := expression.(s.Branch); ok {
		expression = expr.Right()
		return expression, true
	}
	return nil, false
}

func filterByName(expression s.PathExpression, nodes []s.Element) []s.Element {
	var res []s.Element

	if expr, ok := expression.(s.Name); ok {
		name := expr.Name()

		// fmt.Println("Name :", name, len(nodes))

		for _, v := range nodes {
			if v.Name() == name {
				res = append(res, v)
			}
		}
	}
	return res
}

func filterByIndex(expression s.PathExpression, nodes []s.Element) []s.Element {
	var res []s.Element

	if expr, ok := expression.(s.Index); ok {
		index := expr.Index()

		// fmt.Println("Index :", index, len(nodes))

		if num := len(nodes); index >= 0 && index < num {
			res = append(res, nodes[index])
		}
	}

	return res
}
