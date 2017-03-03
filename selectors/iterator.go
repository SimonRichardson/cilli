package selectors

type PathLexerIterator interface {
	HasNext() bool
	Next() (PathToken, error)
}
