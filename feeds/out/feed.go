package out

import (
	"time"

	"github.com/beevik/etree"
	c "github.com/cbojar/feed-munger/configuration"
	"github.com/cbojar/feed-munger/feeds"
	"github.com/cbojar/feed-munger/feeds/in"
)

type Feed struct {
	document *etree.Document
}

func NewFeed(config *c.Configuration, firstDocument *in.Feed) *Feed {
	document := etree.NewDocument()
	document.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	root := document.CreateElement("feed")
	root.CreateAttr("xmlns", "http://www.w3.org/2005/Atom")

	id := root.CreateElement("id")
	if config.ID() != "" {
		id.SetText(config.ID())
	} else if firstDocumentID, found := firstDocument.ID(); found {
		id.SetText(firstDocumentID)
	}

	title := root.CreateElement("title")
	if config.Name() != "" {
		title.SetText(config.Name())
	} else if firstDocumentTitle, found := firstDocument.Title(); found {
		title.SetText(firstDocumentTitle)
	}

	alternate := root.CreateElement("link")
	alternate.CreateAttr("rel", "alternate")
	alternate.CreateAttr("type", "text/html")
	if config.Alternate() != "" {
		alternate.CreateAttr("href", config.Alternate())
	} else if firstDocumentAlternate, found := firstDocument.Alternate(); found {
		alternate.CreateAttr("href", firstDocumentAlternate)
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
