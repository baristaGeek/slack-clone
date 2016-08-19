package main

import(
  "fmt"
  "time"
  "math/rand"
)

type Message struct {
	Name string `json:"name"` //package we're providing special tagging for and name
	Data interface{} `json:"data"`
}

type Client struct{
  send chan Message
}

func (client *Client) write(){
  for msg := range client.send{
    //TODO: socket.sendJSON(msg)
    fmt.Printf("%#v\n", msg)
  }
}

func (client *Client) subscribeChannels(){
  //TODO: Changefeed RethinkDB
  for {
    time.Sleep(r())
    client.send <- Message{"channel add", ""} //data doesn't need a realy payload for now, so let's just set it to an empty stringc
  }
}

func (client *Client) subscribeMessages(){
  //TODO: Changefeed RethinkDB
  for {
    time.Sleep(r())
    client.send <- Message{"message add", ""}
  }
}

func r() time.Duration{
    return time.Millisecond * time.Duration(rand.Intn(1000))
}

//How to create objects in a non-OOP language, such as Golang
//It allows us to easily instantiate a new client in func main
func NewClient() *Client{
  return &Client{ //pointer returning newly instantiated client
    send: make(chan Message),
  }
}

func main(){
  client := NewClient()
  go client.subscribeChannels()
  go client.subscribeMessages()
  client.write()
}
