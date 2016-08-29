package main

import(
	"log"
	"net/http"
	r "gopkg.in/dancannon/gorethink.v2"
)

type Channel struct{
	Id string `json:"id" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}

type User struct{
  Id string `gorethink:"id,omitempty"`
  Name string `gorethink:"name"` //tags because that's what main.go expects
}

func main(){
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
		Database: "rtsupport",
	})

	if err!= nil{
		log.Panic(err.Error())
	}

	router := NewRouter(session)
	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subscribeChannel)
	router.Handle("channel unsubscribe", unsubscribeChannel)

	//TODO: user edit, user subscribe, user unsubscribe, message add, message subscribe, message unsubscribe
	router.Handle("user edit", editUser)
	router.Handle("user subscribe", subscribeUser)
	router.Handle("user unsubscribe", unsubscribeUser)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
