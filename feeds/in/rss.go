package in

import (
	"regexp"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/cbojar/feed-munger/feeds"
)

type RSSFormat struct {
	document *etree.Document
}

func (format *RSSFormat) channel() *etree.Element {
	return format.document.Root().SelectElement("channel")
}

func (format *RSSFormat) Namespaces() feeds.Namespaces {
	namespaces := feeds.Namespaces{}

	for _, attr := range format.document.Root().Attr {
		if strings.HasPrefix(attr.FullKey(), "xmlns") && attr.Value != AtomNS {
			namespaces = append(namespaces, feeds.NewNamespace(attr.FullKey(), attr.Value))
		}
	}

	return namespaces
}

func (format *RSSFormat) ID() string {
	return format.channel().SelectElement("atom:id").NotNil().Text()
}

func (format *RSSFormat) Title() string {
	return format.channel().SelectElement("title").NotNil().Text()
}

func (format *RSSFormat) Alternate() string {
	for linkXML := range format.channel().SelectElementsSeq("link") {
		if linkXML.NamespaceURI() == "" {
			return linkXML.Text()
		}
	}

	return ""
}

func (format *RSSFormat) Self() string {
	return format.channel().FindElement("./atom:link[@rel='self']").NotNil().SelectAttrValue("href", "")
}

func (format *RSSFormat) Authors() []*feeds.Person {
	authors := []*feeds.Person{}

	if editorXML := format.channel().SelectElement("managingEditor"); editorXML != nil {
		pattern := regexp.MustCompile(`^(.*@.*)\s+\((.*)\)\s*$`)
		if matches := pattern.FindStringSubmatch(editorXML.Text()); matches != nil {
			editor := &feeds.Person{
				Name:       matches[2],
				URI:        "",
				Email:      matches[1],
				Extensions: []*etree.Element{},
			}

			authors = append(authors, editor)
		} else {
			editor := &feeds.Person{
				Name:       "",
				URI:        "",
				Email:      editorXML.Text(),
				Extensions: []*etree.Element{},
			}

			authors = append(authors, editor)
		}
	}

	for authorXML := range format.channel().SelectElementsSeq("atom:author") {
		author := &feeds.Person{
			Name:       authorXML.SelectElement("name").NotNil().Text(),
			URI:        authorXML.SelectElement("uri").NotNil().Text(),
			Email:      authorXML.SelectElement("email").NotNil().Text(),
			Extensions: extensionsToRSSUnder(authorXML),
		}

		authors = append(authors, author)
	}

	return authors
}

func (format *RSSFormat) Contributors() []*feeds.Person {
	// Contributors do not exist in RSS
	return []*feeds.Person{}
}

func (format *RSSFormat) Categories() []*feeds.Category {
	categories := []*feeds.Category{}

	for categoryXML := range format.channel().SelectElementsSeq("category") {
		category := &feeds.Category{
			Term:   categoryXML.Text(),
			Scheme: categoryXML.SelectAttrValue("domain", ""),
			Label:  "",
		}

		categories = append(categories, category)
	}

	return categories
}

func (format *RSSFormat) Copyright() string {
	return format.channel().SelectElement("copyright").NotNil().Text()
}

func (format *RSSFormat) Description() string {
	return format.channel().SelectElement("description").NotNil().Text()
}

func (format *RSSFormat) Image() string {
	return format.channel().SelectElement("image").NotNil().Text()
}

func (format *RSSFormat) Language() string {
	return format.channel().SelectElement("language").NotNil().Text()
}

