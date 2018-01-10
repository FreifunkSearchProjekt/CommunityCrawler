package utils

type Transaction struct {
	BasicWebpages []WebpageBasic `json:"basic_webpages"`
	RssFeed       []FeedBasic    `json:"rss_feeds"`
}

type WebpageBasic struct {
	URL         string `json:"url"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Description string `json:"description"`
}

type FeedBasic struct {
	URL         string `json:"url"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
