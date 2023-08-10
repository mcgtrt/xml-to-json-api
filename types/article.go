package types

import "encoding/xml"

type Article struct {
	ID                string   `bson:"_id,omitempty" json:"id,omitempty"`
	XMLName           xml.Name `xml:"NewsletterNewsItem"`
	ArticleURL        string   `xml:"ArticleURL" bson:"articleURL" json:"articleURL"`
	NewsArticleID     int      `xml:"NewsArticleID" bson:"newsArticleID" json:"newsArticleID"`
	PublishDate       string   `xml:"PublishDate" bson:"publishDate" json:"publishDate"`
	Taxonomies        string   `xml:"Taxonomies" bson:"taxonomies" json:"taxonomies"`
	TeaserText        string   `xml:"TeaserText" bson:"teaserText" json:"teaserText"`
	ThumbnailImageURL string   `xml:"ThumbnailImageURL" bson:"thumbnailImageURL" json:"thumbnailImageURL"`
	Title             string   `xml:"Title" bson:"title" json:"title"`
	OptaMatchId       string   `xml:"OptaMatchId" bson:"optaMatchId" json:"optaMatchId"`
	LastUpdateDate    string   `xml:"LastUpdateDate" bson:"lastUpdateDate" json:"lastUpdateDate"`
	IsPublished       bool     `xml:"IsPublished" bson:"isPublished" json:"isPublished"`
}

type ArticleList struct {
	XMLName xml.Name  `xml:"NewsletterNewsItems"`
	Values  []Article `xml:"NewsletterNewsItem"`
}

// Top level XML class. Skipped other XML values for this level
// of nesting according to the task requirements.
type NewListInformation struct {
	XMLName     xml.Name    `xml:"NewListInformation"`
	ArticleList ArticleList `xml:"NewsletterNewsItems"`
}