func (format *RSSFormat) Links() []*feeds.Link {
	links := []*feeds.Link{}

	for linkXML := range format.channel().SelectElementsSeq("atom:link") {
		if rel := linkXML.SelectAttrValue("rel", ""); rel == "" || rel == "alternative" || rel == "self" {
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

func (format *RSSFormat) Published() *time.Time {
	return elementTimeRSS(format.channel(), "pubDate")
}

func (format *RSSFormat) Updated() *time.Time {
	return elementTimeRSS(format.channel(), "lastBuildDate")
}

func (format *RSSFormat) Entries() []*feeds.Entry {
	feedTitle := format.Title()
	entriesXML := format.channel().SelectElementsSeq("item")
	entries := []*feeds.Entry{}

	for entryXML := range entriesXML {
		entries = append(entries, entryRSS(entryXML, feedTitle))
	}

	return entries
}

func entryRSS(entryXML *etree.Element, feedTitle string) *feeds.Entry {
	links := []*feeds.Link{}
	// TODO: enclosures
	if linkXML := entryXML.SelectElement("link"); linkXML != nil {
		link := &feeds.Link{
			Relationship: "alternate",
			ContentType:  "text/html",
			HrefLang:     "",
			Href:         linkXML.Text(),
		}
		links = append(links, link)
	}

	for linkXML := range entryXML.SelectElementsSeq("atom:link") {
		link := &feeds.Link{
			Relationship: linkXML.SelectAttrValue("rel", ""),
			ContentType:  linkXML.SelectAttrValue("type", ""),
			HrefLang:     linkXML.SelectAttrValue("hreflang", ""),
			Href:         linkXML.SelectAttrValue("href", ""),
		}

		links = append(links, link)
	}

	authors := []*feeds.Person{}
	for authorXML := range entryXML.SelectElementsSeq("atom:author") {
		author := &feeds.Person{
			Name:       elementText(authorXML, "name"),
			URI:        elementText(authorXML, "uri"),
			Email:      elementText(authorXML, "email"),
			Extensions: extensionsToAtomUnder(authorXML),
		}

		authors = append(authors, author)
	}

	categories := []*feeds.Category{}
	for categoryXML := range entryXML.SelectElementsSeq("category") {
		category := &feeds.Category{
			Term:   categoryXML.Text(),
			Scheme: categoryXML.SelectAttrValue("domain", ""),
			Label:  "",
		}

		categories = append(categories, category)
	}

	extensions := []*etree.Element{}
	for _, element := range extensionsToRSSUnder(entryXML) {
		extensions = append(extensions, element)
	}

	entry := &feeds.Entry{
		ID:           elementText(entryXML, "guid"),
		Title:        elementText(entryXML, "title"),
		Updated:      nil,
		Published:    elementTimeRSS(entryXML, "pubDate"),
		Summary:      elementText(entryXML, "description"),
		Content:      "", //TODO
		Rights:       elementText(entryXML, "rights"),
		Source:       feedTitle,
		Links:        links,
		Authors:      authors,
		Contributors: []*feeds.Person{},
		Categories:   categories,
		Extensions:   extensions,
	}

	return entry
}

func elementTimeRSS(xml *etree.Element, tag string) *time.Time {
	if text := elementText(xml, tag); text == "" {
		return nil
	} else if parsed, err := time.Parse(time.RFC822, text); err == nil {
		return &parsed
	} else if parsed, err := time.Parse(time.RFC1123, text); err == nil {
		return &parsed
	} else if parsed, err := time.Parse("Mon, 02 Jan 2006 15:04 MST", text); err == nil {
		return &parsed
	} else if parsed, err := time.Parse(time.RFC850, text); err == nil {
		return &parsed
	} else if parsed, err := time.Parse(time.ANSIC, text); err == nil {
		return &parsed
	} else if parsed, err := time.Parse(time.DateTime, text); err == nil {
		return &parsed
	} else if parsed, err := time.Parse(time.Layout, text); err == nil {
		return &parsed
	} else if parsed, err := time.Parse(time.RubyDate, text); err == nil {
		return &parsed
	} else if parsed, err := time.Parse(time.UnixDate, text); err == nil {
		return &parsed
	} else {
		return nil
	}
}

func extensionsToRSSUnder(xml *etree.Element) []*etree.Element {
	elements := []*etree.Element{}

	for element := range xml.ChildElementsSeq() {
		if ns := element.NamespaceURI(); ns == "" {
			continue
		} else if ns == AtomNS && (element.Tag == "id" || element.Tag == "author" || element.Tag == "link") {
			continue
		} else {
			elements = append(elements, element)
		}
	}

	return elements
}
