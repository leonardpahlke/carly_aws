package main

import (
	spider_downloader "carly_aws/internal/spider-downloader"
	"carly_aws/pkg"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const NewspaperNameTAZ = "taz"
const NewspaperNameZeitOnline = "zeitonline"
const NewspaperNameFrankfurterRundschau = "frankfurterrundschau"

const name = "Main" + pkg.EnvSpiderName

func main() {
	_ = os.Setenv("AWS_REGION", pkg.AWSDeployRegion)
	_ = os.Setenv(pkg.EnvSpiderName, "SpiderDownloader")
	_ = os.Setenv(pkg.EnvArticleBucket, "carly-dev-bucket-article-dom-store")
	assetsPath := "assets/"
	_ = os.Setenv(pkg.EnvFilePrefix, assetsPath)
	_ = os.Setenv(pkg.EnvLogLevel, strconv.Itoa(int(log.InfoLevel)))

	err := os.Mkdir(assetsPath, 0755)
	if err == nil {
		pkg.LogInfo(name, "folder created")
	}

	ok := reqSpiderDownloader(pkg.SpiderDownloaderEvent{
		ArticleReference: fmt.Sprintf("%s-5751873", NewspaperNameTAZ),
		ArticleUrl:       "https://taz.de/Neue-All-Parteien-Koalition-in-Italien/!5751873/",
		Newspaper:        NewspaperNameTAZ,
	}, NewspaperNameTAZ)
	if !ok {
		return
	}

	ok = reqSpiderDownloader(pkg.SpiderDownloaderEvent{
		ArticleReference: fmt.Sprintf("%s-202102-5751873", NewspaperNameZeitOnline),
		ArticleUrl:       "https://www.zeit.de/wirtschaft/2021-02/umweltverschmutzung-shell-nigeria-bauern-klage-oel-leck",
		Newspaper:        NewspaperNameZeitOnline,
	}, NewspaperNameTAZ)
	if !ok {
		return
	}

	ok = reqSpiderDownloader(pkg.SpiderDownloaderEvent{
		ArticleReference: fmt.Sprintf("%s-202102-90200529", NewspaperNameFrankfurterRundschau),
		ArticleUrl:       "https://www.fr.de/politik/corona-grenze-kontrolle-seehofer-oesterreich-tirol-tschechien-bayern-sachsen-seehofer-soeder-merkel-eu-zr-90200529.html",
		Newspaper:        NewspaperNameFrankfurterRundschau,
	}, NewspaperNameTAZ)
	if !ok {
		return
	}
}


func reqSpiderDownloader(event pkg.SpiderDownloaderEvent, newspaper string) bool {
	resp, err := spider_downloader.Main(event)
	if err != nil {
		pkg.LogError(name, fmt.Sprintf("%s-downloader processing error", newspaper), err)
		return false
	}
	pkg.LogInfo(name, fmt.Sprintf("%s-downloader resp: %s", newspaper, resp))
	return true
}