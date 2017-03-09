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

type PathPredicate struct {
	Equality   func(s.Element, string, interface{}) bool
	Inequality func(s.Element, string, interface{}) bool
}

type Path struct {
	expression s.PathExpression
	predicate  PathPredicate
}

func NewPath(expression s.PathExpression) *Path {
	return &Path{
		expression: expression,
	}
}

func (p *Path) With(predicate PathPredicate) *Path {
	p.predicate = predicate
	return p
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
		case s.PETNameDescendants, s.PETInstance, s.PETbranch:
			if x, ok := left(expression); ok {
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
				case s.PETGroup:
					if n, _, ok := group(p.predicate, x, nodes); ok {
						nodes = n
						break
					}
					return nil, ErrUnexpectedExpression
				default:
					return nil, ErrUnexpectedExpression
				}

				if y, ok := right(expression); ok {
					switch y.Type() {
					case s.PETName:
						expression = expressions.MakePathDescendants(s.PDTContext, y)
					case s.PETNameDescendants, s.PETInstance:
						nodes = getContextChildren(nodes)
						expression = y
					case s.PETWildcard, s.PETbranch:
						expression = y
					case s.PETIndexAccess:
						expression = expressions.MakePathDescendants(
							s.PDTContext,
							expressions.MakePathNameDescendants(y, expressions.MakePathWildcard()),
						)
					case s.PETGroup:
						if n, expr, ok := group(p.predicate, y, nodes); ok {
							nodes = n
							expression = expr
							continue loop
						}
						return nil, ErrUnexpectedExpression
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

func list(expression s.PathExpression) ([]s.PathExpression, bool) {
	if expr, ok := expression.(s.List); ok {
		return expr.List(), true
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

func group(predicates PathPredicate, expr s.PathExpression, nodes []s.Element) ([]s.Element, s.PathExpression, bool) {
	if exprs, ok := list(expr); ok {
	loop:
		for k, v := range exprs {
			switch v.Type() {
			case s.PETAttribute:
				if next, ok := peek(exprs, k+1); ok && validAttribute(next) {
					continue loop
				}
				return nil, nil, false
			case s.PETEquality:
				if x, ok := left(v); ok {
					if y, ok := right(v); ok {
						nodes = filterByPredicate(predicates.Equality, x, y, nodes)
						continue loop
					}
				}
			default:
				return nil, nil, false
			}
		}

		return nodes, expressions.MakePathWildcard(), true
	}

	return nil, nil, false
}

func peek(exprs []s.PathExpression, pos int) (s.PathExpression, bool) {
	if num := len(exprs); pos >= 0 && pos < num {
		return exprs[pos], true
	}
	return nil, false
}

func validAttribute(expr s.PathExpression) bool {
	// A valid attribute should always have a left hand side of name.
	switch expr.Type() {
	case s.PETEquality:
		if x, ok := left(expr); ok && x.Type() == s.PETName {
			return true
		}
	}
	return false
}

func filterByName(expression s.PathExpression, nodes []s.Element) []s.Element {
	var res []s.Element

	if expr, ok := expression.(s.Name); ok {
		name := expr.Name()

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

		if num := len(nodes); index >= 0 && index < num {
			res = append(res, nodes[index])
		}
	}

	return res
}

func filterByPredicate(predicate func(s.Element, string, interface{}) bool,
	left, right s.PathExpression,
	nodes []s.Element,
) []s.Element {
	var res []s.Element

	if x, ok := left.(s.Name); ok {
		prop := x.Name()

		if y, ok := right.(s.Value); ok {
			value := y.Value()

			for _, v := range nodes {
				if predicate(v, prop, value) {
					res = append(res, v)
				}
			}
		}
	}

	return res
}
