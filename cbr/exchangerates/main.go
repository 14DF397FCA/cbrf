package exchangerates

import (
	"cbrf/common"
	"encoding/xml"
	"fmt"
	"net/http"
)

type Currency struct {
	Text     string `xml:",chardata" json:"Text"`
	ID       string `xml:"ID,attr" json:"ID"`
	NumCode  string `xml:"NumCode" json:"NumCode"`
	CharCode string `xml:"CharCode" json:"CharCode"`
	Nominal  string `xml:"Nominal" json:"Nominal"`
	Name     string `xml:"Name" json:"Name"`
	Value    string `xml:"Value" json:"Value"`
}
type ExchangeRates struct {
	XMLName    xml.Name   `xml:"ValCurs" json:"-"`
	Text       string     `xml:",chardata" json:"Text"`
	Date       string     `xml:"Date,attr" json:"Date"`
	Name       string     `xml:"name,attr" json:"Name"`
	Currencies []Currency `xml:"Valute" json:"Valute"`
}

func Do(r *http.Request) interface{} {
	url := makeUrl(r)
	return common.GetRates(url, ExchangeRates{})
}

func getUrl(r *http.Request) string {
	if l := r.FormValue("lang"); l == "eng" {
		return fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily_%s.asp", l)
	} else {
		return "https://www.cbr.ru/scripts/XML_daily.asp"
	}
}

func getDate(r *http.Request) string {
	if d := r.FormValue("date_req"); len(d) > 0 {
		return fmt.Sprintf("?date_req=%s", d)
	}
	return ""
}

func makeUrl(r *http.Request) string {
	return fmt.Sprintf("%s%s", getUrl(r), getDate(r))
}
