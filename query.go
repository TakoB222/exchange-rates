package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	value atomic.Value
	rate []Rate
	response []*respJSON
)


func query(){
	//ticker := time.NewTicker(5 * time.Minute)
	//for _ = range ticker.C {
	for {
		resp, err := http.Get("https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5")
		if err != nil {
			fmt.Println("http get error")
			os.Exit(1)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&response); err != nil {
			fmt.Println("error with decode")
			os.Exit(1)
		}

		for _, obj := range response {
			rate = append(rate, Rate{ccy: obj.Ccy, base_ccy: obj.Base_ccy, buy: obj.Buy, sale: obj.Sale})
		}
		value.Store(rate)
		fmt.Printf("%+v\n", rate)

		time.Sleep(5 * time.Minute)
	}

	//}
}


