package main

import(
  "github.com/gorilla/websocket"
)

type FindHandler func(string) (Handler, bool)

type Message struct {
	Name string `json:"name"` //package we're providing special tagging for and name
	Data interface{} `json:"data"`
}

type Client struct{
  send chan Message
  socket *websocket.Conn
  findHandler FindHandler
}

func (client *Client) Read(){
  var message Message
  for{
    if err := client.socket.ReadJSON(&message); err != nil{
      break
    }
    if handler, found := client.findHandler(message.Name); found{
      handler(client, message.Data)
    }
  }
  client.socket.Close()
}

func (client *Client) Write(){ //capitalizing a func name makes the access modifier public
  for msg := range client.send{
    if err := client.socket.WriteJSON(msg); err != nil {
      break
    }
  }
  client.socket.Close()
}



//How to create objects in a non-OOP language, such as Golang
//It allows us to easily instantiate a new client in func main
func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client{
  return &Client{ //pointer returning newly instantiated client
    send: make(chan Message),
    socket: socket, //set socket-field to the past websocket
    findHandler: findHandler,
  }
}
