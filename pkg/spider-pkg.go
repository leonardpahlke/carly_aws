package pkg

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

const EnvSpiderName = "NAME"
const EnvArticleBucket = "ARTICLE_BUCKET"
const EnvFilePrefix = "FILE_PREFIX"
const EnvLogLevel = "LOG_LEVEL"

type SpiderDownloaderEvent struct {
	ArticleReference string
	ArticleUrl       string
	Newspaper        string
}

type SpiderDownloaderResponse struct {
	S3ArticleDomLink string
	ArticleReference string
	ArticleUrl string
	Newspaper string
}


func LogInfo(methodName string, logMessage string) {
	log.Infof("%s: %s \n", strings.ToUpper(methodName), logMessage)
}


func LogError(methodName string, logMessage string, err error) {
	log.Errorf("%s: %s \n  ERROR: %s \n", strings.ToUpper(methodName), logMessage, err)
}

func LogWarning(methodName string, logMessage string) {
	log.Warnf("%s: %s \n", strings.ToUpper(methodName), logMessage)
}


func SetLogLevel() {
	logLevel := os.Getenv(EnvLogLevel)
	if logLevel == "" {
		log.SetLevel(log.ErrorLevel)
	} else {
		convertedEnv, err := strconv.ParseInt(logLevel, 10, 64)
		if err != nil {
			LogWarning("SetLogLevel", "could not set log level")
			log.SetLevel(log.ErrorLevel)
		}
		log.SetLevel(log.Level(convertedEnv))
	}
}