package spider_translator

import (
	"carly_aws/pkg"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/translate"
)

/*
	1. check if document is german 									(comprehend)
	2. create article json document to store it in the s3-bucket 	(s3)
	3. store document in s3 bucket 									(s3)
	4. if article is in german translate it to english and store it (translate, s3)
*/

/* Article Bucket Analytics Files:
FOLDER:<newspaper>
	FOLDER:<ARTICLE-REFERENCE>
		FILE:<LANGUAGE-CODE> (de or en)
		FILE:comprehend
		FILE:DOM
*/

func Handler(request pkg.SpiderTranslatorEvent) (pkg.SpiderTranslatorResponse, error) {
	rfc5646German := "de"
	rfc5646English := "en"

	spiderName, _ := pkg.CheckEnvNotEmpty(pkg.EnvSpiderName)
	bucketNameAnalytics, _ := pkg.CheckEnvNotEmpty(pkg.EnvArticleBucketAnalytics)
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
		pkg.LogError(spiderName, "Language detected is not supported", errors.New("language code not supported"))
	}

	/*
		2. --- create json article document of the article to store it
	*/

	// create json document - SpiderMLTextDocument
	spiderMLTextDocumentJsonByteArray := pkg.MarshalStruct(pkg.BucketAnalytics_TEXT{
		ArticleReference: request.ArticleReference,
		ArticleText:      request.ArticleText,
		Language:         detectedLanguage,
		Newspaper:        request.Newspaper,
	})
	strSpiderMLTextDocumentJsonByteArray := string(spiderMLTextDocumentJsonByteArray)

	/*
		3. --- store article document in s3 bucket
	*/

	// Create a S3 client - to store the article text in the s3 bucket
	uploader := s3manager.NewUploader(mySession)
	documentUploadS3Key := pkg.GetBucketKeyForAnalyticsBucket(request.Newspaper, request.ArticleReference, detectedLanguage, "json")

	// store article text in s3 bucket
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: &bucketNameAnalytics,
		Key:    &documentUploadS3Key,
		Body:   strings.NewReader(strSpiderMLTextDocumentJsonByteArray),
	})
	if err != nil {
		pkg.LogError(spiderName, "s3 upload error", err)
	}

	/*
		4. --- if article is in german translate it to english and store it
	*/

	if *dominantLanguage.Languages[0].LanguageCode == rfc5646German {
		pkg.LogInfo(spiderName, "Language code is german, translate document into english")

		// Create a Translate client - to translate the text german -> english
		clientTranslate := translate.New(mySession)
		targetTranslatedBucketKey := pkg.GetBucketFileName(bucketNameAnalytics, request.Newspaper, rfc5646English, "json")

		// translate text to english
		textTranslationResult, err := clientTranslate.Text(&translate.TextInput{
			SourceLanguageCode: &rfc5646German,
			TargetLanguageCode: &rfc5646English,
			Text:               &strSpiderMLTextDocumentJsonByteArray,
		})
		if err != nil {
			pkg.LogError(spiderName, "clientTranslate.Text error", err)
		}

		// format translated text
		translatedText := *textTranslationResult.TranslatedText
		translatedText = strings.ReplaceAll(translatedText, "\\ n", "")
		translatedText = strings.ReplaceAll(translatedText, "  ", " ")

		pkg.LogInfo(spiderName, fmt.Sprintf("Text translated.. %s", translatedText))

		// upload translated text to s3
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: &bucketNameAnalytics,
			Key:    &targetTranslatedBucketKey,
			Body:   strings.NewReader(translatedText),
		})
		if err != nil {
			pkg.LogError(spiderName, "s3 upload error", err)
		}

		return pkg.SpiderTranslatorResponse{
			ArticleReference: request.ArticleReference,
			Newspaper:        request.Newspaper,
			ArticleText:      translatedText,
		}, nil
	}
	return pkg.SpiderTranslatorResponse{
		ArticleReference: request.ArticleReference,
		Newspaper:        request.Newspaper,
		ArticleText:      request.ArticleText,
	}, nil
}

/* DOWNLOAD S3 File example

// create downlaoder manager to get translated article text
		downloader := s3manager.NewDownloader(mySession)

		// Download the item from the bucket. If an error occurs, log it and exit.
		// 		Otherwise, notify the user that the download succeeded.
		fileName := targetTranslatedBucketKey
		file, err := os.Create(fileName)
		if err != nil {
			pkg.LogError(spiderName, "Unable to create local file", err)
		}

		// Download s3 file and save it in the created local file
		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: &bucketNameAnalytics,
				Key:    &targetTranslatedBucketKey,
			})
		if err != nil {
			pkg.LogError(spiderName, fmt.Sprintf("Unable to download item: %q, %v", targetTranslatedBucketKey, err), err)
		}

		pkg.LogInfo(spiderName, fmt.Sprintf("Downloaded %s, %x, bytes", file.Name(), numBytes))

		// read downloaded file
		fileData, err := ioutil.ReadFile(fileName)
		if err != nil {
			pkg.LogError(spiderName, "ioutil.ReadFile unable to read file", err)
		}

		// decode payload to
		var downloadedFile pkg.BucketAnalytics_TEXT

		err = json.Unmarshal(fileData, &downloadedFile)
		if err != nil {
			pkg.LogError(spiderName, "json.Unmarshal unmarshal of translated file failed", err)
		}
*/
