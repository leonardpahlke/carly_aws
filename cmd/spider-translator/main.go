package main

import (
	spider_translator "carly_aws/internal/spider-translator/handler"
	"carly_aws/pkg"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const localTestName = "LocalTest-" + pkg.SpiderNameTranslator

func main() {
	_ = os.Setenv("AWS_REGION", pkg.AWSDeployRegion)
	_ = os.Setenv(pkg.EnvSpiderName, pkg.SpiderNameTranslator)
	_ = os.Setenv(pkg.EnvLogLevel, strconv.Itoa(int(log.InfoLevel)))
	_ = os.Setenv(pkg.EnvArticleBucketAnalytics, "carly-dev-bucket-article-analytics-store")

	pkg.LogInfo(localTestName, "Starting local translaor test")

	spider_translator.Handler(pkg.SpiderTranslatorEvent{
		ArticleReference: "german-article-test",
		Newspaper:        pkg.NewspaperTESTING,
		ArticleText:      testArticle,
	})

}

const testArticle = `
Renaissance der Samstags-Konferenz Auf einer Wellenlänge Pandemie ist nicht nur schlimm: Mit der samstäglichen Radio-Konferenz kann man sich sogar schlimmen Standfußball schönhören.
 Ich pack das nicht, ich halt das nicht mehr aus, ich will das nicht mehr sehen.
“ Günther Kochs Stimme überschlägt sich.
 Es ist der 29.
 Mai 1999 und die Bundesliga erlebt einen ihrer spannendsten letzten Spieltage seit ihrer Gründung.
 Insgesamt fünf Mannschaften können in diesen zwei Stunden noch absteigen.
 In der ARD-Radiokonferenz dürften so viele Zuhörer hängen wie lange nicht.
 Hätten wir in diesen letzten Minuten vor Spielende in einem fahrenden Auto gesessen, wir hätten angehalten.
 Zugegeben, am aktuellen Spieltag war weit weniger Dynamik in der 15.
30-Uhr-Schalte als in der Frisur von Günter Netzer, aber selbst diese spärlichen Ereignisse hält den/die Ra­dio­re­por­te­r*in nicht davon ab, mittels energischer Zwischenrufe kleine Momente der Spannung zu erzeugen.
 Dass die ein oder andere Szene größer aufgeblasen wird, als sie tatsächlich in der Nachschau war, das muss verziehen sein.
 Das Radio lebt von Bildern, die mit Worten erzeugt werden.
 Man will doch bitte hochschrecken, wenn Julia Metzners Stimme sich plötzlich aus Stuttgart meldet.
 Wer nicht minütlich an die Dramatik des Fußballs erinnert werden will, kann ja einen Liveticker lesen.
 Der Reporter (früher waren da ja tatsächlich ausschließlich Männer) ist das letzte Bindeglied zwischen technisiertem Hochleistungsfußball und der Erinnerung an eine Kindheit vor dem Gartenradio.
 Er hat seine großen Momente („Aus dem Hintergrund müsste Rahn schießen …“) vielleicht hinter sich gelassen, erlebt aber gerade wieder Zulauf.
 Umrundet die Radiokonferenz doch jedes Interessenspektrum.
 Sei es mit halbem Interesse als Hintergrundgeräusch aus dem Nebenzimmer oder als Ersatz für den glühenden Anhänger, der aktuell nicht ins Stadion kann.
 Möglich, dass genau dieser Punkt den Reiz ausmacht.
 Das Erlebnis Fußball, wie wir es gewohnt sind, hat sich auf das sportliche Ereignis reduziert.
 Der Glamour ist verschwunden, aber eben auch das umfangreiche soziale Gefüge, das den Fußball umgibt, stillgelegt.
 Ganze Samstagabende hat der Freundeskreis einst zur Spielanalyse genutzt, nicht selten stand der ganze Tag im Zeichen des Fußballs.
 Grob fahrlässiges Mittagsbetrinken und für einen Tag Experte sein.
 Doch jetzt findet das Zuschauererlebnis vor allem zu Hause statt.
 Wie der Bundesliga-Wahnsinn konsumiert wird, das bleibt jedem selbst überlassen – und da König Fußball gar nicht mehr so sehr im Fokus steht, wie die Führungsetage einiger Klubs es gerne hätte, reicht manchmal nur das Küchenradio.
 Oder das Autoradio, das Handyradio beim Spazieren, das Radio am Arbeitsplatz.
 Die Verlässlichkeit, dass jeden Samstag zur selben Uhrzeit die Konferenz beginnt, die einzelnen Re­por­ter*in­nen sich knapp vorstellen und 90 Minuten keinen Informationsteppich über ein Spiel legen, sondern die wichtigen Spielszenen im Blick haben, ist geradezu wohltuend.
 Radio ist Ergebnis.
 Vielleicht will ich nicht wissen, ob der Innenverteidiger zum zweiten Mal Vater geworden ist, sondern ob er hoch genug springt, wenn er muss.
 Diese Szene möchte ich knackig erklärt bekommen.
 Danach zurück nach Bremen, wo 90 Minuten lang Standfußball gegen Freiburg gespielt wird.
 Dieses Spiel beispielsweise hat an diesem vergangenen Bundesliga-Samstag mit der Tatsache gewonnen, dass das Radio keine Bilder senden kann.
 „Radio, Radio, Radio – das schnellste Medium der Welt.
“ Was Günther Koch völlig außer Atem ins Mikro haucht, als sich am letzten Spieltag 1999 die Eintracht und Hansa Rostock in den Klassenerhalt retten, mag vor den heutigen technischen Leistungen und Möglichkeiten am Mitleid rütteln, aber vermutlich gilt diese Aussage immer noch.
 Einzig um die Bilder, wie Hansas Slawomir Majak nach Abpfiff nur noch mit Unterhose bekleidet Richtung Kurve jubelt, beneiden wir alle, die das Spiel doch live sehen konnten.
 .
`
