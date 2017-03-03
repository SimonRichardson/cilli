package cilli

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	s "github.com/SimonRichardson/cilli/selectors"
)

func Test_PathLexerForWildcard(t *testing.T) {
	lex := NewPathLexer("*").With(s.PathTokenTypes())
	val, err := lex.Iter().Next()
	if err != nil {
		t.Error(err)
	}

	if val.Val() != "*" {
		t.Fail()
	}
}

func Test_PathLexerForFloat(t *testing.T) {
	var (
		f = func(a float64) string {
			lex := NewPathLexer(fmt.Sprintf("%f", a))
			val, err := lex.Iter().Next()
			if err != nil {
				t.Fatal(err)
			}
			return val.Val()
		}
		g = func(a float64) string {
			return fmt.Sprintf("%f", a)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathLexerForInt(t *testing.T) {
	var (
		f = func(a int64) string {
			lex := NewPathLexer(fmt.Sprintf("%d", a))
			val, err := lex.Iter().Next()
			if err != nil {
				t.Fatal(err)
			}
			return val.Val()
		}
		g = func(a int64) string {
			return fmt.Sprintf("%d", a)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathLexerForUint(t *testing.T) {
	var (
		f = func(a int64) string {
			lex := NewPathLexer(fmt.Sprintf("%d", a))
			val, err := lex.Iter().Next()
			if err != nil {
				t.Fatal(err)
			}
			return val.Val()
		}
		g = func(a int64) string {
			return fmt.Sprintf("%d", a)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathLexerForString(t *testing.T) {
	var (
		f = func(a string) string {
			lex := NewPathLexer(fmt.Sprintf("%q", a))
			val, err := lex.Iter().Next()
			if err != nil {
				t.Fatal(err)
			}
			return val.Val()
		}
		g = func(a string) string {
			return fmt.Sprintf("%q", a)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathLexerForNamed(t *testing.T) {
	var (
		f = func(a Named) string {
			lex := NewPathLexer(a.String())
			val, err := lex.Iter().Next()
			if err != nil {
				t.Fatal(err)
			}
			return val.Val()
		}
		g = func(a Named) string {
			return a.String()
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

type Named []string

func (s Named) Generate(r *rand.Rand, size int) reflect.Value {
	m := []string{GenerateNamedWithRand(r, size)}
	return reflect.ValueOf(m)
}

func (s Named) String() string {
	return s[0]
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateNamedWithRand(r *rand.Rand, size int) string {
	result := make([]byte, size)
	// Make sure the first char is always alpha.
	result[0] = chars[r.Intn(len(chars)-11)]
	for i := 1; i < size; i++ {
		result[i] = chars[r.Intn(len(chars)-1)]
	}
	return string(result)
}

func Test_PathLexerForStringFloatNamed(t *testing.T) {
	var (
		f = func(a string, b float64, c Named) string {
			var (
				lex  = NewPathLexer(fmt.Sprintf("%q%f%s", a, b, c.String()))
				iter = lex.Iter()

				x = next(t, iter).Val()
				y = next(t, iter).Val()
				z = next(t, iter).Val()
			)
			return fmt.Sprintf("%s%s%s", x, y, z)
		}
		g = func(a string, b float64, c Named) string {
			return fmt.Sprintf("%q%f%s", a, b, c.String())
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func next(t *testing.T, iter s.PathLexerIterator) s.PathToken {
	res, err := iter.Next()
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func Test_PathLexerForFloatNamedIntString(t *testing.T) {
	var (
		f = func(a float64, b Named, c int, d string) string {
			var (
				// We need to insert space here, because Named values are greedy.
				lex  = NewPathLexer(fmt.Sprintf("%f%s %d%q", a, b.String(), c, d))
				iter = lex.Iter()

				v = next(t, iter).Val()
				x = next(t, iter).Val()
				y = next(t, iter).Val()
				z = next(t, iter).Val()
			)
			return fmt.Sprintf("%s%s%s%s", v, x, y, z)
		}
		g = func(a float64, b Named, c int, d string) string {
			return fmt.Sprintf("%f%s%d%q", a, b.String(), c, d)
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathLexerForStringFloatNamedWithSpaces(t *testing.T) {
	var (
		f = func(a string, b float64, c Named) string {
			var (
				lex  = NewPathLexer(fmt.Sprintf("%q           %f %s", a, b, c.String()))
				iter = lex.Iter()

				x = next(t, iter).Val()
				y = next(t, iter).Val()
				z = next(t, iter).Val()
			)
			return fmt.Sprintf("%s%s%s", x, y, z)
		}
		g = func(a string, b float64, c Named) string {
			return fmt.Sprintf("%q%f%s", a, b, c.String())
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}

func Test_PathLexerWithTypesForForwardSlashesAndNamed(t *testing.T) {
	var (
		f = func(a Named, b Named) string {
			var (
				lex  = NewPathLexer(fmt.Sprintf("/%s/%s", a.String(), b.String())).With(s.PathTokenTypes())
				iter = lex.Iter()

				w = next(t, iter).Val()
				x = next(t, iter).Val()
				y = next(t, iter).Val()
				z = next(t, iter).Val()
			)

			return fmt.Sprintf("%s%s%s%s", w, x, y, z)
		}
		g = func(a Named, b Named) string {
			return fmt.Sprintf("/%s/%s", a.String(), b.String())
		}
	)

	if err := quick.CheckEqual(f, g, nil); err != nil {
		t.Error(err)
	}
}
