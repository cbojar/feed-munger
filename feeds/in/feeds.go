package infeeds

import (
	"slices"
	"strings"

	"github.com/beevik/etree"
	"github.com/cbojar/feed-munger/feeds"
)

type Feeds struct {
	feeds []*Feed
}

func (in *Feeds) First() *Feed {
	if len(in.feeds) > 0 {
		return in.feeds[0]
	} else {
		return nil
	}
}

func (in *Feeds) All() []*Feed {
	return in.feeds
}

func (in *Feeds) Namespaces() feeds.Namespaces {
	namespaces := feeds.Namespaces{}

	for _, feed := range in.feeds {
		namespaces = append(namespaces, feed.Namespaces()...)
	}

	return namespaces
}

func (in *Feeds) Entries() []*etree.Element {
	entries := []*etree.Element{}
	for _, feed := range in.feeds {
		entries = append(entries, feed.Entries()...)
	}

	slices.SortFunc(entries, func(entry1 *etree.Element, entry2 *etree.Element) int {
		date1 := entry1.SelectElement("updated").Text()
		date2 := entry2.SelectElement("updated").Text()

		return strings.Compare(date2, date1)
	})

	return entries
}

func (in *Feeds) append(feed *Feed) {
	in.feeds = append(in.feeds, feed)
}
