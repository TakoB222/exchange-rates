package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

func main(){

	go query()
	http.HandleFunc("/rate/", handle)
	fmt.Println("Listening on port 8080...")

	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request){
	fmt.Println("iam work")
	format := strings.TrimPrefix(r.URL.String(), "/rate/")
	fmt.Println(format)
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(value.Load())
		json.NewEncoder(w).Encode(string(body))
		break
	case "xml":
		w.Header().Set("Content-Type", "application/xml")
		body, _ := xml.Marshal(value.Load())
		xml.NewEncoder(w).Encode(string(body))
		break
	case "text":
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, value.Load())
		break
	default:
		fmt.Fprint(w, "wrong format")
	}
}
