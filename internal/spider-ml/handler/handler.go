package spider_ml

import (
	"carly_aws/pkg"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/*

1. analyze entities, key phrases and sentiment 					(comprehend)
2. store document in s3 bucket 									(s3)

*/

/* Article Bucket Analytics Files:
FOLDER:<newspaper>
	FOLDER:<ARTICLE-REFERENCE>
		FILE:<LANGUAGE-CODE> (de or en)
		FILE:comprehend
		FILE:DOM
*/

func Handler(request pkg.SpiderMLEvent) (pkg.SpiderMLResponse, error) {
	rfc5646English := pkg.RFC_5646_ENGLISH

	spiderName, _ := pkg.CheckEnvNotEmpty(pkg.EnvSpiderName)
	bucketNameAnalytics, _ := pkg.CheckEnvNotEmpty(pkg.EnvArticleBucketAnalytics)
	mySession := session.Must(session.NewSession())

	// Create a Comprehend client - to analyse sentiment, entities, key phrases,
	clientComprehend := comprehend.New(mySession)

	/*
		1. --- analyze entities, key phrases and sentiment
	*/

	detectedKeyPhrases, err := clientComprehend.DetectKeyPhrases(&comprehend.DetectKeyPhrasesInput{
		LanguageCode: &rfc5646English,
		Text:         &request.ArticleText,
	})
	if err != nil {
		pkg.LogError(spiderName, "clientComprehend.DetectKeyPhrases error", err)
	}
	pkg.LogInfo(spiderName, "detectedKeyPhrases complete..")

	detectedEntites, err := clientComprehend.DetectEntities(&comprehend.DetectEntitiesInput{
		LanguageCode: &rfc5646English,
		Text:         &request.ArticleText, //todo - detectedEntites
	})
	if err != nil {
		pkg.LogError(spiderName, "clientComprehend.DetectEntities error", err)
	}
	pkg.LogInfo(spiderName, "detectedEntites complete..")

	var sentimentBySentence []pkg.BucketAnalytics_COMPREHEND_sentiment
	splittedTextBySentence := strings.Split(request.ArticleText, ".")
	splittedTextBySentence = pkg.TrimStringAry(splittedTextBySentence)

	for _, sentence := range splittedTextBySentence {
		detectedSentiment, err := clientComprehend.DetectSentiment(&comprehend.DetectSentimentInput{
			LanguageCode: &rfc5646English,
			Text:         &sentence,
		})
		if err != nil {
			pkg.LogError(spiderName, "clientComprehend.DetectSentiment error", err)
		}
		sentimentBySentence = append(sentimentBySentence, pkg.ConvBucketAnalytics_COMPREHEND_sentiment(sentence, detectedSentiment))
	}
	pkg.LogInfo(spiderName, fmt.Sprintf("detectedSentiment complete.. %v", sentimentBySentence))

	/*
		2. --- store document in s3 bucket
	*/

	spiderMLComprehendDocumentJsonByteArray := pkg.MarshalStruct(pkg.BucketAnalytics_COMPREHEND{
		KeyPhrases:       detectedKeyPhrases.KeyPhrases,
		Entities:         detectedEntites.Entities,
		Sentiment:        sentimentBySentence,
		ArticleReference: request.ArticleReference,
	})

	// Create a S3 client - to store the article text in the s3 bucket
	uploader := s3manager.NewUploader(mySession)
	documentS3Key := pkg.GetBucketKeyForAnalyticsBucket(request.Newspaper, request.ArticleReference, "comprehend", "json")

	// store article comprehend analytics results in s3 bucket
	articleUploadResponse, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &bucketNameAnalytics,
		Key:    &documentS3Key,
		Body:   strings.NewReader(string(spiderMLComprehendDocumentJsonByteArray)),
	})
	if err != nil {
		pkg.LogError(spiderName, "s3 upload error", err)
	}

	return pkg.SpiderMLResponse{
		ArticleReference: request.ArticleReference,
		Newspaper:        request.Newspaper,
		S3ArticleFileUrl: articleUploadResponse.Location,
	}, nil
}
