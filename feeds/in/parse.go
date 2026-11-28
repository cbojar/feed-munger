package in

import (
	"io"

	"github.com/beevik/etree"
)

func ParseAll(datas []io.ReadCloser) (*Feeds, error) {
	documents := &Feeds{feeds: make([]*Feed, 0, len(datas))}

	for _, data := range datas {
		if document, err := parse(data); err != nil {
			return nil, err
		} else {
			documents.append(document)
		}
	}

	return documents, nil
}

func parse(data io.ReadCloser) (*Feed, error) {
	document := etree.NewDocument()
	if _, err := document.ReadFrom(data); err != nil {
		return nil, err
	} else if document.Root().Tag == "rss" {
		return NewFeed(&RSSFormat{document: document}), nil
	} else if document.Root().Tag == "feed" {
		return NewFeed(&AtomFormat{document: document}), nil
	} else {
		panic("Unknown feed type")
	}
}
