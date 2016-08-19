package main

import(
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type Message struct {
	Name string `json:"name"` //package we're providing special tagging for and name
	Data interface{} `json:"data"`
}

/*
type Speaker interface{
	Speak() //structs can have methods that have access to its fields
}

func (m Message) Speak(){
	fmt.Println("I'm a " + m.name + "event!")
}

func someFunc(speaker Speaker){
	speaker.Speak()
}
*/

type Channel struct{
	Id string `json:"id"`
	Name string `json:"name"`
}

func main(){
	recRawMsg := []byte(`{"name":"channel add",` +
	`"data":{"name":"Hardware Support"}}`)

	var recMessage Message
	err := json.Unmarshal(recRawMsg, &recMessage) //&returns a pointer to the value, so that the value can be modified in memmory
	if err != nil{
		fmt.Println(err)
		return 
	}

	fmt.Printf("%#v\n", recMessage) //print values in Go syntax

	if recMessage.Name == "channel add"{
		channel, err := addChannel(recMessage.Data)
		var sendMessage Message
		sendMessage.Name = "channel add"
		sendMessage.Data = channel 
		sendRawMsg, err :=json.Marshal(sendMessage) //returns the encoded JSON as a byte array
		if err != nil{
			fmt.Println(err)
			return
		}
		fmt.Println(string(sendRawMsg))
	}
}




