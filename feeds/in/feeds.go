package in

import (
	"slices"

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

func (in *Feeds) Entries() []*feeds.Entry {
	entries := []*feeds.Entry{}
	for _, feed := range in.feeds {
		entries = append(entries, feed.Entries()...)
	}

	slices.SortFunc(entries, func(entry1 *feeds.Entry, entry2 *feeds.Entry) int {
		return entry2.CompareUpdated(entry1)
	})

	return entries
}

func (in *Feeds) append(feed *Feed) {
	in.feeds = append(in.feeds, feed)
}
