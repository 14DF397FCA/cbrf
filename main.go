package main

import (
	"cbrf/cbrf/currency"
	"cbrf/cbrf/currencyCode"
	"cbrf/cbrf/dynamic"
	"cbrf/cbrf/metal"
	"cbrf/common"
	"log"
	"net/http"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	s := "CBRF to json/xml in UTF-8\n"
	w.Write([]byte(s))
}

func CBRF(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}

	CBRFResp := make([]byte, 0)
	switch r.URL.Path {
	case "/cbrf/json":
		w.Header().Set("Content-Type", "application/json")
		data := currency.GetRates(r)
		CBRFResp = common.ToJson(data)
	case "/cbrf/xml":
		w.Header().Set("Content-Type", "application/xml")
		data := currency.GetRates(r)
		CBRFResp = common.ToXML(data)
	case "/cbrf/metals/json":
		w.Header().Set("Content-Type", "application/json")
		data := metal.GetRates(r)
		CBRFResp = common.ToJson(data)
	case "/cbrf/metals/xml":
		w.Header().Set("Content-Type", "application/xml")
		data := metal.GetRates(r)
		CBRFResp = common.ToXML(data)
	case "/cbrf/dynamic/json":
		w.Header().Set("Content-Type", "application/json")
		data := dynamic.GetRates(r)
		CBRFResp = common.ToJson(data)
	case "/cbrf/dynamic/xml":
		w.Header().Set("Content-Type", "application/xml")
		data := dynamic.GetRates(r)
		CBRFResp = common.ToXML(data)
	}
	_, err := w.Write(CBRFResp)
	if err != nil {
		log.Println(err)
	}
}

func CurCode(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}

	allCurs := currencyCode.MergeCurrencies(currencyCode.GetCurrenciesMonthly(), currencyCode.GetCurrenciesDaily())

	CBRFResp := make([]byte, 0)

	var cur currencyCode.Currency
	if f := r.FormValue("id"); f != "" {
		cur = allCurs.SearchByID(f)
	} else if f := r.FormValue("isonum"); f != "" {
		cur = allCurs.SearchByISONum(f)
	} else if f := r.FormValue("isocode"); f != "" {
		cur = allCurs.SearchByISOCharCode(f)
	}

	if r.Form.Has("json") {
		w.Header().Set("Content-Type", "application/json")
		CBRFResp = common.ToJson(cur)
	} else {
		w.Header().Set("Content-Type", "application/xml")
		CBRFResp = common.ToXML(cur)
	}
	_, err := w.Write(CBRFResp)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexPage)
	mux.HandleFunc("/index", IndexPage)
	mux.HandleFunc("/cbrf/json", CBRF)
	mux.HandleFunc("/cbrf/xml", CBRF)
	mux.HandleFunc("/cbrf/metals/json", CBRF)
	mux.HandleFunc("/cbrf/metals/xml", CBRF)
	mux.HandleFunc("/cbrf/dynamic/json", CBRF)
	mux.HandleFunc("/cbrf/dynamic/xml", CBRF)
	mux.HandleFunc("/cbrf/currency", CurCode)

	address := "0.0.0.0:8000"
	serv := http.Server{
		Addr:    address,
		Handler: mux,
	}
	log.Printf("Start listening %s...", address)
	serv.ListenAndServe()
}
