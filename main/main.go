package main

import (
	"DrawAndGuess/server"
)

func main() {
	svr := server.NewServer()
	svr.ListenAndServe(":8090")
}
