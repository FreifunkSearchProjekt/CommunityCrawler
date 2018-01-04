package common

type transaction struct {
	BasicWebpages []WebpageBasic `json:"basic_webpages"`
}

type WebpageBasic struct {
	URL         string `json:"url"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Description string `json:"description"`
}
