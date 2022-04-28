package currencyCode

import (
	"bytes"
	"cbrf/common"
	"encoding/xml"
	"log"
	"strings"
)

type Currency struct {
	Text        string `xml:",chardata" json:"Text"`
	ID          string `xml:"ID,attr" json:"ID"`
	Name        string `xml:"Name" json:"Name"`
	EngName     string `xml:"EngName" json:"EngName"`
	Nominal     string `xml:"Nominal" json:"Nominal"`
	ParentCode  string `xml:"ParentCode" json:"ParentCode"`
	ISONumCode  string `xml:"ISO_Num_Code" json:"ISONumCode"`
	ISOCharCode string `xml:"ISO_Char_Code" json:"ISOCharCode"`
}

type Currencies struct {
	XMLName xml.Name   `xml:"Valuta" json:"-"`
	Text    string     `xml:",chardata" json:"Text"`
	Name    string     `xml:"name,attr" json:"Name"`
	Item    []Currency `xml:"Item" json:"Item"`
}

func (cs *Currencies) SearchByID(d string) Currency {
	for _, c := range cs.Item {
		if c.ID == strings.ToUpper(d) {
			return c
		}
	}
	return Currency{}
}

func (cs *Currencies) SearchByISONum(d string) Currency {
	for _, c := range cs.Item {
		if c.ISONumCode == d {
			return c
		}
	}
	return Currency{}
}

func (cs *Currencies) SearchByISOCharCode(d string) Currency {
	for _, c := range cs.Item {
		if c.ISOCharCode == strings.ToUpper(d) {
			return c
		}
	}
	return Currency{}
}

func DecodeRates(buf []byte) (Currencies, error) {
	rates := Currencies{}
	d := xml.NewDecoder(bytes.NewReader(buf))
	d.CharsetReader = common.Decode
	err := d.Decode(&rates)
	if err != nil {
		return Currencies{}, err
	}
	return rates, nil
}

func GetCurrencies(url string) Currencies {
	xmlData, err := common.GetXML(url)
	if err != nil {
		log.Println(err)
	}
	data, err := DecodeRates(xmlData)
	if err != nil {
		log.Println(err)
	}
	return data
}

func GetCurrenciesMonthly() Currencies {
	//	https://www.cbr.ru/development/SXML/ Example 1
	return GetCurrencies("https://www.cbr.ru/scripts/XML_valFull.asp?d=1")
}
func GetCurrenciesDaily() Currencies {
	//	https://www.cbr.ru/development/SXML/ Example 1
	return GetCurrencies("https://www.cbr.ru/scripts/XML_valFull.asp?d=0")
}

func MergeCurrencies(m Currencies, d Currencies) Currencies {
	for _, r := range d.Item {
		m.Item = append(m.Item, r)
	}
	return m
}
