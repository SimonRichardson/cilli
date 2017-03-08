package expressions

import (
	"bufio"

	s "github.com/SimonRichardson/cilli/selectors"
)

type attributeType struct {
}

func MakePathAttribute() s.PathExpression {
	return attributeType{}
}

func (p attributeType) Type() s.PathExpressionType {
	return s.PETAttribute
}

func (p attributeType) Describe(w *bufio.Writer) error {
	_, err := w.WriteString(p.Type().String())
	return err
}

type methodCallType struct {
	method     s.PathExpression
	parameters []s.PathExpression
}

func MakePathMethodCall(method s.PathExpression,
	params []s.PathExpression,
) s.PathExpression {
	return methodCallType{method, params}
}

func (p methodCallType) Type() s.PathExpressionType {
	return s.PETMethodCall
}

func (p methodCallType) Method() s.PathExpression {
	return p.method
}

func (p methodCallType) Parameters() []s.PathExpression {
	return p.parameters
}

func (p methodCallType) Describe(w *bufio.Writer) error {
	if x, ok := p.method.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	for k, v := range p.parameters {
		if x, ok := v.(s.Describe); ok {
			if err := x.Describe(w); err != nil {
				return err
			}
		}
		if k < len(p.parameters)-1 {
			w.WriteString(", ")
		}
	}

	return nil
}

type groupType struct {
	expressions []s.PathExpression
}

func MakePathGroup(expressions []s.PathExpression) s.PathExpression {
	return groupType{expressions}
}

func (p groupType) Type() s.PathExpressionType {
	return s.PETGroup
}

func (p groupType) List() []s.PathExpression {
	return p.expressions
}

func (p groupType) Describe(w *bufio.Writer) error {
	if _, err := w.WriteRune('('); err != nil {
		return err
	}

	for k, v := range p.expressions {
		if x, ok := v.(s.Describe); ok {
			if err := x.Describe(w); err != nil {
				return err
			}
		}
		if k < len(p.expressions)-1 {
			w.WriteString(", ")
		}
	}

	if _, err := w.WriteRune(')'); err != nil {
		return err
	}

	return nil
}

type infixAttributeType struct {
	left, right s.PathExpression
}

func MakePathInfixAttribute(left, right s.PathExpression) s.PathExpression {
	return infixAttributeType{
		left:  left,
		right: right,
	}
}

func (p infixAttributeType) Left() s.PathExpression {
	return p.left
}

func (p infixAttributeType) Right() s.PathExpression {
	return p.right
}

func (p infixAttributeType) Type() s.PathExpressionType {
	return s.PETInfixAttribute
}

func (p infixAttributeType) Describe(w *bufio.Writer) error {
	if x, ok := p.left.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("@"); err != nil {
		return err
	}

	if x, ok := p.right.(s.Describe); ok {
		if err := x.Describe(w); err != nil {
			return err
		}
	}

	return nil
}
