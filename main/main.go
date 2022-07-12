package main

import (
	"DrawAndGuess/server"
	"DrawAndGuess/storage"
	_ "mysql"
)

func main() {
	svr := server.NewServer()
	err := storage.Init()
	if err != nil {
		panic(err)
	}
	svr.ListenAndServe(":8090")
}
