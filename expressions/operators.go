package expressions

import (
	"bufio"

	s "github.com/SimonRichardson/cilli/selectors"
)

type equalityType struct {
	left, right s.PathExpression
}

func MakePathEquality(left, right s.PathExpression) s.PathExpression {
	return equalityType{
		left:  left,
		right: right,
	}
}

func (p equalityType) Type() s.PathExpressionType {
	return s.PETEquality
}

func (p equalityType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("=="); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p equalityType) Left() s.PathExpression {
	return p.left
}

func (p equalityType) Right() s.PathExpression {
	return p.right
}

type greaterThanType struct {
	left, right s.PathExpression
}

func MakePathGreaterThan(left, right s.PathExpression) s.PathExpression {
	return greaterThanType{
		left:  left,
		right: right,
	}
}

func (p greaterThanType) Type() s.PathExpressionType {
	return s.PETGreaterThan
}

func (p greaterThanType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p greaterThanType) Left() s.PathExpression {
	return p.left
}

func (p greaterThanType) Right() s.PathExpression {
	return p.right
}

type greaterThanOrEqualToType struct {
	left, right s.PathExpression
}

func MakePathGreaterThanOrEqualTo(left, right s.PathExpression) s.PathExpression {
	return greaterThanOrEqualToType{
		left:  left,
		right: right,
	}
}

func (p greaterThanOrEqualToType) Type() s.PathExpressionType {
	return s.PETGreaterThanOrEqualTo
}

func (p greaterThanOrEqualToType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString(">="); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p greaterThanOrEqualToType) Left() s.PathExpression {
	return p.left
}

func (p greaterThanOrEqualToType) Right() s.PathExpression {
	return p.right
}

type inequalityType struct {
	left, right s.PathExpression
}

func MakePathInequality(left, right s.PathExpression) s.PathExpression {
	return inequalityType{
		left:  left,
		right: right,
	}
}

func (p inequalityType) Type() s.PathExpressionType {
	return s.PETInequality
}

func (p inequalityType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("!="); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p inequalityType) Left() s.PathExpression {
	return p.left
}

func (p inequalityType) Right() s.PathExpression {
	return p.right
}

type lessThanType struct {
	left, right s.PathExpression
}

func MakePathLessThan(left, right s.PathExpression) s.PathExpression {
	return lessThanType{
		left:  left,
		right: right,
	}
}

func (p lessThanType) Type() s.PathExpressionType {
	return s.PETLessThan
}

func (p lessThanType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p lessThanType) Left() s.PathExpression {
	return p.left
}

func (p lessThanType) Right() s.PathExpression {
	return p.right
}

type lessThanOrEqualToType struct {
	left, right s.PathExpression
}

func MakePathLessThanOrEqualTo(left, right s.PathExpression) s.PathExpression {
	return lessThanOrEqualToType{
		left:  left,
		right: right,
	}
}

func (p lessThanOrEqualToType) Type() s.PathExpressionType {
	return s.PETLessThanOrEqualTo
}

func (p lessThanOrEqualToType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString(">="); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p lessThanOrEqualToType) Left() s.PathExpression {
	return p.left
}

func (p lessThanOrEqualToType) Right() s.PathExpression {
	return p.right
}

type logicalAndType struct {
	left, right s.PathExpression
}

func MakePathLogicalAnd(left, right s.PathExpression) s.PathExpression {
	return logicalAndType{
		left:  left,
		right: right,
	}
}

func (p logicalAndType) Type() s.PathExpressionType {
	return s.PETLogicalAnd
}

func (p logicalAndType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("&&"); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p logicalAndType) Left() s.PathExpression {
	return p.left
}

func (p logicalAndType) Right() s.PathExpression {
	return p.right
}

type logicalOrType struct {
	left, right s.PathExpression
}

func MakePathLogicalOr(left, right s.PathExpression) s.PathExpression {
	return logicalOrType{
		left:  left,
		right: right,
	}
}

func (p logicalOrType) Type() s.PathExpressionType {
	return s.PETLogicalOr
}

func (p logicalOrType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("||"); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}

func (p logicalOrType) Left() s.PathExpression {
	return p.left
}

func (p logicalOrType) Right() s.PathExpression {
	return p.right
}
