package common

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"
)

func GetData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: %v", resp.StatusCode)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func Decode(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unknown charset: %s", charset)
	}
}

func GetToday() string {
	return time.Now().Format("02/01/2006")
}

func GetYesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("02/01/2006")
}

func ToXML(in interface{}) []byte {
	if data, err := xml.Marshal(in); err != nil {
		log.Println(err)
		return nil
	} else {
		return data
	}
}

func ToJson(in interface{}) []byte {
	if data, err := json.Marshal(in); err != nil {
		log.Println(err)
		return nil
	} else {
		return data
	}
}

func GetRates(url string, t interface{}) interface{} {
	log.Println("URL", url)

	dataBytes, err := GetData(url)
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	data, err := DecodeRates(dataBytes, t)
	if err != nil {
		log.Println(err)
	}
	return data
}

func DecodeRates(buf []byte, s interface{}) (interface{}, error) {
	out := reflect.New(reflect.TypeOf(s)).Interface()
	d := xml.NewDecoder(bytes.NewReader(buf))
	d.CharsetReader = Decode
	err := d.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
