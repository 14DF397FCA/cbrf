package dynamic

import (
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
	XMLName    xml.Name   `xml:"ValCurs"`
	Text       string     `xml:",chardata"`
	ID         string     `xml:"ID,attr"`
	DateRange1 string     `xml:"DateRange1,attr"`
	DateRange2 string     `xml:"DateRange2,attr"`
	Name       string     `xml:"name,attr"`
	Record     []Currency `xml:"Record"`
}

func GetRates(r *http.Request) interface{} {
	url := makeURL(r)
	log.Println("URL", url)

	xmlBytes, err := common.GetXML(url)
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	//data, err := DecodeRates(xmlBytes)
	data, err := common.DecodeRates(xmlBytes, &Currencies{})
	if err != nil {
		log.Println("Error in decode:", err)
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
