package server

import (
	"encoding/json"
	"fmt"
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
	s.HandleFunc("/users", s.listOnlineUsers()).Methods("GET")                     //Added
	s.HandleFunc("/users/{name}", s.newUserLogin()).Methods("POST")                //Added
	s.HandleFunc("/users", s.userLogout()).Methods("DELETE")                       //Added
	s.HandleFunc("/games", s.listGames()).Methods("GET")                           //Added
	s.HandleFunc("/games/{gameId}/players", s.listPlayersInGame()).Methods("GET")  //Added
	s.HandleFunc("/games/{gameId}/lines", s.getLinesInGame()).Methods("GET")       //Added
	s.HandleFunc("/games/{gameId}/lines", s.appendNewLineInGame()).Methods("POST") //Added
	s.HandleFunc("/games/{gameId}/join", s.userJoinGame()).Methods("POST")         //Added
	s.HandleFunc("/games/{gameId}/messages", s.listMessagesInGame()).Methods("GET")
	s.HandleFunc("/games/{gameId}/messages", s.appendMessageInGame()).Methods("POST")
	s.HandleFunc("/games/create", s.newGame()).Methods("POST") //Added
}

func (s *Server) ListenAndServe(port string) {
	err := http.ListenAndServe(port, s)
	if err != nil {
		return
	}
}

func (s *Server) listOnlineUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.userSet.users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) newUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["name"]
		u := &User{
			UserName: username,
			UserId:   uuid.New(),
		}
		if err := s.userSet.appendUser(u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		var batch GameCreationBatch
		if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := s.userSet.findUserById(batch.User.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		g := NewGame(user, batch.Answer)
		user.GameId = g.Id
		if err = s.gameSet.appendGame(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(string(g.Answer))
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(g); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) userLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err := s.userSet.findUserById(u.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.gameSet.deletePlayerInGame(&u)
		s.userSet.deleteUser(&u)

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(u); err != nil {
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
		user.GameId = gameUUID

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

		g.Messages = append(g.Messages, m)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(g.Messages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
