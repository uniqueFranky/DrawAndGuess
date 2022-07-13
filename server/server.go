package server

import (
	"DrawAndGuess/game"
	"DrawAndGuess/user"
	"gorilla/mux"
	"net/http"
)

type Server struct {
	*mux.Router
	allGameSet game.AllGameSet
	userSet    user.UserSet
}

func NewServer() *Server {
	server := Server{
		Router:     mux.NewRouter(),
		userSet:    user.UserSet{Users: []*user.User{}},
		allGameSet: game.NewAllGameSet(),
	}
	server.routes()
	return &server
}

func (s *Server) routes() {
	//s.HandleFunc("/users", s.listOnlineUsers()).Methods("GET")      //Added
	//s.HandleFunc("/users/{name}", s.newUserLogin()).Methods("POST") //Added
	//s.HandleFunc("/users/{name}", s.getUser()).Methods("GET")
	//s.HandleFunc("/users", s.userLogout()).Methods("DELETE")                          //Added
	s.HandleFunc("/games", s.listGames()).Methods("GET")                           //Added
	s.HandleFunc("/games/{gameId}/players", s.listPlayersInGame()).Methods("GET")  //Added
	s.HandleFunc("/games/{gameId}/lines", s.getLinesInGame()).Methods("GET")       //Added
	s.HandleFunc("/games/{gameId}/lines", s.appendNewLineInGame()).Methods("POST") //Added
	s.HandleFunc("/games/{gameId}/lines", s.setLinesInGame()).Methods("PUT")       //Added
	s.HandleFunc("/games/{gameId}/join", s.userJoinGame()).Methods("POST")         //Added
	s.HandleFunc("/games/{gameId}/leave/{name}/{userId}", s.userLeaveGame()).Methods("DELETE")
	s.HandleFunc("/games/{gameId}/ended", s.getEndedGame()).Methods("GET")
	s.HandleFunc("/games/{gameId}", s.getGame()).Methods("GET")
	s.HandleFunc("/games/{gameId}/messages", s.listMessagesInGame()).Methods("GET")   //Added
	s.HandleFunc("/games/{gameId}/messages", s.appendMessageInGame()).Methods("POST") //Added
	s.HandleFunc("/games/create/{answer}", s.newGame()).Methods("POST")               //Added

	s.HandleFunc("/users/reg/{name}/{psw}", s.newUserRegister()).Methods("POST")
	s.HandleFunc("/users/login/{name}/{psw}", s.newUserLogin()).Methods("POST")
	s.HandleFunc("/users/logout/{name}/{id}", s.userLogout()).Methods("DELETE")
	s.HandleFunc("/users/list", s.listOnlineUsers()).Methods("GET")
	s.HandleFunc("/users/hasreg/{name}", s.hasReg()).Methods("GET")

	s.HandleFunc("/vocabs", s.listVocabs()).Methods("GET")
}

func (s *Server) ListenAndServe(port string) {
	err := http.ListenAndServe(port, s)
	if err != nil {
		return
	}
}
