package spider_ml

import (
	"carly_aws/pkg"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/translate"
	"strings"
)

/*

1. check if document is german 									(comprehend)
2. create article json document to store it in the s3-bucket 	(s3)
3. store document in s3 bucket 									(s3)
4. if article is in german translate it to english and store it (translate, s3)
5. analyze entities, key phrases and sentiment 					(comprehend)
6. store document in s3 bucket 									(s3)

 */

/* Article Bucket Analytics Files:
FOLDER:<newspaper>
	FOLDER:<ARTICLE-REFERENCE>
		FILE:<LANGUAGE-CODE> (de or en)
		FILE:comprehend
 */

// handler is a simple function that takes a string and does a ToUpper.
func Handler(request pkg.SpiderMLEvent) (pkg.SpiderMLResponse, error) {
	rfc5646German := "de"
	rfc5646English := "en"

	spiderName, _ := pkg.CheckEnvNotEmpty(pkg.EnvSpiderName)
	bucketNameAnalytics, _ := pkg.CheckEnvNotEmpty(pkg.EnvArticleBucketAnalytics)
	spiderRoleArn, _ := pkg.CheckEnvNotEmpty(pkg.EnvSpiderRoleArn)
	mySession := session.Must(session.NewSession())

	/*
	1. --- check if document is german (comprehend)
	 */

	// Create a Comprehend client - to analyse sentiment, entities, key phrases,
	clientComprehend := comprehend.New(mySession)

	// check language of document
	dominantLanguage, err := clientComprehend.DetectDominantLanguage(&comprehend.DetectDominantLanguageInput{
		Text: &request.ArticleText,
	})
	if err != nil {
		pkg.LogError(spiderName, "clientComprehend.DetectDominantLanguage error", err)
	}

	// get language code of the most dominant language
	detectedLanguage := *dominantLanguage.Languages[0].LanguageCode
	pkg.LogInfo(spiderName, fmt.Sprintf("Language detected: %s", detectedLanguage))

	// if the detected language is not german or english -> error
	if detectedLanguage != rfc5646German && detectedLanguage != rfc5646English {
		pkg.LogError(spiderName, "Language detected is not supported", errors.New("Language code not supported"))
	}

	/*
	2. --- create json document of the article to store it
	 */

	// create json document - SpiderMLTextDocument
	spiderMLTextDocumentJsonByteArray := pkg.MarshalStruct(pkg.SpiderMLTextDocument{
		ArticleReference: request.ArticleReference,
		ArticleText:      request.ArticleText,
		Language:         detectedLanguage,
		Newspaper:        request.Newspaper,
	})

	/*
	3. --- store document in s3 bucket
	 */

	// Create a S3 client - to store the article text in the s3 bucket
	uploader := s3manager.NewUploader(mySession)

	// set meta data of this article in a separate struct to use it later on
	storeFileMetaStruct := storeFileMetaStruct{
		spiderName:       spiderName,
		bucketName:       bucketNameAnalytics,
		articleReference: request.ArticleReference,
		newspaper:        request.Newspaper,
		uploader:         *uploader,
	}

	// store article text in s3 bucket
	articleUploadResponse := storeFileInArticleAnalyticsBucket(storeFileStruct{
		filename:   detectedLanguage,
		fileEnding: "json",
		file:       string(spiderMLTextDocumentJsonByteArray),
	}, storeFileMetaStruct)


	/*
	4. --- if article is in german translate it to english and store it
	 */

	if *dominantLanguage.Languages[0].LanguageCode == rfc5646German {
		pkg.LogInfo(spiderName, "Language code is german, translate document into english")

		// Create a Translate client - to translate the text german -> english
		clientTranslate := translate.New(mySession)
		jsonContentType := "json"

		// create uri where to store the result document
		translateBucketUri := getBucketFileUri(bucketNameAnalytics, request.Newspaper, fmt.Sprintf("%s", rfc5646English), "json")

		textTranslationJob, err := clientTranslate.StartTextTranslationJob(&translate.StartTextTranslationJobInput{
			DataAccessRoleArn: &spiderRoleArn,
			InputDataConfig: &translate.InputDataConfig{
				ContentType: &jsonContentType,
				S3Uri:       &articleUploadResponse.Location,
			},
			JobName:             &request.ArticleReference,
			OutputDataConfig:    &translate.OutputDataConfig{
				S3Uri: &translateBucketUri,
			},
			SourceLanguageCode:  &rfc5646German,
			TargetLanguageCodes: []*string{&rfc5646English},
		})
		if err != nil {
			pkg.LogError(spiderName, "clientTranslate.StartTextTranslationJob error", err)
		}
		pkg.LogInfo(spiderName, fmt.Sprintf("TextTranslationJob status response: %s", *textTranslationJob.JobStatus))

	}

	/*
	5. --- analyze entities, key phrases and sentiment
	 */

	detectedKeyPhrases, err := clientComprehend.DetectKeyPhrases(&comprehend.DetectKeyPhrasesInput{
		LanguageCode: &rfc5646English,
		Text:         nil,  //todo - detectedKeyPhrases
	})
	if err != nil {
		pkg.LogError(spiderName, "clientComprehend.DetectKeyPhrases error", err)
	}

	detectedEntites, err := clientComprehend.DetectEntities(&comprehend.DetectEntitiesInput{
		LanguageCode: &rfc5646English,
		Text:         nil,  //todo - detectedEntites
	})
	if err != nil {
		pkg.LogError(spiderName, "clientComprehend.DetectEntities error", err)
	}

	detectedSentiment, err := clientComprehend.DetectSentiment(&comprehend.DetectSentimentInput{
		LanguageCode: &rfc5646English,
		Text:         nil,  //todo - detectedSentiment
	})
	if err != nil {
		pkg.LogError(spiderName, "clientComprehend.DetectSentiment error", err)
	}

	/*
	6. --- store document in s3 bucket
	 */

	spiderMLComprehendDocumentJsonByteArray := pkg.MarshalStruct(pkg.SpiderMLComprehendDocument{
		KeyPhrases:     detectedKeyPhrases.KeyPhrases,
		Entities:       detectedEntites.Entities,
		Sentiment:      detectedSentiment.Sentiment,
		SentimentScore: detectedSentiment.SentimentScore,
	})

	// store article comprehend analytics results in s3 bucket
	articleComprehendUploadResponse := storeFileInArticleAnalyticsBucket(storeFileStruct{
		filename:   "comprehend",
		fileEnding: "json",
		file:       string(spiderMLComprehendDocumentJsonByteArray),
	}, storeFileMetaStruct)

	return pkg.SpiderMLResponse{
		ArticleReference: request.ArticleReference,
		Newspaper:        request.Newspaper,
		// todo check other response parameters
	}, nil
}

// store file in article analytics bucket
func storeFileInArticleAnalyticsBucket(fileInfoIn storeFileStruct, fileMetaIn storeFileMetaStruct) *s3manager.UploadOutput {
	result, err := fileMetaIn.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(fileMetaIn.bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s/%s.%s",
			fileMetaIn.newspaper,
			fileMetaIn.articleReference,
			fileInfoIn.filename,
			fileInfoIn.fileEnding)),
		Body:   strings.NewReader(fileInfoIn.file),
	})
	if err != nil {
		pkg.LogError(fileMetaIn.spiderName, "s3 upload error", err)
		return &s3manager.UploadOutput{}
	}
	return result
}

func getBucketFileUri(bucketName string, newspaper string, filename string, fileEnding string) string {
	return fmt.Sprintf("s3://%s/%s/%s.%s", bucketName, newspaper, filename, fileEnding)
}

type storeFileMetaStruct struct {
	spiderName string
	bucketName string
	articleReference string
	newspaper string
	uploader s3manager.Uploader
}

type storeFileStruct struct {
	filename string
	fileEnding string
	file string
}