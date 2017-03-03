package cilli

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"testing/quick"

	s "github.com/SimonRichardson/cilli/selectors"
)

func Test_PathDescribe(t *testing.T) {
	var (
		f = func(a string) string {
			var (
				lex       = NewPathLexer(fmt.Sprintf("%q", a))
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			var (
				path = NewPath(expr)

				buffer = new(bytes.Buffer)
				writer = bufio.NewWriter(buffer)
			)

			path.Describe(writer)
			writer.Flush()

			res, _ := strconv.Unquote(buffer.String())
			return res
		}
		g = func(a string) string {
			return fmt.Sprintf("%q", a)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteWildcardZero(t *testing.T) {
	var (
		types     = s.PathTokenTypes()
		lex       = NewPathLexer("*").With(types)
		parser    = NewPathParser(lex.Iter())
		expr, err = parser.ParseExpression()
	)
	if err != nil {
		t.Error(err)
	}

	path := NewPath(expr)
	res, err := path.Execute(MakeElement(func() []s.Element {
		return []s.Element{}
	}))
	if err != nil {
		t.Error(err)
	}

	if len(res) != 0 {
		t.Error("Expected zero element")
	}
}

type element struct {
	children func() []s.Element
}

func MakeElement(children func() []s.Element) s.Element {
	return element{children}
}

func (e element) Children() []s.Element {
	return e.children()
}

func (e element) Name() string {
	return "event"
}

func MakeElements(amount uint) []s.Element {
	var (
		x   = int(amount)
		res = make([]s.Element, x, x)
	)
	for i := 0; i < x; i++ {
		res[i] = element{func() []s.Element {
			return []s.Element{}
		}}
	}
	return res
}

func MakeElementsWithChildren(amount, numOfChildren uint) []s.Element {
	var (
		x   = int(amount)
		res = make([]s.Element, x, x)
	)
	for i := 0; i < x; i++ {
		res[i] = element{func() []s.Element {
			return MakeElements(numOfChildren)
		}}
	}
	return res
}

func Test_PathExecuteWildcard(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("*").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement(func() []s.Element {
				return MakeElementsWithChildren(a, 8)
			}))
			if err != nil {
				t.Error(err)
			}

			return len(res) == (int(a) * 9)
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}

}

func Test_PathExecuteForwardSlashWithWildcard(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/*").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement(func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}

			return len(res) == int(a)
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func clamp(fn func(uint) bool) func(uint) bool {
	return func(x uint) bool {
		return fn(x % 1000)
	}
}

func Test_PathExecuteName(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("event").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement(func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}

			return len(res) == int(a)
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteForwardSlashWithName(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/event").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement(func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}

			return len(res) == int(a)
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteForwardSlashWithNameThenName(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/event/event").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement(func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			return len(res) == int(a)*10
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
