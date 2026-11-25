package entries

import (
	"time"

	"github.com/beevik/etree"
)

type Entry struct {
	id                 string
	title              string
	updated            *time.Time
	published          *time.Time
	links              []*Link
	summary            string
	content            string
	rights             string
	authors            []*Person
	contributors       []*Person
	source             string
	categories         []*Category
	additionalElements []*etree.Element
}

func (entry *Entry) CompareUpdated(other *Entry) int {
	if entry.updated != nil && other.updated != nil {
		return entry.updated.Compare(*other.updated)
	} else if entry.updated == nil && other.updated != nil {
		return -1
	} else if entry.updated != nil && other.updated == nil {
		return 1
	} else {
		return 0
	}
}

func (entry *Entry) ID() string {
	return entry.id
}

func (entry *Entry) Title() string {
	return entry.title
}

func (entry *Entry) Updated() *time.Time {
	return entry.updated
}

func (entry *Entry) Published() *time.Time {
	return entry.published
}

func (entry *Entry) Summary() string {
	return entry.summary
}

func (entry *Entry) Content() string {
	return entry.content
}

func (entry *Entry) Rights() string {
	return entry.rights
}

func (entry *Entry) Source() string {
	return entry.source
}

func (entry *Entry) Links() []*Link {
	return entry.links
}

func (entry *Entry) Authors() []*Person {
	return entry.authors
}

func (entry *Entry) Contributors() []*Person {
	return entry.contributors
}

func (entry *Entry) Categories() []*Category {
	return entry.categories
}

func (entry *Entry) AdditionalElements() []*etree.Element {
	return entry.additionalElements
}
