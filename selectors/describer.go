package selectors

import "bufio"

type Describe interface {
	Describe(*bufio.Writer) error
}
