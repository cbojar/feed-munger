package infeeds

import (
	"strings"

	"github.com/beevik/etree"
	"github.com/cbojar/feed-munger/feeds"
)

type Feed struct {
	document *etree.Document
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

func (feed *Feed) ID() (string, bool) {
	if id := feed.document.FindElement("/feed/id"); id == nil {
		return "", false
	} else {
		return id.Text(), true
	}
}

func (feed *Feed) Title() (string, bool) {
	if title := feed.document.FindElement("/feed/title"); title == nil {
		return "", false
	} else {
		return title.Text(), true
	}
}

func (feed *Feed) Alternate() (string, bool) {
	if link := feed.document.FindElement("/feed/link[@rel='alternate']"); link == nil {
		return "", false
	} else {
		return link.SelectAttrValue("href", ""), true
	}
}

func (feed *Feed) Entries() []*etree.Element {
	feedTitle, titleFound := feed.Title()

	entries := feed.document.Root().SelectElements("entry")
	extracted := make([]*etree.Element, 0, len(entries))

	for _, entry := range entries {
		extractedEntry := entry.Copy()

		if titleFound {
			source := extractedEntry.CreateElement("source")
			sourceTitle := source.CreateElement("title")
			sourceTitle.SetText(feedTitle)
		}

		extracted = append(extracted, extractedEntry)
	}

	return extracted
}
