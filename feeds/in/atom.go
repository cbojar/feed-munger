package in

import (
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/cbojar/feed-munger/feeds"
)

const AtomNS = "http://www.w3.org/2005/Atom"
const AtomNSS = "https://www.w3.org/2005/Atom"

type AtomFormat struct {
	document *etree.Document
}

func (format *AtomFormat) Namespaces() feeds.Namespaces {
	namespaces := feeds.Namespaces{}

	for _, attr := range format.document.Root().Attr {
		if strings.HasPrefix(attr.FullKey(), "xmlns") {
			namespaces = append(namespaces, feeds.NewNamespace(attr.FullKey(), attr.Value))
		}
	}

	return namespaces
}

func (format *AtomFormat) ID() string {
	return format.document.FindElement("/feed/id").NotNil().Text()
}

func (format *AtomFormat) Title() string {
	return format.document.FindElement("/feed/title").NotNil().Text()
}

func (format *AtomFormat) Alternate() string {
	return format.document.FindElement("/feed/link[@rel='alternate']").NotNil().SelectAttrValue("href", "")
}

func (format *AtomFormat) Self() string {
	return format.document.FindElement("/feed/link[@rel='self']").NotNil().SelectAttrValue("href", "")
}

func (format *AtomFormat) Authors() []*feeds.Person {
	authors := []*feeds.Person{}

	for authorXML := range format.document.Root().SelectElementsSeq("author") {
		author := &feeds.Person{
			Name:       authorXML.SelectElement("name").NotNil().Text(),
			URI:        authorXML.SelectElement("uri").NotNil().Text(),
			Email:      authorXML.SelectElement("email").NotNil().Text(),
			Extensions: extensionsToAtomUnder(authorXML),
		}

		authors = append(authors, author)
	}

	return authors
}

func (format *AtomFormat) Contributors() []*feeds.Person {
	contributors := []*feeds.Person{}

	for contributorXML := range format.document.Root().SelectElementsSeq("contributor") {
		contributor := &feeds.Person{
			Name:       contributorXML.SelectElement("name").NotNil().Text(),
			URI:        contributorXML.SelectElement("uri").NotNil().Text(),
			Email:      contributorXML.SelectElement("email").NotNil().Text(),
			Extensions: extensionsToAtomUnder(contributorXML),
		}

		contributors = append(contributors, contributor)
	}

	return contributors
}

func (format *AtomFormat) Categories() []*feeds.Category {
	categories := []*feeds.Category{}

	for categoryXML := range format.document.Root().SelectElementsSeq("category") {
		category := &feeds.Category{
			Term:   categoryXML.SelectAttrValue("term", ""),
			Scheme: categoryXML.SelectAttrValue("scheme", ""),
			Label:  categoryXML.SelectAttrValue("label", ""),
		}

		categories = append(categories, category)
	}

	return categories
}

func (format *AtomFormat) Copyright() string {
	return format.document.Root().SelectElement("rights").NotNil().Text()
}

func (format *AtomFormat) Description() string {
	return format.document.Root().SelectElement("subtitle").NotNil().Text()
}

func (format *AtomFormat) Image() string {
	return format.document.Root().SelectElement("logo").NotNil().Text()
}

func (format *AtomFormat) Language() string {
	return format.document.Root().SelectAttrValue("xml:lang", "")
}

func (format *AtomFormat) Links() []*feeds.Link {
	links := []*feeds.Link{}

	for linkXML := range format.document.Root().SelectElementsSeq("link") {
		if rel := linkXML.SelectAttrValue("rel", ""); rel == "" || rel == "alternate" || rel == "self" {
			continue
		}

		link := &feeds.Link{
			Relationship: linkXML.SelectAttrValue("rel", ""),
			ContentType:  linkXML.SelectAttrValue("type", ""),
			HrefLang:     linkXML.SelectAttrValue("hreflang", ""),
			Href:         linkXML.SelectAttrValue("href", ""),
		}

		links = append(links, link)
	}

	return links
}

func (format *AtomFormat) Published() *time.Time {
	return elementTimeAtom(format.document.Root(), "published")
}

func (format *AtomFormat) Updated() *time.Time {
	return elementTimeAtom(format.document.Root(), "updated")
}

func (format *AtomFormat) Entries() []*feeds.Entry {
	feedTitle := format.Title()
	entriesXML := format.document.Root().SelectElementsSeq("entry")
	entries := []*feeds.Entry{}

	for entryXML := range entriesXML {
		entries = append(entries, entryAtom(entryXML, feedTitle))
	}

	return entries
}

func entryAtom(entryXML *etree.Element, feedTitle string) *feeds.Entry {
	links := []*feeds.Link{}
	for linkXML := range entryXML.SelectElementsSeq("link") {
		link := &feeds.Link{
			Relationship: linkXML.SelectAttrValue("rel", ""),
			ContentType:  linkXML.SelectAttrValue("type", ""),
			HrefLang:     linkXML.SelectAttrValue("hreflang", ""),
			Href:         linkXML.SelectAttrValue("href", ""),
		}

		links = append(links, link)
	}

	authors := []*feeds.Person{}
	for authorXML := range entryXML.SelectElementsSeq("author") {
		author := &feeds.Person{
			Name:       elementText(authorXML, "name"),
			URI:        elementText(authorXML, "uri"),
			Email:      elementText(authorXML, "email"),
			Extensions: extensionsToAtomUnder(authorXML),
		}

		authors = append(authors, author)
	}

	contributors := []*feeds.Person{}
	for contributorXML := range entryXML.SelectElementsSeq("contributor") {
		contributor := &feeds.Person{
			Name:       elementText(contributorXML, "name"),
			URI:        elementText(contributorXML, "uri"),
			Email:      elementText(contributorXML, "email"),
			Extensions: extensionsToAtomUnder(contributorXML),
		}

		contributors = append(contributors, contributor)
	}

	categories := []*feeds.Category{}
	for categoryXML := range entryXML.SelectElementsSeq("category") {
		category := &feeds.Category{
			Term:   categoryXML.SelectAttrValue("term", ""),
			Scheme: categoryXML.SelectAttrValue("scheme", ""),
			Label:  categoryXML.SelectAttrValue("label", ""),
		}

		categories = append(categories, category)
	}

	extensions := []*etree.Element{}
	for _, element := range extensionsToAtomUnder(entryXML) {
		extensions = append(extensions, element)
	}

	entry := &feeds.Entry{
		ID:           elementText(entryXML, "id"),
		Title:        elementText(entryXML, "title"),
		Updated:      elementTimeAtom(entryXML, "updated"),
		Published:    elementTimeAtom(entryXML, "published"),
		Summary:      elementText(entryXML, "summary"),
		Content:      elementText(entryXML, "content"),
		Rights:       elementText(entryXML, "rights"),
		Source:       feedTitle,
		Links:        links,
		Authors:      authors,
		Contributors: contributors,
		Categories:   categories,
		Extensions:   extensions,
	}

	return entry
}

func elementTimeAtom(xml *etree.Element, tag string) *time.Time {
	if text := elementText(xml, tag); text == "" {
		return nil
	} else if parsed, err := time.Parse(time.RFC3339, text); err != nil {
		return nil
	} else {
		return &parsed
	}
}

func extensionsToAtomUnder(xml *etree.Element) []*etree.Element {
	elements := []*etree.Element{}

	for element := range xml.ChildElementsSeq() {
		if ns := element.NamespaceURI(); ns != AtomNS && ns != AtomNSS {
			elements = append(elements, element)
		}
	}

	return elements
}
