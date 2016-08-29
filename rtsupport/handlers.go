package main

import(
  "github.com/mitchellh/mapstructure"
  r "gopkg.in/dancannon/gorethink.v2"
  "fmt"
)

const (
  ChannelStop = iota //automatically assigns const values to 0,1,...,n
  UserStop
  MessageStop
)

func addChannel(client *Client, data interface{}){
  var channel Channel
  err := mapstructure.Decode(data, &channel)
  if err != nil{
    client.send <- Message{"error", err.Error()}
    return
  }
  go func(){ //Own goroutine in order to run independtly from the Client.Read() goroutine, because the addchannel func has a slow IO process
    err = r.Table("channel").
      Insert(channel).
      Exec(client.session)
    if err != nil{
      client.send <- Message{"error", err.Error()}
    }
  }()
}

func subscribeChannel(client *Client, data interface{}){
  stop := client.NewStopChannel(ChannelStop)
  result := make(chan r.ChangeResponse)
  cursor, err := r.Table("channel").
    Changes(r.ChangesOpts{IncludeInitial: true}). //returns existing channel records
    Run(client.session)
  if err != nil{
    client.send <- Message{"error", err.Error()}
    return
  }
  go func(){
    var change r.ChangeResponse
    for cursor.Next(&change){
      result <- change
    }
  }()
  changeFeedHelper(cursor, "channel", client.send, stop)
}

func unsubscribeChannel(client *Client, data interface{}){
  client.StopForKey(ChannelStop)
}

func editUser(client *Client, data interface{}){
  var user User
  err := mapstructure.Decode(data, &user)
  if err != nil{
    client.send <- Message{"error", err.Error()}
    return
  }
  client.userName = user.Name //variable used to pass as the author's name
  go func(){
    _, err := r.Table("user").
      Get(client.id).
      Update(user).
      RunWrite(client.session)
    if err != nil{
      client.send <- Message{"error", err.Error()}
    }
  }()
}

func subscribeUser(client *Client, data interface{}){
  go func(){
    stop := client.NewStopChannel(UserStop)
    cursor, err := r.Table("user").
      Changes(r.ChangesOpts{IncludeInitial: true}). //returns existing channel records
      Run(client.session)
    if err != nil{
      client.send <- Message{"error", err.Error()}
      return
    }
    changeFeedHelper(cursor, "user", client.send, stop)
  }()
}

func unsubscribeUser(client *Client, data interface{}){
  client.StopForKey(UserStop)
}

func changeFeedHelper(cursor *r.Cursor, changeEventName string,
  send chan<- Message, stop <-chan bool){
  change := make make (chan r.ChangeResponse)
  cursor.Listen(change)
  for{
    eventName := ""
    var data interface{}
    select{
    case <-stop:
      cursor.Close()
      return
    case val := <-change:
      if val.NewValue != nil && val.OldValue == nil{
        eventName = changeEventName + " add"
        data = val.NewValue
      }else if val.NewValue == nil && val.OldValue != nil{
        eventName = changeEventName + " remove"
        data = val.OldValue
      }else if val.NewValue != nil && val.OldValue != nil{
        eventName = changeEventName + " edit"
        data = val.NewValue
      }
      send <- Message{eventName, data}
    }
  }
}
