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
	res, err := path.Execute(MakeElement("root", func() []s.Element {
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
	name     string
	children func() []s.Element
}

func MakeElement(name string, children func() []s.Element) s.Element {
	return element{name, children}
}

func (e element) Children() []s.Element {
	return e.children()
}

func (e element) Name() string {
	return e.name
}

func MakeElements(name string, amount uint) []s.Element {
	var (
		x   = int(amount)
		res = make([]s.Element, x, x)
	)
	for i := 0; i < x; i++ {
		res[i] = element{name, func() []s.Element {
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
		res[i] = element{"node", func() []s.Element {
			return MakeElements("subnode", numOfChildren)
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
			res, err := path.Execute(MakeElement("root", func() []s.Element {
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
			res, err := path.Execute(MakeElement("root", func() []s.Element {
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
				lex       = NewPathLexer("root").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}

			return len(res) == 1
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
				lex       = NewPathLexer("/node").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
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
				lex       = NewPathLexer("/node/subnode").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
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

func Test_PathExecuteNoForwardSlashWithNameThenName(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("root/node/subnode").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
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

func Test_PathExecuteForwardSlashWithNameThenNameAndWildcard(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/node/subnode/*").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
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

func Test_PathExecuteNoForwardSlashWithNameThenNameAndWildcard(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("root/node/subnode/*").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
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

func Test_PathExecuteForwardSlashWithIndexThenName(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/node[0]/subnode").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			for _, v := range res {
				if v.Name() != "subnode" {
					return false
				}
			}
			return len(res) == 10
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteForwardSlashWithIndexThenNameAndIndex(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/node[0]/subnode[0]").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			for _, v := range res {
				if v.Name() != "subnode" {
					return false
				}
			}
			return len(res) == 1
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteNoForwardSlashWithIndexThenName(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("root/node[0]/subnode").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			for _, v := range res {
				if v.Name() != "subnode" {
					return false
				}
			}
			return len(res) == 10
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteNoForwardSlashWithIndexThenNameAndIndex(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("root/node[0]/subnode[0]").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			for _, v := range res {
				if v.Name() != "subnode" {
					return false
				}
			}
			return len(res) == 1
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteForwardSlashWithIndexThenNameAndEmptyGroup(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/node[0]/subnode.()").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			for _, v := range res {
				if v.Name() != "subnode" {
					return false
				}
			}
			return len(res) == 10
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteNoForwardSlashWithIndexThenNameAndEmptyGroup(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("root/node[0]/subnode.()").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr)
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			for _, v := range res {
				if v.Name() != "subnode" {
					return false
				}
			}
			return len(res) == 10
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteForwardSlashWithIndexThenNameAndGroupEquality(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/node[0]/subnode.(@Name==\"subnode\")").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr).With(PathPredicate{
				Equality: func(elem s.Element, prop string, value interface{}) bool {
					if prop == "Name" {
						if str, err := strconv.Unquote(value.(string)); err == nil {
							return elem.Name() == str
						}
					}
					return false
				},
			})
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			for _, v := range res {
				if v.Name() != "subnode" {
					return false
				}
			}
			return len(res) == 10
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathExecuteForwardSlashWithDoubleGroupEquality(t *testing.T) {
	var (
		f = clamp(func(a uint) bool {
			var (
				types     = s.PathTokenTypes()
				lex       = NewPathLexer("/node.(@Name==\"node\")/subnode.(@Name==\"subnode\")").With(types)
				parser    = NewPathParser(lex.Iter())
				expr, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}

			path := NewPath(expr).With(PathPredicate{
				Equality: func(elem s.Element, prop string, value interface{}) bool {
					if prop == "Name" {
						if str, err := strconv.Unquote(value.(string)); err == nil {
							return elem.Name() == str
						}
					}
					return false
				},
			})
			res, err := path.Execute(MakeElement("root", func() []s.Element {
				return MakeElementsWithChildren(a, 10)
			}))
			if err != nil {
				t.Error(err)
			}
			for _, v := range res {
				if v.Name() != "subnode" {
					return false
				}
			}
			return len(res) == int(a)*10
		})
	)

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
