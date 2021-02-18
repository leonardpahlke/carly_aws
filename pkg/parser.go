package pkg

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"strings"
)

// handler is a simple function that takes a string and does a ToUpper.
func ParseArticle(request ArticleParseRequest) (ArticleParseResponse, error) {
	methodName := "pkg.ParserArticle"

	tag := ""
	enter := false

	article := ""

	tokenizer := html.NewTokenizer(strings.NewReader(request.ArticleHtmlDom))
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tokenType {
		case html.ErrorToken:
			LogError(methodName, "Token error", err)

		case html.StartTagToken, html.SelfClosingTagToken:
			enter = false
			tagAttribute := token.Attr
			tag = token.Data
			for _, textTag := range request.TextTagsToParse {
				if tag == textTag && !CheckElementAttribute(request.WhitelistAttributeValues, tagAttribute) {
					enter = true
					break
				}
			}

		case html.TextToken:
			if enter {
				data := strings.TrimSpace(token.Data)
				if len(data) > 0 {
					article += fmt.Sprintf("%s\n", data)
				}
			}
		}
	}

	return ArticleParseResponse{
		ArticleText: article,
	}, nil
}


func ParseDomElement(htmlDom string, element string, requiredAttribute []html.Attribute) (string, error) {
	doc, _ := html.Parse(strings.NewReader(htmlDom))
	bn, err := getElement(doc, element, requiredAttribute)
	if err != nil {
		return "", err
	}
	log.Info(bn.Attr)
	return renderNode(bn), nil
}


func getElement(doc *html.Node, element string, requiredAttribute []html.Attribute) (*html.Node, error) {
	var elementNode *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == element && CheckElementAttribute(requiredAttribute, node.Attr) {
			elementNode = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if elementNode != nil {
		return elementNode, nil
	}
	return nil, errors.New(fmt.Sprintf("Missing <%s> in the node tree", element))
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	_ = html.Render(w, n)
	return buf.String()
}

func BundleSentences(articleText string) string {
	updatedArticleText := ""
	splitedByDot := strings.Split(articleText, ".")
	for _, sentence := range splitedByDot {
		updatedArticleText += strings.ReplaceAll(sentence, "\n", " ") + ".\n"
	}
	return updatedArticleText
}

func CheckElementAttribute(attrToCheck []html.Attribute, elementAttr []html.Attribute) bool {
	for _, checkAttr := range attrToCheck {
		for _, elemAttr := range elementAttr {
			if checkAttr.Key == elemAttr.Key && checkAttr.Val == elemAttr.Val {
				return true
			}
		}
	}
	return false
}

type ArticleParseRequest struct {
	ArticleHtmlDom           string
	TextTagsToParse          []string
	WhitelistAttributeValues []html.Attribute
}

type ArticleParseResponse struct {
	ArticleText string
}
