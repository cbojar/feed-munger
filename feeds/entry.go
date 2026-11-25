package feeds

import "github.com/cbojar/feed-munger/feeds/entries"

type Entry = entries.Entry
type EntryBuilder = entries.EntryBuilder
type Category = entries.Category
type Link = entries.Link
type Person = entries.Person

func BuildEntry(build func(*EntryBuilder)) *Entry {
	return entries.BuildEntry(build)
}
