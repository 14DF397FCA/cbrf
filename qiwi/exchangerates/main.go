package exchangerates

import (
	"cbrf/common"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
)

type Currency struct {
	Set  string  `xml:"set" json:"set"`
	From string  `xml:"from" json:"from"`
	To   string  `xml:"to" json:"to"`
	Rate float64 `xml:"rate" json:"rate"`
}

type ExchangeRate struct {
	XMLName xml.Name   `xml:"QIWIExchangeRates" json:"-"`
	Result  []Currency `xml:"item" json:"result"`
}

func Do(r *http.Request) interface{} {
	url := "https://edge.qiwi.com/sinap/crossRates"
	log.Println("URL", url)

	data, err := common.GetData(url)
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	t := ExchangeRate{}
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Printf("Error unmarshal JSON: %s", err)
	}
	return t
}
