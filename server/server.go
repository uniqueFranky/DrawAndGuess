package server

import (
	"encoding/json"
	"google/uuid"
	"gorilla/mux"
	"net/http"
)

func NewServer() *Server {
	server := Server{
		Router:  mux.NewRouter(),
		userSet: UserSet{users: []*User{}},
		gameSet: GameSet{games: []*Game{}},
	}
	server.routes()
	return &server
}

func (s *Server) routes() {
	//s.HandleFunc("/users", s.listOnlineUsers()).Methods("GET")      //Added
	//s.HandleFunc("/users/{name}", s.newUserLogin()).Methods("POST") //Added
	//s.HandleFunc("/users/{name}", s.getUser()).Methods("GET")
	//s.HandleFunc("/users", s.userLogout()).Methods("DELETE")                          //Added
	//s.HandleFunc("/games", s.listGames()).Methods("GET")                              //Added
	//s.HandleFunc("/games/{gameId}/players", s.listPlayersInGame()).Methods("GET")     //Added
	//s.HandleFunc("/games/{gameId}/lines", s.getLinesInGame()).Methods("GET")          //Added
	//s.HandleFunc("/games/{gameId}/lines", s.appendNewLineInGame()).Methods("POST")    //Added
	//s.HandleFunc("/games/{gameId}/join", s.userJoinGame()).Methods("POST")            //Added
	//s.HandleFunc("/games/{gameId}/messages", s.listMessagesInGame()).Methods("GET")   //Added
	//s.HandleFunc("/games/{gameId}/messages", s.appendMessageInGame()).Methods("POST") //Added
	//s.HandleFunc("/games/create/{answer}", s.newGame()).Methods("POST")               //Added

	s.HandleFunc("/users/reg/{name}/{psw}", s.newUserRegister()).Methods("POST")
	s.HandleFunc("/users/login/{name}/{psw}", s.newUserLogin()).Methods("POST")
	s.HandleFunc("/users/logout/{name}/{id}", s.userLogout()).Methods("DELETE")
	s.HandleFunc("/users/list", s.listOnlineUsers()).Methods("GET")
}

func (s *Server) ListenAndServe(port string) {
	err := http.ListenAndServe(port, s)
	if err != nil {
		return
	}
}

func (s *Server) newUserRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		psw := mux.Vars(r)["psw"]
		_, err := s.userSet.userReg(name, psw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", " application/json")
		u := User{
			UserName: name,
			UserId:   uuid.Nil,
		}
		if err = json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) newUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		psw := mux.Vars(r)["psw"]
		id, err := s.userSet.userLogin(name, psw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		u, err := s.userSet.findUserById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listOnlineUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.userSet.getUserNames()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) userLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		idStr := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := s.userSet.findUserById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if u.UserName != name {
			http.Error(w, "username and uuid dismatch", http.StatusBadRequest)
			return
		}
		err = s.userSet.deleteUserById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["name"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := s.userSet.findUserById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (s *Server) listGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.gameSet.games); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listPlayersInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.gameSet.findGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g.PlayerSet.users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getLinesInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		g, err := s.gameSet.findGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g.Lines); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) appendNewLineInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var newLine Line
		if err = json.NewDecoder(r.Body).Decode(&newLine); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.gameSet.findGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		g.appendLine(newLine)
		if err = json.NewEncoder(w).Encode(g.Lines); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) newGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ans := mux.Vars(r)["answer"]
		user, err := s.userSet.findUserById(u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		g := NewGame(user, ans)
		//user.GameId = g.Id
		if err = s.gameSet.appendGame(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) userJoinGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var u User
		if err = json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.gameSet.findGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := s.userSet.findUserById(u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = g.appendPlayer(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//user.GameId = gameUUID

		if err = json.NewEncoder(w).Encode(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listMessagesInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		g, err := s.gameSet.findGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g.Messages); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) appendMessageInGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameIdStr := mux.Vars(r)["gameId"]
		gameUUID, err := uuid.Parse(gameIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.gameSet.findGameById(gameUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var m Message
		err = json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := s.userSet.findUserById(m.From.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		m.From = u
		g.Messages = append(g.Messages, m)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(g.Messages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
