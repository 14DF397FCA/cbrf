package dynamic

import (
	"bytes"
	"cbrf/common"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type Currency struct {
	Text    string `xml:",chardata"`
	Date    string `xml:"Date,attr"`
	ID      string `xml:"Id,attr"`
	Nominal string `xml:"Nominal"`
	Value   string `xml:"Value"`
}

type Currencies struct {
	XMLName    xml.Name   `xml:"Currencies"`
	Text       string     `xml:",chardata"`
	ID         string     `xml:"ID,attr"`
	DateRange1 string     `xml:"DateRange1,attr"`
	DateRange2 string     `xml:"DateRange2,attr"`
	Name       string     `xml:"name,attr"`
	Record     []Currency `xml:"Record"`
}

func GetRates(r *http.Request) Currencies {
	url := makeURL(r)
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

func makeURL(r *http.Request) string {
	var dateReq1, dateReq2, cur string
	if dateReq1 = r.FormValue("date_req1"); dateReq1 == "" {
		dateReq1 = common.GetYesterday()
	}
	if dateReq2 = r.FormValue("date_req2"); dateReq2 == "" {
		dateReq2 = common.GetToday()
	}
	if cur = r.FormValue("VAL_NM_RQ"); cur == "" {
		cur = "R01010"
	}
	log.Println("date_req1:", dateReq1)
	log.Println("date_req2:", dateReq2)
	log.Println("VAL_NM_RQ:", cur)

	return fmt.Sprintf("https://www.cbr.ru/scripts/XML_dynamic.asp?date_req1=%s&date_req2=%s&VAL_NM_RQ=%s", dateReq1, dateReq2, cur)
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
