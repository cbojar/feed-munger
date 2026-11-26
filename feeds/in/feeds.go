package in

import (
	"slices"

	"github.com/cbojar/feed-munger/feeds"
)

type Feeds struct {
	feeds []*Feed
}

func (in *Feeds) Masthead() *feeds.Masthead {
	if len(in.feeds) == 0 {
		return nil
	}

	first := in.feeds[0]

	return &feeds.Masthead{
		ID:          first.ID,
		Title:       first.Title,
		Alternate:   first.Alternate,
		Self:        first.Self,
		Authors:     in.authors(),
		Categories:  in.categories(),
		Copyright:   first.Copyright,
		Description: first.Description,
		Image:       first.Image,
		Language:    first.Language,
		Links:       in.links(),
		Published:   first.Published,
		Updated:     first.Updated,
	}
}

func (in *Feeds) authors() []*feeds.Person {
	return allOf(in, func(feed *Feed) []*feeds.Person { return feed.Authors })
}

func (in *Feeds) categories() []*feeds.Category {
	return allOf(in, func(feed *Feed) []*feeds.Category { return feed.Categories })
}

func (in *Feeds) links() []*feeds.Link {
	return allOf(in, func(feed *Feed) []*feeds.Link { return feed.Links })
}

func allOf[T any](in *Feeds, these func(*Feed) []T) []T {
	all := []T{}

	for _, feed := range in.feeds {
		all = append(all, these(feed)...)
	}

	return all
}

func (in *Feeds) All() []*Feed {
	return in.feeds
}

func (in *Feeds) Namespaces() feeds.Namespaces {
	namespaces := feeds.Namespaces{}

	for _, feed := range in.feeds {
		namespaces = append(namespaces, feed.Namespaces...)
	}

	return namespaces
}

func (in *Feeds) Entries() []*feeds.Entry {
	entries := []*feeds.Entry{}
	for _, feed := range in.feeds {
		entries = append(entries, feed.Entries...)
	}

	slices.SortFunc(entries, func(entry1 *feeds.Entry, entry2 *feeds.Entry) int {
		return entry2.CompareUpdated(entry1)
	})

	return entries
}

func (in *Feeds) append(feed *Feed) {
	in.feeds = append(in.feeds, feed)
}
