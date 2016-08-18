package main

import(
	"fmt"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

func main(){
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}

func handler(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hello from go")
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		fmt.Println(err)
		return
	}
	for{
		/*msgType, msg, err := socket.ReadMessage()
		if err != nil{
			fmt.Println(err)
			return
		}*/
		var inMessage Message
		var outMessage Message
		if err := socket.ReadJSON(&inMessage); err != nil{
			fmt.Println(err)
			break
		}
		fmt.Printf("%#v\n", inMessage)
		switch inMessage.Name{
		case "channel add":
			err := addChannel(inMessage.Data)
			if err != nil{
				outMessage = Message{"error", err}
				if err := socket.WriteJSON(outMessage); err != nil{
					fmt.Println(err)
					break
				}
			}
		case "channel subscribe":
			go subscribeChannel(socket)
		}
		/*fmt.Println(string(msg))
		if err = socket.WriteMessage(msgType, msg); err != nil{
			fmt.Println(err)
			return
		}*/
	}
}

func addChannel(data interface{}) error{ //Don't return channel because only people subscribed to a certain channel want to see the messages in that channel
	var channel Channel 
	/* manual type assertion
	channelMap := data.(map[]string interface{})
	channel.Name = channelMap["name"].(string)
	*/
	err := mapstructure.Decode(data, &channel)
	if err != nil{
		fmt.Println(err)
		return err
	}

	channel.Id = "1" //naturally set by RethinkDB but hardcoded for the moment
	fmt.Println("added channel")
	return nil 
}

func subscribeChannel(socket *websocket.Conn){
	//TODO: Query RethinkDB / changefeed
	for{
		time.Sleep(time.Second * 1)
		message := Message{"channel add",
			Channel{"1", "Software Support"}}
		socket.WriteJSON(message)
		fmt.Println("sent new channel")
	}
}






