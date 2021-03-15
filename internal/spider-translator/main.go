package main

import (
	spider_translator "carly_aws/internal/spider-translator/handler"
	"carly_aws/pkg"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(event pkg.SpiderTranslatorEvent) (pkg.SpiderTranslatorResponse, error) {
	return spider_translator.Handler(event)
}

func main() {
	lambda.Start(Handler)
}
