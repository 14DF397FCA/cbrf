package metal

import (
	"bytes"
	"cbrf/common"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type Metal struct {
	Text string `xml:",chardata" json:"Text"`
	Date string `xml:"Date,attr" json:"Date"`
	Code string `xml:"Code,attr" json:"Code"`
	Buy  string `xml:"Buy" json:"Buy"`
	Sell string `xml:"Sell" json:"Sell"`
}

type Metals struct {
	XMLName  xml.Name `xml:"Metall" json:"-"`
	Text     string   `xml:",chardata" json:"Text"`
	FromDate string   `xml:"FromDate,attr" json:"FromDate"`
	ToDate   string   `xml:"ToDate,attr" json:"ToDate"`
	Name     string   `xml:"name,attr" json:"Name"`
	Record   []Metal  `xml:"Record" json:"Record"`
}

func DecodeRates(buf []byte) (Metals, error) {
	out := Metals{}
	d := xml.NewDecoder(bytes.NewReader(buf))
	d.CharsetReader = common.Decode
	err := d.Decode(&out)
	if err != nil {
		return Metals{}, err
	}
	return out, nil
}

func GetRates(r *http.Request) Metals {
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
	var dateReq1, dateReq2 string
	if dateReq1 = r.FormValue("date_req1"); dateReq1 == "" {
		dateReq1 = common.GetToday()
	}
	if dateReq2 = r.FormValue("date_req2"); dateReq2 == "" {
		dateReq2 = common.GetToday()
	}
	log.Println("date_req1:", dateReq1)
	log.Println("date_req2:", dateReq2)

	return fmt.Sprintf("https://www.cbr.ru/scripts/xml_metall.asp?date_req1=%s&date_req2=%s", dateReq1, dateReq2)
}
