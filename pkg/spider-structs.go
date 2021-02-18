package pkg

const EnvSpiderName = "NAME"
const EnvArticleBucket = "ARTICLE_BUCKET"
const EnvFilePrefix = "FILE_PREFIX"
const EnvLogLevel = "LOG_LEVEL"

const SpiderNameParser = "SpiderParser"
const SpiderNameMl = "SpiderMl"
const SpiderNameDownloader = "SpiderDownloader"

// Downloader
type SpiderDownloaderEvent struct {
	ArticleReference string `json:"article_reference"`
	ArticleUrl       string `json:"article_url"`
	Newspaper        string `json:"newspaper"`
}
type SpiderDownloaderResponse struct {
	ArticleDom 		 string `json:"article_dom"`
	S3ArticleDomLink string `json:"s_3_article_dom_link"`
	ArticleReference string `json:"article_reference"`
	ArticleUrl       string `json:"article_url"`
	Newspaper        string `json:"newspaper"`
}

// Parser
type SpiderParserEvent struct {
	ArticleDom 		 string `json:"article_dom"`
	ArticleReference string `json:"article_reference"`
	S3ArticleDomLink string `json:"s_3_article_dom_link"`
	Newspaper        string `json:"newspaper"`
}
type SpiderParserResponse struct {
	ArticleReference  string            `json:"article_reference"`
	S3ArticleDomLink  string            `json:"s_3_article_dom_link"`
	Newspaper         string            `json:"newspaper"`

	ArticleText       string            `json:"article_text"`
	// ArticleAttributes map[string]string `json:"article_attributes"`
}

// ML
type SpiderMLEvent struct {
	ArticleReference  string            `json:"article_reference"`
	Newspaper         string            `json:"newspaper"`
	ArticleText       string            `json:"article_text"`
	ArticleAttributes map[string]string `json:"article_attributes"`
}
type SpiderMLResponse struct {
	ArticleReference  string            `json:"article_reference"`
	Newspaper         string            `json:"newspaper"`
	ArticleText       string            `json:"article_text"`
	ArticleAttributes map[string]string `json:"article_attributes"`
}
