package cilli

import (
	"bytes"
	"fmt"
	"strings"

	s "github.com/SimonRichardson/cilli/selectors"
)

type PathLexer struct {
	source string
	types  map[rune]s.PathTokenType
}

func NewPathLexer(source string) *PathLexer {
	return &PathLexer{
		source: source,
		types:  make(map[rune]s.PathTokenType),
	}
}

func (p *PathLexer) With(types []s.PathTokenType) *PathLexer {
	values := make(map[rune]s.PathTokenType, len(types))
	for _, v := range types {
		values[v.Rune()] = v
	}

	res := NewPathLexer(p.source)
	res.types = values
	return res
}

func (p *PathLexer) Iter() s.PathLexerIterator {
	return newPathIterator(p.source, p.types)
}

type pathLexerIterator struct {
	reader *strings.Reader
	types  map[rune]s.PathTokenType
}

func newPathIterator(source string, types map[rune]s.PathTokenType) *pathLexerIterator {
	return &pathLexerIterator{
		reader: strings.NewReader(source),
		types:  types,
	}
}

func (i *pathLexerIterator) HasNext() bool {
	return i.reader.Len() > 0
}

func (i *pathLexerIterator) Next() (s.PathToken, error) {
	var (
		token  = s.PTTNull
		buffer = bytes.NewBufferString("")

		lastChar rune
	)

loop:
	for {
		char, _, err := i.reader.ReadRune()
		if err != nil {
			return s.PathToken{}, err
		}

		// Number!
		if token == s.PTTNull || token == s.PTTNumber {
			// Include exponential numbers
			if (char >= 48 && char <= 57) || char == 45 || char == 46 || (token == s.PTTNumber && (char == 43 || char == 101)) {
				if token == s.PTTNull {
					token = s.PTTNumber
				}

				buffer.WriteRune(char)
				if i.HasNext() {
					continue loop
				}
				return s.MakePathToken(token, buffer.String()), nil
			}

			if token == s.PTTNumber {
				if buffer.Len() == 1 {
					// Work out if it's just an invalid number (., -, etc)
					nan, _, err := buffer.ReadRune()
					if err != nil {
						return s.PathToken{}, err
					}
					if !(nan >= 48 && nan <= 57) {
						if tokenType, ok := i.types[nan]; ok {
							// We gone too far, we're not what we think we are,
							// so rewind the rune, so we can move forward again.
							if err := i.reader.UnreadRune(); err != nil {
								return s.PathToken{}, err
							}

							return s.MakePathToken(tokenType, string(nan)), nil
						}
					}
					if err := buffer.UnreadRune(); err != nil {
						return s.PathToken{}, err
					}
				}

				if err := i.reader.UnreadRune(); err != nil {
					return s.PathToken{}, err
				}

				return s.MakePathToken(token, buffer.String()), nil
			}
		}

		// Strings
		if (token == s.PTTNull && char == 34) || token == s.PTTString {
			buffer.WriteRune(char)
			lastChar = char

			if token == s.PTTString && (char == 34 && lastChar != 92) {
				return s.MakePathToken(token, buffer.String()), nil
			}

			if token == s.PTTNull {
				token = s.PTTString
			}

			if i.HasNext() {
				continue loop
			}
			return s.MakePathToken(token, buffer.String()), nil
		}

		// Custom types
		if token == s.PTTNull {
			if tokenType, ok := i.types[char]; ok {
				return s.MakePathToken(tokenType, string(char)), nil
			}
		}

		// Named properties that are not strings.
		if token == s.PTTNull || token == s.PTTName {
			if (char >= 48 && char <= 57) || (char >= 65 && char <= 90) || (char >= 97 && char <= 122) || char == 95 {
				buffer.WriteRune(char)

				if token == s.PTTNull {
					token = s.PTTName
				}

				if i.HasNext() {
					continue loop
				}
				return s.MakePathToken(token, buffer.String()), nil
			}

			if token == s.PTTName {
				if err := i.reader.UnreadRune(); err != nil {
					return s.PathToken{}, err
				}

				return s.MakePathToken(token, buffer.String()), nil
			}
		}

		// Ignore spaces
		if char <= 32 {
			continue loop
		}

		// No idea what it is!
		break
	}

	return s.MakePathToken(s.PTTNull, ""), fmt.Errorf("Nothing found.")
}
