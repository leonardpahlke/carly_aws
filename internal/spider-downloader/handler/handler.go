package spider_downloader

import (
	"carly_aws/pkg"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

/*
Lambda SpiderDownloader
1. set up lambda clients
2. create file to store article html
3. retrieve article html from article website
4. store html to local file
5. upload file to s3
*/

// Handler - gets executed when lambda is run
func Handler(event pkg.SpiderDownloaderEvent) (pkg.SpiderDownloaderResponse, error) {
	// 1. set up
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	s3BucketName, _ := pkg.CheckEnvNotEmpty(pkg.EnvArticleBucket)
	spiderName, _ := pkg.CheckEnvNotEmpty(pkg.EnvSpiderName)
	pkg.SetLogLevel()

	fileName := fmt.Sprintf("%s.html", event.ArticleReference)
	fileFullPathName := "/tmp/" + fileName

	pkg.LogInfo(spiderName, fmt.Sprintf("setup %s", event))

	// 2. create file
	f, err := os.Create(fileFullPathName)
	if err != nil {
		pkg.LogError(spiderName, "create file error", err)
		return pkg.SpiderDownloaderResponse{}, err
	}

	// 3. get article DOM frm article url
	resp, err := http.Get(event.ArticleUrl)
	if err != nil {
		pkg.LogError(spiderName, "Website GET request error", err)
		return pkg.SpiderDownloaderResponse{}, err
	}
	defer resp.Body.Close()

	pkg.LogInfo(spiderName, "Website data retrieved")

	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	// 4. write html body bytes to file
	l, err := f.WriteString(string(html))
	if err != nil {
		pkg.LogError(spiderName, "write data to file error", err)
		return pkg.SpiderDownloaderResponse{}, err
	}

	err = f.Close()
	if err != nil {
		pkg.LogError(spiderName, "closing file error", err)
		return pkg.SpiderDownloaderResponse{}, err
	}
	f, err = os.Open(fileFullPathName)

	// 5. Upload file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", event.Newspaper, fileName)),
		Body:   f,
	})
	if err != nil {
		pkg.LogError(spiderName, "s3 upload error", err)
		return pkg.SpiderDownloaderResponse{}, err
	}

	// show the HTML code as a string %s
	log.Infof("Written bytes to file %x\n", l)

	return pkg.SpiderDownloaderResponse{
		ArticleReference: event.ArticleReference,
		ArticleUrl:       event.ArticleUrl,
		Newspaper:        event.Newspaper,
		S3ArticleDomLink: result.Location,
	}, nil
}
