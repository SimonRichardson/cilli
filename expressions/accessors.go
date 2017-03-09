package expressions

import (
	"bufio"
	"fmt"

	s "github.com/SimonRichardson/cilli/selectors"
)

type nameDescendantsType struct {
	left, right s.PathExpression
}

func MakePathNameDescendants(left, right s.PathExpression) s.PathExpression {
	return nameDescendantsType{
		left:  left,
		right: right,
	}
}

func (p nameDescendantsType) Type() s.PathExpressionType {
	return s.PETNameDescendants
}

func (p nameDescendantsType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune('('); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune(')'); err != nil {
		return err
	}

	return nil
}

func (p nameDescendantsType) Left() s.PathExpression {
	return p.left
}

func (p nameDescendantsType) Right() s.PathExpression {
	return p.right
}

type branchType struct {
	left, right s.PathExpression
}

func MakePathBranch(left, right s.PathExpression) s.PathExpression {
	return branchType{
		left:  left,
		right: right,
	}
}

func (p branchType) Type() s.PathExpressionType {
	return s.PETbranch
}

func (p branchType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune('('); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune(')'); err != nil {
		return err
	}

	return nil
}

func (p branchType) Left() s.PathExpression {
	return p.left
}

func (p branchType) Right() s.PathExpression {
	return p.right
}

type instanceType struct {
	left, right s.PathExpression
}

func MakePathInstance(left, right s.PathExpression) s.PathExpression {
	return instanceType{
		left:  left,
		right: right,
	}
}

func (p instanceType) Type() s.PathExpressionType {
	return s.PETInstance
}

func (p instanceType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune('.'); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p instanceType) Left() s.PathExpression {
	return p.left
}

func (p instanceType) Right() s.PathExpression {
	return p.right
}

type indexAccessType struct {
	left, right s.PathExpression
}

func MakePathIndexAccess(left, right s.PathExpression) s.PathExpression {
	return indexAccessType{
		left:  left,
		right: right,
	}
}

func (p indexAccessType) Type() s.PathExpressionType {
	return s.PETIndexAccess
}

func (p indexAccessType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune('['); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune(']'); err != nil {
		return err
	}

	return nil
}

func (p indexAccessType) Left() s.PathExpression {
	return p.left
}

func (p indexAccessType) Right() s.PathExpression {
	return p.right
}

type indexAccessDescendantsType struct {
	left, right s.PathExpression
}

func MakePathIndexAccessDescendants(left, right s.PathExpression) s.PathExpression {
	return indexAccessDescendantsType{
		left:  left,
		right: right,
	}
}

func (p indexAccessDescendantsType) Type() s.PathExpressionType {
	return s.PETIndexAccessDescendants
}

func (p indexAccessDescendantsType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune('['); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune(']'); err != nil {
		return err
	}

	return nil
}

func (p indexAccessDescendantsType) Left() s.PathExpression {
	return p.left
}

func (p indexAccessDescendantsType) Right() s.PathExpression {
	return p.right
}

type descendantsType struct {
	expType     s.PathDescendantsType
	descendants s.PathExpression
}

func MakePathDescendants(expType s.PathDescendantsType, descendants s.PathExpression) s.PathExpression {
	return descendantsType{
		expType:     expType,
		descendants: descendants,
	}
}

func (p descendantsType) Type() s.PathExpressionType {
	if p.ExpType() == s.PDTAll {
		return s.PETAllDescendants
	}
	return s.PETDescendants
}

func (p descendantsType) Describe(w *bufio.Writer) error {
	if _, err := w.WriteRune('('); err != nil {
		return err
	}

	var (
		slash = s.PTTForwardSlash.String()
		val   = slash
	)
	if p.expType == s.PDTAll {
		val = fmt.Sprintf("%s%s", slash, slash)
	}

	if _, err := w.WriteString(val); err != nil {
		return err
	}

	if x, ok := p.descendants.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteRune(')'); err != nil {
		return err
	}

	return nil
}

func (p descendantsType) ExpType() s.PathDescendantsType {
	return p.expType
}

func (p descendantsType) Descendants() s.PathExpression {
	return p.descendants
}
