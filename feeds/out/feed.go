package out

import (
	"time"

	"github.com/beevik/etree"
	"github.com/cbojar/feed-munger/feeds"
)

type Feed struct {
	document *etree.Document
}

func NewFeed(identity *feeds.Identity) *Feed {
	document := etree.NewDocument()
	document.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := document.CreateElement("feed")
	root.CreateAttr("xmlns", "http://www.w3.org/2005/Atom")

	if identity.ID() != "" {
		id := root.CreateElement("id")
		id.SetText(identity.ID())
	}

	if identity.Title() != "" {
		title := root.CreateElement("title")
		title.SetText(identity.Title())
	}

	if identity.Alternate() != "" {
		alternate := root.CreateElement("link")
		alternate.CreateAttr("rel", "alternate")
		alternate.CreateAttr("type", "text/html")
		alternate.CreateAttr("href", identity.Alternate())
	}

	if identity.Self() != "" {
		self := root.CreateElement("link")
		self.CreateAttr("rel", "self")
		self.CreateAttr("type", "application/atom+xml")
		self.CreateAttr("href", identity.Self())
	}

	updated := root.CreateElement("updated")
	updated.SetText(time.Now().Format(time.RFC3339))

	return &Feed{document: document}
}

func (feed *Feed) AddAllNamespaces(namespaces feeds.Namespaces) {
	for _, namespace := range namespaces {
		feed.AddNamespace(namespace)
	}
}

func (feed *Feed) AddNamespace(namespace *feeds.Namespace) {
	feed.document.Root().CreateAttr(namespace.Key(), namespace.Value())
}

func (feed *Feed) AddAllEntries(entries []*feeds.Entry) {
	for _, entry := range entries {
		feed.AddEntry(entry)
	}
}

func (feed *Feed) AddEntry(entry *feeds.Entry) {
	entryXML := feed.document.Root().CreateElement("entry")
	entryXML.CreateElement("id").SetText(entry.ID())
	entryXML.CreateElement("title").SetText(entry.Title())
	entryXML.CreateElement("updated").SetText(entry.Updated().Format(time.RFC3339))

	if entry.Published() != nil {
		entryXML.CreateElement("published").SetText(entry.Published().Format(time.RFC3339))
	}

	for _, link := range entry.Links() {
		linkXML := entryXML.CreateElement("link")

		if link.Relationship() != "" {
			linkXML.CreateAttr("rel", link.Relationship())
		} else {
			linkXML.CreateAttr("rel", "alternate")
		}

		if link.ContentType() != "" {
			linkXML.CreateAttr("type", link.ContentType())
		}

		if link.HrefLang() != "" {
			linkXML.CreateAttr("hreflang", link.HrefLang())
		}

		if link.Href() != nil {
			linkXML.CreateAttr("href", link.Href().String())
		}
	}

	for _, author := range entry.Authors() {
		authorXML := entryXML.CreateElement("author")
		if author.Name() != "" {
			authorXML.CreateElement("name").SetText(author.Name())
		}

		if author.URI() != nil {
			authorXML.CreateElement("uri").SetText(author.URI().String())
		}

		if author.Email() != "" {
			authorXML.CreateElement("email").SetText(author.Email())
		}

		for _, element := range author.AdditionalElements() {
			authorXML.AddChild(element)
		}
	}

	for _, contributor := range entry.Contributors() {
		contributorXML := entryXML.CreateElement("contributor")
		if contributor.Name() != "" {
			contributorXML.CreateElement("name").SetText(contributor.Name())
		}

		if contributor.URI() != nil {
			contributorXML.CreateElement("uri").SetText(contributor.URI().String())
		}

		if contributor.Email() != "" {
			contributorXML.CreateElement("email").SetText(contributor.Email())
		}

		for _, element := range contributor.AdditionalElements() {
			contributorXML.AddChild(element)
		}
	}

	if entry.Summary() != "" {
		entryXML.CreateElement("summary").SetText(entry.Summary())
	}

	if entry.Content() != "" {
		entryXML.CreateElement("content").SetText(entry.Content())
	}

	if entry.Rights() != "" {
		entryXML.CreateElement("rights").SetText(entry.Rights())
	}

	if entry.Source() != "" {
		entryXML.CreateElement("source").CreateElement("name").SetText(entry.Source())
	}

	for _, category := range entry.Categories() {
		categoryXML := entryXML.CreateElement("category")

		categoryXML.CreateAttr("term", category.Term())

		if category.Scheme() != "" {
			categoryXML.CreateAttr("scheme", category.Scheme())
		}

		if category.Label() != "" {
			categoryXML.CreateAttr("label", category.Label())
		}
	}

	for _, element := range entry.AdditionalElements() {
		entryXML.AddChild(element)
	}
}

func (feed *Feed) Bytes() ([]byte, error) {
	feed.document.IndentTabs()
	return feed.document.WriteToBytes()
}
