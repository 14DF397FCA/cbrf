package metal

import (
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

func Do(r *http.Request) interface{} {
	url := makeURL(r)
	return common.GetRates(url, Metals{})
}

func makeURL(r *http.Request) string {
	var dateReq1, dateReq2 string
	if dateReq1 = r.FormValue("date_req1"); dateReq1 == "" {
		dateReq1 = common.GetToday()
	}
	if dateReq2 = r.FormValue("date_req2"); dateReq2 == "" {
		dateReq2 = common.GetToday()
	}
	url := fmt.Sprintf("https://www.cbr.ru/scripts/xml_metall.asp?date_req1=%s&date_req2=%s", dateReq1, dateReq2)
	log.Printf("Metal URL: %s", url)
	return url
}
