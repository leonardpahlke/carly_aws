package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	splittedTextBySentence := strings.Split(text, ".")
	log.Infof("splittedTextBySentence.. %v", splittedTextBySentence)
}

const text = `
The United States is willing to sit down with Iran “tomorrow” and jointly agree to full compliance with the nuclear accord they and five other world powers signed in 2015, according to a senior Biden administration official.
“We’ve made clear that we’re not talking about renegotiating the deal,” the official said of the agreement that curbed Iran’s nuclear program in exchange for lifting U.S. and other sanctions.
Iran has made equally clear it shares the goal of going back to the terms of the original agreement, before President Donald Trump pulled out of it. Trump reinstituted the sanctions and added what Biden officials estimate were at least 1,500 new ones. In response, Iran reactivated key elements of the program the United States and others say could produce nuclear weapons. Iran denies any such ambition.
But nearly two months into Biden’s presidency, with Iran’s own contentious presidential election approaching in June, the two sides have been unable even to talk to each other about what both say they want.
Biden has vowed to quickly restore the Iran nuclear deal, but that may be easier said than done
There was a near miss more than three weeks ago, when the administration said it would attend a meeting called by the European Union with Iran and the other original signatories still party to the agreement — Britain, France, Germany, Russia and China. Iran said no, indicating it wanted to know more about what was on the table.
Since then, the United States and Iran have issued sometimes contradictory, often intransigent statements that reflect mutual suspicion and agendas that are far broader than the simple reactivation of an agreement that many opponents of their efforts say was flawed to begin with.
This report is drawn from public pronouncements from Washington and Tehran and interviews with a half-dozen senior U.S. and European officials and with experts familiar with the issue. The officials spoke on the condition of anonymity about what one called the sensitive, and halting, diplomatic “dance.”
Iran wants all Trump sanctions lifted and an immediate influx of cash from the release of blocked international loans and frozen funds, along with foreign investment and removal of bans on oil sales. It seeks assurances that the next U.S. administration won’t jettison the deal again.
Even when the nuclear agreement was in force, Iran complained that U.S. threats limited foreign investment.
`
