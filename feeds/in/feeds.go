package infeeds

import (
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

func (in *Feeds) append(feed *Feed) {
	in.feeds = append(in.feeds, feed)
}
