package in

import (
	"net/url"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/cbojar/feed-munger/feeds"
)

const AtomNS = "http://www.w3.org/2005/Atom"
const AtomNSS = "https://www.w3.org/2005/Atom"

type Feed struct {
	document *etree.Document
}

func (feed *Feed) Identity() *feeds.Identity {
	id, _ := feed.id()
	title, _ := feed.title()
	alternate, _ := feed.alternate()
	self, _ := feed.self()

	return feeds.NewIdentity(id, title, alternate, self)
}

func (feed *Feed) Namespaces() feeds.Namespaces {
	namespaces := feeds.Namespaces{}

	for _, attr := range feed.document.Root().Attr {
		if strings.HasPrefix(attr.FullKey(), "xmlns") {
			namespaces = append(namespaces, feeds.NewNamespace(attr.FullKey(), attr.Value))
		}
	}

	return namespaces
}

func (feed *Feed) id() (string, bool) {
	if id := feed.document.FindElement("/feed/id"); id == nil {
		return "", false
	} else {
		return id.Text(), true
	}
}

func (feed *Feed) title() (string, bool) {
	if title := feed.document.FindElement("/feed/title"); title == nil {
		return "", false
	} else {
		return title.Text(), true
	}
}

func (feed *Feed) alternate() (string, bool) {
	if link := feed.document.FindElement("/feed/link[@rel='alternate']"); link == nil {
		return "", false
	} else {
		return link.SelectAttrValue("href", ""), true
	}
}

func (feed *Feed) self() (string, bool) {
	if link := feed.document.FindElement("/feed/link[@rel='self']"); link == nil {
		return "", false
	} else {
		return link.SelectAttrValue("href", ""), true
	}
}

func (feed *Feed) Entries() []*feeds.Entry {
	feedTitle, _ := feed.title()
	entriesXML := feed.document.Root().SelectElementsSeq("entry")
	entries := []*feeds.Entry{}

	for entryXML := range entriesXML {
		entries = append(entries, entry(entryXML, feedTitle))
	}

	return entries
}

func entry(entryXML *etree.Element, feedTitle string) *feeds.Entry {
	return feeds.BuildEntry(func(builder *feeds.EntryBuilder) {
		builder.ID = elementText(entryXML, "id")
		builder.Title = elementText(entryXML, "title")
		builder.Updated = elementTime(entryXML, "updated")
		builder.Published = elementTime(entryXML, "published")
		builder.Summary = elementText(entryXML, "summary")
		builder.Content = elementText(entryXML, "content")
		builder.Rights = elementText(entryXML, "rights")
		builder.Source = feedTitle

		for linkXML := range entryXML.SelectElementsSeq("link") {
			rel := linkXML.SelectAttrValue("rel", "")
			contentType := linkXML.SelectAttrValue("type", "")
			hreflang := linkXML.SelectAttrValue("hreflang", "")
			href, _ := url.Parse(linkXML.SelectAttrValue("href", ""))
			builder.AddLink(rel, contentType, hreflang, href)
		}

		for authorXML := range entryXML.SelectElementsSeq("author") {
			name := elementText(authorXML, "name")
			uri := elementURL(authorXML, "uri")
			email := elementText(authorXML, "email")
			builder.AddAuthor(name, uri, email, additionalElementsUnder(authorXML))
		}

		for contributorXML := range entryXML.SelectElementsSeq("contributor") {
			name := elementText(contributorXML, "name")
			uri := elementURL(contributorXML, "uri")
			email := elementText(contributorXML, "email")
			builder.AddContributor(name, uri, email, additionalElementsUnder(contributorXML))
		}

		for categoryXML := range entryXML.SelectElementsSeq("category") {
			term := categoryXML.SelectAttrValue("term", "")
			scheme := categoryXML.SelectAttrValue("scheme", "")
			label := categoryXML.SelectAttrValue("label", "")
			builder.AddCategory(term, scheme, label)
		}

		for _, element := range additionalElementsUnder(entryXML) {
			builder.AddAdditionalElement(element)
		}
	})
}

func elementText(xml *etree.Element, tag string) string {
	return xml.SelectElement(tag).NotNil().Text()
}

func elementURL(xml *etree.Element, tag string) *url.URL {
	if text := elementText(xml, tag); text == "" {
		return nil
	} else if uri, err := url.Parse(text); err != nil {
		return nil
	} else {
		return uri
	}
}

func elementTime(xml *etree.Element, tag string) *time.Time {
	if text := elementText(xml, tag); text == "" {
		return nil
	} else if parsed, err := time.Parse(time.RFC3339, text); err != nil {
		return nil
	} else {
		return &parsed
	}
}

func additionalElementsUnder(xml *etree.Element) []*etree.Element {
	elements := []*etree.Element{}

	for element := range xml.ChildElementsSeq() {
		if ns := element.NamespaceURI(); ns != AtomNS && ns != AtomNSS {
			elements = append(elements, element)
		}
	}

	return elements
}
