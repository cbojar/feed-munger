package entries

type Category struct {
	term   string
	scheme string
	label  string
}

func (category *Category) Term() string {
	return category.term
}

func (category *Category) Scheme() string {
	return category.scheme
}

func (category *Category) Label() string {
	return category.label
}
