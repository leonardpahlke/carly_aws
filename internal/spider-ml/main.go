package main

import (
	"carly_aws/internal/spider-ml/handler"
	"carly_aws/pkg"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(event pkg.SpiderMLEvent) (pkg.SpiderMLResponse, error) {
	return spider_ml.Handler(event)
}

func main() {
	lambda.Start(Handler)
}
