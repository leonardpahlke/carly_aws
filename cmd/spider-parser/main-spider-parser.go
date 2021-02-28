package main

import (
	spider_parser "carly_aws/internal/spider-parser/handler"
	"carly_aws/pkg"
	"math/rand"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const name = "LocalTest-" + pkg.SpiderNameParser

func main() {
	_ = os.Setenv("AWS_REGION", pkg.AWSDeployRegion)
	_ = os.Setenv(pkg.EnvSpiderName, pkg.SpiderNameParser)
	_ = os.Setenv(pkg.EnvLogLevel, strconv.Itoa(int(log.InfoLevel)))
	rand.Seed(time.Now().UnixNano())

	folderPath := fmt.Sprintf("tmp/%s/", time.Now().String())

	err := os.Mkdir(folderPath, 0755)
	if err == nil {
		pkg.LogInfo(name, "folder created")
	}

	testArticles := []string{
		pkg.NewspaperNameFrankfurterRundschau,
		// pkg.NewspaperNameZeitOnline,
		pkg.NewspaperNameTAZ,
	}

	testArticleUrls := map[string][]string{
		pkg.NewspaperNameFrankfurterRundschau: {
			"https://www.fr.de/kultur/tv-kino/markus-lanz-zdf-corona-talk-lockdown-regierung-kritik-pranger-tv-90209129.html",
		},
		pkg.NewspaperNameZeitOnline: {
			"https://www.zeit.de/wissen/gesundheit/2021-02/astrazeneca-corona-impfstoff-wirksamkeit-infektionsschutz-biontech-moderna",
		},
		pkg.NewspaperNameTAZ: {
			// "https://taz.de/Medienstreit-in-Australien/!5753122/",
			// "https://taz.de/Extreme-Kaelte-in-den-USA/!5747388/",
			// "https://taz.de/Biobranche-erzielt-Rekordumsatz/!5747387/",
			// "https://taz.de/Trainerwechsel-bei-Borussia-Dortmund/!5747257/",
			"https://taz.de/Renaissance-der-Samstags-Konferenz/!5746942/",
			// "https://taz.de/Waffenrecht-in-Deutschland/!5747097/",
		},
	}

	for _, testArticle := range testArticles {
		articleUrls := testArticleUrls[testArticle]
		if articleUrls == nil {
			pkg.LogWarning(name, "testArticle reference not found in testArticleUrls")
			return
		}
		testNewspaper(
			testArticle,
			articleUrls,
			true, true, folderPath)
	}
}


func testNewspaper(newspaper string, articleUrls []string, forceWriteArticleText bool, writeArticleDom bool, folderPath string) {
	for _, articleUrl := range articleUrls {
		bodyBytes := sendRequestToUrl(articleUrl)
		articleReference := fmt.Sprintf("%s-%x", newspaper, rand.Int())
		if writeArticleDom {
			writeTextToFile("txt", fmt.Sprintf("dom-%s", articleReference), string(bodyBytes), folderPath)
		}
		ok := reqSpiderParser(pkg.SpiderParserEvent{
			ArticleReference: articleReference,
			ArticleDom:       string(bodyBytes),
			Newspaper:        newspaper,
			S3ArticleDomLink: "dummy-link",
		}, forceWriteArticleText, folderPath)
		if !ok {
			return
		}
	}
}

func sendRequestToUrl(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	return bodyBytes
}

func reqSpiderParser(event pkg.SpiderParserEvent, forceStore bool, folderPath string) bool {
	resp, err := spider_parser.Handler(event)
	if err != nil {
		pkg.LogError(name, fmt.Sprintf("%s-parser processing error", event.Newspaper), err)
		return false
	}

	ok := checkArticleText(resp.ArticleText)
	pkg.LogInfo(name, fmt.Sprintf("Article {%s} parsed check {%v}", resp.ArticleReference, ok))

	if !ok || forceStore {
		return writeTextToFile("txt", resp.ArticleReference, resp.ArticleText, folderPath)
	}
	return true
}

func writeTextToFile(fileEnding string, fileName string, text string, folderPath string) bool {
	fileFullName := fmt.Sprintf("%s.%s", fileName, fileEnding)
	fileFullPathName := folderPath + fileFullName
	f, err := os.Create(fileFullPathName)
	if err != nil {
		pkg.LogError(name, "create file error", err)
		return false
	}

	_, err = f.WriteString(text)
	if err != nil {
		pkg.LogError(name, "write file error", err)
		return false
	}
	return true
}


func checkArticleText(articleText string) bool {
	checkChars := "<>"
	splittedArticleText := strings.Split(articleText, " ")
	articleOk := true
	for _, word := range splittedArticleText {
	 	if strings.ContainsAny(word, checkChars) {
	 		pkg.LogWarning(
	 			"cmd.SpiderParser.CheckArticleText",
	 			fmt.Sprintf("word {%s} contains any of these unwanted chars {%s}", word, checkChars))
	 		articleOk = false
		}
	 }
	 return articleOk
}