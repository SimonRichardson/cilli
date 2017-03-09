package cilli

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/SimonRichardson/cilli/expressions"
	s "github.com/SimonRichardson/cilli/selectors"
)

func Test_PathParserForString(t *testing.T) {
	var (
		f = func(a string) s.PathExpression {
			var (
				lex      = NewPathLexer(fmt.Sprintf("%q", a))
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		}
		g = func(a string) s.PathExpression {
			return expressions.MakePathString(fmt.Sprintf("%q", a))
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathParserForName(t *testing.T) {
	var (
		f = func(a Named) s.PathExpression {
			var (
				lex      = NewPathLexer(a.String())
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		}
		g = func(a Named) s.PathExpression {
			return expressions.MakePathName(a.String())
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathParserForNumber(t *testing.T) {
	var (
		f = func(a float64) s.PathExpression {
			var (
				lex      = NewPathLexer(fmt.Sprintf("%f", a))
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		}
		g = func(a float64) s.PathExpression {
			return expressions.MakePathNumber(a)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathParserWithTypesForOneForwardSlash(t *testing.T) {
	var (
		f = func(a string) s.PathExpression {
			var (
				lex      = NewPathLexer(fmt.Sprintf("/%q", a)).With(s.PathTokenTypes())
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		}
		g = func(a string) s.PathExpression {
			return expressions.MakePathDescendants(
				s.PDTContext,
				expressions.MakePathString(fmt.Sprintf("%q", a)),
			)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathParserWithTypesForTwoForwardSlashes(t *testing.T) {
	var (
		f = func(a string) s.PathExpression {
			var (
				lex      = NewPathLexer(fmt.Sprintf("//%q", a)).With(s.PathTokenTypes())
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		}
		g = func(a string) s.PathExpression {
			return expressions.MakePathDescendants(
				s.PDTAll,
				expressions.MakePathString(fmt.Sprintf("%q", a)),
			)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathParserWithTypesForNamedIndexAccess(t *testing.T) {
	var (
		clamp = func(fn func(Named, uint) s.PathExpression) func(Named, uint) s.PathExpression {
			return func(a Named, b uint) s.PathExpression {
				return fn(a, b%1000)
			}
		}
		f = clamp(func(a Named, b uint) s.PathExpression {
			var (
				lex      = NewPathLexer(fmt.Sprintf("%s[%d]", a.String(), b)).With(s.PathTokenTypes())
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		})
		g = clamp(func(a Named, b uint) s.PathExpression {
			return expressions.MakePathIndexAccess(
				expressions.MakePathName(a.String()),
				expressions.MakePathNumber(float64(b)),
			)
		})
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathParserWithTypesForNamedGroup(t *testing.T) {
	var (
		f = func(a Named) s.PathExpression {
			var (
				lex      = NewPathLexer(fmt.Sprintf("%s.()", a.String())).With(s.PathTokenTypes())
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		}
		g = func(a Named) s.PathExpression {
			return expressions.MakePathInstance(
				expressions.MakePathName(a.String()),
				expressions.MakePathGroup([]s.PathExpression{}),
			)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathParserWithTypesForNamedGroupWithEquality(t *testing.T) {
	var (
		f = func(a Named) s.PathExpression {
			var (
				lex      = NewPathLexer(fmt.Sprintf("%s.(@Name==%q)", a.String(), a.String())).With(s.PathTokenTypes())
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		}
		g = func(a Named) s.PathExpression {
			return expressions.MakePathInstance(
				expressions.MakePathName(a.String()),
				expressions.MakePathGroup([]s.PathExpression{
					expressions.MakePathAttribute(),
					expressions.MakePathEquality(
						expressions.MakePathName("Name"),
						expressions.MakePathString(fmt.Sprintf("%q", a.String())),
					),
				}),
			)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathParserWithTypesForDoubleGroupWithEquality(t *testing.T) {
	var (
		f = func(a Named) s.PathExpression {
			var (
				dsl      = fmt.Sprintf("/%s.(@Name==%q)/%s.(@Name==%q)", a.String(), a.String(), a.String(), a.String())
				lex      = NewPathLexer(dsl).With(s.PathTokenTypes())
				parser   = NewPathParser(lex.Iter())
				res, err = parser.ParseExpression()
			)
			if err != nil {
				t.Error(err)
			}
			return res
		}
		g = func(a Named) s.PathExpression {
			return expressions.MakePathDescendants(
				s.PDTContext,
				expressions.MakePathInstance(
					expressions.MakePathName(a.String()),
					expressions.MakePathBranch(
						expressions.MakePathGroup([]s.PathExpression{
							expressions.MakePathAttribute(),
							expressions.MakePathEquality(
								expressions.MakePathName("Name"),
								expressions.MakePathString(fmt.Sprintf("%q", a.String())),
							),
						}),
						expressions.MakePathInstance(
							expressions.MakePathName(a.String()),
							expressions.MakePathGroup([]s.PathExpression{
								expressions.MakePathAttribute(),
								expressions.MakePathEquality(
									expressions.MakePathName("Name"),
									expressions.MakePathString(fmt.Sprintf("%q", a.String())),
								),
							}),
						),
					),
				),
			)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}
