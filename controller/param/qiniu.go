package param

type GetQiniuTokenRequest struct {
	URL         string `json:"url"`
	ResourceKey string `json:"resource_key"`
	Token       string `json:"token"`
}

type QiniuAuthentication struct {
	Url         string ``
	Method      string
	Path        string
	RawQuery    string
	Host        string
	ContentType string
	BodyStr     string
}
