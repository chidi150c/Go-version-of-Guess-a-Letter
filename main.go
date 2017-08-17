package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"user-apiv2/apiapp"
	"user-apiv2/apigame"
	"user-apiv2/apiuser"
	"user-apiv2/data"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8006", "used to chose listening address port")
	flag.Parse()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//[Start User Session]
	//Implemented User in apiuser and currently the only option
	var dbUser = make(data.DBType)
	var dbGame = make(apigame.GDBType)

	gmg := apigame.NewGameGuiService()

	us := apiuser.NewSession(dbUser)
	gs := apigame.NewSession(dbGame, us, gmg)

	gh := apigame.NewGameHandler(gs)
	uh := apiuser.NewUserHandler(us)
	_ = uh.Session.Userservice.AddUser(&data.User{Username: "chidi", Password: "cc", Level: "Admin"})
	//initial the handler
	h := apiapp.NewHandler(uh, gh)
	//open apiuser server
	server := apiapp.NewServer(addr, h)
	//defer apiuser server close
	//defer server.Close()
	if err := server.Open(done, sigs); err != nil {
		log.Fatal(err)
	}

	//    fresh      gukjj\
	fmt.Println("Listening on: ", server.Port())
	<-done
	fmt.Println("exiting")
}
