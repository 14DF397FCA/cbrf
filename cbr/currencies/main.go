package currencies

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

func (cs *Currencies) SearchByID(d string) (bool, Currency) {
	log.Printf("Search by ID: %s", d)
	for _, c := range cs.Item {
		if c.ID == strings.ToUpper(d) {
			return true, c
		}
	}
	return false, Currency{}
}

func (cs *Currencies) SearchByISONum(d string) (bool, Currency) {
	//	true if currency found
	//log.Printf("Search by ISO Num: %s", d)
	for _, c := range cs.Item {
		if c.ISONumCode == d {
			return true, c
		}
	}
	return false, Currency{}
}

func (cs *Currencies) SearchByISOCharCode(d string) (bool, Currency) {
	log.Printf("Search by ISO Char code: %s", d)
	for _, c := range cs.Item {
		if c.ISOCharCode == strings.ToUpper(d) {
			return true, c
		}
	}
	return false, Currency{}
}

func DecodeRates(buf []byte) (Currencies, error) {
	out := Currencies{}
	d := xml.NewDecoder(bytes.NewReader(buf))
	d.CharsetReader = common.Decode
	err := d.Decode(&out)
	if err != nil {
		return Currencies{}, err
	}
	return out, nil
}

func getCurrencies(url string) Currencies {
	dataBytes, err := common.GetData(url)
	if err != nil {
		log.Println(err)
	}
	//data, err := common.DecodeRates(dataBytes, Currencies{})
	data, err := DecodeRates(dataBytes)
	if err != nil {
		log.Println(err)
	}
	return data
}

func getCurrenciesMonthly() Currencies {
	//	https://www.cbr.ru/development/SXML/ Example 1
	return getCurrencies("https://www.cbr.ru/scripts/XML_valFull.asp?d=1")
}
func getCurrenciesDaily() Currencies {
	//	https://www.cbr.ru/development/SXML/ Example 1
	return getCurrencies("https://www.cbr.ru/scripts/XML_valFull.asp?d=0")
}

func mergeCurrencies(m Currencies, d Currencies) Currencies {
	for _, r := range d.Item {
		m.Item = append(m.Item, r)
	}
	return m
}

func GetAllCurrencies() Currencies {
	curs := mergeCurrencies(getCurrenciesMonthly(), getCurrenciesDaily())
	return curs
}
