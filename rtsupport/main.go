package main

import(
	"fmt"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Message struct {
	Name string `json:"name"` //package we're providing special tagging for and name
	Data interface{} `json:"data"`
}

type Channel struct{
	Id string `json:"id"`
	Name string `json:"name"`
}

func main(){
	router := NewRoute()
	router.Handle("channel add", addChannel)
	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
