package expressions

import (
	"bufio"
	"fmt"

	s "github.com/SimonRichardson/cilli/selectors"
)

type emptyType struct {
	expressionType s.PathExpressionType
}

func MakePathWildcard() s.PathExpression {
	return emptyType{s.PETWildcard}
}

func (p emptyType) Type() s.PathExpressionType {
	return p.expressionType
}

func (p emptyType) Describe(w *bufio.Writer) error {
	_, err := w.WriteRune('*')
	return err
}

type booleanType struct {
	value bool
}

func MakePathBoolean(value bool) s.PathExpression {
	return booleanType{value}
}

func (p booleanType) Type() s.PathExpressionType {
	return s.PETBoolean
}

func (p booleanType) Val() bool {
	return p.value
}

func (p booleanType) Describe(w *bufio.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("%v", p.value))
	return err
}

type integerType struct {
	value int
}

func MakePathInteger(value int) s.PathExpression {
	return integerType{value}
}

func (p integerType) Type() s.PathExpressionType {
	return s.PETInteger
}

func (p integerType) Val() int {
	return p.value
}

func (p integerType) Describe(w *bufio.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("%d", p.value))
	return err
}

type nameType struct {
	value string
}

func MakePathName(value string) s.PathExpression {
	return nameType{value}
}

func (p nameType) Type() s.PathExpressionType {
	return s.PETName
}

func (p nameType) Val() string {
	return p.value
}

func (p nameType) Describe(w *bufio.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("%s", p.value))
	return err
}

func (p nameType) Name() string {
	return p.value
}

type numberType struct {
	value float64
}

func MakePathNumber(value float64) s.PathExpression {
	return numberType{value}
}

func (p numberType) Type() s.PathExpressionType {
	return s.PETNumber
}

func (p numberType) Val() float64 {
	return p.value
}

func (p numberType) Index() int {
	return int(p.value)
}

func (p numberType) Describe(w *bufio.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("%f", p.value))
	return err
}

type stringType struct {
	value string
}

func MakePathString(value string) s.PathExpression {
	return stringType{value}
}

func (p stringType) Type() s.PathExpressionType {
	return s.PETString
}

func (p stringType) Val() string {
	return p.value
}

func (p stringType) Describe(w *bufio.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("%q", p.value))
	return err
}
