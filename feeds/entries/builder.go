package entries

import (
	"net/url"
	"time"

	"github.com/beevik/etree"
)

type EntryBuilder struct {
	ID                 string
	Title              string
	Updated            *time.Time
	Published          *time.Time
	Links              []*Link
	Summary            string
	Content            string
	Rights             string
	Authors            []*Person
	Contributors       []*Person
	Source             string
	Categories         []*Category
	AdditionalElements []*etree.Element
}

func BuildEntry(build func(*EntryBuilder)) *Entry {
	builder := &EntryBuilder{
		Links:              []*Link{},
		Authors:            []*Person{},
		Contributors:       []*Person{},
		Categories:         []*Category{},
		AdditionalElements: []*etree.Element{},
	}

	build(builder)

	return builder.Done()
}

func (entry *EntryBuilder) AddLink(rel string, contentType string, hreflang string, href *url.URL) {
	entry.Links = append(entry.Links, &Link{
		relationship: rel,
		contentType:  contentType,
		hreflang:     hreflang,
		href:         href,
	})
}

func (entry *EntryBuilder) AddAuthor(name string, uri *url.URL, email string, additionalElements []*etree.Element) {
	entry.Authors = append(entry.Authors, &Person{
		name:               name,
		uri:                uri,
		email:              email,
		additionalElements: additionalElements,
	})
}

func (entry *EntryBuilder) AddContributor(name string, uri *url.URL, email string, additionalElements []*etree.Element) {
	entry.Contributors = append(entry.Contributors, &Person{
		name:               name,
		uri:                uri,
		email:              email,
		additionalElements: additionalElements,
	})
}

func (entry *EntryBuilder) AddCategory(term string, scheme string, label string) {
	entry.Categories = append(entry.Categories, &Category{
		term:   term,
		scheme: scheme,
		label:  label,
	})
}

func (entry *EntryBuilder) AddAdditionalElement(element *etree.Element) {
	entry.AdditionalElements = append(entry.AdditionalElements, element)
}

func (entry *EntryBuilder) Done() *Entry {
	return &Entry{
		id:                 entry.ID,
		title:              entry.Title,
		updated:            entry.Updated,
		published:          entry.Published,
		links:              entry.Links,
		summary:            entry.Summary,
		content:            entry.Content,
		rights:             entry.Rights,
		authors:            entry.Authors,
		contributors:       entry.Contributors,
		source:             entry.Source,
		categories:         entry.Categories,
		additionalElements: entry.AdditionalElements,
	}
}
