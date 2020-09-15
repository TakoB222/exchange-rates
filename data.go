package main

type respJSON struct{
	Ccy string `json:"ccy"`
	Base_ccy string `json:"base_ccy"`
	Buy string `json:"buy"`
	Sale string `json:"sale"`
}
type Rate struct{
	ccy string
	base_ccy string
	buy string
	sale string
}




