package main

import (
	"carly_aws/internal/spider-downloader/handler"
	"carly_aws/pkg"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(event pkg.SpiderDownloaderEvent) (pkg.SpiderDownloaderResponse, error) {
	return spider_downloader.Handler(event)
}

func main() {
	lambda.Start(Handler)
}
