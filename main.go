package main

import (
	CBRCurs "cbrf/cbr/currencies"
	"cbrf/cbr/dynamic"
	CBRRates "cbrf/cbr/exchangerates"
	"cbrf/cbr/metal"
	"cbrf/common"
	QIWIRates "cbrf/qiwi/exchangerates"
	"log"
	"net/http"
)

var AllCBRCurrencies = CBRCurs.GetAllCurrencies()

func IndexPage(w http.ResponseWriter, r *http.Request) {
	s := "Rates to json/xml in UTF-8\n"
	w.Write([]byte(s))
}

func Rates(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}

	resp := make([]byte, 0)
	switch r.URL.Path {
	case "/cbr/json":
		w.Header().Set("Content-Type", "application/json")
		data := CBRRates.Do(r)
		resp = common.ToJson(data)
	case "/cbr/xml":
		w.Header().Set("Content-Type", "application/xml")
		data := CBRRates.Do(r)
		resp = common.ToXML(data)
	case "/cbr/metals/json":
		w.Header().Set("Content-Type", "application/json")
		data := metal.Do(r)
		resp = common.ToJson(data)
	case "/cbr/metals/xml":
		w.Header().Set("Content-Type", "application/xml")
		data := metal.Do(r)
		resp = common.ToXML(data)
	case "/cbr/dynamic/json":
		w.Header().Set("Content-Type", "application/json")
		data := dynamic.Do(r)
		resp = common.ToJson(data)
	case "/cbr/dynamic/xml":
		w.Header().Set("Content-Type", "application/xml")
		data := dynamic.Do(r)
		resp = common.ToXML(data)
	case "/qiwi/json":
		w.Header().Set("Content-Type", "application/json")
		data := QIWIRates.Do(r.Form.Has("rich"), AllCBRCurrencies)
		resp = common.ToJson(data)
	case "/qiwi/xml":
		w.Header().Set("Content-Type", "application/xml")
		data := QIWIRates.Do(r.Form.Has("rich"), AllCBRCurrencies)
		resp = common.ToXML(data)
	}
	_, err := w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

func CurCode(w http.ResponseWriter, r *http.Request) {
	//if err := r.ParseForm(); err != nil {
	//	log.Println(err)
	//}
	//
	//allCurs := currencies.MergeCurrencies(currencies.GetCurrenciesMonthly(), currencies.GetCurrenciesDaily())
	//
	//CBRFResp := make([]byte, 0)
	//
	//var cur currencies.Currency
	//if f := r.FormValue("id"); f != "" {
	//	cur = allCurs.SearchByID(f)
	//} else if f := r.FormValue("isonum"); f != "" {
	//	cur = allCurs.SearchByISONum(f)
	//} else if f := r.FormValue("isocode"); f != "" {
	//	cur = allCurs.SearchByISOCharCode(f)
	//}
	//
	//if r.Form.Has("json") {
	//	w.Header().Set("Content-Type", "application/json")
	//	CBRFResp = common.ToJson(cur)
	//} else {
	//	w.Header().Set("Content-Type", "application/xml")
	//	CBRFResp = common.ToXML(cur)
	//}
	//_, err := w.Write(CBRFResp)
	//if err != nil {
	//	log.Println(err)
	//}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexPage)
	mux.HandleFunc("/index", IndexPage)
	mux.HandleFunc("/cbr/json", Rates)
	mux.HandleFunc("/cbr/xml", Rates)
	mux.HandleFunc("/cbr/metals/json", Rates)
	mux.HandleFunc("/cbr/metals/xml", Rates)
	mux.HandleFunc("/cbr/dynamic/json", Rates)
	mux.HandleFunc("/cbr/dynamic/xml", Rates)
	mux.HandleFunc("/qiwi/json", Rates)
	mux.HandleFunc("/qiwi/xml", Rates)
	mux.HandleFunc("/currency/json", CurCode)

	address := "0.0.0.0:8000"
	serv := http.Server{
		Addr:    address,
		Handler: mux,
	}
	log.Printf("Start listening %s...", address)
	serv.ListenAndServe()
}
