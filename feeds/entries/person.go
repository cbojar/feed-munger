package entries

import (
	"net/url"

	"github.com/beevik/etree"
)

type Person struct {
	name               string
	uri                *url.URL
	email              string
	additionalElements []*etree.Element
}

func (person *Person) Name() string {
	return person.name
}

func (person *Person) URI() *url.URL {
	return person.uri
}

func (person *Person) Email() string {
	return person.email
}

func (person *Person) AdditionalElements() []*etree.Element {
	return person.additionalElements
}
