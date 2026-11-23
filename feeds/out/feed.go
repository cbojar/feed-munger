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

func (feed *Feed) AddAllEntries(entries []*etree.Element) {
	for _, entry := range entries {
		feed.AddEntry(entry)
	}
}

func (feed *Feed) AddEntry(entry *etree.Element) {
	feed.document.Root().AddChild(entry)
}

func (feed *Feed) Bytes() ([]byte, error) {
	feed.document.IndentTabs()
	return feed.document.WriteToBytes()
}
