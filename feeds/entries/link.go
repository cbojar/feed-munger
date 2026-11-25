package entries

import "net/url"

type Link struct {
	relationship string
	contentType  string
	hreflang     string
	href         *url.URL
}

func (link *Link) Relationship() string {
	return link.relationship
}

func (link *Link) ContentType() string {
	return link.contentType
}

func (link *Link) HrefLang() string {
	return link.hreflang
}

func (link *Link) Href() *url.URL {
	return link.href
}
