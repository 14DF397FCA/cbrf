package exchangerates

import (
	CBRCurs "cbrf/cbr/currencies"
	"cbrf/common"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
)

type Currency struct {
	Set          string  `xml:"set" json:"set"`
	From         string  `xml:"from" json:"from"`
	To           string  `xml:"to" json:"to"`
	FromCharCode string  `xml:"from_char_code,omitempty" json:"from_char_code,omitempty"`
	ToCharCode   string  `xml:"to_char_code,omitempty" json:"to_char_code,omitempty"`
	FromNameENG  string  `xml:"from_name_eng,omitempty" json:"from_name_eng,omitempty"`
	ToNameENG    string  `xml:"to_name_eng,omitempty" json:"to_name_eng,omitempty"`
	Rate         float64 `xml:"rate" json:"rate"`
}

type ExchangeRate struct {
	XMLName xml.Name   `xml:"QIWIExchangeRates" json:"-"`
	Result  []Currency `xml:"item" json:"result"`
}

func (c *Currency) ToString() string {
	return fmt.Sprintf("%s: From %s (%s) to %s (%s) rate: %f", c.Set, c.FromCharCode, c.From, c.ToCharCode, c.To, c.Rate)
}

func Do(enrich bool, all CBRCurs.Currencies) interface{} {
	url := "https://edge.qiwi.com/sinap/crossRates"
	log.Println("URL", url)

	bytes, err := common.GetData(url)
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	data := ExchangeRate{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Printf("Error unmarshal JSON: %s", err)
	}

	if enrich {
		data = Enrich(data, all)
	}

	return data
}

func Enrich(in ExchangeRate, all CBRCurs.Currencies) ExchangeRate {
	var curFrom, curTo CBRCurs.Currency
	var out ExchangeRate
	for _, i := range in.Result {
		curFrom = getNewCur(i.From, curFrom, all)
		curTo = getNewCur(i.To, curTo, all)
		out.Result = append(out.Result, Currency{
			Set:          i.Set,
			From:         i.From,
			FromCharCode: curFrom.ISOCharCode,
			FromNameENG:  curFrom.EngName,
			To:           i.To,
			ToCharCode:   curTo.ISOCharCode,
			ToNameENG:    curTo.EngName,
			Rate:         i.Rate,
		})
	}
	return out
}

func getNewCur(in string, prev CBRCurs.Currency, all CBRCurs.Currencies) CBRCurs.Currency {
	if in != prev.ISONumCode {
		ok, c := all.SearchByISONum(in)
		if !ok {
			log.Printf("Currency with ISO code %v not found", in)
		}
		return c
	}
	return prev
}
