package selectors

type Element interface {
	Name() string
	Children() []Element
}
