package currency

import (
	"bytes"
	"cbrf/common"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
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
	Currencies []Currency `xml:"Valute" json:"ValCurs"`
}

func (rates *ExchangeRates) ToJson() []byte {
	if data, err := json.Marshal(rates); err != nil {
		log.Println(err)
		return nil
	} else {
		return data
	}
}

func (rates *ExchangeRates) ToXML() []byte {
	if data, err := xml.Marshal(rates); err != nil {
		log.Println(err)
		return nil
	} else {
		return data
	}
}

func DecodeRates(buf []byte) (ExchangeRates, error) {
	rates := ExchangeRates{}
	d := xml.NewDecoder(bytes.NewReader(buf))
	d.CharsetReader = common.Decode

	err := d.Decode(&rates)
	if err != nil {
		return ExchangeRates{}, err
	}
	return rates, nil
}

func GetRates(r *http.Request) ExchangeRates {
	url := makeUrl(r)
	log.Println("URL", url)

	xmlBytes, err := common.GetXML(url)
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	data, err := DecodeRates(xmlBytes)
	if err != nil {
		log.Println(err)
	}
	return data
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
		return fmt.Sprintf("date_req=%s", d)
	}
	return ""
}

func makeUrl(r *http.Request) string {
	return fmt.Sprintf("%s?%s", getUrl(r), getDate(r))
}
