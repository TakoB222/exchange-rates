package main

import (
	"./models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	value atomic.Value
)

func GetExchangeRatesFromPrivatBankAsJSON(){
	ticker := time.NewTicker(5 * time.Minute)
	for _ = range ticker.C {
		var (
			rate []models.Rate
			respJSON []*models.RespJSON
		)
		resp, err := http.Get("https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5")
		if err != nil {
			fmt.Println("http get error:", err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil{
				fmt.Println("error with decode:", err)
				continue
			}
		}
		for _, obj := range respJSON {
			rate = append(rate, models.Rate{Сcy: obj.Ccy, Base_ccy: obj.Base_ccy, Buy: obj.Buy, Sale: obj.Sale})
		}
		value.Store(rate)
		fmt.Printf("%+v\n", rate)
	}
}


