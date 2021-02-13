package main

import (
	"carly_aws/internal/spider-taz-parser/handler"
	"carly_aws/pkg"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(event pkg.SpiderParserEvent) (pkg.SpiderParserResponse, error) {
	return spider_taz_parser.Handler(event)
}

func main() {
	lambda.Start(Handler)
}
