package pkg

import "github.com/aws/aws-sdk-go/service/comprehend"

/*
	This file contains the structure of the S3 Bucket - Bucket-Article-Analytics-Store
*/

// TOP STRUCTURES

// document structure stored in s3 bucket
type BucketAnalytics_TEXT struct {
	ArticleReference string `json:"article_reference"`
	ArticleText      string `json:"article_text"`
	Language         string `json:"language"`
	Newspaper        string `json:"newspaper"`
}

// document to store comprehend analytics results
type BucketAnalytics_COMPREHEND struct {
	KeyPhrases       []*comprehend.KeyPhrase                `json:"key_phrases"`
	Entities         []*comprehend.Entity                   `json:"entities"`
	Sentiment        []BucketAnalytics_COMPREHEND_sentiment `json:"sentiment"`
	ArticleReference string                                 `json:"article_reference"`
}

// SUB STRUCTURES

// struct that describes how the stored sentiment results are structured
type BucketAnalytics_COMPREHEND_sentiment struct {
	Sentiment      string
	SentimentScore bucketAnalytics_COMPREHEND_sentiment_scoredetails
	Sentence       string
}

type bucketAnalytics_COMPREHEND_sentiment_scoredetails struct {
	Mixed    float64
	Negative float64
	Neutral  float64
	Positive float64
}

// SUPPORT METHODS

// convert comprehend seniment analysis output to internal format
func ConvBucketAnalytics_COMPREHEND_sentiment(sentence string, detectedSentimentOut *comprehend.DetectSentimentOutput) BucketAnalytics_COMPREHEND_sentiment {
	return BucketAnalytics_COMPREHEND_sentiment{
		Sentence:  sentence,
		Sentiment: *detectedSentimentOut.Sentiment,
		SentimentScore: bucketAnalytics_COMPREHEND_sentiment_scoredetails{
			Mixed:    *detectedSentimentOut.SentimentScore.Mixed,
			Negative: *detectedSentimentOut.SentimentScore.Negative,
			Neutral:  *detectedSentimentOut.SentimentScore.Neutral,
			Positive: *detectedSentimentOut.SentimentScore.Positive,
		},
	}
}
