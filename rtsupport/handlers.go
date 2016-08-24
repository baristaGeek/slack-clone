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

  go func(){
    for{
      select{
      case <-stop:
        cursor.Close()
        return
      case change := <-result:
        if change.NewValue != nil && change.OldValue == nil{ //new record has been added
          client.send <- Message{"channel add", change.NewValue}
          fmt.Println("sent channel add msg")
        }
      }
    }
  }()
}

func unsubscribeChannel(client *Client, data interface{}){
  client.StopForKey(ChannelStop)
}
